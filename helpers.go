/*
Helper functions
 */
package main

import (
    "encoding/json"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    "github.com/go-errors/errors"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "runtime/debug"
    "strconv"
    "strings"
    "unicode/utf16"
)

func format(i int64) string {
    return strconv.FormatInt(i, 10)
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
            text += "\n" + pp.Sprint(arg)
        }
    }
    msg := tgbotapi.NewMessage(AdminID, text)
    bot.Send(msg)
}

func WarnErrorAdmin(err error) {
    msg := tgbotapi.NewMessage(AdminID, err.Error() + "\n\n" + string(debug.Stack()))
    bot.Send(msg)
}

// cutString cut string using runes by limit
func cutStringUTF16 (text string, limit int) string {
    points := utf16.Encode([]rune(text))
    if len(points) > limit {
        return string(utf16.Decode(points[:limit]))
    }
    return text
}

func in(arr []string, k string) bool {
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
    pointers := make(map[int]string)

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

        pointers[entity.Offset] += startTag


        //startTag = strings.TrimPrefix(startTag, "<")
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
        pointers[entity.Offset+entity.Length] += endTag
    }

    var out = make([]uint16, 0, len(encoded))

    for i, ch := range encoded {
        if m, ok := pointers[i]; ok {
            out = append(out, utf16.Encode([]rune(m))...)
        }
        out = append(out, ch)

        if i == len(encoded) - 1 {
            if m, ok := pointers[i+1]; ok {
                out = append(out, utf16.Encode([]rune(m))...)
            }
        }
    }
    ret := string(utf16.Decode(out))
    ret = strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "").Replace(ret)
    ret = strings.ReplaceAll(ret, `<br>`, "\n")
    return ret
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

func makeArticle(id string, title, description, messageText string) tgbotapi.InlineQueryResultArticle {
    //keyboard := tgbotapi.NewInlineKeyboardMarkup(
    //    tgbotapi.NewInlineKeyboardRow(
    //        tgbotapi.InlineKeyboardButton{
    //            Text:                         "translate",
    //            SwitchInlineQueryCurrentChat: &description,
    //        }))
    return tgbotapi.InlineQueryResultArticle{
        Type:                "article",
        ID:                  id,
        Title:               title,
        InputMessageContent: map[string]interface{}{
            "message_text": messageText,
            "disable_web_page_preview":false,
        },
        //ReplyMarkup: &keyboard,
        //URL:                 "https://t.me/TransloBot?start=from_inline",
        //HideURL:             true,
        Description:         description,
    }
}

func inMapValues(m map[string]string, values ...string) bool {
    for _, v := range values {
        var ok bool
        for _, v1 := range m {
            if v  == v1 {
                ok = true
                break
            }
        }
        if !ok {
            return false
        }
    }
    return true
}

func sendSpeech(lang, text string, callbackID string, user User) error {
    sdec, err := translate.TTS(lang, text)
    if err != nil {
        if err == translate.ErrTTSLanguageNotSupported {
            call := tgbotapi.NewCallback(callbackID, user.Localize("%s language is not supported", iso6391.Name(lang)))
            call.ShowAlert = true
            bot.AnswerCallbackQuery(call)
            return nil
        }
        if e, ok := err.(translate.HTTPError); ok {
            if e.Code == 500 || e.Code == 414 {
                call := tgbotapi.NewCallback(callbackID, "Too big text")
                call.ShowAlert = true
                bot.AnswerCallbackQuery(call)
                return nil
            }
        }
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "Iternal error"))
        pp.Println(err)
        return err
    }
    bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "⏳"))

    f, err := os.CreateTemp("", "")
    if err != nil {
        return err
    }
    defer func() {
        if err = f.Close(); err != nil {
            WarnAdmin(err)
        }
    }()
    _, err = f.Write(sdec)
    if err != nil {
        return err
    }
    audio := tgbotapi.NewAudio(user.ID, f.Name())
    audio.Title = text
    kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
    audio.ReplyMarkup = kb
    bot.Send(audio)
    return nil
}


func WitAiSpeech(wav io.Reader, lang string, bits int) (string, error) {
    var key, ok = WitAPIKeys[lang]
    if !ok {
        return "", errors.New("no wit.ai key for lang " + lang)
    }
    req, err := http.NewRequest("POST", "https://api.wit.ai/speech?v=20210928&bits=" + strconv.Itoa(bits), wav)
    if err != nil {
        panic(err)
    }
    req.Header.Set("Authorization", "Bearer " + key)
    req.Header.Set("Content-Type", "audio/wave")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    parts := strings.Split(string(body), "\r\n")
    if len(parts) == 0 {
        return "", errors.New("empty parts: " + string(body))
    }
    var result struct{
        Text string `json:"text"`
        Error string `json:"error"`
        Code string `json:"code"`
    }
    if err = json.Unmarshal([]byte(parts[len(parts)-1]), &result); err != nil {
        return "", errors.WrapPrefix(err, string(body), 0)
    }
    if result.Error != "" {
        return "", errors.New(result.Error)
    }
    return result.Text, nil
}

func BuildSupportedLanguagesKeyboard() (tgbotapi.InlineKeyboardMarkup, error) {
    keyboard := tgbotapi.NewInlineKeyboardMarkup()
    for i, code := range BotLocalizedLangs {
        lang, ok := langs[code]
        if !ok {
            return tgbotapi.InlineKeyboardMarkup{}, errors.New("no such code "+ code + " in langs")
        }

        if i % 2 == 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_bot_lang_and_register:"  + code)))
        } else {
            l := len(keyboard.InlineKeyboard)-1
            keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name, "set_bot_lang_and_register:"  + code))
        }
    }
    return keyboard, nil
}