package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/helpers"
	"github.com/armanokka/translobot/pkg/lingvo"
	"github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"html"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func (app *App) onMessage(ctx context.Context, message tgbotapi.Message) {
	log := app.log.With(zap.Int64("id", message.From.ID))
	defer func() {
		if err := recover(); err != nil {
			_, f, line, ok := runtime.Caller(2)
			if ok {
				log = log.With(zap.String("caller", f+":"+strconv.Itoa(line)))
			}
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
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Excuses")))
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
				TTS:          true,
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
	log = log.With(zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang), zap.Stringp("act", user.Act), zap.String("message", message.Text), zap.String("caption", message.Caption))
	log.Debug("new message")
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
				btn := tgbotapi.NewInlineKeyboardButtonData(name, "set_bot_lang:"+code)
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
			if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
				warn(err)
				return
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
			"ar": "عربي🇪🇬",
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
			btn := tgbotapi.NewInlineKeyboardButtonData(name, "set_bot_lang:"+code)
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
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
		}
		return
	case "id":
		log.Info("/id")
		msg := tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10))
		app.bot.Send(msg)
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
		}
		return
	case "mailing":
		log.Info("/mailing")
		mailing := "mailing"
		if err = app.db.UpdateUser(message.From.ID, tables.Users{Act: &mailing}); err != nil {
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
	case "tts_on":
		app.bot.Send(tgbotapi.NewMessage(message.From.ID, user.Localize(`Бот будет озвучивать переводы`)))
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"tts": true}); err != nil {
			warn(err)
		}
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
		}
		return
	case "tts_off":
		app.bot.Send(tgbotapi.NewMessage(message.From.ID, user.Localize(`Бот больше не будет озвучивать переводы`)))
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"tts": false}); err != nil {
			warn(err)
		}
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
		}
		return
	}
	if user.Lang == nil {
		user.Lang = &message.From.LanguageCode
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
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[*user.Lang][user.ToLang]+" "+flags[user.ToLang].Emoji))),
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
		if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
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
			Text: user.Localize("Choose language or send its name"),
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
			Text: user.Localize("Choose language or send its name"),
		}
		// TODO handle app.bot.send errors that are not 403
		app.bot.Send(msg)
		if err = app.analytics.Bot(msg, "to_lang"); err != nil {
			app.notifyAdmin(err)
		}
		return
	}

	if user.Act != nil {
		switch *user.Act {
		case "set_my_lang":
			app.bot.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))

			filter := "" // TODO implement filter_to_lang as I did filter_my_lang
			// TODO понимать другую раскладку
			for _, ch := range message.Text {
				if !unicode.IsLetter(ch) && ch != '-' {
					continue
				}
				filter += string(ch)
			}
			filter = strings.ToLower(filter)
			searchSet := []string{*user.Lang, message.From.LanguageCode, user.MyLang, user.ToLang}
			results := make([]string, 0, 2)
			keyboard := tgbotapi.NewInlineKeyboardMarkup()
			for i, set := range searchSet {
				if i < len(searchSet)-1 && in(searchSet[i+1:], set) {
					continue // уже искали в этом сете
				}
				for code, name := range langs[set] { // TODO filter and sort by less differnece
					if !hasPrefix(name, filter, 1) && filter != code || in(results, code) {
						continue
					}
					results = append(results, code)
				}
			}
			sort.Slice(results, func(i, j int) bool {
				return diff(filter, results[i]) < diff(filter, results[j])
			})
			for _, code := range results {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(flags[code].Emoji+" "+langs[*user.Lang][code], "filtered_set_my_lang:"+code)))
			}
			msgIDsBytes, err := app.bc.Get([]byte(strconv.FormatInt(message.From.ID, 10)))
			if err != nil {
				warn(err)
				return
			}
			msgIDs := strings.Split(string(msgIDsBytes), ";") // msgIDs[0] - search query message. msgIDs[1] - languages pagination message.
			if len(msgIDs) != 2 {
				warn(fmt.Errorf("strings.Split(app.bc.Get(message.From.ID), \";\") is not 2 chunks"))
				return
			}
			msgID, err := strconv.ParseInt(msgIDs[0], 10, 64)
			if err != nil {
				warn(err)
				log.Error("couldn't parse int64: app.bc.Get(message.From.ID)", zap.Error(err), zap.String("result", string(msgIDsBytes)))
				return
			}
			// TODO send new message and delete previous instead of editing
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(user.Localize(`❌ Cancel`), "close_type_language_name_menu:my_lang"),
				tgbotapi.NewInlineKeyboardButtonData(user.Localize(`🔄Try again`), `try_again_to_search_my_lang`)))
			if len(results) == 0 {
				keyboard.InlineKeyboard = keyboard.InlineKeyboard[len(keyboard.InlineKeyboard)-1:]
				app.bot.Send(tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      message.From.ID,
						MessageID:   int(msgID),
						ReplyMarkup: &keyboard,
					},
					Text:      user.Localize(`No languages found starting with <b>%s</b>`, filter),
					ParseMode: tgbotapi.ModeHTML,
				})
				return
			}

			app.bot.Send(tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      message.From.ID,
					MessageID:   int(msgID),
					ReplyMarkup: &keyboard,
				},
				Text: user.Localize(`🔎 Found %d languages starting with <b>%s</b>. Tap on language to choose it`, len(results), filter),
				//Text:      fmt.Sprintf("%s\n%s", user.Localize(`Не найдены языки, начинающиеся на %s`, fmt.Sprintf(`<b>%s</b>`, filter)), user.Localize(`Напишите название языка, который вы хотите выбрать:`)),
				ParseMode: tgbotapi.ModeHTML,
			})
			return
		case "set_to_lang":
			app.bot.Send(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))

			filter := "" // TODO implement filter_to_lang as I did filter_my_lang
			// TODO понимать другую раскладку
			for _, ch := range message.Text {
				if !unicode.IsLetter(ch) && ch != '-' {
					continue
				}
				filter += string(ch)
			}
			filter = strings.ToLower(filter)
			searchSet := []string{*user.Lang, message.From.LanguageCode, user.MyLang, user.ToLang}
			usedLangs := make([]string, 0, 2)
			keyboard := tgbotapi.NewInlineKeyboardMarkup()
			for i, set := range searchSet {
				if i < len(searchSet)-1 && in(searchSet[i+1:], set) {
					continue // уже искали в этом сете
				}
				for code, name := range langs[set] { // TODO filter and sort by less differnece
					if !hasPrefix(name, filter, 1) && filter != code || in(usedLangs, code) {
						continue
					}
					usedLangs = append(usedLangs, code)
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(flags[code].Emoji+" "+langs[*user.Lang][code], "filtered_set_to_lang:"+code)))
				}
			}
			msgIDsBytes, err := app.bc.Get([]byte(strconv.FormatInt(message.From.ID, 10)))
			if err != nil {
				warn(err)
				return
			}
			msgIDs := strings.Split(string(msgIDsBytes), ";") // msgIDs[0] - search query message. msgIDs[1] - languages pagination message.
			if len(msgIDs) != 2 {
				warn(fmt.Errorf("strings.Split(app.bc.Get(message.From.ID), \";\") is not 2 chunks"))
				return
			}
			msgID, err := strconv.ParseInt(msgIDs[0], 10, 64)
			if err != nil {
				warn(err)
				log.Error("couldn't parse int64: app.bc.Get(message.From.ID)", zap.Error(err), zap.String("result", string(msgIDsBytes)))
				return
			}
			// TODO send new message and delete previous instead of editing
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(user.Localize(`❌ Cancel`), "close_type_language_name_menu:to_lang"),
				tgbotapi.NewInlineKeyboardButtonData(user.Localize(`🔄Try again`), `try_again_to_search_to_lang`)))
			if len(usedLangs) == 0 {
				keyboard.InlineKeyboard = keyboard.InlineKeyboard[len(keyboard.InlineKeyboard)-1:]
				app.bot.Send(tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:      message.From.ID,
						MessageID:   int(msgID),
						ReplyMarkup: &keyboard,
					},
					Text:      user.Localize(`No languages found starting with <b>%s</b>`, filter),
					ParseMode: tgbotapi.ModeHTML,
				})
				return
			}

			app.bot.Send(tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      message.From.ID,
					MessageID:   int(msgID),
					ReplyMarkup: &keyboard,
				},
				Text: user.Localize(`🔎 Found %d languages starting with <b>%s</b>. Tap on language to choose it`, len(usedLangs), filter),
				//Text:      fmt.Sprintf("%s\n%s", user.Localize(`Не найдены языки, начинающиеся на %s`, fmt.Sprintf(`<b>%s</b>`, filter)), user.Localize(`Напишите название языка, который вы хотите выбрать:`)),
				ParseMode: tgbotapi.ModeHTML,
			})
			return
		case "mailing":
			if err = app.bc.Put([]byte("mailing_message_id"), []byte(strconv.Itoa(message.MessageID))); err != nil {
				warn(err)
				return
			}
			if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": "mailing_keyboards"}); err != nil {
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

			if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
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

	detection, err := app.translo.Translate(ctx, "auto", "en", text)
	if err != nil {
		warn(err)
		return
	}
	from := strings.ToLower(detection.TextLang)
	if from == "" {
		log.Error("from is auto")
	} else if user.MyLang == "auto" {
		if err = app.db.UpdateUser(message.From.ID, tables.Users{MyLang: from}); err != nil {
			warn(err)
			return
		}
		user.MyLang = from
	}
	log = log.With(zap.String("from", from))

	// Подбираем язык перевода, зная язык сообщения
	var to string
	if from == user.ToLang {
		to = user.MyLang
	} else if from == user.MyLang {
		to = user.ToLang // TODO: fix inline
	} else { // никакой из
		if user.ToLang == *user.Lang {
			to = user.ToLang
		} else {
			to = user.MyLang
		}
	}
	log = log.With(zap.String("to", to))

	entities := message.Entities
	if len(message.CaptionEntities) > 0 {
		entities = message.CaptionEntities
	}
	log = log.With(zap.String("source", text))
	text = norm.NFKC.String(helpers.ApplyEntitiesHtml(text, entities))
	log = log.With(zap.String("source_with_html", text))

	var trMylangTolang, trTolangMylang, trFromTo, trDict string
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		tr, err := app.translo.Translate(ctx, user.MyLang, user.ToLang, text)
		trMylangTolang = tr.TranslatedText
		return err
	})
	g.Go(func() error {
		tr, err := app.translo.Translate(ctx, user.ToLang, user.MyLang, text)
		trTolangMylang = tr.TranslatedText
		return err
	})
	g.Go(func() error {
		tr, err := app.translo.Translate(ctx, from, to, text)
		trFromTo = tr.TranslatedText
		return err
	})
	if from == user.MyLang {
		_, ok1 := lingvo.Lingvo[from]
		_, ok2 := lingvo.Lingvo[to]
		if ok1 && ok2 && len(text) < 50 && !strings.ContainsAny(text, " \r\n") {
			g.Go(func() error {
				l, err := lingvo.GetDictionary(ctx, from, to, strings.ToLower(text))
				if err != nil {
					if IsCtxError(err) {
						return nil
					}
					log.Error("lingvo err", zap.Error(err))
					return err
				}
				tr := strings.TrimSpace(writeLingvo(l))
				if tr != "" {
					trDict = tr + "\n❤️ @TransloBot"
				}
				return nil
			})
		}
	}

	if err = g.Wait(); err != nil {
		log.Error("", zap.Error(err))
		app.notifyAdmin(err)
	}

	// request translations: user.from-user.to, user.to-user.from, from-to

	//tr, from, err := app.translate(ctx, from, to, text) // examples мы сохраняем, чтобы соединить с keyboard.Examples и положить в кэш
	//if err != nil {
	//	log.Error("", zap.Error(err))
	//	warn(err)
	//	return
	//}
	//// TODO каждые 24ч капча
	//if from == to && user.MyLang == user.ToLang && tr == text {
	//	app.bot.Send(tgbotapi.NewMessage(message.From.ID, user.Localize(`You translate from one language to the same language`)))
	//}

	//app.bc.PutWithTTL() УПОМЯНУТЬ ОБ АПИ
	pp.Println([]string{trMylangTolang, trTolangMylang, trFromTo, trDict})
	tr := maxDiff(text, []string{trMylangTolang, trTolangMylang, trFromTo, trDict})
	if strings.TrimSpace(tr) == "" {
		log.Error("empty translation", zap.String("trMylangTolang", trMylangTolang), zap.String("trTolangMylang", trFromTo), zap.String("trFromTo", trFromTo), zap.String("trDict", trDict))
	} // TODO fix closeUnclosedTagsAndClearUnsupported
	log.Debug(tr)
	chunks := translate.SplitIntoChunksBySentences(tr, 4000)
	lastMsgID := 0
	for _, chunk := range chunks {
		chunk = closeUnclosedTagsAndClearUnsupported(chunk) /* + "\n❤️ @TransloBot" */ // я не знаю, почему это не работает
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(langs[*user.Lang][user.MyLang]+" "+flags[user.MyLang].Emoji),
				tgbotapi.NewKeyboardButton("↔"),
				tgbotapi.NewKeyboardButton(langs[*user.Lang][user.ToLang]+" "+flags[user.ToLang].Emoji)))
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
			app.log.Error("couldn't send translation to user", zap.String("text", text), zap.String("translation", chunk), zap.Error(err))
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
	if user.TTS {
		tr, err = removeHtml(tr)
		if err != nil {
			warn(err)
			return
		}
		data, _ := translate.TTS(to, tr)
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
	}

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

	if user.Usings == 3 || user.Usings == 6 || user.Usings == 10 {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize(`Закрепите чат с ботом, чтобы не искать его! 📌`)))
	}
	if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": nil}); err != nil {
		warn(err)
		return
	}

}
