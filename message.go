package main

import (
    "database/sql"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "strconv"
    "strings"
)

func handleMessage(update *tgbotapi.Update) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", update.Message.From.LanguageCode)))
        WarnAdmin(err)
    }
    analytics.User(update.Message.Text, update.Message.From)
    
    if update.Message.Chat.ID < 0 {
        return
    }
    
    var UserLang string
    err := db.Model(&Users{}).Select("lang").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&UserLang).Error
    if err != nil {
        warn(err)
        return
    }
    if UserLang == "" {
        UserLang = "en"
    }
    
    if strings.HasPrefix(update.Message.Text, "/start") || inArray(update.Message.Text, []string{"⬅Back", "⬅️Zurück","⬅️Atrás","⬅️Kembali","⬅️Indietro","⬅️Back","⬅️Назад","⬅️Arka", "⬅Zurück","⬅Atrás","⬅Kembali","⬅Indietro","⬅Back","⬅Назад","⬅Назад","⬅Arka"}) {
        var userExists bool
        err = db.Raw("SELECT EXISTS(SELECT id FROM users WHERE id=?)", update.Message.Chat.ID).Find(&userExists).Error
        if err != nil {
            warn(err)
            return
        }
        
        parts := strings.Fields(update.Message.Text)
        if len(parts) == 2 && !userExists { // Рефка
            var referrerExists bool // Check for exists
            err = db.Raw("SELECT EXISTS(SELECT id FROM referrers WHERE code=?)", strings.ToLower(parts[1])).Find(&referrerExists).Error
            if err != nil {
                WarnAdmin(err)
            }
            if referrerExists {
                err = db.Exec("UPDATE referrers SET users=users+1 WHERE code=?", strings.ToLower(parts[1])).Error
                if err != nil {
                    WarnAdmin(err)
                }
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
            
            lang := update.Message.From.LanguageCode // язык бота
            if lang == "" {
                lang = "en"
            }
            
            err = db.Create(&Users{
                ID:     update.Message.Chat.ID,
                MyLang: fromLang,
                ToLang: translateLang,
                Lang: lang,
                Act:    sql.NullString{},
            }).Error
            if err != nil {
                warn(err)
                return
            }
        }
    
        var user Users
        err = db.Model(&Users{}).Select("my_lang", "to_lang", "lang").Where("id = ?", update.Message.Chat.ID).Find(&user).Error
        if err != nil {
            warn(err)
            return
        }
        if user.Lang == "" {
            user.Lang = "en"
        }
    
        user.MyLang = iso6391.Name(user.MyLang)
        user.ToLang = iso6391.Name(user.ToLang)
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("/start", user.Lang, user.MyLang, user.ToLang))
        keyboard := tgbotapi.NewReplyKeyboard(
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(Localize("🙎‍♂️Profile", user.Lang)),
                ),
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(Localize("💬 Bot language", user.Lang)),
                tgbotapi.NewKeyboardButton(Localize("💡 Instruction", user.Lang)),
                ),
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(Localize("My Language", user.Lang)),
                tgbotapi.NewKeyboardButton(Localize("Translate Language", user.Lang)),
                ),
            )
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
    case "🙎‍♂️Profile", "🙍‍♂️Profil", "🙍‍♂️Perfil", "🙍‍♂️Профиль", "🙍‍♂️Profilo", "🙍‍♂️Профіль":
        var user Users
        err = db.Model(&Users{}).Select("my_lang", "to_lang", "lang").Where("id = ?", update.Message.Chat.ID).Find(&user).Error
        if err != nil {
            warn(err)
            return
        }
        user.ToLang = iso6391.Name(user.ToLang)
        user.MyLang = iso6391.Name(user.MyLang)
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("/start", user.Lang, user.MyLang, user.ToLang))
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Profile")
    case "My Language", "/my_lang", "Мой Язык","Mi Idioma","Моя Мова","A Minha Língua","Bahasa Saya","La mia lingua","Tilimni","Meine Sprache":
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        var i int
        for code, lang := range langs {
            if i >= 10 {
                break
            }
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_my_lang_by_callback:"  + code)))
            i++
        }
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("10/"+strconv.Itoa(len(langs)), "none"),
            tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:10")))
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Выберите ваш родной язык", UserLang))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Set my lang")
    case "Translate Language", "/to_lang", "Sprache zum Übersetzen","Idioma para traducir","Bahasa untuk menerjemahkan","Lingua per tradurre","Língua para tradução","Язык перевода","Мова перекладу","Tarjima qilish uchun til":
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        var i int
        for code, lang := range langs {
            if i >= 10 {
                break
            }
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_translate_lang_by_callback:"  + code)))
            i++
        }
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("10/"+strconv.Itoa(len(langs)), "none"),
            tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:10")))
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Выберите язык, на котором хотите переводить текст", UserLang))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Set translate lang")
    case "💡 Instruction", "/help", "💡 Инструкция", "💡 Instrucción","💡 Інструкція","💡 Instrucao","💡 Instruksi","💡 Istruzione","💡 Yo'riqnoma","💡 Anweisung":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("/help", UserLang))
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Help")

    case "💬 Bot language", "💬 Bot-Sprache","💬 Lenguaje bot","💬 Linguagem de bot", "💬 Бот-мова", "💬 Bot tili", "💬 Bahasa bot", "💬 Linguaggio Bot", "💬 Язык бота":
        langs := map[string]string{"en": "🇬🇧 English", "it": "🇮🇹 Italiano", "uz":"🇺🇿 O'zbek tili", "de":"🇩🇪 Deutsch", "ru":"🇷🇺 Русский", "es":"🇪🇸 Español", "uk":"🇺🇦 Український", "pt":"🇵🇹 Português", "id":"🇮🇩 Indonesia"}
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for code, name := range langs {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(name, "set_bot_lang:"+code)))
        }
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Please, select bot language", UserLang))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Choose bot lang")
    case "/sponsorship":
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("sponsorship", UserLang))
        keyboard := tgbotapi.NewReplyKeyboard(
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(Localize("⬅Back", UserLang))))
        msg.ReplyMarkup = keyboard
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
        
        if err = setUserStep(update.Message.Chat.ID, "sponsorship_set_text"); err != nil {
            warn(err)
        }
        analytics.Bot(update.Message.Chat.ID, msg.Text, "Look at sponsorship")
    case "/stats":
        var users int
        err = db.Model(&Users{}).Raw("SELECT COUNT(*) FROM users").Find(&users).Error
        if err != nil {
            warn(err)
            return
        }
        bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Всего " + strconv.Itoa(users) + " юзеров"))

    default: // Сообщение не является командой.
    
        userStep, err := getUserStep(update.Message.Chat.ID)
        if err != nil {
            warn(err)
            return
        }
        switch userStep {
        case "sponsorship_set_text":
            if len(update.Message.Text) > 130 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Too big text", UserLang)))
                return
            }
            var sponsorshipExists bool
            err = db.Model(&SponsorshipsOffers{}).Raw("SELECT EXISTS(SELECT id FROM sponsorships_offers WHERE id=?)", update.Message.From.ID).Find(&sponsorshipExists).Error
            if err != nil {
                warn(err)
                return
            }
            if sponsorshipExists {
                if err = db.Model(&SponsorshipsOffers{}).Where("id = ?", update.Message.From.ID).Update("text", update.Message.Text).Error; err != nil {
                    warn(err)
                    return
                }
            } else {
                if err = db.Create(&SponsorshipsOffers{
                    ID:      update.Message.Chat.ID,
                    Text:    update.Message.Text,
                }).Error; err != nil {
                    warn(err)
                    return
                }
            }
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("sponsorship_set_days", UserLang))
            keyboard := tgbotapi.NewInlineKeyboardMarkup(
                tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData("1д - 9р", "sponsorship_set_days:1"),
                    tgbotapi.NewInlineKeyboardButtonData("2д - 20р", "sponsorship_set_days:2"),
                    tgbotapi.NewInlineKeyboardButtonData("7д - 60р", "sponsorship_set_days:7")),
                tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData("10д - 90р", "sponsorship_set_days:10"),
                    tgbotapi.NewInlineKeyboardButtonData("15д -130р", "sponsorship_set_days:15"),
                    tgbotapi.NewInlineKeyboardButtonData("🔥 30д - 270р", "sponsorship_set_days:30"),
                ))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
            
            if err = setUserStep(update.Message.Chat.ID, ""); err != nil {
                warn(err)
            }
        case "sponsorship_set_days":
            days, err := strconv.Atoi(update.Message.Text)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Введите целое число без лишних символов", UserLang)))
                return
            }
            if days < 1 || days > 30 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Количество дней должно быть от 1 до 30 включительно", UserLang)))
                return
            }
            if err = db.Model(&SponsorshipsOffers{}).Where("id = ?", update.Message.Chat.ID).Updates(&SponsorshipsOffers{
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
            if err = setUserStep(update.Message.Chat.ID, "sponsorship_set_langs"); err != nil {
                warn(err)
            }
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
                    tgbotapi.NewKeyboardButton(Localize("⬅Back", UserLang))))
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
            err = db.Model(&Users{}).Select("my_lang", "to_lang", "usings").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&user).Error
            if err != nil {
                warn(err)
                return
            }
            x := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("⏳ Translating...", UserLang))
            x.ReplyToMessageID = update.Message.MessageID
            msg, err := bot.Send(x)
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
            } else if lang == user.MyLang {
                to = user.ToLang
            } else { // никакой из
                to = user.MyLang
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
            keyboard := tgbotapi.NewInlineKeyboardMarkup(
                tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData(Localize("To voice", UserLang), "speech:"+to),
                    ),
                )
            if len(tr.Variants) > 0 {
                keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(Localize("Variants", UserLang), "variants:"+lang+":"+to)))
            }
    
            edit := tgbotapi.NewEditMessageTextAndMarkup(update.Message.Chat.ID, msg.MessageID, tr.Text, keyboard)
            edit.ParseMode = tgbotapi.ModeHTML
            edit.DisableWebPagePreview = true
            
            var sponsorship Sponsorships
            err = db.Model(&Sponsorships{}).Select("text", "to_langs").Where("start <= current_timestamp AND finish >= current_timestamp").Limit(1).Find(&sponsorship).Error
            if err != nil {
                WarnAdmin(err)
            } else { // no error
                langs := strings.Split(sponsorship.ToLangs, ",")
                if inArray(UserLang, langs) {
                    edit.Text += "\n⚡️" + sponsorship.Text
                }
            }

            _, err = bot.Send(edit)
            if err != nil {
                pp.Println(err)
            }
            
            // if len(tr.Images) > 0 {
            //     msg = tgbotapi.NEwph
            //     bot.Send()
            // }
            
            analytics.Bot(update.Message.Chat.ID, tr.Text, "Translated")

            // if user.Usings == 10 || user.Usings % 20 == 0 && user.Usings != 20 {
            //     msg := tgbotapi.NewMessage(update.Message.Chat.ID, Localize("Rate", UserLang))
            //     kb := tgbotapi.NewInlineKeyboardMarkup(
            //         tgbotapi.NewInlineKeyboardRow(
            //             tgbotapi.NewInlineKeyboardButtonData("👍", "rate")),
            //         tgbotapi.NewInlineKeyboardRow(
            //             tgbotapi.NewInlineKeyboardButtonData("👎", "delete")))
            //     msg.ReplyMarkup = kb
            //     bot.Send(msg)
            // }
            
            
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