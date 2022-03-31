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
	"gorm.io/gorm"
	"html"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func (app *App) onMessage(ctx context.Context, message tgbotapi.Message) {
	user := tables.Users{Lang: message.From.LanguageCode}

	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Произошла ошибка")))
		app.notifyAdmin(err)
	}
	app.analytics.User(message.Text, message.From)

	if message.Chat.ID < 0 {
		return
	}

	defer func() {
		if err := app.db.UpdateUserMetrics(message.From.ID, message.Text); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	var err error
	user, err = app.db.GetUserByID(message.From.ID)
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
	user.SetLang(message.From.LanguageCode)

	switch message.Command() {
	case "start":
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              tgbotapi.NewRemoveKeyboard(true),
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text: user.Localize("Просто напиши мне текст, а я его переведу"),
		})
		query := user.Localize("пример текста")
		app.bot.Send(tgbotapi.VideoConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           message.From.ID,
					ChannelUsername:  "",
					ReplyToMessageID: 0,
					ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.InlineKeyboardButton{
								Text:                         user.Localize("Chat mode"),
								URL:                          nil,
								LoginURL:                     nil,
								CallbackData:                 nil,
								SwitchInlineQuery:            &query,
								SwitchInlineQueryCurrentChat: nil,
								CallbackGame:                 nil,
								Pay:                          false,
							})),
					DisableNotification:      false,
					AllowSendingWithoutReply: false,
				},
				File: tgbotapi.FilePath("inline.mp4"),
			},
			Thumb:             nil,
			Duration:          0,
			Caption:           user.Localize("Как переводить сообщения, не выходя из чата"),
			ParseMode:         "",
			CaptionEntities:   nil,
			SupportsStreaming: false,
		})
		if err = app.db.UpdateUser(message.From.ID, tables.Users{Act: "setup_langs"}); err != nil {
			warn(err)
		}
		return
	case "users":
		if message.From.ID != config.AdminID {
			return
		}
		f, err := os.CreateTemp("", "")
		if err != nil {
			warn(err)
			return
		}
		defer f.Close()
		defer os.Remove(f.Name())

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
		doc := tgbotapi.NewInputMediaDocument(tgbotapi.FilePath(f.Name()))
		group := tgbotapi.NewMediaGroup(message.From.ID, []interface{}{doc})
		app.bot.Send(group)
		if err = app.db.LogBotMessage(message.From.ID, "pm_users", "shared users' ids"); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
		return
	case "id":
		msg := tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10))
		app.bot.Send(msg)
		if err = app.db.LogBotMessage(message.From.ID, "pm_id", msg.Text); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}

		return
	case "mailing":
		if err = app.db.UpdateUser(message.From.ID, tables.Users{Act: "mailing"}); err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("отменить", "cancel_mailing_act"),
			),
		)

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         message.MessageID,
				ReplyMarkup:              keyboard,
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			Text:                  "отправь сообщение для рассылки",
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		})
		return
	}

	switch user.Act {
	case "mailing":
		if err := app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": ""}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewMessage(message.From.ID, "рассылка начата"))
		if err = app.db.DropMailings(); err != nil {
			warn(err)
			return
		}
		if err = app.db.CreateMailingTable(); err != nil {
			warn(err)
			return
		}
		if err = app.bc.Put([]byte("mailing_message_id"), []byte(strconv.Itoa(message.MessageID))); err != nil {
			warn(err)
			return
		}

		rows, err := app.db.GetMailersRows()
		if err != nil {
			warn(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var id int64
			if err = rows.Scan(&id); err != nil {
				warn(err)
				return
			}
			if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   id,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              nil,
					DisableNotification:      false,
					AllowSendingWithoutReply: false,
				},
				FromChatID:          config.AdminID,
				FromChannelUsername: "",
				MessageID:           message.MessageID,
				Caption:             "",
				ParseMode:           "",
				CaptionEntities:     nil,
			}); err != nil {
				pp.Println(err)
			}
			if err = app.db.DeleteMailuser(id); err != nil {
				warn(err)
			}
			time.Sleep(time.Second / 20)
		}
		err = app.db.DropMailings()
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewMessage(message.From.ID, "рассылка закончена"))
		return
	case "setup_langs":
		fromLang, err := translate.GoogleTranslate("auto", "en", cutStringUTF16(message.Text, 100))
		if err != nil {
			warn(err)
		}
		from := fromLang.FromLang

		keyboard, err := buildLangsPagination(0, 18, fromLang.FromLang, fmt.Sprintf("setup_langs:%s:%s", from, "%s"), fmt.Sprintf("setup_langs_pagination:%s:0", from), fmt.Sprintf("setup_langs_pagination:%s:18", from))
		if err != nil {
			warn(err)
		}
		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
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
		}); err != nil {
			pp.Println(err)
		}
		return
	}

	go app.bot.Send(tgbotapi.NewChatAction(message.From.ID, "typing"))

	var text = message.Text
	message.Text = ""
	if message.Caption != "" {
		text = message.Caption
		message.Caption = ""
	}

	if text == "" {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Отправь текстовое сообщение, чтобы я его перевел")))
		app.analytics.Bot(message.Chat.ID, "Please, send text message", "Message is not text message")
		return
	}

	from, err := translate.DetectLanguageGoogle(cutStringUTF16(text, 100))
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

	ret, err := app.SuperTranslate(user, from, to, text, message.Entities)
	if err != nil {
		warn(err)
		return
	}
	//ret.TranslatedText, err = url.QueryUnescape(ret.TranslatedText)
	//if err != nil {
	//	warn(err)
	//	return
	//} хуй

	ret.TranslatedText = html.UnescapeString(ret.TranslatedText)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:                         user.Localize("Chat mode"),
				URL:                          nil,
				LoginURL:                     nil,
				CallbackData:                 nil,
				SwitchInlineQuery:            nil,
				SwitchInlineQueryCurrentChat: &text,
				CallbackGame:                 nil,
				Pay:                          false,
			}), // inline
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔊 "+user.Localize("Озвучить"), fmt.Sprintf("speech_this_message_and_replied_one:%s:%s", from, to))))
	if ret.Examples {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("💬 "+user.Localize("Примеры"), fmt.Sprintf("exm:%s:%s", from, to)))
	}
	if ret.Translations {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📚 "+user.Localize("Переводы"), fmt.Sprintf("trs:%s:%s", from, to))))
	}
	if ret.Dictionary {
		l := len(keyboard.InlineKeyboard) - 1
		if l < 0 {
			l = 0
		}
		keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData("ℹ️"+user.Localize("Словарь"), fmt.Sprintf("dict:%s", from)))
	}

	if _, err = app.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
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

	if err = app.db.LogBotMessage(message.From.ID, "pm_translate", ret.TranslatedText); err != nil {
		app.notifyAdmin(fmt.Errorf("%w", err))
	}

	if user.Usings == 5 || (user.Usings > 0 && user.Usings%20 == 0) {
		link := strings.ReplaceAll(user.Localize("Я рекомендую @translobot"), " ", "+")
		link = url.PathEscape(link)
		defer func() {
			if _, err := app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           message.From.ID,
					ChannelUsername:  "",
					ReplyToMessageID: 0,
					ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonURL(user.Localize("Рассказать про нас"), "http://t.me/share/url?url="+link))),
					DisableNotification:      true,
					AllowSendingWithoutReply: false,
				},
				Text:                  user.Localize("Понравился бот? 😎 Поделись с друзьями, нажав на кнопку"),
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			}); err != nil {
				pp.Println(err)
			}
		}()
	}
}
