package main

import (
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "github.com/sirupsen/logrus"
    "strconv"
    "strings"
    "sync"
)

func handleCallback(callback tgbotapi.CallbackQuery) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewCallback(callback.ID, "Error, sorry"))
        WarnAdmin(err)
        logrus.Error(err)
    }

    user := NewUser(callback.From.ID, warn)
    user.Fill()
    
    switch callback.Data {
    case "none":
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
        user.WriteBotLog("cb_none", "")
        return
    case "delete":
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          callback.From.ID,
            MessageID:       callback.Message.MessageID,
        })
        user.WriteBotLog("cb_delete", "")
        return
    }

    arr := strings.Split(callback.Data, ":")
    if len(arr) == 0 {
        return
    }
    switch arr[0] {
    case "speech_this_message_and_replied_one": // arr[1] - from, arr[2] - to
        text := callback.Message.Text
        if callback.Message.Caption != "" {
            text = callback.Message.Caption
        }
        if err := sendSpeech(arr[1], callback.Message.ReplyToMessage.Text, callback.ID, user); err != nil { // озвучиваем непереведенное сообщение
            warn(err)
            return
        }
        if err := sendSpeech(arr[2], text, callback.ID, user); err != nil { // озвучиваем переведенное сообщение
            warn(err)
            return
        }
        user.WriteUserLog("cb_voice")
    case "dictonary": // arr[1] - lang, arr[2] - text
        meaning, err := translate.GoogleDictionary(arr[1], arr[2])
        if err != nil {
            warn(err)
            return
        }
        text := ""
        for _, data := range meaning.DictionaryData {
            for _, entry := range data.Entries {
                for _, family := range entry.SenseFamilies {
                    for _, sense := range family.Senses {
                        text += "\n✅ " + sense.Definition.Text
                    }
                }
            }
        }
        msg := tgbotapi.NewMessage(callback.From.ID, text)
        msg.ReplyToMessageID = callback.Message.MessageID
        keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
        bot.Send(tgbotapi.NewCallback(callback.ID, ""))
        user.WriteBotLog("cb_meaning", text)
    case "examples": // arr[1], arr[2], arr[3] = from, to, source text. Target text in replied message
        samples, err := translate.ReversoQueryService(arr[3], arr[1], callback.Message.ReplyToMessage.Text, arr[2])
        if err != nil {
            warn(err)
            return
        }

        var text string

        if len(samples.List) > 10 {
            samples.List = samples.List[:10]
        }

        for i, sample := range samples.List {
            sample.TText = strings.ReplaceAll(sample.TText, ` class="both"`, "")
            text += "\n\n<b>" + strconv.Itoa(i+1) + ".</b> " + sample.TText + "\n└" + sample.SText
        }
        text += "\n"
        for _, suggestion := range samples.Suggestions {
            text += "\n>" + suggestion.Suggestion
        }
        text = strings.NewReplacer("<em>", "<b>", "</em>", "</b>").Replace(text)
        msg := tgbotapi.NewMessage(callback.From.ID, text)
        msg.ParseMode = tgbotapi.ModeHTML
        msg.ReplyToMessageID = callback.Message.MessageID
        keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        msg.ReplyMarkup = keyboard
        if _, err = bot.Send(msg); err != nil {
            pp.Println(err)
        }

        bot.Send(tgbotapi.NewCallback(callback.ID, ""))
        user.WriteBotLog("cb_exmp", text)
    case "translations": // arr[1], arr[2] = from, to (in iso6391)
        callback.Message.ReplyToMessage.Text = strings.ToLower(callback.Message.ReplyToMessage.Text)

        var (
            tr translate.GoogleTranslateSingleResult
            trscript translate.YandexTranscriptionResponse
            wg sync.WaitGroup
            err error
            errs = make(chan error, 2)
        )

        wg.Add(1)
        go func() {
            defer wg.Done()
            tr, err = translate.GoogleTranslateSingle(arr[1], arr[2], callback.Message.ReplyToMessage.Text)
            if err != nil {
                errs <- err
            }
        }()

        wg.Add(1)
        go func() {
            defer wg.Done()
            trscript, err = translate.YandexTranscription(arr[1], arr[2], callback.Message.ReplyToMessage.Text)
            if err != nil {
                errs <- err
            }
        }()

        wg.Wait()
        if len(errs) > 0 {
            err = <-errs
            WarnAdmin("ошибка, но юзер не узнал", err, arr[1], arr[2], callback.Message.ReplyToMessage.Text)
        }

        var text string = callback.Message.ReplyToMessage.Text
        if trscript.StatusCode == 200 {
            text += " <b>" + trscript.Transcription + "</b> <i>" + trscript.Pos + "</i>"
        }


        for _, dict := range tr.Dict {
            text += "\n"
            for i, term := range dict.Terms {
                text += "\n<b>" + term + "\n└</b>"
                //l := len(dict.Entry[i].ReverseTranslation) - 1
                for idx, entry := range dict.Entry[i].ReverseTranslation {
                    if idx > 0 {
                        text += ", "
                    }
                    text += entry
                }
            }
        }

        if text == "" {
            call := tgbotapi.NewCallback(callback.ID, user.Localize("Available only for idioms, nouns, verbs and adjectives"))
            call.ShowAlert = true
            bot.Send(call)
            return
        }

        message := tgbotapi.NewMessage(callback.From.ID, text)
        message.ParseMode = tgbotapi.ModeHTML
        keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        message.ReplyMarkup = keyboard
        message.ReplyToMessageID = callback.Message.MessageID

        if _, err = bot.Send(message); err != nil {
            pp.Println(err)
        }
        bot.Send(tgbotapi.NewCallback(callback.ID, ""))
        user.WriteBotLog("cb_dict", text)
    case "set_my_lang_by_callback": // arr[1] - lang
        user.Update(Users{MyLang: arr[1]})

        bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, user.Localize("Ваш язык %s. Выберите Ваш язык.", iso6391.Name(user.MyLang)), *callback.Message.ReplyMarkup))

        call := tgbotapi.NewCallback(callback.ID, user.Localize("Now your language is %s", iso6391.Name(arr[1])))
        call.ShowAlert = true
        bot.AnswerCallbackQuery(call)

        analytics.Bot(callback.From.ID, "callback answered", "My language detected by callback")
        user.WriteBotLog("cb_set_my_lang", call.Text)
    case "set_translate_lang_by_callback": // arr[1] - lang
        user.Update(Users{ToLang: arr[1]})

        bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, user.Localize("Сейчас бот переводит на %s. Выберите язык для перевода", iso6391.Name(user.ToLang)), *callback.Message.ReplyMarkup))

        call := tgbotapi.NewCallback(callback.ID, user.Localize("Now translate language is %s", iso6391.Name(arr[1])))
        call.ShowAlert = true
        bot.AnswerCallbackQuery(call)
    
        analytics.Bot(callback.From.ID, "callback answered", "Translate language detected by callback")
        user.WriteBotLog("cb_set_to_lang", call.Text)
    case "set_translate_lang_pagination": // arr[1] - offset
        offset, err := strconv.Atoi(arr[1])
        if err != nil {
            warn(err)
            return
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for i, code := range codes[offset:] {
            if i >= LanguagesPaginationLimit {
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
            prev = strconv.Itoa(offset - LanguagesPaginationLimit)
        }
        next := strconv.Itoa(offset + LanguagesPaginationLimit)
        if offset + LanguagesPaginationLimit > len(codes) {
            next = strconv.Itoa(offset)
            current = len(codes)
        }
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:" + prev),
            tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(current) + "/"+strconv.Itoa(len(codes)), "none"),
            tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:" + next)))

        bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, keyboard))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
    case "set_my_lang_pagination":
        offset, err := strconv.Atoi(arr[1])
        if err != nil {
            warn(err)
            return
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for i, code := range codes[offset:] {
            if i >= LanguagesPaginationLimit {
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
            prev = strconv.Itoa(offset - LanguagesPaginationLimit)
        }
        next := strconv.Itoa(offset + LanguagesPaginationLimit)
        if offset + LanguagesPaginationLimit > len(codes) {
            next = strconv.Itoa(offset)
            current = len(codes)
        }
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:" + prev),
            tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(current) + "/"+strconv.Itoa(len(codes)), "none"),
            tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:" + next)))

        bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, keyboard))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
    }
}