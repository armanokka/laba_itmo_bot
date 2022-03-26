package bot

import (
	"context"
	"fmt"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"golang.org/x/sync/errgroup"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (app *App) onCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewCallback(callback.ID, "Error, sorry"))
		app.notifyAdmin(err)
		pp.Println(err)
	}

	//user := tables.Users{Lang: callback.From.LanguageCode}

	arr := strings.Split(callback.Data, ":")

	defer func() {
		if err := app.db.UpdateUserLastActivity(callback.From.ID); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	switch arr[0] {
	case "none":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.db.LogBotMessage(callback.From.ID, "cb_none", "")
		return
	case "delete":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.bot.Send(tgbotapi.DeleteMessageConfig{
			ChatID:    callback.From.ID,
			MessageID: callback.Message.MessageID,
		})
		app.db.LogBotMessage(callback.From.ID, "cb_delete", "")
		return
	case "speech_this_message_and_replied_one": // arr[1] - from, arr[2] - to
		text := callback.Message.Text
		if callback.Message.Caption != "" {
			text = callback.Message.Caption
		}
		user, err := app.db.GetUserByID(callback.From.ID)
		if err != nil {
			warn(err)
			return
		}
		go func() {
			if err := app.sendSpeech(user, arr[1], callback.Message.ReplyToMessage.Text, callback.ID); err != nil { // озвучиваем непереведенное сообщение
				warn(err)
				return
			}
		}()
		if err := app.sendSpeech(user, arr[2], text, callback.ID); err != nil { // озвучиваем переведенное сообщение
			warn(err)
			return
		}
		app.db.LogUserMessage(callback.From.ID, "cb_voice")
	case "dict": // arr[1] - from, arr[2] - to, text in replied message
		go app.bot.Send(tgbotapi.NewChatAction(callback.From.ID, "typing"))

		from := arr[1]

		meaning, err := translate2.GoogleDictionary(from, strings.ToLower(callback.Message.ReplyToMessage.Text))
		if err != nil {
			warn(err)
			return
		}

		text := ""

		for _, data := range meaning.DictionaryData {
			for _, entry := range data.Entries {
				for _, senseFamily := range entry.SenseFamilies {
					for _, sense := range senseFamily.Senses {
						text += "\n- " + sense.Definition.Text
					}
				}
			}
		}

		if len(strings.Fields(text)) == 0 {
			app.bot.Send(tgbotapi.CallbackConfig{
				CallbackQueryID: callback.ID,
				Text:            "Empty result :(",
				ShowAlert:       true,
				URL:             "",
				CacheTime:       0,
			})
			return
		}

		if callback.From.LanguageCode != "en" {
			tr, err := translate2.MicrosoftTranslate("en", callback.From.LanguageCode, text)
			if err != nil {
				warn(err)
			}
			text = tr.TranslatedText
		}

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         callback.Message.MessageID,
				ReplyMarkup:              nil,
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			Text:                  text,
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		})
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		app.db.LogBotMessage(callback.From.ID, "cb_meaning", text)
	case "exm": // arr[1] - from, arr[2] - to. Target text in replied message
		tr, err := translate2.ReversoTranslate(translate2.ReversoIso6392(arr[1]), translate2.ReversoIso6392(arr[2]), callback.Message.ReplyToMessage.Text)
		if err != nil {
			warn(err)
		}

		var (
			text string
		)

		if len(tr.ContextResults.Results) > 3 {
			tr.ContextResults.Results = tr.ContextResults.Results[:3]
		}

		idx := 0
		count := 0

		for _, result := range tr.ContextResults.Results {
			idx++
			if len(result.SourceExamples) == 0 {
				continue
			}
			for i := 0; i < len(result.SourceExamples) && count < 3; i++ {
				if i > 0 {
					idx++
				}
				text += "\n\n<b>" + strconv.Itoa(idx) + ".</b> "
				text += result.SourceExamples[i] + "\n<b>└</b>" + result.TargetExamples[i]
				count++
			}
		}
		text = strings.NewReplacer("<em>", "<b>", "</em>", "</b>").Replace(text)

		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         callback.Message.MessageID,
				ReplyMarkup:              nil,
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:                  text,
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			pp.Println(err)
		}

		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		app.db.LogBotMessage(callback.From.ID, "cb_exmp", text)
	case "trs": // arr[1], arr[2] = from, to (in iso6391)
		callback.Message.ReplyToMessage.Text = strings.ToLower(callback.Message.ReplyToMessage.Text)

		var (
			trscript translate2.YandexTranscriptionResponse
			err      error
			rev      translate2.ReversoTranslation
		)
		reversoTo := translate2.ReversoIso6392(arr[1])
		reversoFrom := translate2.ReversoIso6392(arr[2])

		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		keyboard.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, 0, 2)

		g, _ := errgroup.WithContext(ctx)
		g.Go(func() error {
			rev, err = translate2.ReversoTranslate(translate2.ReversoIso6392(arr[1]), translate2.ReversoIso6392(arr[2]), callback.Message.ReplyToMessage.Text)
			return err
		})
		//g.Go(func() error {
		//	if v, err := lingvo.Suggestions(arr[1], arr[2], callback.Message.Text, 6, 0); err == nil && len(v.Items) > 0 {
		//		suggestions = v.Items
		//	}
		//	return nil
		//})
		g.Go(func() error {
			trscript, err = translate2.YandexTranscription(reversoFrom, reversoTo, callback.Message.ReplyToMessage.Text)
			return err
		})

		if err = g.Wait(); err != nil {
			warn(err)
			return
		}

		var text string
		var wg sync.WaitGroup

		for _, result := range rev.ContextResults.Results {
			if result.Translation == "" {
				continue
			}
			result := result
			wg.Add(1)
			go func() {
				defer wg.Done()
				tr, err := translate2.ReversoTranslate(reversoTo, reversoFrom, result.Translation)
				if err != nil {
					warn(err)
					return
				}

				text += "\n<b>" + result.Translation + "\n└</b>"
				if len(tr.ContextResults.Results) > 4 {
					tr.ContextResults.Results = tr.ContextResults.Results[:4]
				}
				for i, result := range tr.ContextResults.Results {
					if result.Translation == "" {
						continue
					}
					if i > 0 {
						text += ", "
					}
					if result.Translation == callback.Message.ReplyToMessage.Text {
						text += "<b>" + result.Translation + "</b>"
						continue
					}
					text += "<code>" + result.Translation + "</code>"
				}
			}()
		}
		wg.Wait()

		if trscript.StatusCode == 200 {
			addition := ""
			if trscript.Transcription != "" {
				addition += " <b>/" + trscript.Transcription + "/</b>"
			}
			if trscript.Pos != "" {
				addition += " <i>" + trscript.Pos + "</i>"
			}
			text = callback.Message.ReplyToMessage.Text + addition + text
		} else {
			text = callback.Message.ReplyToMessage.Text + "\n" + text
		}

		if text == "" {
			warn(fmt.Errorf("callback: trs: Empty text that need to send. From %s to %s. Query: %s", arr[1], arr[2], callback.Message.ReplyToMessage.Text))
			return
		}

		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         callback.Message.MessageID,
				ReplyMarkup:              keyboard,
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:                  text,
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		app.db.LogBotMessage(callback.From.ID, "cb_dict", text)
	//case "l": // lingvo, arr[1] - from, arr[2] - to, arr[3] - text
	//	text := strings.Join(arr[3:], ":")
	//	res, err := lingvo.GetDictionary(arr[1], arr[2], text)
	//	if err != nil {
	//		warn(err)
	//		return
	//	}

	case "setup_langs": // arr[1] - source language, arr[2] - direction to translate, in replied message there is source text
		app.bot.Send(tgbotapi.NewChatAction(callback.From.ID, "typing"))

		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))

		from := arr[1]
		to := arr[2]
		if err := app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{
			"my_lang":       to,
			"to_lang":       from,
			"last_activity": time.Now(),
			"act":           "",
		}); err != nil {
			warn(err)
			return
		}

		user, err := app.db.GetUserByID(callback.From.ID)
		if err != nil {
			warn(err)
			return
		}

		answer, err := app.SuperTranslate(user, from, to, callback.Message.ReplyToMessage.Text, callback.Message.ReplyToMessage.Entities)
		if err != nil {
			warn(err)
			return
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔊 "+user.Localize("Озвучить"), fmt.Sprintf("speech_this_message_and_replied_one:%s:%s", from, to))))
		if answer.Examples {
			keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("💬 "+user.Localize("Примеры"), fmt.Sprintf("exm:%s:%s", from, to)))
		}
		if answer.Translations {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📚 "+user.Localize("Переводы"), fmt.Sprintf("trs:%s:%s", from, to))))
		}
		if answer.Dictionary {
			l := len(keyboard.InlineKeyboard) - 1
			if l < 0 {
				l = 0
			}
			keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData("ℹ️"+user.Localize("Словарь"), fmt.Sprintf("dict:%s", from)))
		}

		//if answer.Suggestions {
		//	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		//		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Lingvo", fmt.Sprintf("lingvo_vars:%s:%s:%d", from, to, 0))))
		//}

		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         callback.Message.ReplyToMessage.MessageID,
				ReplyMarkup:              keyboard,
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:                  answer.TranslatedText,
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			pp.Println(err)
		}
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		if err = app.db.IncreaseUserUsings(callback.From.ID); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
		app.bot.Send(tgbotapi.NewMessage(callback.From.ID, user.Localize("Теперь я буду переводить с %s на %s и обратно. Если захочешь изменить, напишешь /start", langs[arr[1]].Name, langs[arr[2]].Name)))
	case "setup_langs_pagination": // arr[1] - source language of the text, arr[2] - offset
		from := arr[1]
		offset, err := strconv.Atoi(arr[2])
		if err != nil {
			warn(err)
			return
		}
		if offset < 0 || offset > len(codes)-1 {
			warn(fmt.Errorf("offset is too big, len(codes) is %d, offset ois %d", len(codes), offset))
			return
		}

		count := 18
		if offset+count > len(codes)-1 {
			count = len(codes) - 1 - offset
		}

		if count == 0 {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
			return
		}

		back := offset - 18
		if back < 0 {
			back = 0
		}
		kb, err := buildLangsPagination(offset, count, arr[1],
			fmt.Sprintf("setup_langs:%s:%s", from, "%s"),
			fmt.Sprintf("setup_langs_pagination:%s:%d", from, back),
			fmt.Sprintf("setup_langs_pagination:%s:%d", from, offset+count))
		if err != nil {
			warn(err)
		}

		if reflect.DeepEqual(*callback.Message.ReplyMarkup, kb) {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
			return
		}

		_, err = app.bot.Send(tgbotapi.EditMessageReplyMarkupConfig{tgbotapi.BaseEdit{
			ChatID:          callback.From.ID,
			ChannelUsername: "",
			MessageID:       callback.Message.MessageID,
			InlineMessageID: "",
			ReplyMarkup:     &kb,
		}})
		if err != nil {
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
		//case "inline_translate": //arr[1] - lang
		//	pp.Println(callback)
		//	tr, err := translate2.MicrosoftTranslate("", arr[1], callback.Message.Text)
		//	if err != nil {
		//		warn(err)
		//		return
		//	}
		//	app.bot.Send(tgbotapi.NewEditMessageText(callback.From.ID, callback.Message.MessageID, tr.TranslatedText))
		//	default:
		//	app.bot.Send(tgbotapi.NewMessage(callback.From.ID, "Action is expired. Send /start"))
	}
}
