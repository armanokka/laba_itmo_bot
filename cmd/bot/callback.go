package bot

import (
	"database/sql"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	iso6391 "github.com/emvi/iso-639-1"
	"github.com/go-errors/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
)

func (app app) onCallbackQuery(callback tgbotapi.CallbackQuery) {
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewCallback(callback.ID, "Error, sorry"))
		app.notifyAdmin(err)
		logrus.Error(err)
	}

	user := app.loadUser(callback.From.ID, warn)
	user.Fill()
	defer user.UpdateLastActivity()

	arr := strings.Split(callback.Data, ":")
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
	case "I'm developer":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     nil,
			},
			Text:                  user.Localize("Advice Translo API"),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		})
	case "good_bot":
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				MessageID:       callback.Message.MessageID,
			},
			Text:                  user.Localize("Donation"),
		})

	case "show_my_lang_tab":
		user.SendMyLangTab(callback.Message.MessageID)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "show_langs_by_letter_and_set_my_lang": // arr[1] - letter
		keyboard := buildOneLetterKeyboard(arr[1], "set_my_lang_by_callback:%s", "show_my_lang_tab")
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &keyboard,
			},
			Text:                  applyEntitiesHtml(callback.Message.Text, callback.Message.Entities),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		})
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))

	case "show_to_lang_tab":
		user.SendToLangTab(callback.Message.MessageID)
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "show_langs_by_letter_and_set_to_lang":
		keyboard := buildOneLetterKeyboard(arr[1], "set_to_lang_by_callback:%s", "show_to_lang_tab")
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &keyboard,
			},
			Text:                  applyEntitiesHtml(callback.Message.Text, callback.Message.Entities),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		})
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))

	case "set_bot_lang_and_register": // arr[1] - bot lang
		tolang := "en"
		if arr[1] == tolang {
			tolang = "fr"
		}
		var exists bool
		if err := app.db.Model(&tables.Users{}).
			Raw("SELECT EXISTS(SELECT lang FROM users WHERE id=?)", callback.From.ID).
			Find(&exists).Error; err != nil {
			warn(err)
		}

		if exists {
			user.Update(tables.Users{Lang: arr[1]})
		}

		if !exists{
			if err := app.db.Create(&tables.Users{
				ID:         callback.From.ID,
				MyLang:     arr[1],
				ToLang:     tolang,
				Act:        sql.NullString{},
				Mailing:    true,
				Lang:       arr[1],
			}).Error; err != nil {
				warn(err)
			}
		}

		user.Fill()

		msg := user.StartMessage()

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{
				ChatID:                   callback.From.ID,
				ReplyMarkup:              msg.Keyboard,
				DisableNotification:      true,
				AllowSendingWithoutReply: true,
			},
			ParseMode: tgbotapi.ModeHTML,
			Text:                  msg.Text,
		})
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
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
	case "dictonary": // arr[1] - lang, arr[2] - text
		meaning, err := translate2.GoogleDictionary(arr[1], arr[2])
		if err != nil {
			warn(err)
			return
		}
		text := ""
		var i int

		for _, data := range meaning.DictionaryData {
			i += 1 // тут i+1 без проверки индекса
			for idx, entry := range data.Entries {
				if idx > 0 {
					i += 1
				}
				for idx, family := range entry.SenseFamilies {
					if idx > 0 {
						i += 1
					}
					for idx, sense := range family.Senses {
						if idx > 0 {
							i += 1
						}
						text += "\n\n<b>" + strconv.Itoa(i) + ".</b> " + sense.Definition.Text
					}
				}
			}
		}
		text = strings.ReplaceAll(text, "\n", "<br>")
		if arr[1] != user.MyLang {
			tr, err := translate2.GoogleHTMLTranslate(arr[1], user.MyLang, text)
			if err != nil {
				warn(err)
				return
			}
			text = tr.Text
		}
		text = strings.ReplaceAll(text, "<br>", "\n")

		if len(strings.Fields(text)) == 0 {
			app.bot.Send(tgbotapi.CallbackConfig{
				CallbackQueryID: callback.ID,
				Text:            "Error",
				ShowAlert:       true,
				URL:             "",
				CacheTime:       0,
			})
			return
		}

		msg := tgbotapi.NewMessage(callback.From.ID, text)
		msg.ReplyToMessageID = callback.Message.MessageID
		msg.ParseMode = tgbotapi.ModeHTML
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
		msg.ReplyMarkup = keyboard
		app.bot.Send(msg)
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		app.writeBotLog(callback.From.ID, "cb_meaning", text)
	case "examples": // arr[1] - from, arr[2] - to. Target text in replied message
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
		var firstlySourceExamples bool
		if arr[1] == user.MyLang {
			firstlySourceExamples = true
		} else {
			firstlySourceExamples = false
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
				if firstlySourceExamples {
					text += result.SourceExamples[i] + "\n<b>└</b>" + result.TargetExamples[i]
				} else {
					text += result.TargetExamples[i] + "\n<b>└</b>" + result.SourceExamples[i]
				}
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
	case "translations": // arr[1], arr[2] = from, to (in iso6391)
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
		pp.Println(trscript)
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

		message := tgbotapi.NewMessage(callback.From.ID, text)
		message.ParseMode = tgbotapi.ModeHTML
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
		message.ReplyMarkup = keyboard
		message.ReplyToMessageID = callback.Message.MessageID

		if _, err = app.bot.Send(message); err != nil {
			call := tgbotapi.NewCallback(callback.ID, user.Localize("Available only for idioms, nouns, verbs and adjectives"))
			call.ShowAlert = true
			app.bot.Send(call)
		} else {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		}
		app.writeBotLog(callback.From.ID, "cb_dict", text)
	case "set_my_lang_by_callback": // arr[1] - lang
		user.Update(tables.Users{MyLang: arr[1]})
		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     callback.Message.ReplyMarkup,
			},
			Text:                  user.Localize("Ваш язык <b>%s</b>. Выберите Ваш язык.", langs[user.MyLang].Name),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		})

		call := tgbotapi.NewCallback(callback.ID, user.Localize("Now your language is %s", iso6391.Name(arr[1])))
		call.ShowAlert = true
		app.bot.AnswerCallbackQuery(call)

		app.analytics.Bot(callback.From.ID, "callback answered", "My language detected by callback")
		app.writeBotLog(callback.From.ID, "cb_set_my_lang", call.Text)
	case "set_to_lang_by_callback": // arr[1] - lang
		user.Update(tables.Users{ToLang: arr[1]})

		app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit:              tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     callback.Message.ReplyMarkup,
			},
			Text:                  user.Localize("Сейчас бот переводит на <b>%s</b>. Выберите язык, на который хотите переводить", langs[user.ToLang].Name),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		})

		call := tgbotapi.NewCallback(callback.ID, user.Localize("Now translate language is %s", iso6391.Name(arr[1])))
		call.ShowAlert = true
		app.bot.AnswerCallbackQuery(call)

		app.analytics.Bot(callback.From.ID, "callback answered", "Translate language detected by callback")
		app.writeBotLog(callback.From.ID, "cb_set_to_lang", call.Text)
	case "set_translate_lang_pagination": // arr[1] - offset
		offset, err := strconv.Atoi(arr[1])
		if err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		for i, code := range codes[offset:] {
			if i >= config.LanguagesPaginationLimit {
				break
			}
			lang, ok := langs[code]
			if !ok {
				warn(errors.New("no such code "+ code + " in langs"))
				return
			}
			if i % 2 == 0 {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_translate_lang_by_callback:"  + code)))
			} else {
				l := len(keyboard.InlineKeyboard)-1
				keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_translate_lang_by_callback:"  + code))
			}
		}
		current := offset
		prev := "0"
		if offset > 0 {
			prev = strconv.Itoa(offset - config.LanguagesPaginationLimit)
		}
		next := strconv.Itoa(offset + config.LanguagesPaginationLimit)
		if offset + config.LanguagesPaginationLimit > len(codes) {
			next = strconv.Itoa(offset)
			current = len(codes)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:" + prev),
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(current) + "/"+strconv.Itoa(len(codes)), "none"),
			tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:" + next)))

		app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, keyboard))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "set_my_lang_pagination":
		offset, err := strconv.Atoi(arr[1])
		if err != nil {
			warn(err)
			return
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup()
		for i, code := range codes[offset:] {
			if i >= config.LanguagesPaginationLimit {
				break
			}
			lang, ok := langs[code]
			if !ok {
				warn(errors.New("no such code "+ code + " in langs"))
				return
			}
			if i % 2 == 0 {
				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_my_lang_by_callback:"  + code)))
			} else {
				l := len(keyboard.InlineKeyboard)-1
				keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_my_lang_by_callback:"  + code))
			}
		}
		current := offset
		prev := "0"
		if offset > 0 {
			prev = strconv.Itoa(offset - config.LanguagesPaginationLimit)
		}
		next := strconv.Itoa(offset + config.LanguagesPaginationLimit)
		if offset + config.LanguagesPaginationLimit > len(codes) {
			next = strconv.Itoa(offset)
			current = len(codes)
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:" + prev),
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(current) + "/"+strconv.Itoa(len(codes)), "none"),
			tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:" + next)))

		app.bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, keyboard))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	default:
		app.notifyAdmin("неизвестный колбэк: " + callback.Data)
	}
}