package main

import (
    "database/sql"
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "github.com/sirupsen/logrus"
    "os"
    "strconv"
    "strings"
    "sync"
)

func handleMessage(message tgbotapi.Message) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewMessage(message.Chat.ID, localize("Sorry, error caused.\n\nPlease, don't block the bot, I'll fix the bug in near future, the administrator has already been warned about this error ;)", message.From.LanguageCode)))
        WarnAdmin(err)
        logrus.Error(err)
    }
    analytics.User(message.Text, message.From)
    
    if message.Chat.ID < 0 {
        return
    }

    var user = NewUser(message.From.ID, warn)
    if message.Text == "/start" {
        user := NewUser(message.From.ID, warn)
        if !user.Exists() {
            if message.From.LanguageCode == "" {
                message.From.LanguageCode = "en"
            }
            err := db.Create(&Users{
                ID:         message.From.ID,
                MyLang:     message.From.LanguageCode,
                ToLang:     "en",
                Act:        sql.NullString{},
                Mailing:    true,
                Lang:       message.From.LanguageCode,
            }).Error
            if err != nil {
                warn(err)
            }
        }

        user.SendStart()
        return
    }
    user.Fill()


    if low := strings.ToLower(message.Text); low != "" {
        switch {
        case in(command("my language"), low):
            keyboard := tgbotapi.NewInlineKeyboardMarkup()
            for i, code := range codes {
                if i >= 20 {
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
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:20")))
            msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Ваш язык %s. Выберите Ваш язык.", iso6391.Name(user.MyLang)))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)

            analytics.Bot(message.Chat.ID, msg.Text, "Set my lang")
            return
        case in(command("translate language"), low):
            keyboard := tgbotapi.NewInlineKeyboardMarkup()
            for i, code := range codes {
                if i >= 20 {
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
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:20")))
            msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Сейчас бот переводит на %s. Выберите язык для перевода", iso6391.Name(user.ToLang)))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)

            analytics.Bot(message.Chat.ID, msg.Text, "Set translate lang")
            return
        }
    }

    switch message.Text {
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
    case "/id":
        bot.Send(tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10)))
    default: // Сообщение не является командой.
        if user.MyLang == user.ToLang {
            bot.Send(tgbotapi.NewMessage(message.From.ID, user.Localize("The original text language and the target language are the same, please set different")))
            return
        }

        if user.Usings == 20 || user.Usings == 50 || user.Usings == 100 || (user.Usings > 100 && user.Usings % 50 == 0) {
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

        var (
            tr translate.GoogleHTMLTranslation
            single translate.GoogleTranslateSingleResult
            samples translate.ReversoQueryResponse
            wg sync.WaitGroup
            errs = make(chan error, 2)
        )

        // Переводим в гугле, как обычно
        wg.Add(1)
        go func() {
            defer wg.Done()
            tr, err = translate.GoogleHTMLTranslate(from, to, text)
            if err != nil {
                errs <- err
            }
            if tr.Text == "" && text != "" {
                WarnAdmin("короче на " + to + " не переводит")
                errs <- err
                return
            }

            tr.Text = strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "").Replace(tr.Text)
            tr.Text = strings.ReplaceAll(tr.Text, `<br>`, "\n")

        }()

        wg.Add(1)
        go func() {
            defer wg.Done()
            single, err = translate.GoogleTranslateSingle(from, to, text)
            if err != nil {
                errs <- err
            }
        }()

        wg.Wait()

        if len(errs) > 0 {
            warn(<-errs)
            return
        }

        close(errs)

        samples, err = translate.ReversoQueryService(text, from, tr.Text, to)
        if err != nil {
            warn(err)
            return
        }

        otherLanguagesButton := tgbotapi.InlineKeyboardButton{
            Text:                         "🔀",
            SwitchInlineQueryCurrentChat: &text,
        }
        From := langs[from]
        keyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("From " + From.Emoji + " " + From.Name, "none")),
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.NewInlineKeyboardButtonData("🎙",  "speech_this_message_and_replied_one:"+from+":"+to),
                otherLanguagesButton,
            ),
        )

        if len(single.Dict) > 0 {
            l := len(keyboard.InlineKeyboard) - 1
            keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l],
                tgbotapi.NewInlineKeyboardButtonData("📚", "dictionary:"+from+":"+to))
        }

        if len(samples.Suggestions) > 0 || len(samples.List) > 0 {
            l := len(keyboard.InlineKeyboard) - 1
            keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l],
                tgbotapi.NewInlineKeyboardButtonData("💬", "examples:"+from+":"+to+":"+text))
        }

        edit := tgbotapi.NewEditMessageTextAndMarkup(message.From.ID, msg.MessageID, tr.Text, keyboard)
        edit.ParseMode = tgbotapi.ModeHTML
        bot.Send(edit)

        analytics.Bot(user.ID, tr.Text, "Translated")

        //user.WriteLog("pm_translate")
    }
}