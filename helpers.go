/*
Helper functions
 */
package main

import (
    "encoding/json"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "runtime/debug"
    "strconv"
    "strings"
    "unicode/utf16"
)

func SendMenu(user User) {
    msg := tgbotapi.NewMessage(user.ID, user.Localize("/start", iso6391.Name(user.MyLang), iso6391.Name(user.ToLang)))
    keyboard := tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton(user.Localize("My Language")),
            tgbotapi.NewKeyboardButton(user.Localize("Translate Language")),
        ),
    )

    msg.ReplyMarkup = keyboard
    msg.ParseMode = tgbotapi.ModeHTML
    bot.Send(msg)

    analytics.Bot(user.ID, msg.Text, "Start")
}

func SendHelp(user User) {
    msg := tgbotapi.NewMessage(user.ID, user.Localize("/help"))
    msg.ParseMode = tgbotapi.ModeHTML
    query := user.Localize("пишите сюда")
    btn := tgbotapi.InlineKeyboardButton{
        Text:                         "inline",
        SwitchInlineQuery:            &query,
    }
    keyboard := tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(btn))
    msg.ReplyMarkup = keyboard
    bot.Send(msg)

    analytics.Bot(user.ID, msg.Text, "Help")
}

// WarnAdmin send error message to the admin (by AdminID const)
func WarnAdmin(args ...interface{}) {
    var text = "Error caused:\n"
    for _, arg := range args {
        switch v := arg.(type) {
        case string:
            text += "\n" + v
        case error:
            text += "\n"+ v.Error()
            text += "\nTrace:<code>" + string(debug.Stack())+"</code>\n"
        case int:
            text += "\n" + strconv.Itoa(v)
        default:
            text += "[wrong type given]"
        }
    }
    msg := tgbotapi.NewMessage(AdminID, text)
    bot.Send(msg)
}

func WarnErrorAdmin(err error) {
    msg := tgbotapi.NewMessage(AdminID, err.Error() + "\n\n" + string(debug.Stack()))
    bot.Send(msg)
}

// setUserStep set user's step to your. If string is empty "", then step will be null
func setUserStep(chatID int64, step string) error {
    if step == "" {
        return db.Model(&Users{ID: chatID}).Where("id = ?", chatID).Limit(1).Updates(map[string]interface{}{"act":nil}).Error
    }
    return db.Model(&Users{ID: chatID}).Where("id = ?", chatID).Limit(1).Update("act", step).Error
}

// getUserStep return user's step at the moment. Step can be set by setUserStep
func getUserStep(chatID int64) (string, error) {
    var user Users
    err := db.Model(&Users{ID: chatID}).Select("act").Where("id", chatID).Limit(1).Find(&user).Error
    return user.Act.String, err
}

// cutString cut string using runes by limit
func cutString (text string, limit int) string {
    runes := []rune(text)
    if len(runes) > limit {
        return string(runes[:limit])
    }
    return text
}

func inArray(k string, arr []string,) bool {
    for _, v := range arr {
        if k == v {
             return true
        }
    }
    return false
}


func GetTickedCallbacks(keyboard tgbotapi.InlineKeyboardMarkup) []string {
    callbacks := make([]string, 0)
    for _, row := range keyboard.InlineKeyboard {
        for _, button := range row {
            if strings.HasPrefix(button.Text, "✅") {
                callbacks = append(callbacks, *button.CallbackData)
            }
        }
    }
    return callbacks
}

func TickByCallback(uniqueCallbackData string, keyboard *tgbotapi.InlineKeyboardMarkup) {
    var done bool
    for i1, row := range keyboard.InlineKeyboard {
        if done {
            break
        }
        for i2, button := range row {
            if *button.CallbackData == uniqueCallbackData && !strings.HasPrefix(*button.CallbackData, "✅ ") {
                keyboard.InlineKeyboard[i1][i2].Text = "✅ " + button.Text
                done = true
                break
            }
        }
    }
}

func UnTickByCallback(uniqueCallbackData string, keyboard *tgbotapi.InlineKeyboardMarkup) {
    var done bool
    for i1, row := range keyboard.InlineKeyboard {
        if done {
            break
        }
        for i2, button := range row {
            if *button.CallbackData == uniqueCallbackData {
                keyboard.InlineKeyboard[i1][i2].Text = strings.TrimPrefix(button.Text, "✅ ")
                done = true
                break
            }
        }
    }
}

func IsTicked(callback string, keyboard *tgbotapi.InlineKeyboardMarkup) bool {
    for _, row := range keyboard.InlineKeyboard {
        for _, button := range row {
            if *button.CallbackData != callback {
                continue
            }
            if strings.HasPrefix(button.Text, "✅") {
                return true
            }
        }
    }
    return false
}

func applyEntitiesHtml(text string, entities []tgbotapi.MessageEntity) string {
    if len(entities) == 0 {
        return text
    }

    encoded := utf16.Encode([]rune(text))
    out := make([]uint16, 0, len(encoded))

    pointers := make(map[int][]uint16)
    for _, entity := range entities {
        var startTag string
        switch entity.Type {
        case "code", "pre":
            startTag = `<label class="notranslate"><code>`
        case "mention", "hashtag", "cashtag", "bot_command", "url", "email", "phone_number":
            startTag = `<label class="notranslate">` // very important to keep '<label class="notranslate">' strongly correct, without any spaces or another
        case "bold":
            startTag = `<b>`
        case "italic":
            startTag = `<i>`
        case "underline":
            startTag = `<u>`
        case "strikethrough":
            startTag = `<s>`
        case "text_link":
            startTag = `<a href="` + entity.URL + `">`
        case "text_mention":
            startTag = `<a href="tg://user?id=` + strconv.FormatInt(entity.User.ID, 10) + `">`
        }
        pointers[entity.Offset] = append(pointers[entity.Offset], utf16.Encode([]rune(startTag))...)


        startTag = strings.TrimPrefix(startTag, "<")
        var endTag string
        switch entity.Type {
        case "code", "pre":
            endTag = "</code></label>" // very important to keep '</label>' strongly correct, without any spaces or another
        case "mention", "hashtag", "cashtag", "bot_command", "url", "email", "phone_number":
            endTag = `</label>`
        case "bold":
            endTag = `</b>`
        case "italic":
            endTag = `</i>`
        case "underline":
            endTag = `</u>`
        case "strikethrough":
            endTag = `</s>`
        case "text_link", "text_mention":
            endTag = `</a>`
        }
        pointers[entity.Offset+entity.Length] = append(pointers[entity.Offset+entity.Length], utf16.Encode([]rune(endTag))...)
    }



    for i, ch := range encoded {
        if m, ok := pointers[i]; ok {
            out = append(out, m...)
        }
        out = append(out, ch)

        if i == len(encoded) - 1 {
            if m, ok := pointers[i+1]; ok {
                out = append(out, m...)
            }
        }
    }
    ret := utf16.Decode(out)
    return strings.Replace(string(ret), "\n", "<br>", -1)
}

func setMyCommands(langs []string, commands []tgbotapi.BotCommand) error {
    newCommands := make(map[string][]tgbotapi.BotCommand)
    for _, lang := range langs {
        newCommands[lang] = []tgbotapi.BotCommand{}
        for _, command := range commands {
            tr, err := translate.GoogleHTMLTranslate("en", lang, command.Description)
            if err != nil {
                return err
            }
            newCommands[lang]= append(newCommands[lang], tgbotapi.BotCommand{
                Command:     command.Command,
                Description: tr.Text,
            })
        }
    }

    for lang, command := range newCommands {
        data, err := json.Marshal(command)
        if err != nil {
            return err
        }
        params := tgbotapi.Params{}
        params.AddNonEmpty("commands", string(data))
        params.AddNonEmpty("language_code", lang)
        if _, err = bot.MakeRequest("setMyCommands", params); err != nil {
            return err
        }
    }
    return nil
}