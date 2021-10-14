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
)

func handleCallback(callback *tgbotapi.CallbackQuery) {
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
    case "delete":
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          callback.From.ID,
            MessageID:       callback.Message.MessageID,
        })
        return
    }

    arr := strings.Split(callback.Data, ":")
    if len(arr) == 0 {
        return
    }
    switch arr[0] {
    case "tick_country": // when user that want to buy sponsorship clicks on a button, arr[1] - lang code
        if IsTicked(callback.Data, callback.Message.ReplyMarkup) {
            UnTickByCallback(callback.Data, callback.Message.ReplyMarkup)
        } else {
            TickByCallback(callback.Data, callback.Message.ReplyMarkup)
        }
        callbacks := GetTickedCallbacks(*callback.Message.ReplyMarkup)
        langs := make([]string, 0)
        for _, callback := range callbacks {
            langs = append(langs, strings.Split(callback, ":")[1])
        }
        err := db.Model(&AdsOffers{}).Where("id = ?", callback.From.ID).Update("to_langs", strings.Join(langs, ",")).Error
        if err != nil {
            warn(err)
            return
        }
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, *callback.Message.ReplyMarkup))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
    case "register_bot_lang": // arr[1] - lang code
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
        user.Update(Users{Lang: arr[1]})
        SendMenu(user)
    case "set_bot_lang": // arr[1] - lang code
        user.Update(Users{Lang: arr[1]})
        SendMenu(user)
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
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
    //case "translate_pagination": // arr[1] - offset, arr[2] - to
    //    offset, err := strconv.Atoi(arr[1])
    //    if err != nil {
    //        warn(err)
    //        return
    //    }
    //    keyboard := tgbotapi.NewInlineKeyboardMarkup()
    //    for i, code := range codes[offset:] {
    //        if i >= LanguagesPaginationLimit {
    //            break
    //        }
    //        lang, ok := langs[code]
    //        if !ok {
    //            warn(errors.New("no such code "+ code + " in langs"))
    //            return
    //        }
    //        if i % 2 == 0 {
    //            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "translate:"  + code + ":" + arr[2])))
    //        } else {
    //            l := len(keyboard.InlineKeyboard)-1
    //            keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "translate:"  + code + ":" + arr[2]))
    //        }
    //    }
    //
    //    prev := offset - 20
    //    if prev < 0 {
    //        prev = 0
    //    }
    //    next := offset + LanguagesPaginationLimit
    //    if next > len(codes) - 1 {
    //        next = len(codes) - 1
    //    }
    //    keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
    //        tgbotapi.NewInlineKeyboardButtonData("◀", "translate_pagination:" + strconv.Itoa(prev) + ":" + arr[2]),
    //        tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(codes)), "none"),
    //        tgbotapi.NewInlineKeyboardButtonData("▶", "translate_pagination:"+strconv.Itoa(next) + ":" + arr[2])))
    //    keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
    //
    //    msg := tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, user.Localize("Select the source language of your text if it was not defined correctly"), keyboard)
    //    bot.Send(msg)
    //    bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
    //case "translate": // arr[1] - from, arr[2] - to. Text in replied message
    //    if err := SendTranslation(user, arr[1], arr[2], callback.Message.ReplyToMessage.ReplyToMessage.Text, callback.Message.ReplyToMessage.ReplyToMessage.MessageID); err != nil {
    //        warn(err)
    //    }
    //    bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
    case "examples": // arr[1], arr[2], arr[3] = from, to, source text. Target text in replied message
        samples, err := translate.ReversoQueryService(arr[3], arr[1], callback.Message.ReplyToMessage.Text, arr[2])
        if err != nil {
            warn(err)
            return
        }

        pp.Println(arr[1], arr[2], callback.Message.ReplyToMessage.Text)
        pp.Println(samples)

        var text string

        if len(samples.List) > 10 {
            samples.List = samples.List[:10]
        }

        for i, sample := range samples.List {
            sample.TText = strings.ReplaceAll(sample.TText, ` class="both"`, "")
            text += "\n\n<b>" + strconv.Itoa(i+1) + ".</b> " + sample.TText + "\n<b>└</b>"+sample.SText
        }
        text += "\n"
        for _, suggestion := range samples.Suggestions {
            text += "\n>" + suggestion.Suggestion
        }
        msg := tgbotapi.NewMessage(callback.From.ID, text)
        msg.ParseMode = tgbotapi.ModeHTML
        keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        msg.ReplyMarkup = keyboard
        if _, err = bot.Send(msg); err != nil {
            pp.Println(err)
        }

        bot.Send(tgbotapi.NewCallback(callback.ID, ""))
    case "dictionary": // arr[1], arr[2] = from, to (in iso6391)
        callback.Message.ReplyToMessage.Text = strings.ToLower(callback.Message.ReplyToMessage.Text)

        tr, err := translate.GoogleTranslateSingle(arr[1], arr[2], callback.Message.ReplyToMessage.Text)
        if err != nil {
            warn(err)
            return
        }


        var text string

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

        if _, err = bot.Send(message); err != nil {
            pp.Println(err)
        }
        bot.Send(tgbotapi.NewCallback(callback.ID, ""))
    case "set_my_lang_by_callback": // arr[1] - lang
        user.Update(Users{MyLang: arr[1]})

        bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, user.Localize("Ваш язык %s. Выберите Ваш язык.", iso6391.Name(user.MyLang)), *callback.Message.ReplyMarkup))

        call := tgbotapi.NewCallback(callback.ID, user.Localize("Now your language is %s", iso6391.Name(arr[1])))
        call.ShowAlert = true
        bot.AnswerCallbackQuery(call)

        analytics.Bot(callback.From.ID, "callback answered", "My language detected by callback")
    case "set_translate_lang_by_callback": // arr[1] - lang
        user.Update(Users{ToLang: arr[1]})

        bot.Send(tgbotapi.NewEditMessageTextAndMarkup(callback.From.ID, callback.Message.MessageID, user.Localize("Сейчас бот переводит на %s. Выберите язык для перевода", iso6391.Name(user.ToLang)), *callback.Message.ReplyMarkup))

        call := tgbotapi.NewCallback(callback.ID, user.Localize("Now translate language is %s", iso6391.Name(arr[1])))
        call.ShowAlert = true
        bot.AnswerCallbackQuery(call)
    
        analytics.Bot(callback.From.ID, "callback answered", "Translate language detected by callback")

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
        if offset == 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:0"),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(codes)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:"+strconv.Itoa(offset+LanguagesPaginationLimit))))
        } else if offset + LanguagesPaginationLimit < len(codes) {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:"+strconv.Itoa(offset-LanguagesPaginationLimit)),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(codes)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:"+strconv.Itoa(offset+LanguagesPaginationLimit))))
        } else {
            langsLen := strconv.Itoa(len(codes))
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:"+strconv.Itoa(offset-LanguagesPaginationLimit)),
                tgbotapi.NewInlineKeyboardButtonData(langsLen+"/"+langsLen, "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:"+strconv.Itoa(offset)),
                ))
        }

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
        if offset <= 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:0"),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(codes)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:"+strconv.Itoa(LanguagesPaginationLimit))))
        } else if offset + LanguagesPaginationLimit < len(codes) {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:"+strconv.Itoa(offset-LanguagesPaginationLimit)),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(codes)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:"+strconv.Itoa(offset+LanguagesPaginationLimit))))
        } else {
            langsLen := strconv.Itoa(len(codes))
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:"+strconv.Itoa(offset-LanguagesPaginationLimit)),
                tgbotapi.NewInlineKeyboardButtonData(langsLen+"/"+langsLen, "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:"+strconv.Itoa(offset))),
            )
        }
        
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(callback.From.ID, callback.Message.MessageID, keyboard))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
    }
}