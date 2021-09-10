package main

import (
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "os"
    "strconv"
    "strings"
)

func handleMessage(message *tgbotapi.Message) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", message.From.LanguageCode)))
        WarnAdmin(err)
    }
    analytics.User(message.Text, message.From)
    
    if message.Chat.ID < 0 {
        return
    }

    user := NewUser(message.Chat.ID, warn)
    if !user.Exists() {
        if message.From.LanguageCode == "" {
            message.From.LanguageCode = "en"
        }
        var referrerID int64
        if strings.HasPrefix(message.Text, "/start ") {
            fields := strings.Fields(message.Text)
            if len(fields) == 2 {
                var err error
                rID, err := strconv.ParseInt(fields[1], 10, 32)
                if err == nil {
                    referrer := NewUser(referrerID, warn)
                    if referrer.Exists() {
                        referrerID = rID
                    }
                }
            }
        }
        user.Create(Users{
            ID:      message.Chat.ID,
            MyLang:  "en",
            ToLang:  "en",
            Lang:    message.From.LanguageCode,
            ReferrerID: referrerID,
        })

        // Приветственное сообщение
        langs := map[string]string{"en": "🇬🇧 English", "it": "🇮🇹 Italiano", "uz":"🇺🇿 O'zbek tili", "de":"🇩🇪 Deutsch", "ru":"🇷🇺 Русский", "es":"🇪🇸 Español", "uk":"🇺🇦 Український", "pt":"🇵🇹 Português", "id":"🇮🇩 Indonesia"}
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        var i int
        for code, name := range langs {
            if i % 3 == 0 {
                keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(name, "register_bot_lang:"+code)))
            } else {
                l := len(keyboard.InlineKeyboard)-1
                keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(name, "register_bot_lang:"+code))
            }
            i++
        }
        msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Please, select bot language"))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)

        analytics.Bot(message.Chat.ID, msg.Text, "New user, register")

        return
    } else {
        user.Fill()
    }

    if strings.HasPrefix(message.Text, "/start") {
        SendMenu(user)
        if err := setUserStep(message.Chat.ID, ""); err != nil {
            warn(err)
            return
        }
        return
    }
    
    switch message.Text {
    case "/bot_lang":
        langs := map[string]string{"en": "🇬🇧 English", "it": "🇮🇹 Italiano", "uz":"🇺🇿 O'zbek tili", "de":"🇩🇪 Deutsch", "ru":"🇷🇺 Русский", "es":"🇪🇸 Español", "uk":"🇺🇦 Український", "pt":"🇵🇹 Português", "id":"🇮🇩 Indonesia"}
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        var i int
        for code, name := range langs {
            if i % 2 == 0 {
                keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(name, "set_bot_lang:"+code)))
            } else {
                l := len(keyboard.InlineKeyboard)-1
                keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(name, "set_bot_lang:"+code))
            }
            i++
        }
        msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Please, select bot language"))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    case "My Language", "/my_lang", "Мой Язык","Mi Idioma","Моя Мова","A Minha Língua","Bahasa Saya","La mia lingua","Tilimni","Meine Sprache":
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for i, code := range codes {
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
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("◀", "set_my_lang_pagination:0"),
            tgbotapi.NewInlineKeyboardButtonData("0/"+strconv.Itoa(len(langs)), "none"),
            tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:10")))
        msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Ваш язык %s. Выберите Ваш язык.", iso6391.Name(user.MyLang)))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    
        analytics.Bot(message.Chat.ID, msg.Text, "Set my lang")
    case "Translate Language", "/to_lang", "Sprache zum Übersetzen","Idioma para traducir","Bahasa untuk menerjemahkan","Lingua per tradurre","Língua para tradução","Язык перевода","Мова перекладу","Tarjima qilish uchun til":
        keyboard := tgbotapi.NewInlineKeyboardMarkup()
        for i, code := range codes {
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
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("◀", "set_translate_lang_pagination:0"),
        tgbotapi.NewInlineKeyboardButtonData("0/"+strconv.Itoa(len(langs)), "none"),
            tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:10")))
        msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Сейчас бот переводит на %s. Выберите язык для перевода", iso6391.Name(user.ToLang)))
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
    
        analytics.Bot(message.Chat.ID, msg.Text, "Set translate lang")
    case "💡 Instruction", "/help", "💡 Инструкция", "💡 Instrucción","💡 Інструкція","💡 Instrucao","💡 Instruksi","💡 Istruzione","💡 Yo'riqnoma","💡 Anweisung":
        SendHelp(user)
    case "/sponsorship":
        msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("sponsorship"))
        keyboard := tgbotapi.NewReplyKeyboard(
            tgbotapi.NewKeyboardButtonRow(
                tgbotapi.NewKeyboardButton(user.Localize("⬅Back"))))
        msg.ReplyMarkup = keyboard
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
        
        if err := setUserStep(message.Chat.ID, "sponsorship_set_text"); err != nil {
            warn(err)
        }
        analytics.Bot(message.Chat.ID, msg.Text, "Look at sponsorship")
    case "/stats":
        var users int
        err := db.Model(&Users{}).Raw("SELECT COUNT(*) FROM users").Find(&users).Error
        if err != nil {
            warn(err)
            return
        }
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Всего " + strconv.Itoa(users) + " юзеров"))
    case "/users":
        if message.From.ID != AdminID {
            return
        }
        f, err := os.Create("users.txt")
        if err != nil {
            warn(err)
            return
        }
        var users []Users
        if err = db.Model(&Users{}).Find(&users).Error; err != nil {
            warn(err)
            return
        }
        for _, user := range users {
            if _, err = f.WriteString(strconv.FormatInt(user.ID, 10) + "\r\n"); err != nil {
                warn(err)
                return
            }
        }
        doc := tgbotapi.NewInputMediaDocument("users.txt")
        group := tgbotapi.NewMediaGroup(message.From.ID, []interface{}{doc})
        if _, err = bot.Send(group); err != nil {
            warn(err)
        }
    default: // Сообщение не является командой.
    
        userStep, err := getUserStep(message.Chat.ID)
        if err != nil {
            warn(err)
            return
        }
        switch userStep {
        case "sponsorship_set_text":
            if len(message.Text) > 130 {
                bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Too big text")))
                return
            }
            var sponsorshipExists bool
            err = db.Model(&SponsorshipsOffers{}).Raw("SELECT EXISTS(SELECT id FROM sponsorships_offers WHERE id=?)", message.From.ID).Find(&sponsorshipExists).Error
            if err != nil {
                warn(err)
                return
            }
            if sponsorshipExists {
                if err = db.Model(&SponsorshipsOffers{}).Where("id = ?", message.From.ID).Update("text", message.Text).Error; err != nil {
                    warn(err)
                    return
                }
            } else {
                if err = db.Create(&SponsorshipsOffers{
                    ID:      message.Chat.ID,
                    Text:    message.Text,
                }).Error; err != nil {
                    warn(err)
                    return
                }
            }
            msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("sponsorship_set_days"))
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
            
            if err = setUserStep(message.Chat.ID, ""); err != nil {
                warn(err)
            }

        case "sponsorship_set_days":
            days, err := strconv.Atoi(message.Text)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Введите целое число без лишних символов")))
                return
            }
            if days < 1 || days > 30 {
                bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Количество дней должно быть от 1 до 30 включительно")))
                return
            }
            if err = db.Model(&SponsorshipsOffers{}).Where("id = ?", message.Chat.ID).Updates(&SponsorshipsOffers{
                Days:    days,
            }).Error; err != nil {
                warn(err)
                return
            }
            msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Выберите языки пользователей, которые получат вашу рассылку."))
            langs := map[string]string{"en": "🇬🇧 English", "it": "🇮🇹 Italiano", "uz":"🇺🇿 O'zbek tili", "de":"🇩🇪 Deutsch", "ru":"🇷🇺 Русский", "es":"🇪🇸 Español", "uk":"🇺🇦 Український", "pt":"🇵🇹 Português", "id":"🇮🇩 Indonesia"}
            keyboard := tgbotapi.NewInlineKeyboardMarkup()
            for code, name := range langs {
                keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(name, "country:"+code)))
            }
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(user.Localize("Далее"), "sponsorship_pay")))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
            if err = setUserStep(message.Chat.ID, "sponsorship_set_langs"); err != nil {
                warn(err)
            }
        default: // У пользователя нет шага и сообщение не команда

            if user.Usings == 10 || user.Usings % 20 == 0 {
                defer func() {
                    photo := tgbotapi.NewPhoto(message.Chat.ID, "ad.jpg")
                    photo.Caption = user.Localize("bot_advertise")
                    if _, err = bot.Send(photo); err != nil {
                        pp.Println(err)
                    }
                }()
            }
            
            
            msg, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("⏳ Translating...")))
            if err != nil {
                return
            }
            
            var text = message.Text
            if message.Caption != "" {
                text = message.Caption
            }

            if len(message.Entities) > 0 {
                text = applyEntitiesHtml(text, message.Entities)
            } else if len(message.CaptionEntities) > 0 {
                text = applyEntitiesHtml(text, message.CaptionEntities)
            }


            if text == "" {
                bot.Send(tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, user.Localize("Please, send text message")))
                analytics.Bot(message.Chat.ID, msg.Text, "Message is not text message")
                return
            }
            
            from, err := translate.DetectLanguageGoogle(cutString(text, 100))
            if err != nil {
                warn(err)
                return
            }

            if from == "" {
                bot.Send(tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, text))
                return
            }
            
            var to string // language into need to translate
            if from == user.ToLang {
                to = user.MyLang
            } else if from == user.MyLang {
                to = user.ToLang
            } else { // никакой из
                to = user.MyLang
            }
            
            tr, err := translate.GoogleHTMLTranslate(from, to, text)
            if err != nil {
                bot.Send(tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, user.Localize("Unsupported language or internal error")))
                warn(err)
                return
            }

            if tr.Text == "" {
                answer := tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, user.Localize("Empty result"))
                bot.Send(answer)
                return
            }

            tr.Text = strings.Replace(tr.Text, `<label class="notranslate">`, "", -1)
            tr.Text = strings.Replace(tr.Text, `</label>`, "", -1)
            tr.Text = strings.Replace(tr.Text, `<br>`, "\n", -1)


            otherLanguagesButton := tgbotapi.InlineKeyboardButton{
                Text:                         user.Localize("Другие языки"),
                SwitchInlineQueryCurrentChat: &message.Text,
            }
            keyboard := tgbotapi.NewInlineKeyboardMarkup(
                tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData(user.Localize("To voice"), "speech:"+to)),
                tgbotapi.NewInlineKeyboardRow(otherLanguagesButton),
            )
    
            edit := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, msg.MessageID, tr.Text, keyboard)
            edit.ParseMode = tgbotapi.ModeHTML
            edit.DisableWebPagePreview = true
            
            var sponsorship Sponsorships
            err = db.Model(&Sponsorships{}).Select("text", "to_langs").Where("start <= current_timestamp AND finish >= current_timestamp").Limit(1).Find(&sponsorship).Error
            if err != nil {
                WarnAdmin(err)
            } else { // no error
                langs := strings.Split(sponsorship.ToLangs, ",")
                if inArray(user.Lang, langs) {
                    edit.Text += "\n⚡️" + sponsorship.Text
                }
            }

            if _, err = bot.Send(edit); err != nil {
                pp.Println(err)
            }
            
            analytics.Bot(message.Chat.ID, tr.Text, "Translated")
            
            if err = db.Exec("UPDATE users SET usings=usings+1 WHERE id=?", message.Chat.ID).Error; err != nil {
                WarnAdmin(err)
            }
        }
        
    }
    
}