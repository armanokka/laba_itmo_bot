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
	"golang.org/x/text/unicode/norm"
	"gorm.io/gorm"
	"net/url"
	"os"
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
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
		}
	}()

	user := tables.Users{Lang: message.From.LanguageCode}

	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Произошла ошибка")))
		app.notifyAdmin(err)
	}

	if message.Chat.ID < 0 {
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
			if message.From.LanguageCode == "" {
				message.From.LanguageCode = "en"
			}
			err = app.db.CreateUser(tables.Users{
				ID:           message.From.ID,
				MyLang:       "",
				ToLang:       "",
				Act:          "setup_langs",
				Usings:       0,
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
	log = log.With(zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang), zap.Int("usings", user.Usings))

	switch message.Command() {
	case "start":
		fallthrough
	case "help":
		log.Info("/start or /help")
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

		if err = app.db.UpdateUser(message.From.ID, tables.Users{Act: "setup_langs"}); err != nil {
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
		keyboard := &tgbotapi.InlineKeyboardMarkup{}
		if message.Text != "Empty" {
			keyboard = parseKeyboard(message.Text)
			if err = app.bc.Put([]byte("mailing_keyboard_raw_text"), []byte(strconv.Itoa(message.MessageID))); err != nil {
				warn(err)
				return
			}
		}
		var withKeyboard bool
		if len(keyboard.InlineKeyboard) > 0 {
			withKeyboard = true
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
		if err = app.db.DropMailings(); err != nil {
			warn(err)
			return
		}
		if err = app.db.CreateMailingTable(); err != nil {
			warn(err)
			return
		}

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

		if withKeyboard {
			app.bot.Send(tgbotapi.CopyMessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   message.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              keyboard,
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
		} else {
			app.bot.Send(tgbotapi.CopyMessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   message.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              nil,
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
			if withKeyboard {
				if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:                   id,
						ChannelUsername:          "",
						ReplyToMessageID:         0,
						ReplyMarkup:              keyboard,
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
			} else {
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
					MessageID:           mailingMessageIdInt,
					Caption:             "",
					ParseMode:           "",
					CaptionEntities:     nil,
				}); err != nil {
					pp.Println(err)
				}
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
		if message.Text == "" {
			app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Отправь текстовое сообщение, чтобы я его перевел")))
			return
		}
		fromLang, err := translate.GoogleTranslate(ctx, "auto", "en", cutStringUTF16(message.Text, 100))
		if err != nil {
			warn(err)
			return
		}
		from := fromLang.FromLang
		//app.SuperTranslate(ctx, user, message.Chat.ID, from, to, text, entities)

		keyboard, err := buildLangsPagination(user, 0, 18, fromLang.FromLang, fmt.Sprintf("setup_langs:%s:%s", from, "%s"), fmt.Sprintf("setup_langs_pagination:%s:0", from), fmt.Sprintf("setup_langs_pagination:%s:18", from))
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

	from, err := translate.DetectLanguageYandex(ctx, cutStringUTF16(text, 100))
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

	entities := message.Entities
	if len(message.CaptionEntities) > 0 {
		entities = message.CaptionEntities
	}
	text = applyEntitiesHtml(norm.NFKC.String(text), entities)
	if err = app.SuperTranslate(ctx, user, message.Chat.ID, message.MessageID, from, to, text, false); err != nil && !errors.Is(err, context.Canceled) {
		warn(err)
		if e, ok := err.(errors.Error); ok {
			log.Error("", zap.Error(e), zap.String("stack", string(e.Stack())))
		} else {
			log.Error("", zap.Error(err))
		}
		return
	}

	app.analytics.Bot(user.ID, "", "Translated")

	if err = app.db.LogBotMessage(message.From.ID, "pm_translate", ""); err != nil {
		app.notifyAdmin(fmt.Errorf("%w", err))
	}

	if user.Usings == 5 || (user.Usings > 0 && user.Usings%20 == 0) {
		if _, err := app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.From.ID,
				ChannelUsername:  "",
				ReplyToMessageID: 0,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL(user.Localize("Рассказать про нас"), "http://t.me/share/url?url="+url.PathEscape(user.Localize("Я рекомендую @translobot"))))),
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
	}

}
