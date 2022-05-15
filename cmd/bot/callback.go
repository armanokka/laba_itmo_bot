package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func (app *App) onCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	log := app.log.With(zap.Int64("id", callback.From.ID))
	defer func() {
		if err := recover(); err != nil {
			log.Error("%w", zap.Any("error", err), zap.String("stack_trace", string(debug.Stack())))
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
	case "choose_another_lang": // arr[1] - lang
		app.bot.Send(tgbotapi.NewCallbackWithAlert(callback.ID, user.Localize("Твое сообщение и так на %s языке. Выбери другой язык, на который надо переводить", langs[user.Lang][arr[1]])))
	case "setup_langs": // arr[1] - source language, arr[2] - direction to translate, in replied message there is source text
		if err := app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{
			"my_lang":       arr[1],
			"to_lang":       arr[2],
			"last_activity": time.Now(),
			"act":           "",
		}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		if err = app.SuperTranslate(ctx, user, callback.From.ID, arr[1], arr[2], callback.Message.ReplyToMessage.Text, callback.Message.ReplyToMessage.MessageID, *callback.Message.ReplyToMessage); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))

		//if _, err := app.bot.Send(tgbotapi.DocumentConfig{
		//	BaseFile: tgbotapi.BaseFile{
		//		BaseChat: tgbotapi.BaseChat{
		//			ChatID:                   callback.From.ID,
		//			ChannelUsername:          "",
		//			ReplyToMessageID:         0,
		//			ReplyMarkup:              nil,
		//			DisableNotification:      false,
		//			AllowSendingWithoutReply: false,
		//		},
		//		File: tgbotapi.FilePath("inline.gif"),
		//	},
		//	Thumb:                       nil,
		//	Caption:                     user.Localize("Как переводить еще удобнее"),
		//	ParseMode:                   "",
		//	CaptionEntities:             nil,
		//	DisableContentTypeDetection: false,
		//}); err != nil {
		//	pp.Println(err)
		//}
	case "setup_langs_pagination": // arr[1] - source language of the text, arr[2] - offset
		from := arr[1]
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

		if count == 0 {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
			return
		}

		back := offset - 18
		if back < 0 {
			back = 0
		}
		kb, err := buildLangsPagination(user, offset, count, arr[1],
			fmt.Sprintf("setup_langs:%s:%s", from, "%s"),
			fmt.Sprintf("setup_langs_pagination:%s:%d", from, back),
			fmt.Sprintf("setup_langs_pagination:%s:%d", from, offset+count),
			fmt.Sprintf("choose_another_lang:%s", arr[1]))
		if err != nil {
			warn(err)
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
