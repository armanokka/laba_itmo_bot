package dashbot

import (
    "bytes"
    "encoding/json"
    "errors"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "net/http"
)

type DashBot struct {
    APIKey string
    Error func(err error)
}

func NewAPI(APIkey string, errFunc func(error)) *DashBot {
    return &DashBot{APIKey: APIkey, Error: errFunc}
}

func (d *DashBot) User(message string,  user *tgbotapi.User) {
    params := map[string]interface{}{
        "text":message,
        "userId":user.ID,
        "platformUserJson": map[string]string{
            "firstName": user.FirstName,
            "lastName": user.LastName,
            "locale": user.LanguageCode,
        },
        "platformJson": map[string]string{
            "username": user.UserName,
        },
    }
    data, err := json.Marshal(params)
    req, err := http.NewRequest("POST", "https://tracker.dashbot.io/track?platform=universal&v=10.1.1-rest&type=incoming&apiKey=" + d.APIKey, bytes.NewBuffer(data))
    if err != nil {
        d.Error(err)
        return
    }
    req.Header.Set("Content-Type", "application/json")
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        d.Error(err)
        return
    }
    if res.StatusCode != 200 {
        d.Error(errors.New("code non 200"))
    }
}

func (d *DashBot) Bot(id int64, message, intent string) {
    params := map[string]interface{}{
        "text":message,
        "userId":id,
        "intent": map[string]string{
            "name": intent,
        },
    }
    data, err := json.Marshal(params)
    req, err := http.NewRequest("POST", "https://tracker.dashbot.io/track?platform=universal&v=11.1.0-rest&type=outgoing&apiKey=" + d.APIKey, bytes.NewBuffer(data))
    if err != nil {
        d.Error(err)
        return    }
    req.Header.Set("Content-Type", "application/json")
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        d.Error(err)
        return
    }
    if res.StatusCode != 200 {
        d.Error(errors.New("code non 200"))
    }
}