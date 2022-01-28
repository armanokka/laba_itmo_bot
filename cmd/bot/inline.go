package bot

import (
	"fmt"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/translate"
	iso6391 "github.com/emvi/iso-639-1"
	"github.com/go-errors/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
	"html"
	"strconv"
	"sync"
)

func (app App) onInlineQuery(update tgbotapi.InlineQuery) {
	app.analytics.User(update.Query, update.From)

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

	localizer := i18n.NewLocalizer(app.bundle, update.From.LanguageCode)

	if update.Query == "" {

		tryItOutLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Try it out"})
		if err != nil {
			warn(err)
			return
		}
		recommendLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Порекомендовать бота"})
		if err != nil {
			warn(err)
			return
		}
		adLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "inline_ad"})
		if err != nil {
			warn(err)
			return
		}
		clickLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "click to recommend a bot in the chat"})
		if err != nil {
			warn(err)
			return
		}
		kbLoc := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(tryItOutLocale, "https://t.me/translobot?start=from_inline")))

		elemLoc := tgbotapi.InlineQueryResultArticle{
			Type:                "article",
			ID:                  "ad_local",
			Title:               recommendLocale,
			InputMessageContent: map[string]interface{}{
				"message_text": adLocale,
				"disable_web_page_preview":true,
				"parse_mode": tgbotapi.ModeHTML,
			},
			ReplyMarkup:         &kbLoc,
			URL:                 "",
			HideURL:             true,
			Description:         clickLocale,
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


	fromlang := ""
	var wg sync.WaitGroup
	var user tables.Users
	var err error
	results := sync.Map{}
	var userExists bool

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

	for _, code := range codes[offset:offset + count] {
		wg.Add(1)
		code := code
		go func() {
			defer wg.Done()

			tr, err := translate.GoogleHTMLTranslate("auto", code, update.Query)
			if err != nil {
				warn(err)
				// not return
				tr.Text = "error"
			}
			if fromlang == "" {
				fromlang = tr.From
			}
			tr.Text = html.UnescapeString(tr.Text)

			if tr.Text == "" {
				// not return
				tr.Text = "error"
			}

			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.InlineKeyboardButton{
						Text:                         "translate",
						URL:                          nil,
						LoginURL:                     nil,
						CallbackData:                 nil,
						SwitchInlineQuery:            nil,
						SwitchInlineQueryCurrentChat: &tr.Text,
						CallbackGame:                 nil,
						Pay:                          false,
					}))

			results.Store(code, tgbotapi.InlineQueryResultArticle{
				Type:                "article",
				ID:                  "да пох вообще",
				Title:               iso6391.Name(code),
				InputMessageContent: map[string]interface{}{
					"message_text": tr.Text,
					"disable_web_page_preview":false,
				},
				ReplyMarkup: &keyboard,
				HideURL:             true,
				Description:         tr.Text,
			})
		}()
	}

	wg.Wait()

	blocks := make([]interface{}, 0, 50)

	if offset == 0 && userExists {
		if fromlang != user.MyLang {
				block, ok := results.Load(user.MyLang)
				if ok {
					article, ok := block.(tgbotapi.InlineQueryResultArticle)
					if !ok {
						app.notifyAdmin(fmt.Errorf("there is no tgbotapi.InlineQueryResultArticle in block with code %s", user.MyLang))
					}
					if ok {
						article.ID = "my_lang!"
						article.Title += " 👤"
						blocks = append(blocks, article)
					}
				}
		}
		if fromlang != user.ToLang {
			block, ok := results.Load(user.ToLang)
			if ok {
				article, ok := block.(tgbotapi.InlineQueryResultArticle)
				if !ok {
					app.notifyAdmin(fmt.Errorf("there is no tgbotapi.InlineQueryResultArticle in block with code %s", user.ToLang))
				}
				if ok {
					article.ID = "to_lang!"
					article.Title += " 👤"
					blocks = append(blocks, article)
				}

			}
		}
	}

	for i, code := range codes[offset:offset+count] {
		block, ok := results.Load(code)
		if !ok {
			warn(fmt.Errorf("couldn't find code %s in translations", code))
		}
		article, ok := block.(tgbotapi.InlineQueryResultArticle)
		if !ok {
			app.notifyAdmin(fmt.Errorf("there is no tgbotapi.InlineQueryResultArticle in block with code %s", code))
		}
		article.ID = strconv.Itoa(offset + i)
		if offset == 0 && i < 18 {
			article.Title += " 📌"
		}
		blocks = append(blocks, article)
	}

	if len(blocks) > 50 {
		count -= len(blocks) - 50
		blocks = blocks[:50]
	}

	pmtext := "From: " + langs[fromlang].Name
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
