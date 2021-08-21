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
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{}))
        bot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Скоро будет дальше, а пока тут пусто"))
    // case "swap_languages":
    //     var user Users
    //     err = db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.CallbackQuery.From.ID).Find(&user).Error
    //     if err != nil {
    //         warn(err)
    //         return
    //     }
    //     err = db.Model(&Users{}).Where("id = ?", update.CallbackQuery.From.ID).Updates(Users{MyLang: user.ToLang, ToLang: user.MyLang}).Error
    //     if err != nil {
    //         warn(err)
    //         return
    //     }
    //     keyboard := tgbotapi.NewReplyKeyboard(
    //         tgbotapi.NewKeyboardButtonRow(
    //             tgbotapi.NewKeyboardButton(Localize("🙎‍♂️Profile", user.Lang)),
    //         ),
    //         tgbotapi.NewKeyboardButtonRow(
    //             tgbotapi.NewKeyboardButton(Localize("💬 Bot language", user.Lang)),
    //             tgbotapi.NewKeyboardButton(Localize("💡 Instruction", user.Lang)),
    //         ),
    //         tgbotapi.NewKeyboardButtonRow(
    //             tgbotapi.NewKeyboardButton(Localize("My Language", user.Lang)),
    //             tgbotapi.NewKeyboardButton(Localize("Translate Language", user.Lang)),
    //         ),
    //     )
    //     msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, Localize("/start", user.Lang, user.MyLang, user.ToLang), keyboard)
    //
    //     msg.ReplyMarkup = keyboard
    //     msg.ParseMode = tgbotapi.ModeHTML
    //     bot.Send(msg)
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
        text := update.CallbackQuery.Message.Text
        if update.CallbackQuery.Message.Caption != "" {
            text = update.CallbackQuery.Message.Caption
        }
        sdec, err := translate.TTS(arr[1], text)
        if err != nil {
            if e, ok := err.(translate.TTSError); ok {
                if e.Code == 500 || e.Code == 414 {
                    callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Too big text")
                    callback.ShowAlert = true
                    bot.AnswerCallbackQuery(callback)
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
        defer func() {
            if err = f.Close(); err != nil {
                warn(err)
            }
        }()
        _, err = f.Write(sdec)
        if err != nil {
            warn(err)
            return
        }
        audio := tgbotapi.NewAudio(update.CallbackQuery.From.ID, f.Name())
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
        for i, code := range codes[offset:] {
            if i >= 10 {
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
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:"+strconv.Itoa(offset+10))))
        } else if offset + 10 < len(langs) {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:"+strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:"+strconv.Itoa(offset+10))))
        } else {
            langsLen := strconv.Itoa(len(langs))
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:"+strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(langsLen+"/"+langsLen, "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:"+strconv.Itoa(offset)),
                ))
        }

        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, keyboard))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "set_my_lang_pagination":
        offset, err := strconv.Atoi(arr[1])
        if err != nil {
            warn(err)
            return
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for i, code := range codes[offset:] {
            if i >= 10 {
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
        if offset == 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:0"),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:"+strconv.Itoa(offset+10))))
        } else if offset + 10 < len(langs) {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:"+strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(arr[1] + "/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:"+strconv.Itoa(offset+10))))
        } else {
            langsLen := strconv.Itoa(len(langs))
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:"+strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(langsLen+"/"+langsLen, "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:"+strconv.Itoa(offset))),
            )
        }
        
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, keyboard))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
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
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{}))

        msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, Localize("Выберите языки пользователей, которые получат вашу рассылку.", UserLang))
        langs := map[string]string{"en": "🇬🇧 English", "it": "🇮🇹 Italiano", "uz":"🇺🇿 O'zbek tili", "de":"🇩🇪 Deutsch", "ru":"🇷🇺 Русский", "es":"🇪🇸 Español", "uk":"🇺🇦 Український", "pt":"🇵🇹 Português", "id":"🇮🇩 Indonesia"}
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for code, name := range langs {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(name, "country:"+code)))
        }
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Localize("Далее", UserLang), "sponsorship_pay")))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    case "translate_by_callback": // arr[1] - source lang of a text in replied message, arr[2] - to lang
        var text = update.CallbackQuery.Message.ReplyToMessage.Text
        if update.CallbackQuery.Message.ReplyToMessage.Caption != "" {
            text = update.CallbackQuery.Message.ReplyToMessage.Caption
        }
        
        tr, err := translate.GoogleTranslate(arr[1], arr[2], text)
        if err != nil {
            warn(err)
            return
        }
        if tr.Text == "" {
            bot.Send(tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, Localize("Empty result", UserLang), *update.CallbackQuery.Message.ReplyMarkup))
            callback := tgbotapi.NewCallback(update.CallbackQuery.ID, Localize("Empty result", UserLang))
            callback.ShowAlert = true
            bot.AnswerCallbackQuery(callback)
            return
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData(Localize("To voice", UserLang), "speech:"+arr[2]),
            ),
        )
        if len(tr.Variants) > 0 {
            l := len(tr.Variants)-1
            keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(Localize("Variants", UserLang), "variants:"+arr[1]+":"+arr[2]))
        }
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, update.CallbackQuery.Message.ReplyMarkup.InlineKeyboard...)
        bot.Send(tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, tr.Text, keyboard))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "translate_to_other_languages_pagination": // arr[1] - source lang, arr[2] - pagination offset
        offset, err := strconv.Atoi(arr[2])
        if err != nil {
            warn(err)
            return
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for i, code := range codes[offset:] {
            if i >= 10 {
                break
            }
            lang, ok := langs[code]
            if !ok {
                warn(errors.New("no such code "+ code + " in langs"))
                return
            }
            if i % 2 == 0 {
                keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "translate_by_callback:"  + arr[1] + ":" + code)))
            } else {
                l := len(keyboard.InlineKeyboard)-1
                keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "translate_by_callback:"  + arr[1] + ":" + code))
            }
        }
        
        if offset == 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "translate_to_other_languages_pagination:" + arr[1] + ":0"),
                tgbotapi.NewInlineKeyboardButtonData("0/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "translate_to_other_languages_pagination:" + arr[1] + ":"+strconv.Itoa(offset+10))))
        } else if offset + 10 < len(langs) {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "translate_to_other_languages_pagination:" + arr[1] + ":"+strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(arr[2] + "/"+strconv.Itoa(len(langs)), "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "translate_to_other_languages_pagination:" + arr[1] + ":"+strconv.Itoa(offset+10))))
        } else {
            langsLen := strconv.Itoa(len(langs))
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("◀", "translate_to_other_languages_pagination:" + arr[1] + ":" +strconv.Itoa(offset-10)),
                tgbotapi.NewInlineKeyboardButtonData(langsLen+"/"+langsLen, "none"),
                tgbotapi.NewInlineKeyboardButtonData("▶", "translate_to_other_languages_pagination:" + arr[1] + ":"+strconv.Itoa(offset))))
        }
    
        bot.Send(tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, keyboard))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    }
}