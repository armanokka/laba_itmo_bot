package main

import (
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "github.com/sirupsen/logrus"
    "os"
    "strconv"
    "strings"
)

func handleMessage(message *tgbotapi.Message) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", message.From.LanguageCode)))
        WarnAdmin(err)
        logrus.Error(err)
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
            tgbotapi.NewInlineKeyboardButtonData("0/"+strconv.Itoa(len(codes)), "none"),
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
        tgbotapi.NewInlineKeyboardButtonData("0/"+strconv.Itoa(len(codes)), "none"),
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

        defer func() {
            err := db.Model(&Users{}).Exec("UPDATE users SET usings=usings+1 WHERE id=?", message.From.ID).Error
            if err != nil {
                WarnAdmin(err)
            }
        }()

        if user.Usings % 20 == 0 {
            defer func() {
                photo := tgbotapi.NewPhoto(message.Chat.ID, "logo.jpg")
                photo.Caption = user.Localize("bot_advertise")
                photo.ParseMode = tgbotapi.ModeHTML
                if _, err := bot.Send(photo); err != nil {
                    pp.Println(err)
                }
            }()
        }

        msg, err := bot.Send(tgbotapi.MessageConfig{
            BaseChat:              tgbotapi.BaseChat{
                ChatID:                   message.Chat.ID,
                ReplyToMessageID: message.MessageID, // very important to "dictionary" callback
            },
            Text:                  user.Localize("⏳ Translating..."),
            ParseMode:             "",
            Entities:              nil,
            DisableWebPagePreview: false,
        })
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

        from, err := translate.DetectLanguageGoogle(cutStringUTF16(text, 100))
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
        text = strings.ReplaceAll(text, "\n", "<br>")

        tr, err := translate.GoogleHTMLTranslate(from, to, text)
        if err != nil {
            bot.Send(tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, user.Localize("Unsupported language or internal error")))
            warn(err)
            return
        }

        if tr.Text == "" {
            answer := tgbotapi.NewEditMessageText(message.Chat.ID, msg.MessageID, user.Localize("%s language is not supported", iso6391.Name(to)))
            WarnAdmin("короче на " + to + " не переводит")
            bot.Send(answer)
            return
        }

        tr.Text = strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "").Replace(tr.Text)
        tr.Text = strings.ReplaceAll(tr.Text, `<br>`, "\n")

        otherLanguagesButton := tgbotapi.InlineKeyboardButton{
            Text:                         user.Localize("Другие языки"),
            SwitchInlineQueryCurrentChat: &message.Text,
        }
        keyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData(user.Localize("To voice"), "speech_this_message_and_replied_one:"+from+":"+to)),
            tgbotapi.NewInlineKeyboardRow(otherLanguagesButton),
        )
        if inMapValues(translate.ReversoSupportedLangs(), from, to) && from != to {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData(user.Localize("Dictionary"), "dictionary:"+from+":"+to)))
        }
        edit := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, msg.MessageID, langs[from].Name + "->" + langs[to].Name + ": " +tr.Text, keyboard)
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