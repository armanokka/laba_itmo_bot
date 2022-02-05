package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func (app *App) onMessage(ctx context.Context, message tgbotapi.Message) {
	localizer := i18n.NewLocalizer(app.bundle, message.From.LanguageCode)

	warn := func(err error) {
		locale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)"})
		if err != nil {
			app.notifyAdmin(err)
			app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)"))
		} else {
			app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, locale))
		}
		app.notifyAdmin(err)
	}
	app.analytics.User(message.Text, message.From)

	if message.Chat.ID < 0 {
		return
	}


	defer func() {
		if err := app.db.UpdateUserLastActivity(message.From.ID); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	user, err := app.db.GetUserByID(message.From.ID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			if message.From.LanguageCode == "" {
				message.From.LanguageCode = "en"
			}
			err = app.db.CreateUser(tables.Users{
				ID:           message.From.ID,
				MyLang:       "",
				ToLang:       "",
				Act:          "setup_langs",
				Usings:       1,
				Blocked:      false,
				LastActivity: time.Now(),
			})
			if err != nil {
				warn(err)
				return
			}

		} else {
			warn(err)
		}
	}

	if strings.HasPrefix(message.Text, "/start") {
		locale, err := localizer.LocalizeMessage(&i18n.Message{
			ID:          "Просто напиши мне текст, а я его переведу",
		})
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:               tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              tgbotapi.NewRemoveKeyboard(true),
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:                  locale,
		})
		if err = app.db.UpdateUser(message.From.ID, tables.Users{Act: "setup_langs"}); err != nil {
			warn(err)
		}
		return
	}


	if err = app.db.LogUserMessage(message.From.ID, message.Text); err != nil {
		app.notifyAdmin(fmt.Errorf("%w", err))
	}

	switch message.Command() {
	case "users":
		if message.From.ID != config.AdminID {
			return
		}
		f, err := os.Create("users.txt")
		if err != nil {
			warn(err)
			return
		}
		users, err := app.db.GetAllUsers()
		if err != nil {
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
		if err = app.db.LogBotMessage(message.From.ID, "pm_users", "shared users' ids"); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
		return
	case "id":
		msg :=  tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10))
		app.bot.Send(msg)
		if err = app.db.LogBotMessage(message.From.ID, "pm_id", msg.Text); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}

		return
	}

	switch user.Act {
	case "setup_langs":
		fromLang, err := translate.GoogleHTMLTranslate("auto", "en", message.Text)
		if err != nil {
			warn(err)
		}
		from := fromLang.From


		keyboard, err := buildLangsPagination(0, 18, fmt.Sprintf("setup_langs:%s:%s", from, "%s"), fmt.Sprintf("setup_langs_pagination:%s:0", from), fmt.Sprintf("setup_langs_pagination:%s:18", from))
		if err != nil {
			warn(err)
		}
		locale, err := localizer.LocalizeMessage(&i18n.Message{ID: "На какой язык перевести?"})
		if err != nil {
			warn(err)
			return
		}
		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         message.MessageID,
				ReplyMarkup:              keyboard,
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:                  locale,
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: true,
		}); err != nil {
			pp.Println(err)
		}
		return
	}

	if user.Usings == 5 || (user.Usings > 0 && user.Usings % 20 == 0) {
		IrecommendBotLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Я рекомендую @translobot"})
		if err != nil {
			warn(err)
			return
		}

		TellAboutUsLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Рассказать про нас"})
		if err != nil {
			warn(err)
			return
		}

		RecommendationLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "recommendation"})
		if err != nil {
			warn(err)
			return
		}

		link := strings.ReplaceAll(IrecommendBotLocale, " ", "+")
		link = url.PathEscape(link)
		defer func() {
			if _, err := app.bot.Send(tgbotapi.MessageConfig{
				BaseChat:              tgbotapi.BaseChat{
					ChatID:                   message.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonURL(TellAboutUsLocale, "http://t.me/share/url?url=" + link))),
					DisableNotification:      true,
					AllowSendingWithoutReply: false,
				},
				Text:                  RecommendationLocale,
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			}); err != nil {
				pp.Println(err)
			}
		}()
	}






	var text = message.Text
	if message.Caption != "" {
		text = message.Caption
	}

	if text == "" {
		locale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Please, send text message"})
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send( tgbotapi.NewMessage(message.Chat.ID, locale))
		app.analytics.Bot(message.Chat.ID, "Please, send text message", "Message is not text message")
		return
	}


	from, err := translate.DetectLanguageGoogle(text)
	if err != nil {
		warn(err)
		return
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

	ToVoiceLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "To voice"})
	if err != nil {
		warn(err)
		return
	}
	ExamplesLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Examples"})
	if err != nil {
		warn(err)
		return
	}
	TranslationsLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Translations"})
	if err != nil {
		warn(err)
		return
	}
	DictionaryLocale, err := localizer.LocalizeMessage(&i18n.Message{ID: "Dictionary"})
	if err != nil {
		warn(err)
		return
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔊 " + ToVoiceLocale, fmt.Sprintf("speech_this_message_and_replied_one:%s:%s", from, to))))
	if ret.Examples {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("💬 " + ExamplesLocale, fmt.Sprintf("exm:%s:%s", from, to)))
	}
	if ret.Translations {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 " + TranslationsLocale, fmt.Sprintf("trs:%s:%s", from, to))))
	}
	if ret.Dictionary {
		l := len(keyboard.InlineKeyboard) - 1
		if l < 0 {
			l = 0
		}
		keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData("ℹ️" + DictionaryLocale, fmt.Sprintf("dict:%s", from)))
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
	if err = app.db.IncreaseUserUsings(message.From.ID); err != nil {
		app.notifyAdmin(fmt.Errorf("%w", err))
	}
	if err = app.db.LogBotMessage(message.From.ID, "pm_translate", ret.TranslatedText); err != nil {
		app.notifyAdmin(fmt.Errorf("%w", err))
	}
}