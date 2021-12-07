package bot

import (
	"context"
	"database/sql"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
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

	if strings.HasPrefix(message.Text, "/start") {
		if !user.Exists() {
			if message.From.LanguageCode == "" || !in(config.BotLocalizedLangs, message.From.LanguageCode) {
				message.From.LanguageCode = "en"
			}
		} else {
			user.Fill()
		}

		kb, err := BuildSupportedLanguagesKeyboard()
		if err  != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:               tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ReplyMarkup:              kb,
			},
			Text:                  user.Localize("Выберите свой язык"),
		})
		return
	}

	user.Fill()

	defer user.UpdateLastActivity()

	app.writeUserLog(message.From.ID, message.Text)

	if low := strings.ToLower(message.Text); low != "" {
		switch {
		case in(command("my language"), low):
			kb, err := buildLangsPagination(0, 18, "set_my_lang_by_callback:%s", "set_my_lang_pagination:0", "set_my_lang_pagination:18", "set_my_lang_by_callback:" + user.MyLang)
			if err != nil {
				warn(err)
			}
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat:              tgbotapi.BaseChat{
					ChatID:                   message.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              kb,
					DisableNotification:      false,
					AllowSendingWithoutReply: false,
				},
				Text:                  user.Localize("Ваш язык <b>%s</b>. Выберите Ваш язык.", langs[user.MyLang].Name),
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			})
			return
		case in(command("translate language"), low):
			kb, err := buildLangsPagination(0, 18, "set_to_lang_by_callback:%s", "set_to_lang_pagination:0", "set_to_lang_pagination:18", "set_to_lang_by_callback:" + user.ToLang)
			if err != nil {
				warn(err)
			}
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat:              tgbotapi.BaseChat{
					ChatID:                   message.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              kb,
					DisableNotification:      false,
					AllowSendingWithoutReply: false,
				},
				Text:                  user.Localize("Сейчас бот переводит на <b>%s</b>. Выберите язык, на который хотите переводить", langs[user.ToLang].Name),
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			})
			return
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

	if user.MyLang == user.ToLang {
		app.bot.Send( tgbotapi.NewMessage(message.From.ID, user.Localize("The original text language and the target language are the same, please set different")))
		return
	}

	if user.Usings == 5 || (user.Usings > 0 && user.Usings % 20 == 0) {
		// tg://share...
	}

	if user.Usings > 15 && !user.IsDeveloper.Valid {
		defer func() {
			_, err := app.bot.Send( tgbotapi.MessageConfig{
				BaseChat:               tgbotapi.BaseChat{
					ChatID:                   message.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         message.MessageID,
					ReplyMarkup:               tgbotapi.NewInlineKeyboardMarkup(
						 tgbotapi.NewInlineKeyboardRow(
							 tgbotapi.NewInlineKeyboardButtonData(user.Localize("Yes"), "I'm developer"),
							 tgbotapi.NewInlineKeyboardButtonData(user.Localize("No"), "delete"))),
					DisableNotification:      true,
					AllowSendingWithoutReply: false,
				},
				Text:                  user.Localize("Are you developer?"),
				ParseMode:             "",
				Entities:              nil,
				DisableWebPagePreview: false,
			})
			if err != nil {
				pp.Println(err)
			}
		}()
		user.Update(tables.Users{IsDeveloper: sql.NullBool{
			Bool:  false,
			Valid: true,
		}})
	}

	msg, err := app.bot.Send( tgbotapi.MessageConfig{
		BaseChat:  tgbotapi.BaseChat{
			ChatID:                   message.Chat.ID,
			ChannelUsername:          "",
			ReplyToMessageID:         message.MessageID, // very important to "dictionary" callback
			ReplyMarkup:              nil,
			DisableNotification:      true,
			AllowSendingWithoutReply: true,
		},
		Text:                  user.Localize("⏳ Translating..."),
		ParseMode:              tgbotapi.ModeHTML,
		Entities:              nil,
		DisableWebPagePreview: false,
	})
	if err != nil {
		return
	}

	var text = message.Text
	if message.Caption != "" {
		text = message.Caption
	}

	if text == "" {
		app.bot.Send( tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, user.Localize("Please, send text message")))
		app.analytics.Bot(message.Chat.ID, msg.Text, "Message is not text message")
		return
	}

	detect, err := translate.GoogleHTMLTranslate("auto", "en", cutStringUTF16(text, 100))
	if err != nil {
		warn(err)
		return
	}
	from := detect.From

	if from != user.MyLang && from != user.ToLang { // если задетекчен язык ни родной, ни переводный
		trTest, err := translate.GoogleHTMLTranslate(user.MyLang, "en", text) // пробуем перевести с нашего языка
		if err != nil {
			warn(err)
		}
		if trTest.Text != text {
			pp.Println(trTest.Text)
			from = user.MyLang
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

	From := langs[from]
	keyboard :=  tgbotapi.NewInlineKeyboardMarkup(
		 tgbotapi.NewInlineKeyboardRow(
			 tgbotapi.NewInlineKeyboardButtonData("From " + From.Emoji + " " + From.Name, "none")),
		 tgbotapi.NewInlineKeyboardRow(
			 tgbotapi.NewInlineKeyboardButtonData("🔊 " + user.Localize("To voice"),  "speech_this_message_and_replied_one:"+from+":"+to)),
	)

	if ret.Examples {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			 tgbotapi.NewInlineKeyboardRow( tgbotapi.NewInlineKeyboardButtonData("💬 " + user.Localize("Examples"), "examples:"+from+":"+to)))
	}
	if ret.Translations {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			 tgbotapi.NewInlineKeyboardRow( tgbotapi.NewInlineKeyboardButtonData("📚 " + user.Localize("Translations"), "translations:"+from+":"+to)))
	}

	if ret.Dictionary {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			 tgbotapi.NewInlineKeyboardRow( tgbotapi.NewInlineKeyboardButtonData("ℹ️" + user.Localize("Dictionary"), "dictonary:"+from+":"+text)))
	}

	if _, err = app.bot.Send( tgbotapi.EditMessageTextConfig{
		BaseEdit:               tgbotapi.BaseEdit{
			ChatID:          message.From.ID,
			MessageID:       msg.MessageID,
			ReplyMarkup:     &keyboard,
		},
		Text:                  ret.TranslatedText,
		ParseMode:              tgbotapi.ModeHTML,
		Entities:              nil,
		DisableWebPagePreview: false,
	}); err != nil {
		logrus.Error(err)
		pp.Println(err)
	}

	app.analytics.Bot(user.ID, ret.TranslatedText, "Translated")
	user.IncrUsings()
	if user.Exists() {
		app.writeBotLog(message.From.ID, "pm_translate", ret.TranslatedText)
	}
}