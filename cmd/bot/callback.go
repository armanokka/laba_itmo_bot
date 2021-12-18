package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/pkg/lingvo"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	"github.com/go-errors/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func (app *app) onCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewCallback(callback.ID, "Error, sorry"))
		app.notifyAdmin(err)
		logrus.Error(err)
	}

	arr := strings.Split(callback.Data, ":")

	user := app.loadUser(callback.From.ID, warn)
	defer user.UpdateLastActivity()

	if arr[0] == "set_bot_lang_and_register" {
		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              nil,
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			Text:                  user.Localize("Бот был обновлен. Нажмите /start"),
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		})
	} else {
		user.Fill()
	}


	switch arr[0] {
	case "none":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.writeBotLog(callback.From.ID, "cb_none", "")
		return
	case "delete":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.bot.Send(tgbotapi.DeleteMessageConfig{
			ChatID:          callback.From.ID,
			MessageID:       callback.Message.MessageID,
		})
		app.writeBotLog(callback.From.ID, "cb_delete", "")
		return
	case "speech_this_message_and_replied_one": // arr[1] - from, arr[2] - to
		text := callback.Message.Text
		if callback.Message.Caption != "" {
			text = callback.Message.Caption
		}
		if err := app.sendSpeech(arr[1], callback.Message.ReplyToMessage.Text, callback.ID, user); err != nil { // озвучиваем непереведенное сообщение
			warn(err)
			return
		}
		if err := app.sendSpeech(arr[2], text, callback.ID, user); err != nil { // озвучиваем переведенное сообщение
			warn(err)
			return
		}
		app.writeUserLog(callback.From.ID, "cb_voice")
	case "dict": // arr[1] - from, arr[2] - to, text in replied message
		go app.bot.Send(tgbotapi.NewChatAction(callback.From.ID, "typing"))

		from := arr[1]
		to := arr[2]

		var (
			meaning translate2.GoogleDictionaryResponse
			lingv []lingvo.Dictionary
			err error
		)

		g, _ := errgroup.WithContext(ctx)

		g.Go(func() error {
			meaning, err = translate2.GoogleDictionary(from, strings.ToLower(callback.Message.ReplyToMessage.Text))
			return err
		})

		g.Go(func() error {
			_, ok1 := lingvo.Lingvo[from]
			_, ok2 := lingvo.Lingvo[to]
			if !ok1 || !ok2 {
				return nil
			}
			lingv, err = lingvo.GetDictionary(from, to, callback.Message.ReplyToMessage.Text)
			return err
		})

		if err := g.Wait(); err != nil {
			warn(err)
			return
		}

		pp.Println(meaning, lingv)

		text := callback.Message.ReplyToMessage.Text
		var transcription bool
		if lingv != nil {
			for i, dict := range lingv {
				if !transcription && dict.Transcription != "" {
					text += " <b>/" + dict.Transcription + "/</b>\n"
					transcription = true
				}

				text += "\n<b>" + strconv.Itoa(i+1) + ". " + dict.Translations + "</b>"
				if dict.Comment != "" {
					text += " (" + dict.Comment + ")"
				}
				if dict.Examples != "" {
					dict.Examples = strings.ReplaceAll(dict.Examples, "\r\n", "\n<b>└</b>")
					text += "\n<b>└</b>" + dict.Examples
				}
				text += "\n"
			}

		}

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

		if user.Lang != "en" && from != user.Lang {
			tr, err := translate2.GoogleHTMLTranslate("en", user.Lang, text)
			if err != nil {
				warn(err)
			}
			text = tr.Text
		}




		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
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
		app.writeBotLog(callback.From.ID, "cb_meaning", text)
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
		text = strings.NewReplacer("<em>","<b>", "</em>", "</b>").Replace(text)

		msg := tgbotapi.NewMessage(callback.From.ID, text)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyToMessageID = callback.Message.MessageID
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
		msg.ReplyMarkup = keyboard
		if _, err = app.bot.Send(msg); err != nil {
			pp.Println(err)
		}

		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		app.writeBotLog(callback.From.ID, "cb_exmp", text)
	case "trs": // arr[1], arr[2] = from, to (in iso6391)
		callback.Message.ReplyToMessage.Text = strings.ToLower(callback.Message.ReplyToMessage.Text)

		var (
			trscript translate2.YandexTranscriptionResponse
			wg       sync.WaitGroup
			errs = make(chan error, 2)
		)
		from := translate2.ReversoIso6392(arr[1])
		to := translate2.ReversoIso6392(arr[2])

		rev, err := translate2.ReversoTranslate(from, to, callback.Message.ReplyToMessage.Text)
		if err != nil {
			warn(errors.New(err))
			return
		}

		var text string

		for i, result := range rev.ContextResults.Results {
			if result.Translation == "" {
				continue
			}
			result := result
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				tr, err := translate2.ReversoTranslate(to, from, result.Translation)
				if err != nil {
					errs <- errors.New(err)
				}
				text += "\n<b>" + result.Translation + "</b> <i>" + user.Localize(rev.ContextResults.Results[i].PartOfSpeech) + "</i>\n<b>└</b>"
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
					text += result.Translation
				}
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			trscript, err = translate2.YandexTranscription(translate2.ReversoIso6391(from), translate2.ReversoIso6391(to), callback.Message.ReplyToMessage.Text)
			if err != nil {
				errs <- err
			}
		}()

		wg.Wait()
		close(errs)

		if len(errs) > 0 {
			err = <-errs
			app.notifyAdmin("ошибка, но юзер не узнал", err.(*errors.Error), arr[1], arr[2], callback.Message.ReplyToMessage.Text)
		}

		//var text string = callback.Message.ReplyToMessage.Text
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
			call := tgbotapi.NewCallback(callback.ID, user.Localize("Available only for idioms, nouns, verbs and adjectives"))
			call.ShowAlert = true
			app.bot.Send(call)
			return
		}

		if _, err = app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              nil,
				DisableNotification:      true,
				AllowSendingWithoutReply: false,
			},
			Text:                  text,
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			call := tgbotapi.NewCallback(callback.ID, user.Localize("Available only for idioms, nouns, verbs and adjectives"))
			call.ShowAlert = true
			app.bot.Send(call)
		} else {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		}
		app.writeBotLog(callback.From.ID, "cb_dict", text)
	case "translate_to": // arr[1] - source language, arr[2] - direction to translate, in replied message there is source text
		go app.bot.Send(tgbotapi.NewChatAction(callback.From.ID, "typing"))
		from := arr[1]
		to := arr[2]
			pp.Println("translate_to:", arr)
		answer, err := app.SuperTranslate(from, to, callback.Message.ReplyToMessage.Text, callback.Message.ReplyToMessage.Entities)
		if err != nil {
			warn(err)
			return
		}
		pp.Println(answer)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔊 " + user.Localize("To voice"), fmt.Sprintf("speech_this_message_and_replied_one:%s:%s", from, to))))
		if answer.Examples {
			keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("💬 " + user.Localize("Examples"), fmt.Sprintf("exm:%s:%s", from, to)))
		}
		if answer.Translations {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📚 " + user.Localize("Translations"), fmt.Sprintf("trs:%s:%s", from, to))))
		}
		if answer.Dictionary {
			l := len(keyboard.InlineKeyboard) - 1
			if l < 0 {
				l = 0
			}
			keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData("ℹ️" + user.Localize("Dictionary"), fmt.Sprintf("dict:%s:%s", from, to)))
		}

		//if answer.Suggestions {
		//	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		//		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Lingvo", fmt.Sprintf("lingvo_vars:%s:%s:%d", from, to, 0))))
		//}

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{
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
		})
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		user.AddUsedLang(to)
		user.IncrUsings()
	case "translate_pagination": // arr[1] - source language of the text, arr[2] - offset
		from := arr[1]
		offset, err := strconv.Atoi(arr[2])
		if err != nil {
			warn(err)
			return
		}
		if offset < 0 || offset > len(codes) - 1 {
			warn(fmt.Errorf("offset is too big, len(codes) is %d, offset ois %d", len(codes), offset))
			return
		}

		count := 18
		if offset + count > len(codes) - 1 {
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
		kb, err := user.buildLangsPagination(offset, count,
			from,
			fmt.Sprintf("translate_to:%s:%s", from, "%s"),
			fmt.Sprintf("translate_pagination:%s:%d", from, back),
			fmt.Sprintf("translate_pagination:%s:%d", from, offset+count))
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
				BaseChat:              tgbotapi.BaseChat{
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

	//case "lingvo_vars": // arr[1] - from, arr[2] - to, arr[3] - offset. Text in REPLIED message
	//	from := arr[1]
	//	to := arr[2]
	//
	//	offset, err := strconv.Atoi(arr[3])
	//	if err != nil {
	//		warn(err)
	//		return
	//	}
	//	if offset < 0 || offset > len(codes) - 1 {
	//		warn(fmt.Errorf("offset is too big, len(codes) is %d, offset ois %d", len(codes), offset))
	//		return
	//	}
	//
	//	back := offset - 10
	//	if back < 0 {
	//		back = 0
	//	}
	//	next := offset + 10
	//	pp.Println(offset, next)
	//
	//
	//	suggestions, err := lingvo.Suggestions(from, to, callback.Message.ReplyToMessage.Text, 10, offset)
	//	if err != nil {
	//		warn(err)
	//		return
	//	}
	//
	//	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	//	for _, sug := range suggestions.Items {
	//		btn := tgbotapi.NewInlineKeyboardButtonData(sug.Heading, fmt.Sprintf("lingvo_word:%s:%s:%s", from, to, callback.Message.ReplyToMessage.Text))
	//		if len(sug.Heading) < 6 {
	//			l := len(keyboard.InlineKeyboard) - 1
	//			if l < 0 {
	//				l = 0
	//			}
	//			if keyboard.InlineKeyboard == nil {
	//				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
	//				continue
	//			}
	//			if len(keyboard.InlineKeyboard[l]) < 3 {
	//				keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], btn)
	//				continue
	//			}
	//
	//			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
	//		}
	//	}
	//
	//
	//	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
	//		tgbotapi.NewInlineKeyboardRow(
	//			tgbotapi.NewInlineKeyboardButtonData("<-- Back", fmt.Sprintf("lingvo_vars:%s:%s:%d", from, to, back)),
	//			tgbotapi.NewInlineKeyboardButtonData("Next -->", fmt.Sprintf("lingvo_vars:%s:%s:%d", from, to, next)),
	//		))
	//
	//	if reflect.DeepEqual(keyboard, *callback.Message.ReplyMarkup) {
	//		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
	//		return
	//	}
	//
	//	if _, err = app.bot.Send(tgbotapi.EditMessageReplyMarkupConfig{tgbotapi.BaseEdit{
	//		ChatID:          callback.From.ID,
	//		ChannelUsername: "",
	//		MessageID:       callback.Message.MessageID,
	//		InlineMessageID: "",
	//		ReplyMarkup:     &keyboard,
	//	}}); err != nil {
	//		warn(err)
	//	}
	//	app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
	//
	//case "lingvo_word": //arr[1] - from, arr[2] - to, arr[3] - word
	default:
		app.bot.Send(tgbotapi.NewMessage(callback.From.ID, "Action is expired. Send /start"))
	}
}