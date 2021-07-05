package main

import (
    "database/sql"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "strings"
)

func handleMessage(update *tgbotapi.Update) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, error caused.\n\nAdmin: Please, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)"))
        WarnAdmin(err)
    }
    
    if update.Message.Chat.ID < 0 {
        return
    }
    switch update.Message.Text {
    case "/start", "/start from_inline", "⬅Back":
        var userExists bool
        err := db.Raw("SELECT EXISTS(SELECT id FROM users WHERE id=?)", update.Message.Chat.ID).Find(&userExists).Error
        if err != nil {
            warn(err)
            return
        }
        if !userExists {
            fromLang := update.Message.From.LanguageCode
            translateLang := "fr"
            if fromLang == "" { // код языка недоступен
                fromLang = "en" // О, вы из Англии
            }
            if fromLang == "fr" { // перед нами француз, зачем ему переводить на французский?
                translateLang = "es"
            }
            err = db.Create(&Users{
                ID:     update.Message.Chat.ID,
                MyLang: fromLang,
                ToLang: translateLang,
                Act:    sql.NullString{},
            }).Error
            if err != nil {
                warn(err)
                return
            }
        }
        
        var user Users
        err = db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.Message.Chat.ID).Find(&user).Error
        if err != nil {
            warn(err)
            return
        }
        
        user.MyLang = iso6391.Name(user.MyLang)
        user.ToLang = iso6391.Name(user.ToLang)
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your language is - <b>" + user.MyLang + "</b>, and the language for translation is - <b>" + user.ToLang + "</b>.")
        keyboard := tgbotapi.NewReplyKeyboard(
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton("💡Instruction")),
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton("My Language")),
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton("Translate Language")))
        msg.ReplyMarkup = keyboard
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
        
        err = setUserStep(update.Message.Chat.ID, "")
        if err != nil {
            warn(err)
            return
        }
    case "My Language", "/my_lang":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "To setup <b>your language</b>, do <b>one</b> of the following: 👇\n\nℹ️ Send <b>few words</b> in <b>your</b> language, for example: \"<code>Hi, how are you today?</code>\" - language will be English, or \"<code>L'amour ne fait pas d'erreurs</code>\" - language will be French, and so on.\nℹ️ Or send the <b>name</b> of your language <b>in English</b>, e.g. \"<code>Russian</code>\", or \"<code>Japanese</code>\", or  \"<code>Arabic</code>\", e.t.c.")
        msg.ParseMode = tgbotapi.ModeHTML
        keyboard:= tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("⬅Back")))
        msg.ReplyMarkup = &keyboard
        bot.Send(msg)
        
        err := setUserStep(update.Message.Chat.ID, "set_my_lang")
        if err != nil {
            warn(err)
            return
        }
    
    case "Translate Language", "/to_lang":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "To setup <b>translate language</b>, do <b>one</b> of the following: 👇\n\nℹ️ Send <b>few words</b> in language <b>into you want to translate</b>, for example: \"<code>Hi, how are you today?</code>\" - language will be English, or \"<code>L'amour ne fait pas d'erreurs</code>\" - language will be French, and so on.\nℹ️ Or send the <b>name</b> of language <b>into you want to translate, in English</b>, e.g. \"<code>Russian</code>\", or \"<code>Japanese</code>\", or  \"<code>Arabic</code>\", e.t.c.")
        msg.ParseMode = tgbotapi.ModeHTML
        keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("⬅Back")))
        msg.ReplyMarkup = &keyboard
        bot.Send(msg)
        
        err := setUserStep(update.Message.Chat.ID, "set_translate_lang")
        if err != nil {
            warn(err)
            return
        }
    
    case "💡Instruction", "/help":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "<b>What can this bot do?</b>\n▫️ Translo allows you to translate your messages into over than 100 languages. (117)\n<b>How to translate message?</b>\n▫️ Firstly, you have to setup your lang (default: English), then setup translate lang (default; Arabic) then send text messages and bot will translate them quickly.\n<b>How to setup my lang?</b>\n▫️ Send /my_lang then send any message <b>IN YOUR LANGUAGE</b>. Bot will detect and suggest you some variants. Select your lang. Done.\n<b>How to setup translate lang?</b>\n▫️ Send /to_lang then send any message <b>IN LANGUAGE YOU WANT TRANSLATE</b>. Bot will detect and suggest you some variants. Select your lang. Done.\n<b>I have a suggestion or I found bug!</b>\n▫️ 👉 Contact me pls - @armanokka")
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
    
    default: // Сообщение не является командой.
        userStep, err := getUserStep(update.Message.Chat.ID)
        if err != nil {
            warn(err)
            return
        }
        switch userStep {
        case "set_my_lang":
            lowerUserMsg := strings.ToLower(update.Message.Text)
            
            if nameOfLang := iso6391.Name(lowerUserMsg); nameOfLang != "" { // Юзер отправил сразу код языка
                err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("my_lang", lowerUserMsg).Error
                if err != nil {
                    warn(err)
                    return
                }
                
                keyboard := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("⬅Back"),
                    ),
                )
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+nameOfLang)
                msg.ReplyMarkup = keyboard
                bot.Send(msg)
                
                return
            }
            
            if codeOfLang := iso6391.CodeForName(strings.Title(lowerUserMsg)); codeOfLang != "" { // Юзер полное его название языка
                err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("my_lang", codeOfLang).Error
                if err != nil {
                    warn(err)
                    return
                }
                keyboard := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("⬅Back")))
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(codeOfLang))
                msg.ReplyMarkup = keyboard
                bot.Send(msg)
                return
            }
            
            
            lang, err := translate.DetectLanguageGoogle(update.Message.Text)
            if err != nil || lang == "" {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Could not detect language, please send something else again"))
                return
            }
            err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("my_lang", lang).Error
            if err != nil {
                warn(err)
                return
            }
            keyboard := tgbotapi.NewReplyKeyboard(
                tgbotapi.NewKeyboardButtonRow(
                    tgbotapi.NewKeyboardButton("⬅Back")))
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(lang))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
            
            
            return
        case "set_translate_lang":
            lowerUserMsg := strings.ToLower(update.Message.Text)
            
            if nameOfLang := iso6391.Name(lowerUserMsg); nameOfLang != "" { // Юзер отправил сразу код языка
                err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("my_lang", lowerUserMsg).Error
                if err != nil {
                    warn(err)
                    return
                }
                keyboard := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("⬅Back")))
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+nameOfLang)
                msg.ReplyMarkup = keyboard
                bot.Send(msg)
                
                return
            }
            
            if codeOfLang := iso6391.CodeForName(strings.Title(lowerUserMsg)); codeOfLang != "" { // Юзер полное его название языка
                err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("my_lang", codeOfLang).Error
                if err != nil {
                    warn(err)
                    return
                }
                keyboard := tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton("⬅Back")))
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(codeOfLang))
                msg.ReplyMarkup = keyboard
                bot.Send(msg)
                
                return
            }
            
            lang, err := translate.DetectLanguageGoogle(update.Message.Text)
            if err != nil || lang == "" {
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Could not detect language, please send something else again")
                bot.Send(msg)
                
                return
            }
            err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("my_lang", lang).Error
            if err != nil {
                warn(err)
                return
            }
            keyboard := tgbotapi.NewReplyKeyboard(
                tgbotapi.NewKeyboardButtonRow(
                    tgbotapi.NewKeyboardButton("⬅Back")))
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(lang))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
        default: // У пользователя нет шага и сообщение не команда
            var user Users // Contains only MyLang and ToLang
            err = db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&user).Error
            if err != nil {
                warn(err)
                return
            }
            msg, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Translating..."))
            if err != nil {
                return
            }
            
            var text = update.Message.Text
            if update.Message.Caption != "" {
                text = update.Message.Caption
            }
            if text == "" {
                bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Please, send text message"))
                return
            }
            
            cutText := cutString(text, 500)
            lang, err := translate.DetectLanguageGoogle(cutText)
            if err != nil {
                
                return
            }
            
            if lang == "" {
                bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, text))
                return
            }
            
            var to string // language into need to translate
            if lang == user.ToLang {
                to = user.MyLang
            } else {
                to = user.ToLang
            }
            
            tr, err := translate.GoogleTranslate(lang, to, text)
            if err != nil {
                if e, ok := err.(translate.HTTPError); ok {
                    if e.Code == 413 {
                        bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Too big text"))
                    } else if e.Code >= 500 {
                        bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Unsupported language or internal error"))
                    } else {
                        warn(e)
                    }
                    return
                }
                warn(err)
                return
            }
            if tr.Text == "" {
                answer := tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Empty result")
                bot.Send(answer)
                
                return
            }
            pp.Println(tr)
            
            _, err = bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, tr.Text))
            if err != nil {
                pp.Println(err)
            }
        }
        
    }
    
}