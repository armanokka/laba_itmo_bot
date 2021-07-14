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
        bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", update.Message.From.LanguageCode)))
        WarnAdmin(err)
    }
    analytics.Error = warn
    analytics.User(update.Message.Text, update.Message.From)
    
    if update.Message.Chat.ID < 0 {
        return
    }
    
    var userExists bool
    err := db.Raw("SELECT EXISTS(SELECT id FROM users WHERE id=?)", update.Message.Chat.ID).Find(&userExists).Error
    if err != nil {
        warn(err)
        return
    }
    
    var UserLang string
    err = db.Model(&Users{}).Select("lang").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&UserLang).Error
    if err != nil {
        warn(err)
        return
    }
    
    if strings.HasPrefix(update.Message.Text, "/start") || update.Message.Text == "⬅Back" || update.Message.Text == "Let's check" {
        parts := strings.Fields(update.Message.Text)
        pp.Println(parts)
        if len(parts) == 2 && !userExists {
            pp.Println("here 1")
            var referrer bool // Check for exists
            err = db.Raw("SELECT EXISTS(SELECT id FROM referrers WHERE code=?)", strings.ToLower(parts[1])).Find(&referrer).Error
            if err != nil {
                WarnAdmin(err)
            }
            pp.Println("here 2")
            if referrer {
                pp.Println("here 3")
                err = db.Exec("UPDATE referrers SET users=users+1 WHERE code=?", strings.ToLower(parts[1])).Error
                if err != nil {
                    WarnAdmin(err)
                }
                pp.Println("here 4")
            }
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
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Your language is - <b>%s</b>, and the language for translation is - <b>%s</b>.", UserLang, user.MyLang, user.ToLang))
        keyboard := tgbotapi.NewReplyKeyboard(
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(Localize("💡Instruction", UserLang))),
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(Localize("My Language", UserLang))),
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(Localize("Translate Language", UserLang))))
        msg.ReplyMarkup = keyboard
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
    
        err = setUserStep(update.Message.Chat.ID, "")
        if err != nil {
            warn(err)
            return
        }
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Main menu")
        
        return
    }
    
    switch update.Message.Text {
    case "My Language", "/my_lang":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("To setup <b>your language</b>, do <b>one</b> of the following: 👇\n\nℹ️ Send <b>few words</b> in <b>your</b> language, for example: \"<code>Hi, how are you today?</code>\" - language will be English, or \"<code>L'amour ne fait pas d'erreurs</code>\" - language will be French, and so on.\nℹ️ Or send the <b>name</b> of your language <b>in English</b>, e.g. \"<code>Russian</code>\", or \"<code>Japanese</code>\", or  \"<code>Arabic</code>\", e.t.c.", UserLang))
        msg.ParseMode = tgbotapi.ModeHTML
        keyboard:= tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("⬅Back")))
        msg.ReplyMarkup = &keyboard
        bot.Send(msg)
        
        err := setUserStep(update.Message.Chat.ID, "set_my_lang")
        if err != nil {
            warn(err)
            return
        }
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Set my lang")
        
    case "Translate Language", "/to_lang":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("To setup <b>translate language</b>, do <b>one</b> of the following: 👇\n\nℹ️ Send <b>few words</b> in language <b>into you want to translate</b>, for example: \"<code>Hi, how are you today?</code>\" - language will be English, or \"<code>L'amour ne fait pas d'erreurs</code>\" - language will be French, and so on.\nℹ️ Or send the <b>name</b> of language <b>into you want to translate, in English</b>, e.g. \"<code>Russian</code>\", or \"<code>Japanese</code>\", or  \"<code>Arabic</code>\", e.t.c.", UserLang))
        msg.ParseMode = tgbotapi.ModeHTML
        keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("⬅Back")))
        msg.ReplyMarkup = &keyboard
        bot.Send(msg)
        
        err := setUserStep(update.Message.Chat.ID, "set_translate_lang")
        if err != nil {
            warn(err)
            return
        }
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Set translate lang")
        
    case "💡Instruction", "/help":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("<b>What can this bot do?</b>\n▫️ Translo allow you to translate messages into 182+ languages.\n<b>How to translate message?</b>\n▫️ Firstly, you have to setup your language, then setup translate language, next send text messages and bot will translate them quickly.\n<b>How to setup my language?</b>\n▫️ Click on the button below called \"My Language\"\n<b>How to setup language into I want to translate?</b>\n▫️ Click on the button below called \"Translate Language\"\n<b>Are there any other interesting things?\n</b>▫️ Yes, the bot support inline mode. Start typing the nickname @translobot in the message input field and then write there the text you want to translate.\n<b>I have a suggestion or I found bug!</b>\n▫️ 👉 Contact me pls - @armanokka", UserLang))
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Help")
        
    default: // Сообщение не является командой.
    
        userStep, err := getUserStep(update.Message.Chat.ID)
        if err != nil {
            warn(err)
            return
        }
        switch userStep {
        case "set_my_lang":
            name, code, err := DetectLang(update.Message.Text)
            if err != nil {
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Failed to detect the language. Please enter something else", UserLang))
                bot.Send(msg)
        
                analytics.Bot(update.Message.Chat.ID, msg.Text, "My language not detected")
                return
            }
    
            keyboard := tgbotapi.NewReplyKeyboard(
                tgbotapi.NewKeyboardButtonRow(
                    tgbotapi.NewKeyboardButton("⬅Back")))
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Now your language is %s\n\nPress \"⬅Back\" to exit to menu", UserLang, name))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
    
            err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("my_lang", code).Error
            if err != nil {
                warn(err)
                return
            }
            
            analytics.Bot(update.Message.Chat.ID, msg.Text, "My language detected default")
        case "set_translate_lang":
            name, code, err := DetectLang(update.Message.Text)
            if err != nil {
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Failed to detect the language. Please enter something else", UserLang))
                bot.Send(msg)
    
                analytics.Bot(update.Message.Chat.ID, msg.Text, "Translate language not detected")
                return
            }
            err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Update("to_lang", code).Error
            if err != nil {
                warn(err)
                return
            }

            
            keyboard := tgbotapi.NewReplyKeyboard(
                tgbotapi.NewKeyboardButtonRow(
                    tgbotapi.NewKeyboardButton("⬅Back")))
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Now translate language is %s\n\nPress \"⬅Back\" to exit to menu", UserLang, name))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
            
            analytics.Bot(update.Message.Chat.ID, msg.Text, "Translate language detected default")
        default: // У пользователя нет шага и сообщение не команда
            var user Users // Contains only MyLang and ToLang
            err = db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&user).Error
            if err != nil {
                warn(err)
                return
            }
            msg, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Localize("⏳ Translating...", UserLang)))
            if err != nil {
                return
            }
            
            var text = update.Message.Text
            if update.Message.Caption != "" {
                text = update.Message.Caption
            }
            if text == "" {
                bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, Localize("Please, send text message", UserLang)))
    
                analytics.Bot(update.Message.Chat.ID, msg.Text, "Message is not text message")
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
                        bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, Localize("Too big text", UserLang)))
                    } else if e.Code >= 500 {
                        bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, Localize("Unsupported language or internal error", UserLang)))
                    } else {
                        warn(e)
                    }
                    return
                }
                warn(err)
                return
            }
            if tr.Text == "" {
                answer := tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, Localize("Empty result", UserLang))
                bot.Send(answer)
                
                return
            }
            pp.Println(tr)
            
            _, err = bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, tr.Text))
            if err != nil {
                pp.Println(err)
            }
            
            // if len(tr.Images) > 0 {
            //     msg = tgbotapi.NEwph
            //     bot.Send()
            // }
            
            analytics.Bot(update.Message.Chat.ID, tr.Text, "Translated")
            
            
            err = db.Exec("UPDATE users SET usings=usings+1 WHERE id=?", update.Message.Chat.ID).Error
            if err != nil {
                WarnAdmin(err)
            }
            // var usings int
            // err = db.Model(&Users{}).Select("usings").Where("id = ?", update.Message.Chat.ID).Find(&usings).Error
            // if err != nil {
            //     WarnAdmin(err)
            //     return
            // }
            // if usings == 5 {
            //     msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please rate the bot 👇")
            //     keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
            //         ))
            // }
        }
        
    }
    
}