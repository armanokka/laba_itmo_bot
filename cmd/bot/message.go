package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	user := tables.Users{Lang: message.From.LanguageCode}

	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Произошла ошибка")))
		app.notifyAdmin(err)
	}

	if message.Chat.ID < 0 {
		app.onGroupMessage(message)
		return
	}

	defer func() {
		app.analytics.User(message.Text, message.From)
		if message.Caption != "" {
			message.Text = message.Caption
		}
		if err := app.db.UpdateUserMetrics(message.From.ID, message.Text); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	var err error
	user, err = app.db.GetUserByID(message.From.ID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			tolang := ""
			if message.From.LanguageCode == "" || message.From.LanguageCode == "en" {
				message.From.LanguageCode = "en"
				tolang = "ru"
			} else if message.From.LanguageCode == "ru" {
				tolang = "en"
			}
			if err = app.db.CreateUser(tables.Users{
				ID:           message.From.ID,
				MyLang:       message.From.LanguageCode,
				ToLang:       tolang,
				Act:          "",
				Usings:       0,
				Blocked:      false,
				LastActivity: time.Now(),
			}); err != nil {
				warn(err)
				return
			}
			user.MyLang = message.From.LanguageCode
			user.ToLang = tolang
		} else {
			warn(err)
			return
		}
	}
	user.SetLang(message.From.LanguageCode)
	log = log.With(zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang))

	switch message.Command() {
	case "start":
		// TODO: быстрая смена текста hello mi hao ola
		if user.Blocked { // разбанил
			app.analytics.User("{bot_was_UNblocked}", message.From)
			app.bot.Send(tgbotapi.NewSticker(message.From.ID, tgbotapi.FileID("CAACAgIAAxkBAAEP5w5iif1KBEzJZ-6N49pvKBvTcz5BYwACBAEAAladvQreBNF6Zmb3bCQE")))
			app.bot.Send(tgbotapi.NewMessage(message.From.ID, user.Localize("Рад видеть вас снова.\nС возвращением.")))
			if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"blocked": false}); err != nil {
				app.notifyAdmin(fmt.Errorf("%w", err))
			}
		} else {
			hello := map[string]string{
				"en": "Hello",
				"zh": "Nǐ hǎo",
				"uk": "Привіт",
				"ar": "As-salām ‘alaykum",
				"pt": "Olá",
				"fr": "Salut",
				"fy": "Gouden Dai",
			}
			if _, ok := hello[user.Lang]; ok {
				delete(hello, user.Lang)
			}
			var msg tgbotapi.Message
			i := 0
			for _, text := range hello {
				if msg.MessageID == 0 {
					msg, err = app.bot.Send(tgbotapi.NewMessage(message.From.ID, text))
					i++
					continue
				}
				app.bot.Send(tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:          message.From.ID,
						ChannelUsername: "",
						MessageID:       msg.MessageID,
						InlineMessageID: "",
						ReplyMarkup:     nil,
					},
					Text:                  text,
					ParseMode:             "",
					Entities:              nil,
					DisableWebPagePreview: false,
				})
				time.Sleep(time.Second / 3)
				i++
			}
			app.bot.Send(tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:          message.From.ID,
					ChannelUsername: "",
					MessageID:       msg.MessageID,
					InlineMessageID: "",
					ReplyMarkup:     nil,
				},
				Text:                  user.Localize("Привет"),
				ParseMode:             "",
				Entities:              nil,
				DisableWebPagePreview: false,
			})
			if _, err = app.bot.Send(tgbotapi.NewSticker(message.From.ID, tgbotapi.FileID("CAACAgIAAxkBAAEP5rViieLUfMyMYArLNLl12AOggTEAAVAAAgEBAAJWnb0KIr6fDrjC5jQkBA"))); err != nil {
				warn(err)
			}
			time.Sleep(time.Second)
		}

		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.From.ID,
				ChannelUsername:  "",
				ReplyToMessageID: 0,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔️"),
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji))),
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text: user.Localize("Просто напиши мне текст, а я его переведу"),
		}); err != nil {
			warn(err)
		}

		if err = app.db.UpdateUser(message.From.ID, tables.Users{Act: ""}); err != nil {
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
		if err = app.db.LogBotMessage(message.From.ID, "pm_users", "shared users' ids"); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
		return
	case "id":
		log.Info("/id")
		msg := tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10))
		app.bot.Send(msg)
		if err = app.db.LogBotMessage(message.From.ID, "pm_id", msg.Text); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}

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
	case "analytics":
		log.Info("/analytics")
		user, err := app.db.GetRandomUser()
		if err != nil {
			warn(err)
			return
		}
		fmt.Println("Getting logs for", user.ID)
		logs, err := app.db.GetUserLogs(user.ID, 10)
		if err != nil {
			warn(err)
			return
		}
		msg := ""
		for _, log := range logs {
			if log.FromBot {
				msg += "🤖: "
				switch log.Intent.String {
				case "cb_meaning":
					msg += "<i>Lookup meaning</i>"
				case "cb_exmp":
					msg += "<i>Open examples</i>"
				case "cb_dict":
					msg += "<i>Lookup in dictionary</i>"
				case "bot_was_blocked":
					msg += "<i>Bot was blocked</i>"
				case "bot_was_unblocked":
					msg += "<i>Bot was unblocked</i>"
				case "inline_succeeded":
					msg += "<i>Inline query was handled</i>"
				}
				msg += " " + log.Text
			} else {
				msg += "👤:" + log.Text
			}

			msg += "\n"

		}
		if _, err = app.bot.Send(tgbotapi.NewMessage(message.From.ID, msg)); err != nil {
			fmt.Println(err)
		}
		return
	}

	switch message.Text {
	case "↔️":
		if err = app.db.SwapLangs(message.Chat.ID); err != nil {
			warn(err)
			return
		}
		user.MyLang, user.ToLang = user.ToLang, user.MyLang
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: message.From.ID,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔️"),
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji))),
			},
			Text: "OK",
		})
		return
	case concatNonEmpty(" ", langs[message.From.LanguageCode][user.MyLang], flags[user.MyLang].Emoji):
		kb, err := buildLangsPagination(user, 0, 18, "",
			fmt.Sprintf("set_my_lang:%s:%d", "%s", 0),
			fmt.Sprintf("set_my_lang_pagination:%d", len(codes[user.Lang])/18*18),
			fmt.Sprintf("set_my_lang_pagination:%d", 18))
		if err != nil {
			warn(err)
		}
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      message.Chat.ID,
				ReplyMarkup: kb,
			},
			Text: user.Localize("Choose language"),
		})
		return
	case concatNonEmpty(" ", langs[message.From.LanguageCode][user.ToLang], flags[user.ToLang].Emoji):
		kb, err := buildLangsPagination(user, 0, 18, "",
			fmt.Sprintf("set_to_lang:%s:%d", "%s", 0),
			fmt.Sprintf("set_to_lang_pagination:%d", len(codes[user.Lang])/18*18),
			fmt.Sprintf("set_to_lang_pagination:%d", 18))
		if err != nil {
			warn(err)
		}
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      message.Chat.ID,
				ReplyMarkup: kb,
			},
			Text: user.Localize("Choose language"),
		})
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
			if err = app.bc.Put([]byte("mailing_keyboard_raw_text"), []byte(message.Text)); err != nil {
				warn(err)
				return
			}
		}

		if err := app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"act": ""}); err != nil {
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
			Text:                  "рассылка начата",
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

		app.bot.Send(tgbotapi.CopyMessageConfig{
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
		})

		usersNumber, err := app.db.GetUsersNumber()
		if err != nil {
			return
		}

		slice := make([]int64, 0, 100)
		var i int64
		for ; usersNumber/100 < i; i++ { // iterate over each 100 users
			offset := i*100 + 100
			if err = app.db.GetUsersSlice(offset, 100, slice); err != nil {
				warn(err)
				return
			}
			for j := 0; j < 100; j++ {
				if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:                   slice[j],
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
					pp.Println(err)
				}
				log.Info("mailing was sent", zap.Int64("recepient_id", slice[j]), zap.Int64("queue_position", i*100+int64(j)))
			}
		}

		app.bot.Send(tgbotapi.NewMessage(message.From.ID, "рассылка закончена"))
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
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Отправь текстовое сообщение, чтобы я его перевел")))
		app.analytics.Bot(message.Chat.ID, "Please, send text message", "Message is not text message")
		return
	}

	from, err := translate.DetectLanguageGoogle(ctx, cutStringUTF16(text, 100))
	if err != nil {
		warn(err)
		return
	}
	if strings.Contains(from, "-") {
		parts := strings.Split(from, "-")
		from = parts[0]
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

	app.bot.Send(tgbotapi.NewChatAction(message.From.ID, "typing")) // TODO: encode translation to utf-8
	if err = app.SuperTranslate(ctx, user, message.Chat.ID, from, to, text, message); err != nil && !errors.Is(err, context.Canceled) {
		err = fmt.Errorf("%s\nuser's text:%s\ntranslation:%s", err.Error(), text)
		warn(err)
		if e, ok := err.(errors.Error); ok {
			log.Error("", zap.Error(e), zap.String("stack", string(e.Stack())))
		} else {
			log.Error("", zap.Error(err))
		}
		return
	}

	app.analytics.Bot(user.ID, "", "Translated")

}
