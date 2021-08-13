package main

import (
    "errors"
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
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          update.CallbackQuery.From.ID,
            MessageID:       update.CallbackQuery.Message.MessageID,
        })
        return
    // case "rate":
    //     bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    //     bot.Send()
    }
    arr := strings.Split(update.CallbackQuery.Data, ":")
    if len(arr) == 0 {
        return
    }
    switch arr[0] {
    case "set_bot_lang": // arr[1] - lang code
        err := db.Model(&Users{}).Where("id = ?", update.CallbackQuery.From.ID).Limit(1).Update("lang", arr[1]).Error
        if err != nil {
            warn(err)
        }
        bot.Send(tgbotapi.DeleteMessageConfig{
            ChatID:          update.CallbackQuery.From.ID,
            MessageID:       update.CallbackQuery.Message.MessageID,
        })
        bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, Localize("Now press /start 👈", arr[1])))
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "speech": // arr[1] - lang code
        parts := strings.Split(update.CallbackQuery.Message.Text, "\n")
        update.CallbackQuery.Message.Text = strings.Join(parts[:len(parts)-1], "")
        sdec, err := translate.TTS(arr[1], update.CallbackQuery.Message.Text)
        if err != nil {
            if e, ok := err.(translate.TTSError); ok {
                if e.Code == 500 || e.Code == 414 {
                    bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Too big text"))
                    return
                }
            }
            bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Iternal error"))
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
        audio.Title = update.CallbackQuery.Message.Text
        audio.Performer = "@TransloBot"
        audio.ReplyToMessageID = update.CallbackQuery.Message.MessageID
        kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        audio.ReplyMarkup = kb
        bot.Send(audio)
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    case "variants": // arr[1] - from, arr[2] - to
        tr, err := translate.GoogleTranslate(arr[1], arr[2], update.CallbackQuery.Message.ReplyToMessage.Text)
        if err != nil {
            warn(err)
        }
        if len(tr.Variants) == 0 {
            warn(errors.New("не найдено вариантов перевода"))
        }
        var text string
        var limit uint8
        for _, variant := range tr.Variants {
            if limit >= 5 {
                break
            }
            text += "\n<b>" + variant.Word + "</b> - " + variant.Meaning
            limit++
        }
        msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, text)
        msg.ParseMode = tgbotapi.ModeHTML
        msg.ReplyToMessageID = update.CallbackQuery.Message.MessageID
        kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
        msg.ReplyMarkup = kb
        bot.Send(msg)
        bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
    }
}