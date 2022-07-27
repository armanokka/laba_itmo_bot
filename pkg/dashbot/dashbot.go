package dashbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/armanokka/translobot/pkg/helpers"
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

type Message struct {
	Text         string   `json:"text"`
	UserId       string   `json:"userId"`
	Intent       Intent   `json:"intent"`
	Images       []string `json:"images"`
	Buttons      []Button `json:"buttons"`
	Postback     Postback `json:"postback"`
	PlatformJson string   `json:"platformJson"` // any json
	//PlatformUserJson map[string]interface{} `json:"platformUserJson"` // there is paid access to this field
	SessionId string `json:"sessionId"`
}

type Postback struct {
	ButtonClick ButtonClick `json:"buttonClick"`
}

type ButtonClick struct {
	ButtonId string `json:"buttonId"`
}

type Intent struct {
	Name   string `json:"name"`
	Inputs []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"inputs"`
	Confidence float64 `json:"confidence"`
}

type Button struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

func (d DashBot) User(message tgbotapi.Message) error {
	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}
	params := Message{
		Text:         message.Text,
		UserId:       strconv.FormatInt(message.Chat.ID, 10),
		PlatformJson: string(raw),
		SessionId:    strconv.FormatInt(message.Chat.ID, 10),
	}
	if message.Caption != "" {
		params.Text = message.Caption
	}
	data, err := json.Marshal(params)
	req, err := http.NewRequest("POST", "https://tracker.dashbot.io/track?platform=universal&v=10.1.1-rest&type=incoming&apiKey="+d.APIKey, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
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
					val := ""
					switch {
					case btn.CallbackData != nil:
						val = *btn.CallbackData
					case btn.Pay:
						val = ":pay"
					case btn.URL != nil:
						val = *btn.URL
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
						Id:    *btn.CallbackData,
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
		Text:    helpers.ApplyEntitiesHtml(message.Text, message.Entities),
		UserId:  strconv.FormatInt(message.ChatID, 10),
		Buttons: btns,
		Intent: Intent{
			Name:   intent,
			Inputs: nil,
		},
		SessionId: strconv.FormatInt(message.ChatID, 10),
	},
	)
	req, err := http.NewRequest("POST", "https://tracker.dashbot.io/track?platform=universal&v=11.1.0-rest&type=outgoing&apiKey="+d.APIKey, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header["Content-Type"] = []string{"application/json"}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 299 || res.StatusCode < 200 {
		return errors.New("code non 200:" + strconv.Itoa(res.StatusCode))
	}
	return nil
}
