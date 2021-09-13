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
    "time"
    "unicode/utf16"
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

    if strings.HasPrefix(message.Text, "/start") || message.Text == user.Localize("⬅Back") {
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
    case "/my_lang", user.Localize("My Language"):
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
    case "/to_lang", user.Localize("Translate Language"):
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
    case "/help":
        SendHelp(user)
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
        bot.Send(group)
    case "/ad":
        if message.From.ID != AdminID {
            return
        }
        bot.Send(tgbotapi.NewMessage(message.From.ID, "Введите текст сообщения: (до 100 символов)\n\n/cancel для отмены, можно вызвать в любом месте"))
        user.SetStep("ad:accept_text")
    case "/cancel":
        if message.From.ID != AdminID {
            return
        }
        if err := db.Model(&AdsOffers{}).Where("id = ?", message.From.ID).Delete(&AdsOffers{}).Error; err != nil {
            warn(err)
            return
        }

        user.SetStep("")

        bot.Send(tgbotapi.NewMessage(message.From.ID, "Создание рекламы отменено."))
        SendMenu(user)
    case "/id":
        bot.Send(tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10)))
    default: // Сообщение не является командой.
    
        userStep, err := getUserStep(message.Chat.ID)
        if err != nil {
            warn(err)
            return
        }
        switch userStep {
        case "ad:accept_text":
            if len(utf16.Encode([]rune(message.Text))) > 100 {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Слишком длинный текст. Максимум 100 символов."))
            }
            message.Text = applyEntitiesHtml(message.Text, message.Entities)
            msg := tgbotapi.NewMessage(message.From.ID, "Вы собираетесь создать рекламу с текстом:\n\n"+message.Text)
            msg.ParseMode = tgbotapi.ModeHTML
            keyboard := tgbotapi.NewInlineKeyboardMarkup(
                tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData("Да", "ad:confirm_entered_text_for_ad")))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
        case "ad:pass_start_date":
            t, err := time.ParseInLocation(layout, message.Text, loc)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Неверный ввод"))
                return
            }
            if !now.Before(t) {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Время должно быть позже, чем сейчас"))
                return
            }
            if err = db.Model(&AdsOffers{}).Where("id = ?", message.From.ID).Update("start_date", t).Error; err != nil {
                warn(err)
                return
            }

            bot.Send(tgbotapi.NewMessage(message.From.ID, "Теперь введите дату окончания в таком же формате"))
            user.SetStep("ad:pass_finish_date")
        case "ad:pass_finish_date":
            t, err := time.ParseInLocation(layout, message.Text, loc)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Неверный ввод"))
                return
            }
            if !now.Before(t) {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Время должно быть позже, чем сейчас"))
                return
            }
            var offer AdsOffers
            if err = db.Model(&AdsOffers{}).Where("id = ?", message.From.ID).Find(&offer).Error; err != nil {
                warn(err)
                return
            }
            if !offer.StartDate.Before(t) { // время начала позже времени окончания
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Время окончания должно быть позже времени начала"))
                return
            }
            if err = db.Model(&AdsOffers{}).Where("id = ?", message.From.ID).Update("finish_date", t).Error; err != nil {
                warn(err)
                return
            }

            bot.Send(tgbotapi.NewMessage(message.From.ID, "Отправьте id человека, чья создается реклама. id можно узнать командой /id"))
            user.SetStep("ad:pass_ad_owner_id")
        case "ad:pass_ad_owner_id":
            id, err := strconv.ParseInt(message.Text, 10, 32)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Неверный ввод. Отправьте цифры"))
                return
            }
            if id < 0 {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "id > 0"))
            }
            owner := NewUser(id, warn)
            if !owner.Exists() {
                bot.Send(tgbotapi.NewMessage(message.From.ID, "Неверный id, человека нет в базе"))
                return
            }

            if err = db.Model(&AdsOffers{}).Where("id = ?", message.From.ID).Update("id_whose_ad", id).Error; err != nil {
                warn(err)
                return
            }
            var offer AdsOffers
            if err = db.Model(&AdsOffers{}).Where("id = ?", message.From.ID).Find(&offer).Error; err != nil {
                warn(err)
            }
            langs := strings.Split(offer.ToLangs, ",")
            recipients := make([]string, 0, len(langs))
            for _, lang := range langs{
                recipients = append(recipients, iso6391.Name(lang))
            }
            msg := tgbotapi.NewMessage(message.From.ID, "Подтвердите данные:\n\n<b>Текст:</b> " + offer.Content + "\n\n<b>Начало:</b> " + offer.StartDate.Format("2006-01-02 15:01") + "\n\n<b>Конец:</b>: "+ offer.FinishDate.Format("2006-01-02 15:01") + "\n\n<b>Языки получателей:</b> " + strings.Join(recipients, ", ") + `\n\n<b>Владелец рекламы:</b> <a href="tg://user?id=` + strconv.FormatInt(offer.IDWhoseAd, 10) + `">` + "\n\nВсё так? Если нет, нажмите /cancel")
            keyboard := tgbotapi.NewInlineKeyboardMarkup(
                tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.NewInlineKeyboardButtonData("Подтверждаю", "ad:confirm_pass")))
            msg.ParseMode = tgbotapi.ModeHTML
            msg.ReplyMarkup = keyboard
            bot.Send(msg)
            user.SetStep("")
        //case "sponsorship_set_text":
        //    if len(message.Text) > 130 {
        //        bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Too big text")))
        //        return
        //    }
        //    var sponsorshipExists bool
        //    err = db.Model(&SponsorshipsOffers{}).Raw("SELECT EXISTS(SELECT id FROM sponsorships_offers WHERE id=?)", message.From.ID).Find(&sponsorshipExists).Error
        //    if err != nil {
        //        warn(err)
        //        return
        //    }
        //    if sponsorshipExists {
        //        if err = db.Model(&SponsorshipsOffers{}).Where("id = ?", message.From.ID).Update("text", message.Text).Error; err != nil {
        //            warn(err)
        //            return
        //        }
        //    } else {
        //        if err = db.Create(&SponsorshipsOffers{
        //            ID:      message.Chat.ID,
        //            Text:    message.Text,
        //        }).Error; err != nil {
        //            warn(err)
        //            return
        //        }
        //    }
        //    msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("sponsorship_set_days"))
        //    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        //        tgbotapi.NewInlineKeyboardRow(
        //            tgbotapi.NewInlineKeyboardButtonData("1д - 9р", "sponsorship_set_days:1"),
        //            tgbotapi.NewInlineKeyboardButtonData("2д - 20р", "sponsorship_set_days:2"),
        //            tgbotapi.NewInlineKeyboardButtonData("7д - 60р", "sponsorship_set_days:7")),
        //        tgbotapi.NewInlineKeyboardRow(
        //            tgbotapi.NewInlineKeyboardButtonData("10д - 90р", "sponsorship_set_days:10"),
        //            tgbotapi.NewInlineKeyboardButtonData("15д -130р", "sponsorship_set_days:15"),
        //            tgbotapi.NewInlineKeyboardButtonData("🔥 30д - 270р", "sponsorship_set_days:30"),
        //        ))
        //    msg.ReplyMarkup = keyboard
        //    bot.Send(msg)
        //
        //    if err = setUserStep(message.Chat.ID, ""); err != nil {
        //        warn(err)
        //    }
        //
        //case "sponsorship_set_days":
        //    days, err := strconv.Atoi(message.Text)
        //    if err != nil {
        //        bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Введите целое число без лишних символов")))
        //        return
        //    }
        //    if days < 1 || days > 30 {
        //        bot.Send(tgbotapi.NewMessage(message.Chat.ID, user.Localize("Количество дней должно быть от 1 до 30 включительно")))
        //        return
        //    }
        //    if err = db.Model(&SponsorshipsOffers{}).Where("id = ?", message.Chat.ID).Updates(&SponsorshipsOffers{
        //        Days:    days,
        //    }).Error; err != nil {
        //        warn(err)
        //        return
        //    }
        //    msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Выберите языки пользователей, которые получат вашу рассылку."))
        //    langs := map[string]string{"en": "🇬🇧 English", "it": "🇮🇹 Italiano", "uz":"🇺🇿 O'zbek tili", "de":"🇩🇪 Deutsch", "ru":"🇷🇺 Русский", "es":"🇪🇸 Español", "uk":"🇺🇦 Український", "pt":"🇵🇹 Português", "id":"🇮🇩 Indonesia"}
        //    keyboard := tgbotapi.NewInlineKeyboardMarkup()
        //    for code, name := range langs {
        //        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(name, "country:"+code)))
        //    }
        //    keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(user.Localize("Далее"), "sponsorship_pay")))
        //    msg.ReplyMarkup = keyboard
        //    bot.Send(msg)
        //    if err = setUserStep(message.Chat.ID, "sponsorship_set_langs"); err != nil {
        //        warn(err)
        //    }
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
                from = "auto"
                //bot.Send(tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, text))
                //return
            }

            var to string // language into need to translate
            if from == user.ToLang {
                to = user.MyLang
            } else if from == user.MyLang {
                to = user.ToLang
            } else { // никакой из
                to = user.MyLang
            }

            if len(message.Entities) > 0 {
                text = applyEntitiesHtml(text, message.Entities)
            } else if len(message.CaptionEntities) > 0 {
                text = applyEntitiesHtml(text, message.CaptionEntities)
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

            tr.Text = strings.NewReplacer(`<label class="notranslate">`, ``, `</label>`, ``).Replace(tr.Text)
            tr.Text = strings.ReplaceAll(tr.Text, `<br>`, "\n")

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