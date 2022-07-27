package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
)

func (app *App) onCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	log := app.log.With(zap.Int64("id", callback.From.ID))
	defer func() {
		if err := recover(); err != nil {
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
	}

	var err error
	user, err := app.db.GetUserByID(callback.From.ID)
	if err != nil {
		warn(err)
	}
	user.SetLang(callback.From.LanguageCode)

	arr := strings.Split(callback.Data, ":")

	defer func() {
		if err = app.db.UpdateUserActivity(callback.From.ID); err != nil {
			app.notifyAdmin(err)
		}
	}()

	switch arr[0] {
	case "cancel_mailing_act":
		if err := app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": ""}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.ReplyToMessage.MessageID))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "OK"))
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
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Translo"))
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
		if offset < 0 || offset > len(codes[user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[user.Lang]) is %d, offset ois %d", len(codes[user.Lang]), offset))
			return
		}
		count := 18
		if offset+count > len(codes[user.Lang])-1 {
			count = len(codes[user.Lang]) - 1 - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[user.Lang]) / count * count // from end
		}

		next := offset + 18
		if next >= len(codes[user.Lang])-1 {
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
						tgbotapi.NewKeyboardButton(langs[callback.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[callback.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji))),
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
		if offset < 0 || offset > len(codes[user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[user.Lang]) is %d, offset ois %d", len(codes[user.Lang]), offset))
			return
		}
		count := 18
		if offset+count > len(codes[user.Lang])-1 {
			count = len(codes[user.Lang]) - 1 - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[user.Lang]) / count * count // from end
		}

		next := offset + 18
		if next >= len(codes[user.Lang])-1 {
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
						tgbotapi.NewKeyboardButton(langs[callback.From.LanguageCode][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔"),
						tgbotapi.NewKeyboardButton(langs[callback.From.LanguageCode][user.ToLang]+" "+flags[user.ToLang].Emoji))),
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
		if offset < 0 || offset > len(codes[user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[user.Lang]) is %d, offset ois %d", len(codes[user.Lang]), offset))
			return
		}

		count := 18
		if offset+count > len(codes[user.Lang])-1 {
			count = len(codes[user.Lang]) - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[user.Lang]) / count * count // end
		}

		next := offset + 18
		if next >= len(codes[user.Lang])-1 {
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
		if offset < 0 || offset > len(codes[user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[user.Lang]) is %d, offset ois %d", len(codes[user.Lang]), offset))
			return
		}

		count := 18
		if offset+count > len(codes[user.Lang])-1 {
			count = len(codes[user.Lang]) - offset
		}

		back := offset - 18
		if back < 0 {
			back = len(codes[user.Lang]) / count * count // end
		}

		next := offset + 18
		if next >= len(codes[user.Lang])-1 {
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
	case "show_from": //arr[1] - lang
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, user.Localize("Переведено с %s", langs[callback.From.LanguageCode][arr[1]]+" "+flags[arr[1]].Emoji)))
	case "start_mailing":
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

		keyboardBytes, err := app.bc.Get([]byte("mailing_keyboard_raw_text"))
		if err != nil {
			warn(err)
			return
		}
		keyboardText := string(keyboardBytes)
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		if keyboardText != "Empty" {
			keyboard = parseKeyboard(keyboardText)
		}

		if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
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
		usersNumber, err := app.db.GetUsersNumber()
		if err != nil {
			return
		}

		slice := make([]int64, 0, 100)
		for n := usersNumber / 100; n < usersNumber/100; n++ { // iterate over each 100 users
			offset := n*100 + 100
			if err = app.db.GetUsersSlice(offset, 100, slice); err != nil {
				warn(err)
				return
			}
			pp.Println(slice)
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
				log.Info("mailing was sent", zap.Int64("recepient_id", slice[j]))
			}
		}

		app.bot.Send(tgbotapi.NewMessage(callback.From.ID, "рассылка закончена"))
		if err = app.bc.Delete([]byte("mailing_keyboard_raw_text")); err != nil {
			warn(err)
			return
		}
		if err = app.bc.Delete([]byte("mailing_message_id")); err != nil {
			warn(err)
			return
		}
	}
}
