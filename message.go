package main

import (
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    "github.com/go-errors/errors"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "github.com/sirupsen/logrus"
    "html"
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
    if !user.Exists() {
        if message.From.LanguageCode == "" || !in(BotLocalizedLangs, message.From.LanguageCode) {
            message.From.LanguageCode = "en"
        }
        kb, err := BuildSupportedLanguagesKeyboard("register")
        if err  != nil {
            warn(err)
            return
        }
        bot.Send(tgbotapi.MessageConfig{
            BaseChat:              tgbotapi.BaseChat{
                ChatID:                   message.From.ID,
                ReplyMarkup:              kb,
            },
            Text:                  user.Localize("Choose bot language"),
        })
        return
    }
    user.Fill()

    if strings.HasPrefix(message.Text, "/start") {
        bot.Send(tgbotapi.MessageConfig{
            BaseChat:              tgbotapi.BaseChat{
                ChatID:                   message.From.ID,
                DisableNotification:      true,
                AllowSendingWithoutReply: true,
                ReplyMarkup: tgbotapi.NewReplyKeyboard(
                    tgbotapi.NewKeyboardButtonRow(
                        tgbotapi.NewKeyboardButton(user.Localize("My Language")),
                        tgbotapi.NewKeyboardButton(user.Localize("Translate Language")),
                    ),
                ),
            },
            Text:                  user.Localize("Just send me a text and I will translate it"),
        })

        return
    }
    user.Fill()
    defer user.UpdateLastActivity()

    user.WriteUserLog(message.Text)

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
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_my_lang_pagination:" + strconv.Itoa(LanguagesPaginationLimit))))
            msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Ваш язык %s. Выберите Ваш язык.", iso6391.Name(user.MyLang)))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)

            analytics.Bot(message.Chat.ID, msg.Text, "Set my lang")
            user.WriteBotLog("pm_to_lang", msg.Text)
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
                tgbotapi.NewInlineKeyboardButtonData("▶", "set_translate_lang_pagination:" + strconv.Itoa(LanguagesPaginationLimit))))
            msg := tgbotapi.NewMessage(message.Chat.ID, user.Localize("Сейчас бот переводит на %s. Выберите язык для перевода", iso6391.Name(user.ToLang)))
            msg.ReplyMarkup = keyboard
            bot.Send(msg)

            analytics.Bot(message.Chat.ID, msg.Text, "Set my lang")
            user.WriteBotLog("pm_to_lang", msg.Text)
            return
        }
    }

    switch message.Command() {
    case "get":
        arg := message.CommandArguments()
        var id int64
        if strings.HasPrefix(arg, "@") {
            user, err := bot.GetChat(tgbotapi.ChatInfoConfig{tgbotapi.ChatConfig{
                SuperGroupUsername: strings.TrimPrefix(arg, "@"),
            }})
            if err != nil {
                warn(err)
                return
            }
            id = user.ID
        } else {
            v, err := strconv.ParseInt(arg, 10, 32)
            if err != nil {
                warn(err)
            }
            id = v
        }
        var logs []UsersLogs
        if err := db.Model(&UsersLogs{}).Where("id = ?", id).Order("date DESC").Limit(20).Find(&logs).Error; err != nil {
            warn(err)
            return
        }
        text := ""
        for _, log := range logs {
            if log.FromBot {
                text += "\n<b>Bot:</b> <i>[" + log.Intent.String + "]</i> "
            } else {
                text += "\n<b>User:</b> "
            }
            text += log.Text
        }
        bot.Send(tgbotapi.MessageConfig{
            BaseChat:              tgbotapi.BaseChat{
                ChatID:                   message.From.ID,
            },
            ParseMode: tgbotapi.ModeHTML,
            Text:                  text,
        })
    case "stats":
        var users int
        err := db.Model(&Users{}).Raw("SELECT COUNT(*) FROM users").Find(&users).Error
        if err != nil {
            warn(err)
            return
        }
        var stats = make(map[string]string, 20)
        if err = db.Model(&UsersLogs{}).Raw("SELECT intent, COUNT(*) FROM users_logs GROUP BY intent ORDER BY count(*) DESC").Find(&stats).Error; err != nil {
            warn(err)
        }
        text := "Всего " + strconv.Itoa(users) + " юзеров"
        for name, count := range stats {
            text += "\n" + name + ": " + count
        }
        msg := tgbotapi.NewMessage(message.Chat.ID, text)
        bot.Send(msg)
        user.WriteBotLog("pm_stats", msg.Text)
        return
    case "users":
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
        user.WriteBotLog("pm_users", "{document was sended}")
        return
    case "id":
        msg := tgbotapi.NewMessage(message.From.ID, strconv.FormatInt(message.From.ID, 10))
        bot.Send(msg)
        user.WriteBotLog("pm_id", msg.Text)
        return
    }

    if user.MyLang == user.ToLang {
        bot.Send(tgbotapi.NewMessage(message.From.ID, user.Localize("The original text language and the target language are the same, please set different")))
        return
    }

    if user.Usings == 20 || user.Usings == 50 || user.Usings == 100 || (user.Usings > 100 && user.Usings % 50 == 0) {
        bot.Send(tgbotapi.MessageConfig{
            BaseChat:              tgbotapi.BaseChat{
                ChatID:                   message.From.ID,
                ReplyMarkup:              tgbotapi.NewInlineKeyboardMarkup(
                    tgbotapi.NewInlineKeyboardRow(
                        tgbotapi.NewInlineKeyboardButtonData("👍", "good_bot"),
                        tgbotapi.NewInlineKeyboardButtonData("👎", "delete"))),
                DisableNotification:      true,
                AllowSendingWithoutReply: true,
            },
            Text:                  user.Localize("Did you like the bot?"),
        })
    }

    msg, err := bot.Send(tgbotapi.MessageConfig{
        BaseChat: tgbotapi.BaseChat{
            ChatID:                   message.Chat.ID,
            ChannelUsername:          "",
            ReplyToMessageID:         message.MessageID, // very important to "dictionary" callback
            ReplyMarkup:              nil,
            DisableNotification:      true,
            AllowSendingWithoutReply: true,
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
    text = html.EscapeString(text)

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

    var (
        tr = translate.GoogleHTMLTranslation{}
        rev = translate.ReversoTranslation{}
        dict = translate.GoogleDictionaryResponse{}
        wg = sync.WaitGroup{}
        errs = make(chan *errors.Error, 4)
    )

    wg.Add(1)
    go func() {
        defer wg.Done()
        dict, err = translate.GoogleDictionary(from, text)
        if err != nil {
            errs <- errors.New(err)
        }
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        if inMapValues(translate.ReversoSupportedLangs(), from, to) && from != to {
            rev, err = translate.ReversoTranslate(translate.ReversoIso6392(from), translate.ReversoIso6392(to), strings.ToLower(text))
            if err != nil {
                errs <- errors.New(err)
            }
        }
    }()

    if len(message.Entities) > 0 {
        text = applyEntitiesHtml(text, message.Entities)
    } else if len(message.CaptionEntities) > 0 {
        text = applyEntitiesHtml(text, message.CaptionEntities)
    }
    text = strings.ReplaceAll(text, "\n", "<br>")

    // Переводим в гугле, как обычно
    wg.Add(1)
    go func() {
        defer wg.Done()
        tr, err = translate.GoogleHTMLTranslate(from, to, text)
        if err != nil {
            errs <- errors.New(err)
        }
        if tr.Text == "" && text != "" {
            WarnAdmin("короче на " + to + " не переводит")
            errs <- errors.New("короче на " + to + " не переводит")
            return
        }

        tr.Text = strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "").Replace(tr.Text)
        tr.Text = strings.ReplaceAll(tr.Text, `<br>`, "\n")
    }()

    wg.Wait()
    close(errs)

    if len(errs) > 0 {
        err = <-errs
        WarnAdmin("Ошибка зафиксирована, но у пользователя всё ок", err.(*errors.Error).ErrorStack(), "\n"+text, from, to)
        logrus.Error(err)
    }


    From := langs[from]
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("From " + From.Emoji + " " + From.Name, "none")),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("🔊 " + user.Localize("To voice"),  "speech_this_message_and_replied_one:"+from+":"+to)),
    )


    if len(rev.ContextResults.Results) > 0 {
        if len(rev.ContextResults.Results[0].SourceExamples) > 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
                tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("💬 " + user.Localize("Examples"), "examples:"+from+":"+to)))
        }
        if rev.ContextResults.Results[0].Translation != "" {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
                tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("📚 " + user.Localize("Translations"), "translations:"+from+":"+to)))

        }
    }

    if dict.Status == 200 && dict.DictionaryData != nil {
        keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
            tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("ℹ️" + user.Localize("Dictionary"), "dictonary:"+from+":"+text)))
    }

    if _, err = bot.Send(tgbotapi.EditMessageTextConfig{
        BaseEdit:              tgbotapi.BaseEdit{
            ChatID:          message.From.ID,
            MessageID:       msg.MessageID,
            ReplyMarkup:     &keyboard,
        },
        Text:                  tr.Text,
        ParseMode:             tgbotapi.ModeHTML,
        Entities:              nil,
        DisableWebPagePreview: false,
    }); err != nil {
        pp.Println(err)
    }

    analytics.Bot(user.ID, tr.Text, "Translated")
    user.IncrUsings()
    user.WriteBotLog("pm_translate", tr.Text)
}