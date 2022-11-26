package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

func removeArticles(user tables.Users, articles []interface{}, codes ...string) []interface{} {
	if user.Lang == nil {
		en := "en"
		user.Lang = &en
	}
	trash := make([]int, 0, len(inlineCodes)+5)
	for i, v := range articles {
		article := v.(tgbotapi.InlineQueryResultArticle)
		for _, code := range codes {
			if strings.HasPrefix(article.Title, langs[*user.Lang][code]) {
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
		if strings.HasPrefix(article.Title, langs[*user.Lang][code]) {
			return v
		}
	}
	return nil
}

func removeIndex(obj []interface{}, idx int) []interface{} {
	return append(obj[:idx], obj[idx+1:]...)
}

func Ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

func (app App) onInlineQuery(ctx context.Context, update tgbotapi.InlineQuery) {
	update.Query = Ucfirst(strings.TrimSpace(norm.NFKC.String(update.Query)))
	log := app.log.With(zap.String("query", update.Query))

	defer func() {
		if err := recover(); err != nil {
			_, f, line, ok := runtime.Caller(2)
			if ok {
				log = log.With(zap.String("caller", f+":"+strconv.Itoa(line)))
			}
			if e, ok := err.(errors.Error); ok {
				log.Error("inline", zap.Error(e))
			} else {
				log.Error("", zap.Any("error", err))
			}
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
		}
	}()
	//go func() {
	//	if err := app.analytics.InlineUser(update); err != nil {
	//		app.notifyAdmin(err)
	//	}
	//}()

	warn := func(err error) {
		app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID:     update.ID,
			SwitchPMText:      "Error, sorry.",
			SwitchPMParameter: "from_inline",
		})
		app.notifyAdmin(err)
		log.Error("", zap.Error(err))
	}

	if update.Query == "" {
		if _, err := app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID:     update.ID,
			IsPersonal:        true,
			SwitchPMText:      tables.Users{Lang: &update.From.LanguageCode}.Localize("начните вводить текст"),
			SwitchPMParameter: "from_empty_inline",
		}); err != nil {
			log.Error("", zap.Error(err))
		}
		return
	}

	user := tables.Users{Lang: &update.From.LanguageCode}
	log = log.With(zap.String("language_code", *user.Lang))
	var offset int // смещение для пагинации
	if update.Offset != "" {
		var err error
		offset, err = strconv.Atoi(update.Offset)
		if err != nil {
			warn(err)
			return
		}
	}

	if offset > len(inlineCodes[*user.Lang]) {
		warn(fmt.Errorf("слишком большое смещение: %d", offset))
		return
	}

	count := 50
	if offset+count > len(inlineCodes[*user.Lang])-1 {
		count = len(inlineCodes[*user.Lang]) - offset
	}

	nextOffset := offset + count

	tr, err := app.translo.Translate(ctx, "auto", "en", update.Query)
	// TODO: detect lang via yandex if there are all emojis in message except spec. chars
	if err != nil {
		warn(err)
		return
	}
	from := strings.ToLower(tr.TextLang)

	g, _ := errgroup.WithContext(context.Background())
	var mu sync.Mutex

	g.Go(func() error {
		user, err = app.db.GetUserByID(update.From.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if user.Lang == nil || *user.Lang == "" {
			user.Lang = &update.From.LanguageCode
		} else {
			log = log.With(zap.String("user_lang", *user.Lang))
		}
		return nil
	})

	blocks := make([]interface{}, 0, 50)

	//needAudio := strings.HasPrefix(update.Query, "!")

	//cacheKey := []byte(fmt.Sprintf("%s;%s;%d-%d", *user.Lang, update.Query, offset, offset+count))
	//log.Debug("inline bitcask cache key", zap.String("cache_key", string(cacheKey)))
	//if len(cacheKey) < 64 && app.bc.Has(cacheKey) {
	//	cacheData, err := app.bc.Get(cacheKey)
	//	if err != nil {
	//		warn(err)
	//		return
	//	}
	//	if err = json.Unmarshal(cacheData, &blocks); err != nil {
	//		warn(err)
	//		return
	//	}
	//	log.Debug("took blocks from cache")
	//} else {
	for i, code := range inlineCodes[*user.Lang][offset : offset+count] {
		code := code
		i := i
		g.Go(func() error {
			title := langs[*user.Lang][code]
			//if code == "emj" || from == "emj" {
			//	translation, err = translate2.YandexTranslate(ctx, from, code, update.Query)
			//} else {
			tr, err := app.translo.Translate(ctx, from, code, update.Query)
			if err != nil {
				log.Error("inline", zap.Error(err))
				return nil
			}
			//}
			if err != nil {
				log.Error("inline", zap.Error(err))
				return nil
			}
			if tr.TranslatedText == "" {
				log.Error("empty translation in inline mode", zap.String("query", update.Query), zap.String("language_code", update.From.LanguageCode))
				return nil
			}

			btn := tgbotapi.InlineKeyboardButton{
				Text:                         tables.Users{Lang: &code}.Localize("translate"),
				URL:                          nil,
				LoginURL:                     nil,
				CallbackData:                 nil,
				WebApp:                       nil,
				SwitchInlineQuery:            nil,
				SwitchInlineQueryCurrentChat: &tr.TranslatedText,
				CallbackGame:                 nil,
				Pay:                          false,
			}
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(btn))

			mu.Lock()
			defer mu.Unlock()
			blocks = append(blocks, tgbotapi.InlineQueryResultArticle{
				Type:  "article",
				ID:    strconv.Itoa(i + offset),
				Title: title,
				InputMessageContent: map[string]interface{}{
					"message_text":             tr.TranslatedText,
					"disable_web_page_preview": true,
				},
				ReplyMarkup: &keyboard,
				URL:         "",
				HideURL:     true,
				Description: tr.TranslatedText,
				ThumbURL:    "",
				ThumbWidth:  0,
				ThumbHeight: 0,
			})
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		warn(err)
		log.Error("", zap.Error(err))
		return
	}
	//	go func() {
	//		cacheData, err := json.Marshal(blocks)
	//		if err != nil {
	//			app.notifyAdmin(err)
	//		}
	//		if err = app.bc.Put(cacheKey, cacheData); err != nil {
	//			app.notifyAdmin(err)
	//		}
	//	}()
	//}

	if offset == 0 && !in([]string{"", "auto"}, user.MyLang, user.ToLang) {
		for i, lang := range []string{user.MyLang, user.ToLang} {
			if lang == from {
				continue
			}
			block := getArticle(user, blocks, lang)
			if block == nil {
				continue
			}
			article := block.(tgbotapi.InlineQueryResultArticle)
			article.Title = strings.TrimSuffix(article.Title, " 📌") + " 🙍‍♂"
			article.ID = strconv.Itoa(-1 - i)
			blocks = append(blocks, article)
			nextOffset--
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

	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID:     update.ID,
		Results:           blocks,
		IsPersonal:        true,
		CacheTime:         0,
		NextOffset:        strconv.Itoa(nextOffset),
		SwitchPMText:      user.Localize("tap on translation to send it"),
		SwitchPMParameter: "from_inline",
	}
	if _, err = app.bot.AnswerInlineQuery(inlineConfig); err != nil {
		log.Error("", zap.Error(err))
	}

	if user.MyLang != "" { // user exists
		if err = app.db.UpdateUserActivity(update.From.ID); err != nil {
			warn(err)
		}
	}
	//if err = app.analytics.InlineBot(*update.From, inlineConfig); err != nil {
	//	app.notifyAdmin(err)
	//}
}
