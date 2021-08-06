package main

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI struct {
    *tgbotapi.BotAPI
}

// func NewBotAPI(token string) (*BotAPI, error) {
//     return NewBotAPIWithClient(token, "https://api.telegram.org/bot%s/%s", &http.Client{})
// }
//
// func NewBotAPIWithClient(token, endpoint string, client *http.Client) (*BotAPI, error) {
//     bot = &BotAPI{
//         BotAPI:          &tgbotapi.BotAPI{
//             Token:  token,
//             Debug:  false,
//             Buffer: 10,
//             Client: client,
//         },
//         apiEndpoint:     endpoint,
//         shutdownChannel: nil,
//     }
//     if _, err := bot.GetMe(); err != nil {
//         return nil, err
//     }
//     return bot, nil
// }

func (bot *BotAPI) AnswerCallbackQuery(config tgbotapi.CallbackConfig) (*tgbotapi.APIResponse, error){
    var params = make(tgbotapi.Params, 0)
    params.AddNonEmpty("callback_query_id", config.CallbackQueryID)
    params.AddNonEmpty("text", config.Text)
    params.AddBool("show_alert", config.ShowAlert)
    params.AddNonEmpty("url", config.URL)
    params.AddNonZero("cache_time", config.CacheTime)
    return bot.MakeRequest("answerCallbackQuery", params)
}

func (bot *BotAPI) AnswerInlineQuery(config tgbotapi.InlineConfig) (*tgbotapi.APIResponse, error) {
    params := make(tgbotapi.Params, 0)
    params.AddNonEmpty("inline_query_id", config.InlineQueryID)
    params.AddInterface("results", config.Results)
    params.AddNonZero("cache_time", config.CacheTime)
    params.AddBool("is_personal", config.IsPersonal)
    params.AddNonEmpty("next_offset", config.NextOffset)
    params.AddNonEmpty("switch_pm_text", config.SwitchPMText)
    params.AddNonEmpty("switch_pm_parameter", config.SwitchPMParameter)
    return bot.MakeRequest("answerInlineQuery", params)
}