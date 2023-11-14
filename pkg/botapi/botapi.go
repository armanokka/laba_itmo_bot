package botapi

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotAPI struct {
	*tgbotapi.BotAPI
}

func NewBotAPI(token string) (BotAPI, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return BotAPI{}, err
	}
	return BotAPI{
		BotAPI: api,
	}, nil
}

func NewBotAPIWithEndpoint(token, endpoint string) (BotAPI, error) {
	api, err := tgbotapi.NewBotAPIWithAPIEndpoint(token, endpoint)
	if err != nil {
		return BotAPI{}, err
	}
	return BotAPI{
		BotAPI: api,
	}, nil
}

func (bot BotAPI) AnswerCallbackQuery(config tgbotapi.CallbackConfig) (*tgbotapi.APIResponse, error) {
	var params = make(tgbotapi.Params, 0)
	params.AddNonEmpty("callback_query_id", config.CallbackQueryID)
	params.AddNonEmpty("text", config.Text)
	params.AddBool("show_alert", config.ShowAlert)
	params.AddNonEmpty("url", config.URL)
	params.AddNonZero("cache_time", config.CacheTime)
	return bot.MakeRequest("answerCallbackQuery", params)
}

func (bot BotAPI) AnswerInlineQuery(config tgbotapi.InlineConfig) (*tgbotapi.APIResponse, error) {
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

func (bot BotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	resp, err := bot.Request(c)
	if err != nil {
		if cfg, ok := c.(tgbotapi.EditMessageTextConfig); ok {
			return bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   cfg.ChatID,
					ChannelUsername:          cfg.ChannelUsername,
					ReplyMarkup:              cfg.ReplyMarkup,
					DisableNotification:      true,
					AllowSendingWithoutReply: true,
				},
				Text:                  cfg.Text,
				ParseMode:             cfg.ParseMode,
				Entities:              cfg.Entities,
				DisableWebPagePreview: cfg.DisableWebPagePreview,
			})
		}
		return tgbotapi.Message{}, err
	}

	var message tgbotapi.Message
	err = json.Unmarshal(resp.Result, &message)

	return message, err
}
