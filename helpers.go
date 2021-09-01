/*
Helper functions
 */
package main

import (
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "runtime/debug"
    "strconv"
    "strings"
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

// DetectLang detect language of text and return 1) name of language, 2) code of language, 3) error
func DetectLang(text string) (string, string, error) {
    text = strings.ToLower(text)
    if c := iso6391.Name(text); c != "" { // Текст - код языка на английском
        return c, text, nil
    }
    if c := iso6391.CodeForName(text); c != "" { // Текст - это название языка на английском
        return text, c, nil
    }
    lang, err := translate.DetectLanguageGoogle(text)
    if err != nil {
        return "", "", err
    }
    if lang == "" {
        lang = "auto"
    }
    tr, err := translate.GoogleTranslate(lang, "en", text)
    if err != nil {
        return "", "", err
    }
    
    if c := iso6391.Name(strings.ToLower(tr.Text)); c != "" { // Текст - это код языка на английском
        return c, tr.Text, nil
    }
    if c := iso6391.CodeForName(strings.Title(tr.Text)); c != "" { // Текст - это название языка на английском
        return tr.Text, c, nil
    }
    
    if lang != "" && lang != "auto" {
        return iso6391.Name(lang), lang, nil
    }
    
    return "", "", errors.New("could not detect language")
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


func GetTickedCallbacks(keyboard *tgbotapi.InlineKeyboardMarkup) []string {
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