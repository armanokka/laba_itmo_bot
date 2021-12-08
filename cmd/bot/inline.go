package bot

import (
	"fmt"
	"github.com/armanokka/translobot/pkg/translate"
	iso6391 "github.com/emvi/iso-639-1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
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
		logrus.Error(err)
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

	pp.Println("offset", offset, "count", count)


	fromlang := ""
	var wg sync.WaitGroup
	var user User
	results := sync.Map{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		user = app.loadUser(update.From.ID, warn)
		if user.Exists() {
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
			results.Store(code, tgbotapi.InlineQueryResultArticle{
				Type:                "article",
				ID:                  "да пох вообще",
				Title:               iso6391.Name(code),
				InputMessageContent: map[string]interface{}{
					"message_text": tr.Text,
					"disable_web_page_preview":false,
				},
				HideURL:             true,
				Description:         tr.Text,
			})
		}()
	}

	wg.Wait()

	blocks := make([]interface{}, 0, 50)

	if offset == 0 {
		if fromlang == user.MyLang || fromlang != user.MyLang && fromlang != user.ToLang {
			block, ok := results.Load(user.ToLang)
			if !ok {
				tr, err := translate.GoogleHTMLTranslate(fromlang, user.ToLang, update.Query)
				if err != nil {
					warn(err)
				}
				tr.Text = html.UnescapeString(tr.Text)

				block = tgbotapi.InlineQueryResultArticle{
					Type:                "article",
					ID:                  "да пох вообще",
					Title:               iso6391.Name(user.ToLang),
					InputMessageContent: map[string]interface{}{
						"message_text": tr.Text,
						"disable_web_page_preview":false,
					},
					HideURL:             true,
					Description:         tr.Text,
				}
			}
			article := block.(tgbotapi.InlineQueryResultArticle)
			article.ID = "my"
			article.Title += " 👤"
			blocks = append(blocks, article)
		}
		if fromlang == user.ToLang || fromlang != user.MyLang && fromlang != user.ToLang {
			block, ok := results.Load(user.MyLang)
			if !ok {
				tr, err := translate.GoogleHTMLTranslate(fromlang, user.MyLang, update.Query)
				if err != nil {
					warn(err)
				}
				tr.Text = html.UnescapeString(tr.Text)

				block = tgbotapi.InlineQueryResultArticle{
					Type:                "article",
					ID:                  "да пох вообще",
					Title:               iso6391.Name(user.MyLang),
					InputMessageContent: map[string]interface{}{
						"message_text": tr.Text,
						"disable_web_page_preview":false,
					},
					HideURL:             true,
					Description:         tr.Text,
				}
			}
			article := block.(tgbotapi.InlineQueryResultArticle)
			article.ID = "to"
			article.Title += " 👤"
			blocks = append(blocks, article)
		}
	}

	for i, code := range codes[offset:offset+count] {
		if offset == 0 && (code == user.MyLang || code == user.ToLang) {
			continue
		}
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

	if user.Exists() {
		user.UpdateLastActivity()
		app.writeUserLog(update.From.ID, update.Query)
		app.writeBotLog(update.From.ID, "inline_succeeded", "")
	}
}
