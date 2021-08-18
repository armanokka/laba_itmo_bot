package main

import (
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "os"
    "strconv"
    "strings"
)

func handleCallback(update *tgbotapi.Update) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Error, sorry"))
        WarnAdmin(err)
    }
    
    var UserLang string
    err := db.Model(&Users{}).Select("lang").Where("id = ?", update.CallbackQuery.From.ID).Limit(1).Find(&UserLang).Error
    if err != nil {
        warn(err)
        return
    }
    
    switch update.CallbackQuery.Data {
    case "none":
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "delete":
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          update.CallbackQuery.From.ID,
            MessageID:       update.CallbackQuery.Message.MessageID,
        })
        return
    case "sponsorship_pay":
        bot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Скоро будет дальше, а пока тут пусто"))
        if err := setUserStep(update.CallbackQuery.From.ID, ""); err != nil {
            warn(err)
        }
    }

    arr := strings.Split(update.CallbackQuery.Data, ":")
    if len(arr) == 0 {
        return
    }
    switch arr[0] {
    case "set_bot_lang": // arr[1] - lang code
        err := db.Model(&Users{}).Where("id = ?", update.CallbackQuery.From.ID).Limit(1).Update("lang", arr[1]).Error
        if err != nil {
            warn(err)
        }
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          update.CallbackQuery.From.ID,
            MessageID:       update.CallbackQuery.Message.MessageID,
        })
        bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, Localize("Now press /start 👈", arr[1])))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "speech": // arr[1] - lang code
        parts := strings.Split(update.CallbackQuery.Message.Text, "\n")
        update.CallbackQuery.Message.Text = strings.Join(parts[:len(parts)-1], "")
        sdec, err := translate.TTS(arr[1], update.CallbackQuery.Message.Text)
        if err != nil {
            if e, ok := err.(translate.TTSError); ok {
                if e.Code == 500 || e.Code == 414 {
                    bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Too big text"))
                    return
                }
            }
            bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Iternal error"))
            warn(err)
            return
        }
        f, err := os.CreateTemp("", "")
        if err != nil {
            warn(err)
            return
        }
        _, err = f.Write(sdec)
        if err != nil {
            warn(err)
            return
        }
        audio := tgbotapi.NewAudio(update.CallbackQuery.From.ID, f.Name())
        err = f.Close()
        if err != nil {
            warn(err)
            return
        }
        audio.Title = update.CallbackQuery.Message.Text
        audio.Performer = "@TransloBot"
        audio.ReplyToMessageID = update.CallbackQuery.Message.MessageID
        kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        audio.ReplyMarkup = kb
        bot.Send(audio)
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "variants": // arr[1] - from, arr[2] - to
        tr, err := translate.GoogleTranslate(arr[1], arr[2], update.CallbackQuery.Message.ReplyToMessage.Text)
        if err != nil {
            warn(err)
        }
        if len(tr.Variants) == 0 {
            warn(errors.New("не найдено вариантов перевода"))
        }
        var text string
        var limit uint8
        for _, variant := range tr.Variants {
            if limit >= 5 {
                break
            }
            text += "\n<b>" + variant.Word + "</b> - " + variant.Meaning
            limit++
        }
        msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, text)
        msg.ParseMode = tgbotapi.ModeHTML
        msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
        kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        msg.ReplyMarkup = kb
        bot.Send(msg)
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "set_my_lang_by_callback": // arr[1] - lang
        callback := tgbotapi.NewCallback(update.CallbackQuery.ID, Localize("Now your language is %s", UserLang, iso6391.Name(arr[1])))
        callback.ShowAlert = true
        bot.AnswerCallbackQuery(callback)
    
        err = db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Update("my_lang", arr[1]).Error
        if err != nil {
            warn(err)
            return
        }
    
        analytics.Bot(update.CallbackQuery.From.ID, "callback answered", "My language detected by callback")
    case "set_translate_lang_by_callback": // arr[1] - lang
        callback := tgbotapi.NewCallback(update.CallbackQuery.ID, Localize("Now translate language is %s", UserLang, iso6391.Name(arr[1])))
        callback.ShowAlert = true
        bot.AnswerCallbackQuery(callback)
    
        err = db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Update("to_lang", arr[1]).Error
        if err != nil {
            warn(err)
            return
        }
    
        analytics.Bot(update.CallbackQuery.From.ID, "callback answered", "Translate language detected by callback")
    case "country": // when user that want to buy sponsorship clicks on a button, arr[1] - lang code
        if IsTicked(update.CallbackQuery.Data, update.CallbackQuery.Message.ReplyMarkup) {
            UnTickByCallback(update.CallbackQuery.Data, update.CallbackQuery.Message.ReplyMarkup)
        } else {
            TickByCallback(update.CallbackQuery.Data, update.CallbackQuery.Message.ReplyMarkup)
        }
        callbacks := GetTickedCallbacks(update.CallbackQuery.Message.ReplyMarkup)
        langs := make([]string, 0)
        for _, callback := range callbacks {
            langs = append(langs, strings.Split(callback, ":")[1])
        }
        err := db.Model(&SponsorshipsOffers{}).Where("id = ?", update.CallbackQuery.From.ID).Update("to_langs", strings.Join(langs, ",")).Error
        if err != nil {
            warn(err)
            return
        }
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, *update.CallbackQuery.Message.ReplyMarkup))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "set_translate_lang_pagination": // arr[1] - offset
        offset, err := strconv.Atoi(arr[1])
        if err != nil {
            warn(err)
            return
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        var i int
        for code, lang := range langs {
            if i <= offset {
                continue
            }
            if i >= offset + 10 {
                break
            }
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_translate_lang_by_callback:"  + code)))
            i++
        }
        
        if offset + 10 < len(langs) {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:"+strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:"+strconv.Itoa(offset+10))))
        } else {
            langsLen := strconv.Itoa(len(langs))
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:"+strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(langsLen+"/"+langsLen, "none")))
        }

        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, keyboard))
    case "set_my_lang_pagination":
        offset, err := strconv.Atoi(arr[1])
        if err != nil {
            warn(err)
            return
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        var i int
        for code, lang := range langs {
            if i <= offset {
                continue
            }
            if i >= offset + 10 {
                break
            }
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_my_lang_by_callback:"  + code)))
            i++
        }
        if offset > 0 {
            if offset + 10 < len(langs) {
                keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:"+strconv.Itoa(offset-10)),
                    tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(langs)), "none"),
                    tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:"+strconv.Itoa(offset+10))))
            } else {
                langsLen := strconv.Itoa(len(langs))
                keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:"+strconv.Itoa(offset-10)),
                    tgbotapi.NewInlineKeyboardButtonData(langsLen+"/"+langsLen, "none")))
            }
        } else {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("10/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:10")))
        }
    
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, keyboard))
    case "sponsorship_set_days": // arr[1] - days, 1-30
        days, err := strconv.Atoi(arr[1])
        if err != nil {
            warn(err)
            return
        }
        if err := db.Model(&SponsorshipsOffers{}).Where("id = ?", update.CallbackQuery.From.ID).Updates(&SponsorshipsOffers{
            Days:    days,
        }).Error; err != nil {
            warn(err)
            return
        }
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Выберите языки пользователей, которые получат вашу рассылку.", UserLang))
        langs := map[string]string{"en": "🇬🇧 English", "it": "🇮🇹 Italiano", "uz":"🇺🇿 O'zbek tili", "de":"🇩🇪 Deutsch", "ru":"🇷🇺 Русский", "es":"🇪🇸 Español", "uk":"🇺🇦 Український", "pt":"🇵🇹 Português", "id":"🇮🇩 Indonesia"}
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for code, name := range langs {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(name, "country:"+code)))
        }
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Localize("Далее", UserLang), "sponsorship_pay")))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    }
}