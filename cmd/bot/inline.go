package bot

import (
	"fmt"
	"github.com/armanokka/translobot/internal/tables"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	"github.com/go-errors/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
)

func (app App) onInlineQuery(update tgbotapi.InlineQuery) {
	go 	app.analytics.User(update.Query, update.From)
	update.Query = strings.Title(update.Query)
	warn := func(err error) {
		err = errors.Wrap(err, 1)
		app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID:     update.ID,
			SwitchPMText:      "Error, sorry.",
			SwitchPMParameter: "from_inline",
		})
		app.notifyAdmin(err)
		pp.Println("onInlineQuery: error", err)
	}

	user := tables.Users{Lang: update.From.LanguageCode}

	if update.Query == "" {
		kbLoc := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(user.Localize("Попробовать"), "https://t.me/translobot?start=from_inline")))

		elemLoc := tgbotapi.InlineQueryResultArticle{
			Type:                "article",
			ID:                  "ad_local",
			Title:               user.Localize("Порекомендовать бота"),
			InputMessageContent: map[string]interface{}{
				"message_text": user.Localize("inline_ad"),
				"disable_web_page_preview":true,
				"parse_mode": tgbotapi.ModeHTML,
			},
			ReplyMarkup:         &kbLoc,
			URL:                 "",
			HideURL:             true,
			Description:         user.Localize("кликните, чтобы порекомендовать бота в чате"),
			ThumbURL:            "https://i.yapx.ru/PdNIa.png",
			ThumbWidth:          200,
			ThumbHeight:         200,
		}

		kbEn := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Try it out", "https://t.me/translobot?start=from_inline")))


		elemEn := tgbotapi.InlineQueryResultArticle{
			Type:                "article",
			ID:                  "ad_en",
			Title:               "Recommend the bot",
			InputMessageContent: map[string]interface{}{
				"message_text": `🔥 <a href="https://t.me/translobot">Translo</a> 🌐 - <i>The best Telegram translator bot in the whole world</i>`,
				"disable_web_page_preview":true,
				"parse_mode": tgbotapi.ModeHTML,
			},
			ReplyMarkup:         &kbEn,
			URL:                 "",
			HideURL:             true,
			Description:         "click to recommend a bot in the chat",
			ThumbURL:            "https://i.yapx.ru/PdNIa.png",
			ThumbWidth:          200,
			ThumbHeight:         200,
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

	if offset > len(codes) - 1 {
		warn(fmt.Errorf("слишком большое смещение: %d", offset))
		return
	}

	count := 50
	if offset + count > len(codes) - 1 {
		count = len(codes) - 1 - offset
	}


	userExists := false
	from := ""
	var wg sync.WaitGroup
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		user, err = app.db.GetUserByID(update.From.ID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				warn(err)
				return
			}
		} else {
			userExists = true
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tr, err := translate2.GoogleTranslate("auto", "en", update.Query)
		if err != nil {
			warn(err)
		}
		from = tr.FromLang
	}()

	wg.Wait()
	blocks := make([]interface{}, 0, 50)

	if offset == 0 && userExists {
		if from != user.MyLang {
			keyboard := inlineTranslationKeyboard(user.MyLang)
			name := langs[user.MyLang].Name
			blocks = append(blocks, tgbotapi.InlineQueryResultArticle{
				Type:                "article",
				ID:                  "inline_translate:" + user.MyLang,
				Title:               name + " 🔥",
				InputMessageContent: map[string]interface{}{
					"message_text": update.Query,
					"disable_web_page_preview":true,
					"parse_mode": "",
				},
				ReplyMarkup:         &keyboard,
				URL:                 "",
				HideURL:             true,
				Description:         "send message in " + langs[user.MyLang].Name,
				ThumbURL:            "",
				ThumbWidth:          200,
				ThumbHeight:         200,
			})
		}
		if from != user.ToLang {
			keyboard := inlineTranslationKeyboard(user.ToLang)
			name := langs[user.ToLang].Name
			blocks = append(blocks, tgbotapi.InlineQueryResultArticle{
				Type:                "article",
				ID:                  "inline_translate:" + user.ToLang,
				Title:               name + " 🔥",
				InputMessageContent: map[string]interface{}{
					"message_text": update.Query,
					"disable_web_page_preview":true,
				},
				ReplyMarkup:         &keyboard,
				URL:                 "",
				HideURL:             true,
				Description:         "send message in " + name,
				ThumbURL:            "",
				ThumbWidth:          200,
				ThumbHeight:         200,
			})
		}
	}


	for i, code := range codes[offset:offset+count] {
		if code == user.MyLang || code == user.ToLang {
			continue
		}
		lang := langs[code]
		title := lang.Name
		if i < 18 && offset == 0 {
			title += " 📌"
		}
		keyboard := inlineTranslationKeyboard(code)
		blocks = append(blocks, tgbotapi.InlineQueryResultArticle{
			Type:                "article",
			ID:                 	strconv.Itoa(i)+":"+code,
			Title:               title,
			InputMessageContent: map[string]interface{}{
				"message_text": update.Query,
				"disable_web_page_preview":true,
				"parse_mode": "",
			},
			ReplyMarkup:         &keyboard,
			URL:                 "",
			HideURL:             true,
			Description:         "send message in " + lang.Name,
			ThumbURL:            "",
			ThumbWidth:          200,
			ThumbHeight:         200,
		})
	}

	if len(blocks) > 50 {
		count -= len(blocks) - 50
		blocks = blocks[:50]
	}

	pmtext := "Translo"
	if update.Query == "" {
		pmtext = "Enter text"
	}

	if _, err := app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
		InlineQueryID:     update.ID,
		Results:           blocks,
		CacheTime:         0,
		NextOffset: 	     strconv.Itoa(offset + count),
		IsPersonal:        true,
		SwitchPMText:      pmtext,
		SwitchPMParameter: "from_inline",
	}); err != nil {
		warn(errors.WrapPrefix(err, "app.bot.AnswerInlineQuery:", 1))
		pp.Println(blocks)
	}

	app.analytics.Bot(update.From.ID, "Inline succeeded", "Inline succeeded")
	if userExists {
		if err = app.db.IncreaseUserUsings(update.From.ID); err != nil {
			warn(err)
		}
	}
}
