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
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

	user := tables.Users{Lang: message.From.LanguageCode}

	warn := func(err error) {
		app.bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Произошла ошибка")))
		app.notifyAdmin(err)
	}

	if message.Chat.ID < 0 {
		app.onGroupMessage(message)
		return
	}

	app.bot.Send(tgbotapi.NewChatAction(message.Chat.ID, "typing"))

	defer func() {
		if err := app.analytics.User(message); err != nil {
			app.notifyAdmin(err)
		}
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
			if err = app.db.CreateUser(tables.Users{
				ID:           message.From.ID,
				MyLang:       "auto",
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
		if user.Blocked { // разбанил
			if err = app.analytics.User(message); err != nil {
				app.notifyAdmin(err)
				//return
			}
			app.bot.Send(tgbotapi.NewSticker(message.From.ID, tgbotapi.FileID("CAACAgIAAxkBAAEP5w5iif1KBEzJZ-6N49pvKBvTcz5BYwACBAEAAladvQreBNF6Zmb3bCQE")))
			app.bot.Send(tgbotapi.NewMessage(message.From.ID, user.Localize("Рад видеть вас снова.\nС возвращением.")))

			if err = app.db.UpdateUserByMap(message.From.ID, map[string]interface{}{"blocked": false}); err != nil {
				app.notifyAdmin(fmt.Errorf("%w", err))
			}
		} else {
			if _, err = app.bot.Send(tgbotapi.NewSticker(message.From.ID, tgbotapi.FileID("CAACAgIAAxkBAAEQJrlinQ2sIDF1R3cjISx_cEv1pawdSgACQhAAAjPFKUmQDtQRpypKgiQE"))); err != nil {
				warn(err)
			}
			time.Sleep(time.Second)
		}
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           message.From.ID,
				ChannelUsername:  "",
				ReplyToMessageID: 0,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji))),
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text: user.Localize("Send me (+ words, or) posts from foreign channels 📣 and I'll translate them."),
		}
		if _, err = app.bot.Send(msg); err != nil {
			warn(err)
		}

		if err = app.analytics.Bot(msg, "/start"); err != nil {
			app.notifyAdmin(err)
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

	switch message.Text {
	case "↔":
		if err = app.db.SwapLangs(message.Chat.ID); err != nil {
			warn(err)
			return
		}
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

		if err = app.analytics.Bot(msg, "↔"); err != nil {
			app.notifyAdmin(err)
		}
		return
	case concatNonEmpty(" ", langs[message.From.LanguageCode][user.MyLang], flags[user.MyLang].Emoji):
		kb, err := buildLangsPagination(user, 0, 18, "",
			fmt.Sprintf("set_my_lang:%s:%d", "%s", 0),
			fmt.Sprintf("set_my_lang_pagination:%d", len(codes[user.Lang])/18*18),
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
	case concatNonEmpty(" ", langs[message.From.LanguageCode][user.ToLang], flags[user.ToLang].Emoji):
		kb, err := buildLangsPagination(user, 0, 18, "",
			fmt.Sprintf("set_to_lang:%s:%d", "%s", 0),
			fmt.Sprintf("set_to_lang_pagination:%d", len(codes[user.Lang])/18*18),
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

	// Определяем язык сообщения

	// Подбираем язык перевода, зная язык сообщения
	from, to := "", ""
	if user.MyLang == "auto" {
		from, err = translate.DetectLanguageGoogle(ctx, text)
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
		to = user.ToLang
	} else {
		if from == user.ToLang {
			to = user.MyLang
		} else if from == user.MyLang {
			to = user.ToLang
		} else { // никакой из
			//if user.ToLang == message.From.LanguageCode {
			//	to = user.ToLang
			//} else {
			to = user.MyLang
			//}
		}

		if from != user.MyLang {
			tr, err := translate.GoogleTranslate(ctx, from, to, text)
			if err != nil {
				warn(err)
				return
			}
			if diff(text, tr.Text) < 2 {
				from = user.MyLang
				to = user.ToLang
			}
		}
	}

	entities := message.Entities
	if len(message.CaptionEntities) > 0 {
		entities = message.CaptionEntities
	}
	text = helpers.ApplyEntitiesHtml(norm.NFKC.String(text), entities)

	g, ctx := errgroup.WithContext(ctx)
	if message.ReplyMarkup != nil {
		for i1, row := range message.ReplyMarkup.InlineKeyboard {
			i1 := i1
			row := row
			for i2, btn := range row {
				i2 := i2
				btn := btn
				g.Go(func() error {
					tr, err := translate.GoogleTranslate(ctx, from, to, btn.Text)
					if err != nil {
						return errors.Wrap(err)
					}
					message.ReplyMarkup.InlineKeyboard[i1][i2].Text = tr.Text
					return nil
				})
			}
		}
	}

	if message.Poll != nil {
		g.Go(func() error {
			tr, err := translate.GoogleTranslate(ctx, from, to, message.Poll.Question)
			message.Poll.Question = tr.Text
			return errors.Wrap(err)
		})
		g.Go(func() error {
			tr, err := translate.GoogleTranslate(ctx, from, to, helpers.ApplyEntitiesHtml(message.Poll.Explanation, message.Poll.ExplanationEntities))
			message.Poll.Explanation = tr.Text
			return errors.Wrap(err)
		})
		for i, q := range message.Poll.Options {
			i := i
			q := q
			g.Go(func() error {
				tr, err := translate.GoogleTranslate(ctx, from, to, q.Text)
				message.Poll.Options[i].Text = tr.Text
				return errors.Wrap(err)
			})
		}
	}

	pp.Println(user, from, to)

	tr, from, err := app.translate(ctx, from, to, text) // examples мы сохраняем, чтобы соединить с keyboard.Examples и положить в кэш
	if err != nil {
		warn(err)
		return
	}
	if !validHtml(tr) {
		tr = html.EscapeString(tr)
	}

	//app.bot.Send(tgbotapi.NewDeleteMessage(chatID, message.MessageID))
	chunks := translate.SplitIntoChunksBySentences(tr, 4096)
	for _, chunk := range chunks {
		chunk = closeUnclosedTags(chunk)
		switch {
		case message.Poll != nil:
			options := make([]string, 0, len(message.Poll.Options))
			for _, opt := range message.Poll.Options {
				options = append(options, opt.Text)
			}
			_, err = app.bot.Send(tgbotapi.SendPollConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      message.Chat.ID,
					ReplyMarkup: message.ReplyMarkup,
				},
				Question:              message.Poll.Question,
				Options:               options,
				IsAnonymous:           message.Poll.IsAnonymous,
				Type:                  message.Poll.Type,
				AllowsMultipleAnswers: message.Poll.AllowsMultipleAnswers,
				CorrectOptionID:       int64(message.Poll.CorrectOptionID),
				Explanation:           message.Poll.Explanation,
				ExplanationParseMode:  tgbotapi.ModeHTML,
				ExplanationEntities:   nil,
				OpenPeriod:            message.Poll.OpenPeriod,
				CloseDate:             message.Poll.CloseDate,
				IsClosed:              message.Poll.IsClosed,
			})
		case message.Audio != nil:
			thumbnail := ""
			if message.Audio.Thumbnail != nil {
				thumbnail = message.Audio.Thumbnail.FileID
			}
			_, err = app.bot.Send(tgbotapi.AudioConfig{
				BaseFile: tgbotapi.BaseFile{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           message.Chat.ID,
						ReplyToMessageID: 0,
						ReplyMarkup:      message.ReplyMarkup,
					},
					File: tgbotapi.FileID(message.Audio.FileID),
				},
				Thumb:     tgbotapi.FileID(thumbnail),
				Caption:   tr,
				ParseMode: tgbotapi.ModeHTML,
				Duration:  message.Audio.Duration,
				Performer: message.Audio.Performer,
				Title:     message.Audio.Title,
			})
		case len(message.Photo) > 0:
			chunk = helpers.CutStringUTF16(chunk, 1024) // MEDIA_CAPTION_TOO_LONG
			maxResolutionPhoto := message.Photo[len(message.Photo)-1]
			_, err = app.bot.Send(tgbotapi.PhotoConfig{
				BaseFile: tgbotapi.BaseFile{
					BaseChat: tgbotapi.BaseChat{
						ChatID:      message.Chat.ID,
						ReplyMarkup: message.ReplyMarkup,
					},
					File: tgbotapi.FileID(maxResolutionPhoto.FileID),
				},
				Caption:   chunk,
				ParseMode: tgbotapi.ModeHTML,
			})
		default:
			var keyboard interface{}
			if user.MyLang == "auto" {
				if err = app.db.UpdateUser(message.From.ID, tables.Users{MyLang: from}); err != nil {
					warn(err)
					return
				}
				user.MyLang = from
				keyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[message.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji)))
			}
			msg := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      message.Chat.ID,
					ReplyMarkup: keyboard,
				},
				Text:                  chunk,
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			}
			_, err = app.bot.Send(msg)
			app.analytics.Bot(msg, "translate")

		}
		if err != nil {
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   config.AdminID,
					ChannelUsername:          "",
					ProtectContent:           false,
					ReplyToMessageID:         0,
					ReplyMarkup:              nil,
					DisableNotification:      false,
					AllowSendingWithoutReply: false,
				},
				Text:                  fmt.Sprintf("Error: %s\nUser's text:%s\nTranslation:%s", err.Error(), text, tr),
				ParseMode:             "",
				Entities:              nil,
				DisableWebPagePreview: false,
			})
			app.log.Error("couldn't send translation to user", zap.String("text", text), zap.String("translation", chunk))
			warn(err)
			return
		}
	}
	data, err := translate.TTS(to, tr)
	if err != nil {
		//return err
	}
	app.bot.Send(tgbotapi.AudioConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: message.Chat.ID,
				ReplyMarkup: tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[user.Lang][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[user.Lang][user.ToLang]+" "+flags[user.ToLang].Emoji))),
			},
			File: tgbotapi.FileBytes{
				Name:  helpers.CutStringUTF16(tr, 50),
				Bytes: data,
			},
		},
		Thumb:           nil,
		Caption:         "",
		ParseMode:       "",
		CaptionEntities: nil,
		Duration:        0,
		Performer:       "",
		Title:           "",
	})

}
