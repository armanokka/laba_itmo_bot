package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/helpers"
	"github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"html"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func (app *App) onMessage(ctx context.Context, message tgbotapi.Message) {
	log := app.log.With(zap.Int64("id", message.From.ID))
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(errors.Error); ok {
				app.bot.Send(tgbotapi.NewMessage(config.AdminID, "recover:"+fmt.Sprint(err)+"\nstack:"+string(e.Stack())))
				return
			}
			log.Error("recover:", zap.Any("error", err))
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)+"\n\n"+string(debug.Stack())))
		}
	}()

	user := tables.Users{Lang: &message.From.LanguageCode}

	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Произошла ошибка")))
		app.notifyAdmin(err)
	}

	if message.Chat.ID < 0 {
		app.onGroupMessage(message)
		return
	}

	app.bot.Send(tgbotapi.NewChatAction(message.Chat.ID, "typing"))

	if err := app.analytics.User(message); err != nil {
		app.notifyAdmin(err)
	}

	defer func() {
		if err := app.db.UpdateUserActivity(message.From.ID); err != nil {
			app.notifyAdmin(err)
		}
	}()

	var err error
	user, err = app.db.GetUserByID(message.From.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tolang := ""
			if message.From.LanguageCode == "" || message.From.LanguageCode == "en" {
				message.From.LanguageCode = "en"
				tolang = "ru"
			} else if message.From.LanguageCode == "ru" {
				tolang = "en"
			}
			user = tables.Users{
				ID:           message.From.ID,
				MyLang:       message.From.LanguageCode,
				ToLang:       tolang,
				Lang:         nil,
				LastActivity: time.Now(),
			}
			if err = app.db.CreateUser(&user); err != nil {
				warn(err)
				return
			}
		} else {
			warn(err)
			return
		}
	}
	if user.Blocked {
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"blocked": false}); err != nil {
			app.notifyAdmin(err)
		}
	}
	log = log.With(zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang))

	switch message.Command() {
	case "start":
		if user.Lang == nil || *user.Lang == "" {
			languages := map[string]string{
				"en": "🇬🇧English",
				"de": "🇩🇪Deutsch",
				"es": "🇪🇸Español",
				"uk": "🇺🇦Українська",
				"ar": "🇪🇬عربي",
				"ru": "🇷🇺Русский",
				"uz": "🇺🇿O'Zbek",
				"id": "🇮🇩Bahasa Indonesia",
				"it": "🇮🇹Italiano",
				"pt": "🇵🇹Português",
			}
			keyboard := tgbotapi.NewInlineKeyboardMarkup()
			i := -1
			for code, name := range languages {
				i++
				btn := tgbotapi.NewInlineKeyboardButtonData(name, "set_lang:"+code)
				if i%2 == 0 || len(keyboard.InlineKeyboard) == 0 {
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
					continue
				}
				l := len(keyboard.InlineKeyboard) - 1
				if l < 0 {
					l = 0
				}
				keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], btn)
			}
			msg := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:              message.From.ID,
					ReplyMarkup:         keyboard,
					DisableNotification: true,
				},
				Text: user.Localize(`Choose language of the bot`),
			}
			if _, err = app.bot.Send(msg); err != nil {
				warn(err)
				return
			}
			if err = app.analytics.Bot(msg, "/start"+message.CommandArguments()); err != nil {
				app.notifyAdmin(err)
			}
			return
		}
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.From.ID,
				ChannelUsername:  "",
				ReplyToMessageID: 0,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.ToLang]+" "+flags[user.ToLang].Emoji))),
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:      user.Localize("<b>Send text</b>, and bot will translate it"),
			ParseMode: tgbotapi.ModeHTML,
		}
		if _, err = app.bot.Send(msg); err != nil {
			warn(err)
			return
		}
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": ""}); err != nil {
			warn(err)
		}
		return
	case "users":
		log.Info("/users")
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
		return
	case "set_bot_lang":
		languages := map[string]string{
			"en": "🇬🇧English",
			"de": "🇩🇪Deutsch",
			"es": "🇪🇸Español",
			"uk": "🇺🇦Українська",
			"ar": "🇪🇬عربي",
			"ru": "🇷🇺Русский",
			"uz": "🇺🇿O'Zbek",
			"id": "🇮🇩Bahasa Indonesia",
			"it": "🇮🇹Italiano",
			"pt": "🇵🇹Português",
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		i := -1
		for code, name := range languages {
			i++
			btn := tgbotapi.NewInlineKeyboardButtonData(name, "set_lang:"+code)
			if i%2 == 0 || len(keyboard.InlineKeyboard) == 0 {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
				continue
			}
			l := len(keyboard.InlineKeyboard) - 1
			if l < 0 {
				l = 0
			}
			keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], btn)
		}
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:              message.From.ID,
				ReplyMarkup:         keyboard,
				DisableNotification: true,
			},
			Text: user.Localize(`Choose language of the bot`),
		}
		if _, err = app.bot.Send(msg); err != nil {
			warn(err)
			return
		}
		if err = app.analytics.Bot(msg, "/start"+message.CommandArguments()); err != nil {
			app.notifyAdmin(err)
		}
		return
	case "id":
		log.Info("/id")
		msg := tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10))
		app.bot.Send(msg)

		return
	case "mailing":
		log.Info("/mailing")
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
			Text:                  "отправь сообщение для рассылки\b\bбез клавиатуры",
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		})
		return
	}

	switch strings.TrimSpace(message.Text) {
	case "↔":
		if user.MyLang == "auto" {
			app.bot.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
			if err = app.analytics.Bot(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: message.Chat.ID,
				},
				Text: ":delete_message",
			}, "tried_to_swap_autodetect_lang"); err != nil {
				app.notifyAdmin(err)
			}
			return
		}

		user.MyLang, user.ToLang = user.ToLang, user.MyLang
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: message.From.ID,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji))),
			},
			Text: user.Localize("Клавиатура обновлена"),
		}
		app.bot.Send(msg)
		if err = app.db.SwapLangs(message.Chat.ID); err != nil {
			warn(err)
			return
		}
		if err = app.analytics.Bot(msg, "↔"); err != nil {
			app.notifyAdmin(err)
		}
		return
	case concatNonEmpty(" ", langs[*user.Lang][user.MyLang], flags[user.MyLang].Emoji):
		kb, err := buildLangsPagination(user, 0, 18, "",
			fmt.Sprintf("set_my_lang:%s:%d", "%s", 0),
			fmt.Sprintf("set_my_lang_pagination:%d", len(codes[*user.Lang])/18*18),
			fmt.Sprintf("set_my_lang_pagination:%d", 18), true)
		if err != nil {
			warn(err)
		}
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      message.Chat.ID,
				ReplyMarkup: kb,
			},
			Text: user.Localize("Choose language"),
		}
		app.bot.Send(msg)
		if err = app.analytics.Bot(msg, "my_lang"); err != nil {
			app.notifyAdmin(err)
		}
		return
	case concatNonEmpty(" ", langs[*user.Lang][user.ToLang], flags[user.ToLang].Emoji):
		if user.Lang == nil {
			user.Lang = &message.From.LanguageCode
		}
		kb, err := buildLangsPagination(user, 0, 18, "",
			fmt.Sprintf("set_to_lang:%s:%d", "%s", 0),
			fmt.Sprintf("set_to_lang_pagination:%d", len(codes[*user.Lang])/18*18),
			fmt.Sprintf("set_to_lang_pagination:%d", 18), false)
		if err != nil {
			warn(err)
		}
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      message.Chat.ID,
				ReplyMarkup: kb,
			},
			Text: user.Localize("Choose language"),
		}
		app.bot.Send(msg)
		if err = app.analytics.Bot(msg, "to_lang"); err != nil {
			app.notifyAdmin(err)
		}
		return
	}

	switch user.Act {
	case "mailing":
		if err = app.bc.Put([]byte("mailing_message_id"), []byte(strconv.Itoa(message.MessageID))); err != nil {
			warn(err)
			return
		}
		if err := app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": "mailing_keyboards"}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.From.ID,
				ChannelUsername:  "",
				ReplyToMessageID: 0,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Empty"))),
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			Text:                  "теперь отправь мне кнопки\nтекст|ссылка\nтекст|ссылка",
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		})
		return
	case "mailing_keyboards":
		keyboard := tgbotapi.InlineKeyboardMarkup{}
		if message.Text != "Empty" {
			keyboard = parseKeyboard(message.Text)
		}
		if err = app.bc.Put([]byte("mailing_keyboard_raw_text"), []byte(message.Text)); err != nil {
			warn(err)
			return
		}
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": ""}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              tgbotapi.NewRemoveKeyboard(false),
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			Text:                  "проверь",
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		})

		mailingMessageId, err := app.bc.Get([]byte("mailing_message_id"))
		if err != nil {
			warn(err)
			return
		}
		mailingMessageIdInt, err := strconv.Atoi(string(mailingMessageId))
		if err != nil {
			warn(err)
			return
		}

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "start_mailing")))

		if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   message.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              &keyboard,
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			FromChatID:          config.AdminID,
			FromChannelUsername: "",
			MessageID:           mailingMessageIdInt,
			Caption:             "",
			ParseMode:           "",
			CaptionEntities:     nil,
		}); err != nil {
			warn(err)
			return
		}
		return
	}

	var text = message.Text
	message.Text = ""
	if message.Caption != "" {
		text = message.Caption
		message.Caption = ""
	}

	if text == "" {
		//app.bot.Send(tgbotapi.NewDeleteMessage())
		msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Отправь текстовое сообщение, чтобы я его перевел"))
		app.bot.Send(msg)
		app.analytics.Bot(msg, "not_text_message")
		return
	}

	from, err := app.translo.Detect(ctx, text)
	if err != nil {
		warn(err)
		return
	}
	from = strings.ToLower(from)
	if from == "" {
		log.Error("from is auto")
	} else if user.MyLang == "auto" {
		if err = app.db.UpdateUser(message.From.ID, tables.Users{MyLang: from}); err != nil {
			warn(err)
			return
		}
		user.MyLang = from
	}

	// Подбираем язык перевода, зная язык сообщения
	var to string
	if from == user.ToLang {
		to = user.MyLang
	} else if from == user.MyLang {
		to = user.ToLang // TODO: fix inline
	} else { // никакой из
		if user.ToLang == message.From.LanguageCode {
			to = user.ToLang
		} else {
			to = user.MyLang
		}
	}

	entities := message.Entities
	if len(message.CaptionEntities) > 0 {
		entities = message.CaptionEntities
	}
	text = norm.NFKC.String(helpers.ApplyEntitiesHtml(text, entities))
	tr, from, err := app.translate(ctx, from, to, text) // examples мы сохраняем, чтобы соединить с keyboard.Examples и положить в кэш
	if err != nil {
		warn(err)
		return
	}

	chunks := translate.SplitIntoChunksBySentences(tr, 4000)
	lastMsgID := 0
	for _, chunk := range chunks {
		chunk = closeUnclosedTagsAndClearUnsupported(chunk) /* + "\n❤️ @TransloBot" */ // я не знаю, почему это не работает
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
				tgbotapi.NewKeyboardButton("↔"),
				tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji)))
		//keyboard.InputFieldPlaceholder = user.Localize("Text to translate...")
		msgConfig := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.Chat.ID,
				ReplyMarkup:      keyboard,
				ReplyToMessageID: message.MessageID,
			},
			Text:                  chunk,
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}
		msg, err := app.bot.Send(msgConfig)
		lastMsgID = msg.MessageID

		if err != nil {
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: config.AdminID,
				},
				Text: fmt.Sprintf("%s\nerror with %d (%s->%s):\nText:%s", err.Error(), message.From.ID, from, to, chunk),
			})
			app.log.Error("couldn't send translation to user", zap.String("text", text), zap.String("translation", chunk))
			//pp.Println("couldn't send translation to user", chunk)
			msg, err = app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, chunk))
			if err != nil {
				warn(err)
				app.notifyAdmin(err, fmt.Sprintf("translation (%s->%s)\nOriginal text: %S", from, to, text))
				return
			}
			lastMsgID = msg.MessageID
			//warn(err)
			return
		}
		app.analytics.Bot(msgConfig, "translate")
	}
	tr, err = removeHtml(tr)
	if err != nil {
		warn(err)
		return
	}
	data, _ := translate.TTS(to, tr)
	if user.Lang == nil {
		lang := message.From.LanguageCode
		user.Lang = &lang
	}
	app.bot.Send(tgbotapi.AudioConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: message.Chat.ID,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.ToLang]+" "+flags[user.ToLang].Emoji))),
			},
			File: tgbotapi.FileBytes{
				Name:  html.UnescapeString(helpers.CutStringUTF16(tr, 50)),
				Bytes: data,
			},
		},
		Title: helpers.CutStringUTF16(tr, 40),
	})

	text, err = removeHtml(text)
	if err != nil {
		warn(err)
		return
	}
	data, _ = translate.TTS(from, html.UnescapeString(text))
	if user.Lang == nil {
		user.Lang = &message.From.LanguageCode
	}
	app.bot.Send(tgbotapi.AudioConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: message.Chat.ID,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.ToLang]+" "+flags[user.ToLang].Emoji))),
			},
			File: tgbotapi.FileBytes{
				Name:  html.UnescapeString(helpers.CutStringUTF16(text, 50)),
				Bytes: data,
			},
		},
		Title: helpers.CutStringUTF16(text, 40),
	})

	app.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           message.Chat.ID,
			ChannelUsername:  "",
			ProtectContent:   false,
			ReplyToMessageID: lastMsgID,
			ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("❌", "wrong_translation:"+from+":"+to),
					tgbotapi.NewInlineKeyboardButtonData("✅", "correct_translation"),
				)),
			DisableNotification:      false,
			AllowSendingWithoutReply: false,
		},
		Text: user.Localize("Did I translate it correctly?"),
	})

}
