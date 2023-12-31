package dashbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/armanokka/laba_itmo_bot/pkg/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"strconv"
)

type DashBot struct {
	APIKey string
}

func NewAPI(APIkey string) DashBot {
	return DashBot{APIKey: APIkey}
}

func (d DashBot) request(data []byte) error {
	req, err := http.NewRequest("POST", "https://tracker.dashbot.io/track?platform=universal&v=10.1.1-rest&type=incoming&apiKey="+d.APIKey, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if (res.StatusCode < 200 || res.StatusCode > 200) && res.StatusCode < 500 {
		return errors.New("code non 200:" + strconv.Itoa(res.StatusCode))
	}
	return nil
}

func (d DashBot) User(message tgbotapi.Message) error {
	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}
	params := Message{
		Text:         message.Text,
		UserId:       message.Chat.ID,
		PlatformJson: string(raw),
		SessionId:    message.Chat.ID,
	}
	if message.Caption != "" {
		params.Text = message.Caption
	}
	data, err := json.Marshal(params)
	req, err := http.NewRequest("POST", "https://tracker.dashbot.io/track?platform=universal&v=10.1.1-rest&type=incoming&apiKey="+d.APIKey, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if (res.StatusCode < 200 || res.StatusCode > 200) && res.StatusCode < 500 {
		return errors.New("code non 200:" + strconv.Itoa(res.StatusCode))
	}
	return nil
}

func (d DashBot) Bot(message tgbotapi.MessageConfig, intent string) error {
	var btns []Button
	if message.ReplyMarkup != nil {
		switch markup := message.ReplyMarkup.(type) {
		case tgbotapi.InlineKeyboardMarkup:
			n := 0
			for _, row := range markup.InlineKeyboard {
				n += len(row)
			}
			btns = make([]Button, 0, n)
			for _, row := range markup.InlineKeyboard {
				for _, btn := range row {
					val, id := "", ""
					switch {
					case btn.CallbackData != nil:
						val, id = *btn.CallbackData, *btn.CallbackData
					case btn.Pay:
						val = ":pay"
					case btn.URL != nil:
						val, id = *btn.URL, *btn.CallbackData
					case btn.CallbackGame != nil:
						val = ":callback_game"
					case btn.LoginURL != nil:
						val = ":login_url"
					case btn.SwitchInlineQuery != nil:
						val = ":switch_inline_query"
					case btn.SwitchInlineQueryCurrentChat != nil:
						val = ":switch_inline_query"
					case btn.WebApp != nil:
						val = ":web_app"
					}
					btns = append(btns, Button{
						Id:    id,
						Label: btn.Text,
						Value: val,
					})
				}
			}
		case tgbotapi.ReplyKeyboardMarkup:
			n := 0
			for _, row := range markup.Keyboard {
				n += len(row)
			}
			btns = make([]Button, 0, n)
			for _, row := range markup.Keyboard {
				for _, btn := range row {
					val := ""
					switch {
					case btn.WebApp != nil:
						val = ":web_app"
					case btn.RequestContact:
						val = ":request_contact"
					case btn.RequestPoll != nil:
						val = ":request_poll"
					case btn.RequestLocation:
						val = ":request_location"
					}
					btns = append(btns, Button{
						Id:    btn.Text,
						Label: btn.Text,
						Value: val,
					})
				}
			}
		}
	}
	data, err := json.Marshal(Message{
		Text:    helpers.ApplyEntitiesHtml(message.Text, message.Entities, 4098)[0],
		UserId:  message.ChatID,
		Buttons: btns,
		Intent: Intent{
			Name:       intent,
			Inputs:     nil,
			Confidence: 1,
		},
		SessionId: message.ChatID,
	},
	)
	if err != nil {
		return err
	}
	return d.request(data)
}

func (d DashBot) InlineUser(update tgbotapi.InlineQuery) error {
	data, err := json.Marshal(Message{
		Text:   update.Query,
		UserId: update.From.ID,
		Intent: Intent{
			Name: "inline_query",
			Inputs: []Input{{
				Name:  "language_code",
				Value: update.From.LanguageCode,
			}, {
				Name:  "offset",
				Value: update.Offset,
			}, {
				Name:  "chat_type",
				Value: update.ChatType,
			}},
			Confidence: 1,
		},
		SessionId: update.From.ID,
	})
	if err != nil {
		return err
	}
	return d.request(data)
}

func (d DashBot) InlineBot(user tgbotapi.User, config tgbotapi.InlineConfig) error {
	var btns []Button
	if len(config.Results) > 0 {
		btns = make([]Button, 0, len(config.Results))
		for _, result := range config.Results {
			switch r := result.(type) {
			case tgbotapi.InlineQueryResultArticle:
				btns = append(btns, Button{
					Id:    r.ID,
					Label: r.Title + "\n" + r.Description,
					Value: r.URL,
				})
			case tgbotapi.InlineQueryResultPhoto:
				btns = append(btns, Button{
					Id:    r.ID,
					Label: r.Title + "\n" + r.Description + "\n" + r.Caption,
					Value: r.URL,
				})
			case tgbotapi.InlineQueryResultCachedPhoto:
				btns = append(btns, Button{
					Id:    r.ID,
					Label: r.Title + "\n" + r.Description + "\n" + r.Caption,
					Value: r.PhotoID,
				})
				//case tgbotapi.InlineQueryResultAudio:
				//	btns = append(btns, Button{
				//		Id:    r.ID,
				//		Label: r.Title + "\n" + r.Description + "\n" + r.Caption,
				//		Value: r.PhotoID,
				//	})
				//case tgbotapi.InlineQueryResultCachedAudio:
				//case tgbotapi.InlineQueryResultVideo:
			}
		}
	}
	data, err := json.Marshal(Message{
		Text:   config.SwitchPMText,
		UserId: user.ID,
		Intent: Intent{
			Name: "answer_inline_query",
			Inputs: []Input{{
				Name:  "next_offset",
				Value: config.NextOffset,
			}, {
				Name:  "switch_pm_text",
				Value: config.SwitchPMText,
			}, {
				Name:  "switch_pm_parameter",
				Value: config.SwitchPMParameter,
			}, {
				Name:  "is_personal",
				Value: strconv.FormatBool(config.IsPersonal),
			}},
			Confidence: 1,
		},
		SessionId: user.ID,
	})
	if err != nil {
		return err
	}
	return d.request(data)
}

func (d DashBot) InlineChosenInlineResult(update tgbotapi.ChosenInlineResult) error {
	data, err := json.Marshal(Message{
		Text:   update.Query,
		UserId: update.From.ID,
		Postback: Postback{
			ButtonClick: ButtonClick{update.ResultID},
		},
		SessionId: update.From.ID,
	})
	if err != nil {
		return err
	}
	return d.request(data)
}

func (d DashBot) UserStartedBot(user tgbotapi.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	data, err := json.Marshal(Event{
		Name:           "bot_start",
		ConversationId: strconv.FormatInt(user.ID, 10),
		Type:           PageLaunchEvent,
		ExtraInfo:      string(userJson),
	})
	if err != nil {
		return err
	}
	return d.request(data)

}

// UserButtonClick buttonID can be callback.Data or a title of keyboard button
func (d DashBot) UserButtonClick(user tgbotapi.User, buttonID string) error {
	params, err := json.Marshal(Message{
		Text:   ":button_click",
		UserId: user.ID,
		Intent: Intent{
			Name: "button_click",
			//Inputs:     nil,
			Confidence: 1,
		},
		//Images:       nil,
		//Buttons:      nil,
		Postback: Postback{ButtonClick: ButtonClick{
			ButtonId: buttonID,
		}},
		//PlatformJson: "",
		SessionId: user.ID,
	})
	if err != nil {
		return err
	}
	return d.request(params)
}
