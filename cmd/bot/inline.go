package bot

import (
	"fmt"
	"github.com/armanokka/translobot/pkg/translate"
	iso6391 "github.com/emvi/iso-639-1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"html"
	"strconv"
	"sync"
)

func (app app) onInlineQuery(update tgbotapi.InlineQuery) {
	app.analytics.User(update.Query, update.From)

	warn := func(err error) {
		app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID:     update.ID,
			SwitchPMText:      "Error, sorry.",
			SwitchPMParameter: "from_inline",
		})
		app.notifyAdmin(err)
		pp.Println("onInlineQuery: error", err)
	}

	if update.Query == "" {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Try it out", "https://t.me/translobot?start=from_inline")))
		elem := tgbotapi.InlineQueryResultArticle{
			Type:                "article",
			ID:                  "ad",
			Title:               "Recommend the bot",
			InputMessageContent: map[string]interface{}{
				"message_text": `🔥 <a href="https://t.me/translobot">Translo</a> 🌐 - <i>The best Telegram translator bot in the whole world</i>`,
				"disable_web_page_preview":true,
				"parse_mode": tgbotapi.ModeHTML,
			},
			ReplyMarkup:         &keyboard,
			URL:                 "",
			HideURL:             true,
			Description:         "click to recommend a bot in the chat",
			ThumbURL:            "https://i.yapx.ru/PdNIa.png",
			ThumbWidth:          200,
			ThumbHeight:         200,
		}
		app.bot.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID:     update.ID,
			Results:           []interface{}{elem},
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
	var user User
	results := sync.Map{}
	var userExists bool

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = app.loadUser(update.From.ID, warn)
		if user.Exists() {
			userExists = true
			user.Fill()
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
				return
			}
			if fromlang == "" {
				fromlang = tr.From
			}
			tr.Text = html.UnescapeString(tr.Text)

			if tr.Text == "" {
				return // ну не вышло, так не вышло, че бубнить-то
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
		langs := user.GetUsedLangs()
		for i, lang := range langs {
			block, ok := results.Load(lang)
			if !ok {
				tr, err := translate.GoogleHTMLTranslate(fromlang, lang, update.Query)
				if err != nil {
					warn(err)
				}
				tr.Text = html.UnescapeString(tr.Text)

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

				block = tgbotapi.InlineQueryResultArticle{
					Type:                "article",
					ID:                  "да пох вообще",
					Title:               iso6391.Name(lang),
					InputMessageContent: map[string]interface{}{
						"message_text": tr.Text,
						"disable_web_page_preview":false,
					},
					ReplyMarkup: &keyboard,
					HideURL:             true,
					Description:         tr.Text,
				}
			}
			article := block.(tgbotapi.InlineQueryResultArticle)
			article.ID = "used_lang:"+strconv.Itoa(i)
			article.Title += " 👤"
			blocks = append(blocks, article)
		}
	}

	for i, code := range codes[offset:offset+count] {
		block, ok := results.Load(code)
		if !ok {
			warn(fmt.Errorf("couldn't find code %s in translations", code))
		}
		article := block.(tgbotapi.InlineQueryResultArticle)
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
		warn(err)
		pp.Println(blocks)
	}

	app.analytics.Bot(update.From.ID, "Inline succeeded", "Inline succeeded")
	if userExists {
		user.IncrUsings()
	}
}
