package main

import (
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "strings"
)

func handleCallback(update *tgbotapi.Update) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Error, sorry"))
        WarnAdmin(err)
    }
    switch update.CallbackQuery.Data {
    case "back":
        defer bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
        var user Users
        err := db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.CallbackQuery.From.ID).Limit(1).Find(&user).Error
        if err != nil {
            bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "error #034, try again later"))
            warn(err)
            return
        }
        
        err = db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Limit(1).Updates(map[string]interface{}{"act": nil}).Error
        if err != nil {
            warn(err)
            return
        }
        user.MyLang = iso6391.Name(user.MyLang)
        user.ToLang = iso6391.Name(user.ToLang)
        msg := tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Your language is - <b>"+user.MyLang+"</b>, and translate language - <b>"+user.ToLang+"</b>.\n\nChange your lang /my_lang\nChange translate lang /to_lang")
        msg.ParseMode = tgbotapi.ModeHTML
        bot.Send(msg)
    case "none":
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    }
    
    arr := strings.Split(update.CallbackQuery.Data, ":")
    if len(arr) == 0 {
        return
    }
    switch arr[0] {
    case "set_my_lang": // arr[1] - language code
        err := db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Updates(map[string]interface{}{"act": nil, "my_lang": arr[1]}).Error
        if err != nil {
            warn(err)
            return
        }
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
        replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
        edit := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Now your language is "+iso6391.Name(arr[1]), replyMarkup)
        bot.Send(edit)
    case "set_translate_lang": // arr[1] - language code
        err := db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Updates(map[string]interface{}{"act": nil, "to_lang": arr[1]}).Error
        if err != nil {
            bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "error #435"))
            warn(err)
            return
        }
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
        replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
        edit := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Now translate language is "+iso6391.Name(arr[1]), replyMarkup)
        bot.Send(edit)
    }
    
}

