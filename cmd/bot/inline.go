package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

func removeArticles(user tables.Users, articles []interface{}, codes ...string) []interface{} {

	trash := make([]int, 0, len(codes)+5)
	for i, v := range articles {
		article := v.(tgbotapi.InlineQueryResultArticle)
		for _, code := range codes {
			if strings.HasPrefix(article.Title, langs[user.Lang][code]) {
				trash = append(trash, i)
			}
		}
	}
	for m, i := range trash {
		articles = removeIndex(articles, i-m)
	}
	return articles
}

func getArticle(user tables.Users, articles []interface{}, code string) interface{} {
	for _, v := range articles {
		article := v.(tgbotapi.InlineQueryResultArticle)
		if strings.HasPrefix(article.Title, langs[user.Lang][code]) {
			return v
		}
	}
	return nil
}

func removeIndex(obj []interface{}, idx int) []interface{} {
	return append(obj[:idx], obj[idx+1:]...)
}

func Title(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (app App) onInlineQuery(ctx context.Context, update tgbotapi.InlineQuery) {
	defer func() {
		if err := recover(); err != nil {
			app.log.Error("%w", zap.Any("error", err))
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
		}
	}()

	update.Query = norm.NFKC.String(update.Query)
	warn := func(err error) {
		app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID:     update.ID,
			SwitchPMText:      "Error, sorry.",
			SwitchPMParameter: "from_inline",
		})
		app.notifyAdmin(err)
		pp.Println("onInlineQuery: error", err)
		pp.Println(string(debug.Stack()))
	}

	user := tables.Users{Lang: update.From.LanguageCode}

	if update.Query == "" {
		kbLoc := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(user.Localize("Попробовать"), "https://t.me/translobot?start=from_inline")))

		elemLoc := tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    "ad_local",
			Title: user.Localize("Порекомендовать бота"),
			InputMessageContent: map[string]interface{}{
				"message_text":             user.Localize("inline_ad"),
				"disable_web_page_preview": true,
				"parse_mode":               tgbotapi.ModeHTML,
			},
			ReplyMarkup: &kbLoc,
			URL:         "",
			HideURL:     true,
			Description: user.Localize("кликните, чтобы порекомендовать бота в чате"),
			ThumbURL:    "https://i.yapx.ru/PdNIa.png",
			ThumbWidth:  200,
			ThumbHeight: 200,
		}

		kbEn := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Try it out", "https://t.me/translobot?start=from_inline")))

		elemEn := tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    "ad_en",
			Title: "Recommend the bot",
			InputMessageContent: map[string]interface{}{
				"message_text":             `🔥 <a href="https://t.me/translobot">Translo</a> 🌐 - <i>The best Telegram translator bot in the whole world</i>`,
				"disable_web_page_preview": true,
				"parse_mode":               tgbotapi.ModeHTML,
			},
			ReplyMarkup: &kbEn,
			URL:         "",
			HideURL:     true,
			Description: "click to recommend a bot in the chat",
			ThumbURL:    "https://i.yapx.ru/PdNIa.png",
			ThumbWidth:  200,
			ThumbHeight: 200,
		}
		app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID:     update.ID,
			Results:           []interface{}{elemLoc, elemEn},
			CacheTime:         0,
			IsPersonal:        true,
			NextOffset:        "",
			SwitchPMText:      "Type text to translate",
			SwitchPMParameter: "from_inline",
		})
		return
	}

	var offset int // смещение для пагинации
	if update.Offset != "" {
		var err error
		offset, err = strconv.Atoi(update.Offset)
		if err != nil {
			warn(err)
			return
		}
	}

	if offset > len(codes[user.Lang])-1 {
		warn(fmt.Errorf("слишком большое смещение: %d", offset))
		return
	}

	count := 50
	if offset+count > len(codes[user.Lang])-1 {
		count = len(codes[user.Lang]) - 1 - offset
	}

	nextOffset := offset + count

	g, _ := errgroup.WithContext(context.Background())
	var mu sync.Mutex
	var err error

	g.Go(func() error {
		user, err = app.db.GetUserByID(update.From.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		user.SetLang(update.From.LanguageCode)
		return nil
	})

	blocks := make([]interface{}, 0, 50)

	//needAudio := strings.HasPrefix(update.Query, "!")

	from := ""
	for i, code := range codes[user.Lang][offset : offset+count] {
		code := code
		i := i
		g.Go(func() error {
			title := langs[user.Lang][code]
			if offset == 0 && i < 19 {
				title += " 📌"
			}

			tr, err := translate2.GoogleTranslate(ctx, "auto", code, update.Query)
			if err != nil || tr.Text == "" {
				return nil
			}
			if from == "" {
				from = tr.FromLang
			}

			//if needAudio {
			//	audio, err := translate2.TTS(code, tr.Text)
			//	if err != nil {
			//		return err
			//	}
			//	app.bot.UploadFiles()
			//	tgbotapi.InlineQueryResultAudio{}
			//}

			mu.Lock()
			defer mu.Unlock()
			blocks = append(blocks, tgbotapi.InlineQueryResultArticle{
				Type:  "article",
				ID:    strconv.Itoa(i + offset),
				Title: title,
				InputMessageContent: map[string]interface{}{
					"message_text":             tr.Text,
					"disable_web_page_preview": true,
				},
				ReplyMarkup: nil,
				URL:         "",
				HideURL:     true,
				Description: tr.Text,
				ThumbURL:    "",
				ThumbWidth:  0,
				ThumbHeight: 0,
			})
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		warn(err)
		app.log.Error("", zap.Error(err))
		return
	}

	blocks = removeArticles(user, blocks, from)

	if offset == 0 && user.Usings > 0 {
		nextOffset -= 2
		for i, lang := range []string{user.MyLang, user.ToLang} {
			if lang == from {
				continue
			}
			block := getArticle(user, blocks, lang)
			if block == nil {
				continue
			}
			article := block.(tgbotapi.InlineQueryResultArticle)
			blocks = removeArticles(user, blocks, lang)
			article.Title = strings.TrimSuffix(article.Title, " 📌")
			article.Title += " 🙍‍♂️"
			article.ID = strconv.Itoa(-1 - i)
			blocks = append(blocks, article)
		}
	}

	sort.Slice(blocks, func(i, j int) bool {
		block1 := blocks[i].(tgbotapi.InlineQueryResultArticle)

		id1, err := strconv.Atoi(block1.ID)
		if err != nil {
			warn(err)
			return false
		}
		id2, err := strconv.Atoi(blocks[j].(tgbotapi.InlineQueryResultArticle).ID)
		if err != nil {
			warn(err)
			return false
		}
		return id1 < id2
	})

	if len(blocks) > 50 {
		diff := len(blocks) - 50
		blocks = blocks[:len(blocks)-diff]
	}

	pmtext := "Translo"
	if update.Query == "" {
		pmtext = "Enter text"
	}

	if _, err := app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID:     update.ID,
		Results:           blocks,
		CacheTime:         0,
		NextOffset:        strconv.Itoa(nextOffset),
		IsPersonal:        true,
		SwitchPMText:      pmtext,
		SwitchPMParameter: "from_inline",
	}); err != nil {
		warn(errors.Wrap(err))
		pp.Println(blocks)
	}

	app.analytics.User(update.Query, update.From)
	app.analytics.Bot(update.From.ID, "Inline succeeded", "Inline succeeded")
	if user.MyLang != "" { // user exists
		if err = app.db.UpdateUserMetrics(update.From.ID, "inline:"+update.Query); err != nil {
			warn(err)
		}
	}
}
