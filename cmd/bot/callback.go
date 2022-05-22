package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		if err := app.db.UpdateUserMetrics(callback.From.ID, "callback:"+callback.Data); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	switch arr[0] {
	case "delete_cache": // arr[1] - document's ._key
		if _, err = app.cache.RemoveDocument(ctx, arr[1]); err != nil {
			_, err = app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, err.Error()))
		} else {
			_, err = app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "OK"))
		}
		if err != nil {
			warn(err)
			log.Error("", zap.Error(err))
			return
		}
		log.Info("cache was deleted", zap.String("_key", arr[1]))
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
		app.db.LogBotMessage(callback.From.ID, "cb_none", "")
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
		app.db.LogBotMessage(callback.From.ID, "cb_delete", "")
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
		back := offset - 18
		if back < 0 {
			back = 0
		}
		kb, err := buildLangsPagination(user, offset, 18, user.MyLang,
			fmt.Sprintf("set_my_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_my_lang_pagination:%d", back),
			fmt.Sprintf("set_my_lang_pagination:%d", offset+18))
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
						tgbotapi.NewKeyboardButton("↔️"),
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
		back := offset - 18
		if back < 0 {
			back = 0
		}
		kb, err := buildLangsPagination(user, offset, 18, user.ToLang,
			fmt.Sprintf("set_to_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_to_lang_pagination:%d", back),
			fmt.Sprintf("set_to_lang_pagination:%d", offset+18))
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
						tgbotapi.NewKeyboardButton("↔️"),
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
			count = len(codes[user.Lang]) - 1 - offset
		}

		back := offset - 18
		if back < 0 {
			back = 180
		}

		next := offset + 18
		if next >= len(codes[user.Lang])-1 {
			next = 0 // from start
		}

		kb, err := buildLangsPagination(user, offset, count, "",
			fmt.Sprintf("set_my_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_my_lang_pagination:%d", back),
			fmt.Sprintf("set_my_lang_pagination:%d", next))
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
			count = len(codes[user.Lang]) - 1 - offset
		}

		back := offset - 18
		if back < 0 {
			back = 180 // end
		}

		next := offset + 18
		if next >= len(codes[user.Lang])-1 {
			next = 0 // start
		}

		kb, err := buildLangsPagination(user, offset, count, "",
			fmt.Sprintf("set_my_lang:%s:%d", "%s", offset),
			fmt.Sprintf("set_my_lang_pagination:%d", back),
			fmt.Sprintf("set_my_lang_pagination:%d", next))
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

	}
}
