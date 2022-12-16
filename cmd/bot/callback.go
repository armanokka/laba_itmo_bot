package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func (app *App) onCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	log := app.log.With(zap.Int64("id", callback.From.ID))

	//go func() {
	//	if err := app.analytics.UserButtonClick(*callback.From, callback.Data); err != nil {
	//		app.notifyAdmin(err)
	//	}
	//}()
	defer func() {
		if err := recover(); err != nil {
			_, f, line, ok := runtime.Caller(4)
			if ok {
				log = log.With(zap.String("caller", f+":"+strconv.Itoa(line)))
			}
			if e, ok := err.(error); ok {
				log.Error("", zap.Error(e))
			} else {
				log.Error("", zap.Any("error", err), zap.String("stack_trace", string(debug.Stack())))
			}
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
		}
	}()
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewCallback(callback.ID, "Error, sorry"))
		app.notifyAdmin(err)
		log.Error("", zap.Error(err))
	}

	var err error
	user, err := app.db.GetUserByID(callback.From.ID)
	if err != nil {
		warn(err)
	}
	if user.Lang == nil {
		user.Lang = &callback.From.LanguageCode
	}

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(10) == 0 {
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, user.Localize("Too many requests. Try again in 10 seconds")))
		return
	}
	arr := strings.Split(callback.Data, ":")

	log = log.With(zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang), zap.Stringp("act", user.Act), zap.String("callback_data", callback.Data))
	log.Debug("new callback")
	switch arr[0] {
	case "set_bot_lang":
		if err = app.db.UpdateUser(callback.From.ID, tables.Users{Lang: &arr[1]}); err != nil {
			app.notifyAdmin(err)
		}
		user.Lang = &arr[1]
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		app.bot.Send(tgbotapi.NewSticker(callback.From.ID, tgbotapi.FileID(`CAACAgIAAxkBAAESzGBjaqr-iDc1XPlF0LQVKxeApeGbVwACQhAAAjPFKUmQDtQRpypKgisE`)))
		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
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
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "correct_translation":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    callback.From.ID,
				MessageID: callback.Message.MessageID,
			},
			Text:      user.Localize("<i>Thank you for choosing our translator Translo</i>"),
			ParseMode: tgbotapi.ModeHTML,
		})
		//app.analytics.Bot(tgbotapi.MessageConfig{
		//	BaseChat: tgbotapi.BaseChat{
		//		ChatID:           callback.From.ID,
		//		ReplyToMessageID: callback.Message.MessageID,
		//	},
		//	Text:                  user.Localize("<i>Thank you for choosing our translator Translo</i>"),
		//	ParseMode:             "",
		//	Entities:              nil,
		//	DisableWebPagePreview: false,
		//}, "correct_translation")
	case "wrong_translation": // arr[1] - used 'from', arr[2] - used 'to'
		tryingToFixMsg, err := app.bot.Send(tgbotapi.NewEditMessageText(callback.From.ID, callback.Message.MessageID, user.Localize("Попытаюсь исправить...")))
		if err != nil {
			warn(err)
			return
		}
		//app.analytics.Bot(tgbotapi.MessageConfig{
		//	BaseChat: tgbotapi.BaseChat{
		//		ChatID:           callback.From.ID,
		//		ReplyToMessageID: callback.Message.MessageID,
		//	},
		//	Text:                  user.Localize("Попытаюсь исправить..."),
		//	ParseMode:             "",
		//	Entities:              nil,
		//	DisableWebPagePreview: false,
		//}, "wrong_translation")
		// from != user.MyLang && from != user.ToLang: (to = user.MyLang)
		// переводим с from на user.ToLang
		// Переводим на user.MyLang-user.ToLang
		// Переводим на user.ToLang-user.MyLang

		// from == user.MyLang (to = user.ToLang)
		// Переводим наоборот

		// from == user.ToLang (to = user.MyLang)
		// Переводим наоборот

		from, text := arr[1], callback.Message.ReplyToMessage.Text

		if len(callback.Message.ReplyToMessage.Entities) > 0 { // TODO handle all results from helpers.ApplyEntitiesHtml
			text = helpers.ApplyEntitiesHtml(text, callback.Message.ReplyToMessage.Entities, 4098)[0]
		} else if len(callback.Message.ReplyToMessage.CaptionEntities) > 0 {
			text = helpers.ApplyEntitiesHtml(text, callback.Message.ReplyToMessage.CaptionEntities, 4098)[0]
		}

		lastMsgID := 0
		switch from {
		case user.MyLang:
			tr, _, err := app.translate(ctx, user.ToLang, user.MyLang, text)
			if err != nil {
				warn(err)
				return
			}
			cfg := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: callback.From.ID,
				},
				Text:                  tr,
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			}
			msg, err := app.bot.Send(cfg)
			if err != nil {
				msg, err = app.bot.Send(tgbotapi.NewMessage(callback.From.ID, tr))
				if err != nil {
					warn(err)
					return
				}
			}
			lastMsgID = msg.MessageID
			//app.analytics.Bot(cfg, "wrong_translation")
		case user.ToLang:
			tr, _, err := app.translate(ctx, user.MyLang, user.ToLang, text)
			if err != nil {
				warn(err)
				return
			}
			cfg := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: callback.From.ID,
				},
				Text:                  tr,
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			}
			msg, err := app.bot.Send(cfg)
			if err != nil {
				msg, err = app.bot.Send(tgbotapi.NewMessage(callback.From.ID, tr))
				if err != nil {
					warn(err)
					return
				}
			}
			lastMsgID = msg.MessageID
			//app.analytics.Bot(cfg, "wrong_translation")
		default:
			g, ctx := errgroup.WithContext(ctx)
			for from, to := range map[string]string{from: user.ToLang, user.MyLang: user.ToLang, user.ToLang: user.MyLang} {
				from := from
				to := to
				g.Go(func() error {
					tr, _, err := app.translate(ctx, from, to, text)
					if err != nil || tr == text {
						return err
					}
					cfg := tgbotapi.MessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID: callback.From.ID,
						},
						Text:                  tr,
						ParseMode:             tgbotapi.ModeHTML,
						Entities:              nil,
						DisableWebPagePreview: false,
					}
					msg, err := app.bot.Send(cfg)
					if err != nil {
						msg, err = app.bot.Send(tgbotapi.NewMessage(callback.From.ID, tr))
						if err != nil {
							return err
						}
					}
					lastMsgID = msg.MessageID
					//app.analytics.Bot(cfg, "wrong_translation")
					return nil
				})
			}
			if err = g.Wait(); err != nil {
				warn(err)
				return
			}
		}
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, tryingToFixMsg.MessageID))
		cfg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
				ReplyToMessageID: lastMsgID,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("❌", "wrong_translation_eventually"),
						tgbotapi.NewInlineKeyboardButtonData("✅", "correct_translation"),
					)),
			},
			Text: user.Localize("Did I translate it correctly?"),
		}
		app.bot.Send(cfg)
		//app.analytics.Bot(cfg, "wrong_translation")
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "wrong_translation_eventually":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.bot.Send(tgbotapi.NewEditMessageText(callback.From.ID, callback.Message.MessageID, user.Localize("Excuses")))
		//app.analytics.Bot(tgbotapi.MessageConfig{
		//	BaseChat: tgbotapi.BaseChat{
		//		ChatID: callback.From.ID,
		//	},
		//	Text: user.Localize("Excuses"),
		//}, "wrong_translation")
	case "send_mailing":
		// Getting message_id of the mailing
		msgIDBytes, err := app.bc.Get([]byte(`mailing_message_id`))
		if err != nil {
			warn(err)
			return
		}
		mailingMsgID, err := strconv.Atoi(string(msgIDBytes))
		if err != nil {
			warn(err)
			return
		}

		// Getting keyboard of the mailing
		keyboardTextBytes, err := app.bc.Get([]byte(`mailing_keyboard_text`))
		if err != nil {
			warn(err)
			return
		}
		keyboard, valid := parseKeyboard(string(keyboardTextBytes))
		if !valid {
			warn(fmt.Errorf("invalid keyboard for mailing: %s", string(keyboardTextBytes)))
			return
		}

		users, err := app.db.GetAllUsers()
		if err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewEditMessageText(callback.From.ID, callback.Message.MessageID, "Рассылка начата! 🚀"))
		rateLimiter := ratelimit.New(30)
		for i := 0; i < len(users); i++ {
			_, err := app.bot.Send(tgbotapi.CopyMessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:              users[i].ID, // TODO multiple admins
					ReplyMarkup:         keyboard,
					DisableNotification: false,
				},
				FromChatID: config.AdminID,
				MessageID:  mailingMsgID,
			})
			if err != nil {
				// TODO handle
			}
			rateLimiter.Take()
		}
		app.bot.Send(tgbotapi.NewMessage(config.AdminID, fmt.Sprintf(`"Рассылка завершена. Отправлено %d сообщений"`, len(users))))
		if err = app.bc.Delete([]byte("mailing_keyboard_raw_text")); err != nil {
			warn(err)
			return
		}
		if err = app.bc.Delete([]byte("mailing_message_id")); err != nil {
			warn(err)
			return
		}
	case "cancel_mailing":
		if err = app.bc.Delete([]byte(`mailing_message_id`)); err != nil {
			warn(err)
			return
		}
		if err = app.bc.Delete([]byte(`mailing_keyboard_text`)); err != nil {
			warn(err)
			return
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    callback.From.ID,
				MessageID: callback.Message.MessageID,
			},
			Text:      "<b>Рассылка отменена.</b> Чтобы создать рассылку, введите /mailing",
			ParseMode: tgbotapi.ModeHTML,
		})
	case "none":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "delete": // arr[1] - text for callback query
		text := ""
		if len(arr) > 1 {
			text = strings.Join(arr[1:], ":")
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, text))
		app.bot.Send(tgbotapi.DeleteMessageConfig{
			ChatID:    callback.From.ID,
			MessageID: callback.Message.MessageID,
		})
		return
	case "set_my_lang": // arr[1] - lang, arr[2] - keyboard offset
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		if err = app.db.UpdateUser(callback.From.ID, tables.Users{MyLang: arr[1]}); err != nil {
			warn(err)
			return
		}
		user.MyLang = arr[1]
		offset, err := strconv.Atoi(arr[2])
		if err != nil {
			warn(err)
			return
		}
		if offset < 0 || offset > len(codes[*user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[user.Lang]) is %d, offset ois %d", len(codes[*user.Lang]), offset))
			return
		}
		count := 18
		if offset+count > len(codes[*user.Lang])-1 {
			count = len(codes[*user.Lang]) - 1 - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[*user.Lang]) / count * count // from end
		}

		next := offset + 18
		if next >= len(codes[*user.Lang])-1 {
			next = 0 // from start
		}
		kb, err := buildLangsPagination(user, offset, count, user.MyLang,
			fmt.Sprintf("set_my_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_my_lang_pagination:%d", back),
			fmt.Sprintf("set_my_lang_pagination:%d", next), true)
		if err != nil {
			log.Error("", zap.Error(err))
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.EditMessageReplyMarkupConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &kb,
			},
		})
		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
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
			Text: user.Localize("Клавиатура обновлена"),
		}); err != nil {
			warn(err)
			return
		}
	case "set_to_lang": // arr[1] - lang, arr[2] - keyboard offset
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Translo"))
		if err = app.db.UpdateUser(callback.From.ID, tables.Users{ToLang: arr[1]}); err != nil {
			warn(err)
			return
		}
		user.ToLang = arr[1]
		offset, err := strconv.Atoi(arr[2])
		if err != nil {
			warn(err)
			return
		}
		if offset < 0 || offset > len(codes[*user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[*user.Lang]) is %d, offset ois %d", len(codes[*user.Lang]), offset))
			return
		}
		count := 18
		if offset+count > len(codes[*user.Lang])-1 {
			count = len(codes[*user.Lang]) - 1 - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[*user.Lang]) / count * count // from end
		}

		next := offset + 18
		if next >= len(codes[*user.Lang])-1 {
			next = 0 // from start
		}
		kb, err := buildLangsPagination(user, offset, count, user.ToLang,
			fmt.Sprintf("set_to_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_to_lang_pagination:%d", back),
			fmt.Sprintf("set_to_lang_pagination:%d", next), false)
		if err != nil {
			log.Error("", zap.Error(err))
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.EditMessageReplyMarkupConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &kb,
			},
		})
		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
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
			Text: user.Localize("Клавиатура обновлена"),
		}); err != nil {
			warn(err)
			return
		}
	case "set_my_lang_pagination": // arr[1] - offset to show
		offset, err := strconv.Atoi(arr[1])
		if err != nil {
			warn(err)
			return
		}
		if offset < 0 || offset > len(codes[*user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[*user.Lang]) is %d, offset ois %d", len(codes[*user.Lang]), offset))
			return
		}

		count := 18
		if offset+count > len(codes[*user.Lang])-1 {
			count = len(codes[*user.Lang]) - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[*user.Lang]) / count * count // end
		}

		next := offset + 18
		if next >= len(codes[*user.Lang])-1 {
			next = 0 // start
		}

		kb, err := buildLangsPagination(user, offset, count, "",
			fmt.Sprintf("set_my_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_my_lang_pagination:%d", back),
			fmt.Sprintf("set_my_lang_pagination:%d", next), true)
		if err != nil {
			log.Error("", zap.Error(err))
			warn(err)
			return
		}

		if reflect.DeepEqual(*callback.Message.ReplyMarkup, kb) {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
			return
		}

		if _, err = app.bot.Send(tgbotapi.EditMessageReplyMarkupConfig{tgbotapi.BaseEdit{
			ChatID:          callback.From.ID,
			ChannelUsername: "",
			MessageID:       callback.Message.MessageID,
			InlineMessageID: "",
			ReplyMarkup:     &kb,
		}}); err != nil {
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   callback.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              &kb,
					DisableNotification:      true,
					AllowSendingWithoutReply: false,
				},
				Text:                  callback.Message.Text,
				ParseMode:             "",
				Entities:              nil,
				DisableWebPagePreview: false,
			})
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "set_to_lang_pagination": // arr[2] - offset to show
		offset, err := strconv.Atoi(arr[1])
		if err != nil {
			warn(err)
			return
		}
		if offset < 0 || offset > len(codes[*user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[*user.Lang]) is %d, offset ois %d", len(codes[*user.Lang]), offset))
			return
		}

		count := 18
		if offset+count > len(codes[*user.Lang])-1 {
			count = len(codes[*user.Lang]) - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[*user.Lang]) / count * count // end
		}

		next := offset + 18
		if next >= len(codes[*user.Lang])-1 {
			next = 0 // start
		}

		kb, err := buildLangsPagination(user, offset, count, "",
			fmt.Sprintf("set_to_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_to_lang_pagination:%d", back),
			fmt.Sprintf("set_to_lang_pagination:%d", next), false)
		if err != nil {
			log.Error("", zap.Error(err))
			warn(err)
			return
		}

		if reflect.DeepEqual(*callback.Message.ReplyMarkup, kb) {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
			return
		}

		if _, err = app.bot.Send(tgbotapi.EditMessageReplyMarkupConfig{tgbotapi.BaseEdit{
			ChatID:          callback.From.ID,
			ChannelUsername: "",
			MessageID:       callback.Message.MessageID,
			InlineMessageID: "",
			ReplyMarkup:     &kb,
		}}); err != nil {
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   callback.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              &kb,
					DisableNotification:      true,
					AllowSendingWithoutReply: false,
				},
				Text:                  callback.Message.Text,
				ParseMode:             "",
				Entities:              nil,
				DisableWebPagePreview: false,
			})
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "type_my_lang_name":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		for i1, row := range callback.Message.ReplyMarkup.InlineKeyboard {
			for i2, btn := range row {
				if btn.CallbackData == nil {
					continue
				}
				if *btn.CallbackData == callback.Data && strings.HasPrefix(btn.Text, "✅") { // наша кнопка
					btn.Text = strings.TrimSpace(strings.TrimPrefix(btn.Text, "✅"))
					callback.Message.ReplyMarkup.InlineKeyboard[i1][i2] = btn
					app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, *callback.Message.ReplyMarkup))

					data, err := app.bc.Get([]byte(strconv.FormatInt(callback.From.ID, 10)))
					if err != nil {
						warn(err)
						return
					}
					chunks := strings.Split(string(data), ";")
					if len(chunks) != 2 {
						warn(err)
						log.Error("len(chunks) is not 2", zap.String("result", string(data)))
						return
					}
					msgID, err := strconv.ParseInt(chunks[0], 10, 64)
					if err != nil {
						warn(err)
						return
					}
					app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, int(msgID)))

					if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": nil}); err != nil {
						warn(err)
						return
					}
					return
				}
			}
		}
		setMyLang := "set_my_lang"
		if err = app.db.UpdateUser(callback.From.ID, tables.Users{Act: &setMyLang}); err != nil {
			app.notifyAdmin(err)
		}

		keyboard := callback.Message.ReplyMarkup
		UntickAll(keyboard.InlineKeyboard)
		Tick(callback.Data, keyboard.InlineKeyboard)
		app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, *keyboard))
		msg, err := app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
				ReplyToMessageID: callback.Message.MessageID,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(user.Localize(`❌ Cancel`), "close_type_language_name_menu:my_lang"))),
				DisableNotification: true,
			},
			Text: user.Localize(`🔎 Write the name of the language you want to search for:`),
		})
		if err != nil {
			warn(err)
			return
		}
		key := []byte(strconv.FormatInt(callback.From.ID, 10))
		value := append([]byte(strconv.Itoa(msg.MessageID)), []byte(";")...)
		value = append(value, []byte(strconv.Itoa(callback.Message.MessageID))...) // value = "msg.MessageID;callback.Message.MessageID"
		if err = app.bc.Put(key, value); err != nil {
			warn(err)
		}
		log.Debug("set act and put msg_id in bitcask", zap.String("act", setMyLang), zap.Int("msg_id", msg.MessageID))
	case "type_to_lang_name":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		for i1, row := range callback.Message.ReplyMarkup.InlineKeyboard {
			for i2, btn := range row {
				if btn.CallbackData == nil {
					continue
				}
				if *btn.CallbackData == callback.Data && strings.HasPrefix(btn.Text, "✅") { // наша кнопка
					btn.Text = strings.TrimSpace(strings.TrimPrefix(btn.Text, "✅"))
					callback.Message.ReplyMarkup.InlineKeyboard[i1][i2] = btn
					app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, *callback.Message.ReplyMarkup))

					data, err := app.bc.Get([]byte(strconv.FormatInt(callback.From.ID, 10)))
					if err != nil {
						warn(err)
						return
					}
					chunks := strings.Split(string(data), ";")
					if len(chunks) != 2 {
						warn(err)
						log.Error("len(chunks) is not 2", zap.String("result", string(data)))
						return
					}
					msgID, err := strconv.ParseInt(chunks[0], 10, 64)
					if err != nil {
						warn(err)
						return
					}
					app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, int(msgID)))

					if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": nil}); err != nil {
						warn(err)
						return
					}
					return
				}
			}
		}
		setToLang := "set_to_lang"
		if err = app.db.UpdateUser(callback.From.ID, tables.Users{Act: &setToLang}); err != nil {
			app.notifyAdmin(err)
		}

		keyboard := callback.Message.ReplyMarkup
		UntickAll(keyboard.InlineKeyboard)
		Tick(callback.Data, keyboard.InlineKeyboard)
		app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, *keyboard))
		msg, err := app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
				ReplyToMessageID: callback.Message.MessageID,
				ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(user.Localize(`❌ Cancel`), "close_type_language_name_menu:to_lang"))),
				DisableNotification: true,
			},
			Text: user.Localize(`🔎 Write the name of the language you want to search for:`),
		})
		if err != nil {
			warn(err)
			return
		}
		key := []byte(strconv.FormatInt(callback.From.ID, 10))
		value := append([]byte(strconv.Itoa(msg.MessageID)), []byte(";")...)
		value = append(value, []byte(strconv.Itoa(callback.Message.MessageID))...) // value = "msg.MessageID;callback.Message.MessageID"
		if err = app.bc.Put(key, value); err != nil {
			warn(err)
		}
		log.Debug("set act and put msg_id in bitcask", zap.String("act", setToLang), zap.Int("msg_id", msg.MessageID))
	case "try_again_to_search_my_lang":
		app.bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, user.Localize(`🔎 Write the name of the language you want to search for:`), tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(user.Localize(`❌ Cancel`), "close_type_language_name_menu:my_lang")))))
		if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": "set_my_lang"}); err != nil {
			warn(err)
			return
		}
	case "try_again_to_search_to_lang":
		app.bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, user.Localize(`🔎 Write the name of the language you want to search for:`), tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(user.Localize(`❌ Cancel`), "close_type_language_name_menu:to_lang")))))
		if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": "set_to_lang"}); err != nil {
			warn(err)
			return
		}
	case "close_type_language_name_menu": // arr[1] -  my_lang/to_lang
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		//if callback.
		keyboard := callback.Message.ReplyToMessage.ReplyMarkup
		UntickAll(keyboard.InlineKeyboard)
		app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.ReplyToMessage.MessageID, *keyboard))
		if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
		}
	case "cancel_set_op":
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": nil}); err != nil {
			warn(err)
			return
		}
	case "filtered_set_my_lang": // arr[1] - set my_lang
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))

		// Building new keyboard for pagination message
		codePosition := -1
		for _, code := range codes[*user.Lang] {
			codePosition++
			if code == arr[1] {
				break
			}
		}
		offset := codePosition / 18 * 18
		count := 18
		if offset+count > len(codes[*user.Lang])-1 {
			count = len(codes[*user.Lang]) - offset
		}
		back := offset - 18
		if back < 0 {
			back = len(codes[*user.Lang]) / count * count // end
		}
		next := offset + 18
		if next >= len(codes[*user.Lang])-1 {
			next = 0 // start
		}
		keyboard, err := buildLangsPagination(user, codePosition/18*18, 18, arr[1],
			fmt.Sprintf("set_my_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_my_lang_pagination:%d", back),
			fmt.Sprintf("set_my_lang_pagination:%d", next), true)
		if err != nil {
			warn(err)
			return
		}

		// Getting pagination message
		msgIDsBytes, err := app.bc.Get([]byte(strconv.FormatInt(callback.From.ID, 10)))
		if err != nil {
			warn(err)
			return
		}
		msgIDs := strings.Split(string(msgIDsBytes), ";") // msgIDs[0] - search query message. msgIDs[1] - languages pagination message.
		if len(msgIDs) != 2 {
			warn(fmt.Errorf("strings.Split(app.bc.Get(message.From.ID), \";\") is not 2 chunks"))
			return
		}
		paginationMsgID, err := strconv.ParseInt(msgIDs[1], 10, 64)
		if err != nil {
			warn(err)
			log.Error("couldn't parse int64: app.bc.Get(message.From.ID)", zap.Error(err), zap.String("result", string(msgIDsBytes)))
			return
		}
		// Updating pagination message
		app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, int(paginationMsgID), keyboard))

		// Updating user's my_lang
		user.MyLang = arr[1]
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
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
			Text: user.Localize("Клавиатура обновлена"),
		}); err != nil {
			warn(err)
			return
		}
		if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"my_lang": arr[1], "act": nil}); err != nil {
			warn(err)
			return
		}
	case "filtered_set_to_lang":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))

		// Building new keyboard for pagination message
		codePosition := -1
		for _, code := range codes[*user.Lang] {
			codePosition++
			if code == arr[1] {
				break
			}
		}
		offset := codePosition / 18 * 18
		count := 18
		if offset+count > len(codes[*user.Lang])-1 {
			count = len(codes[*user.Lang]) - offset
		}
		back := offset - 18
		if back < 0 {
			back = len(codes[*user.Lang]) / count * count // end
		}
		next := offset + 18
		if next >= len(codes[*user.Lang])-1 {
			next = 0 // start
		}
		keyboard, err := buildLangsPagination(user, codePosition/18*18, 18, arr[1],
			fmt.Sprintf("set_to_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_to_lang_pagination:%d", back),
			fmt.Sprintf("set_to_lang_pagination:%d", next), true)
		if err != nil {
			warn(err)
			return
		}

		// Getting pagination message
		msgIDsBytes, err := app.bc.Get([]byte(strconv.FormatInt(callback.From.ID, 10)))
		if err != nil {
			warn(err)
			return
		}
		msgIDs := strings.Split(string(msgIDsBytes), ";") // msgIDs[0] - search query message. msgIDs[1] - languages pagination message.
		if len(msgIDs) != 2 {
			warn(fmt.Errorf("strings.Split(app.bc.Get(message.From.ID), \";\") is not 2 chunks"))
			return
		}
		paginationMsgID, err := strconv.ParseInt(msgIDs[1], 10, 64)
		if err != nil {
			warn(err)
			log.Error("couldn't parse int64: app.bc.Get(message.From.ID)", zap.Error(err), zap.String("result", string(msgIDsBytes)))
			return
		}
		// Updating pagination message
		app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, int(paginationMsgID), keyboard))

		// Updating user's my_lang
		user.ToLang = arr[1]
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           callback.From.ID,
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
			Text: user.Localize("Клавиатура обновлена"),
		}); err != nil {
			warn(err)
			return
		}
		if err = app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"to_lang": arr[1], "act": nil}); err != nil {
			warn(err)
			return
		}
	}
}
