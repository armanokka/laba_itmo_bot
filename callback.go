package main

import (
    "github.com/armanokka/translobot/translate"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "os"
    "strings"
)

func handleCallback(update *tgbotapi.Update) {
    warn := func(err error) {
        bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Error, sorry"))
        WarnAdmin(err)
    }
    switch update.CallbackQuery.Data {
    case "delete":
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          update.CallbackQuery.From.ID,
            MessageID:       update.CallbackQuery.Message.MessageID,
        })
        return
    }
    arr := strings.Split(update.CallbackQuery.Data, ":")
    if len(arr) == 0 {
        return
    }
    switch arr[0] {
    case "set_bot_lang": // arr[1] - lang code
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
        err := db.Model(&Users{}).Where("id = ?", update.CallbackQuery.From.ID).Limit(1).Update("lang", arr[1]).Error
        if err != nil {
            warn(err)
        }
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          update.CallbackQuery.From.ID,
            MessageID:       update.CallbackQuery.Message.MessageID,
        })
        bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, Localize("Now press /start 👈", arr[1])))
    case "speech": // arr[1] - lang code
        sdec, err := translate.TTS(arr[1], update.CallbackQuery.Message.Text)
        if err != nil {
            bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Too big text of iternal error"))
            warn(err)
            return
        }
        f, err := os.CreateTemp("", "")
        if err != nil {
            warn(err)
            return
        }
        _, err = f.Write(sdec)
        if err != nil {
            warn(err)
            return
        }
        audio := tgbotapi.NewAudio(update.CallbackQuery.From.ID, f.Name())
        err = f.Close()
        if err != nil {
            warn(err)
            return
        }
        audio.Title = cutString(update.CallbackQuery.Message.Text, 25)
        audio.Performer = "Google"
        audio.ReplyToMessageID = update.CallbackQuery.Message.MessageID
        bot.Send(audio)
    }
}