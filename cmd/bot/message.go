package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/levenshtein"
	"github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (app *app) onMessage(ctx context.Context, message tgbotapi.Message) {
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", message.From.LanguageCode)))
		app.notifyAdmin(err)
		logrus.Error(err)
	}
	app.analytics.User(message.Text, message.From)

	if message.Chat.ID < 0 {
		return
	}

	var user = app.loadUser(message.From.ID, warn)
	defer user.UpdateLastActivity()

	if strings.HasPrefix(message.Text, "/start") {
		user.SendStart(message)
		return
	} else {
		user.Fill()
	}

	app.writeUserLog(message.From.ID, message.Text)

	if low := strings.ToLower(message.Text); low != "" {
		switch {
		case in(command("my language"), low): // добавить обратную совместимось
			user.SendStart(message)
			return
		case in(command("translate language"), low):
			user.SendStart(message) // backward compatibility
			return
		}
	}

	switch message.Command() {
	case "stats":
		var users int
		err := app.db.Model(&tables.Users{}).Raw("SELECT COUNT(*) FROM users").Find(&users).Error
		if err != nil {
			warn(err)
			return
		}
		var stats = make(map[string]string, 20)
		if err = app.db.Model(&tables.UsersLogs{}).Raw("SELECT intent, COUNT(*) FROM users_logs GROUP BY intent ORDER BY count(*) DESC").Find(&stats).Error; err != nil {
			warn(err)
		}
		text := "Всего " + strconv.Itoa(users) + " юзеров"
		for name, count := range stats {
			text += "\n" + name + ": " + count
		}
		msg :=  tgbotapi.NewMessage(message.Chat.ID, text)
		app.bot.Send(msg)
		app.writeBotLog(message.From.ID, "pm_stats", msg.Text)
		return
	case "users":
		if message.From.ID != config.AdminID {
			return
		}
		f, err := os.Create("users.txt")
		if err != nil {
			warn(err)
			return
		}
		var users []tables.Users
		if err = app.db.Model(&tables.Users{}).Find(&users).Error; err != nil {
			warn(err)
			return
		}
		for _, user := range users {
			if _, err = f.WriteString(strconv.FormatInt(user.ID, 10) + "\r\n"); err != nil {
				warn(err)
				return
			}
		}
		doc :=  tgbotapi.NewInputMediaDocument("users.txt")
		group :=  tgbotapi.NewMediaGroup(message.From.ID, []interface{}{doc})
		app.bot.Send(group)
		app.writeBotLog(message.From.ID, "pm_users", "{document was sended}")
		return
	case "id":
		msg :=  tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10))
		app.bot.Send(msg)
		app.writeBotLog(message.From.ID, "pm_id", msg.Text)
		return
	}

	switch user.Act {
	case "setup_langs":
		fromLang, err := translate.GoogleHTMLTranslate("auto", "en", message.Text)
		if err != nil {
			warn(err)
		}
		from := fromLang.From


		keyboard, err := user.buildLangsPagination(0, 18, from, fmt.Sprintf("setup_langs:%s:%s", from, "%s"), fmt.Sprintf("setup_langs_pagination:%s:0", from), fmt.Sprintf("setup_langs_pagination:%s:18", from))
		if err != nil {
			warn(err)
		}
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         message.MessageID,
				ReplyMarkup:              keyboard,
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:                  user.Localize("На какой язык перевести?"),
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: true,
		})
		return
	}

	if user.Usings == 5 || (user.Usings > 0 && user.Usings % 20 == 0) {
		text := user.Localize("Я рекомендую @translobot")
		link := strings.ReplaceAll(text, " ", "+")
		link = url.PathEscape(link)
		defer func() {
			if _, err := app.bot.Send(tgbotapi.MessageConfig{
				BaseChat:              tgbotapi.BaseChat{
					ChatID:                   message.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonURL(user.Localize("Рассказать про нас"), "http://t.me/share/url?url=" + link))),
					DisableNotification:      true,
					AllowSendingWithoutReply: false,
				},
				Text:                  user.Localize("recommendation"),
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			}); err != nil {
				pp.Println(err)
			}
		}()
		// tg://share...
	}






	var text = message.Text
	if message.Caption != "" {
		text = message.Caption
	}

	if text == "" {
		app.bot.Send( tgbotapi.NewMessage(message.Chat.ID, user.Localize("Please, send text message")))
		app.analytics.Bot(message.Chat.ID, "Please, send text message", "Message is not text message")
		return
	}

	detect, err := translate.GoogleHTMLTranslate("auto", "en", text)
	if err != nil {
		warn(err)
		return
	}
	from := detect.From

	if from != user.MyLang && from != user.ToLang {

		var (
			trFromMyLang translate.GoogleHTMLTranslation
			trFromToLang translate.GoogleHTMLTranslation
		)
		g, _ := errgroup.WithContext(ctx)
		g.Go(func() error {
			trFromMyLang, err = translate.GoogleHTMLTranslate(user.MyLang, "en", text)
			return err
		})
		g.Go(func() error {
			trFromToLang, err = translate.GoogleHTMLTranslate(user.ToLang, "en", text)
			return err
		})

		if err = g.Wait(); err != nil {
			warn(err)
			return
		}

		leven1 := levenshtein.ComputeLevenshteinPercentage(text, detect.Text)
		leven2 := levenshtein.ComputeLevenshteinPercentage(text, trFromMyLang.Text)
		leven3 := levenshtein.ComputeLevenshteinPercentage(text, trFromToLang.Text)

		min := min(leven1, leven2, leven3)
		if min == leven1 {
			// detect.From правильный
		} else if min == leven2 {
			from = trFromMyLang.From
		} else if min == leven3 {
			from = trFromToLang.From
		}
	}

	if from == "" {
		from = "auto"
	}


	var to string // language into need to translate
	if from == user.ToLang {
		to = user.MyLang
	} else if from == user.MyLang {
		to = user.ToLang
	} else { // никакой из
		to = user.MyLang
	}



	ret, err := app.SuperTranslate(from, to, text, message.Entities)
	if err != nil {
		warn(err)
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔊 " + user.Localize("To voice"), fmt.Sprintf("speech_this_message_and_replied_one:%s:%s", from, to))))
	if ret.Examples {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("💬 " + user.Localize("Examples"), fmt.Sprintf("exm:%s:%s", from, to)))
	}
	if ret.Translations {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 " + user.Localize("Translations"), fmt.Sprintf("trs:%s:%s", from, to))))
	}
	if ret.Dictionary {
		l := len(keyboard.InlineKeyboard) - 1
		if l < 0 {
			l = 0
		}
		keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData("ℹ️" + user.Localize("Dictionary"), fmt.Sprintf("dict:%s", from)))
	}

	if _, err = app.bot.Send(tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{
			ChatID:                   message.Chat.ID,
			ChannelUsername:          "",
			ReplyToMessageID:         message.MessageID,
			ReplyMarkup:              keyboard,
			DisableNotification:      true,
			AllowSendingWithoutReply: false,
		},
		Text:                  ret.TranslatedText,
		ParseMode:             tgbotapi.ModeHTML,
		Entities:              nil,
		DisableWebPagePreview: false,
	}); err != nil {
		pp.Println(err)
	}

	app.analytics.Bot(user.ID, ret.TranslatedText, "Translated")
	user.IncrUsings()
	app.writeBotLog(message.From.ID, "pm_translate", ret.TranslatedText)
}