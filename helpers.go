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

// DetectLang detect language of text and return 1) name of language, 2) code of language, 3) error
func DetectLang(text string) (string, string, error) {
    if c := iso6391.Name(text); c != "" { // Текст - код языка, c - его название
        return c, text, nil
    }
    if c := iso6391.CodeForName(text); c != "" { // Текст - это название языка, c - код
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
    
    if c := iso6391.Name(strings.ToLower(tr.Text)); c != "" { // Текст - это код языка, c - его название
        return c, tr.Text, nil
    }
    if c := iso6391.CodeForName(strings.Title(tr.Text)); c != "" { // Текст - это название языка, c - код
        return tr.Text, c, nil
    }
    
    if lang != "" && lang != "auto" {
        return iso6391.Name(lang), lang, nil
    }
    
    return "", "", errors.New("could not detect language...")
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

