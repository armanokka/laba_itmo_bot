package main

import (
	"encoding/json"
	"github.com/detectlanguage/detectlanguage-go"
	"github.com/k0kubun/pp"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TranslateAPIResponse struct {
	Align []string `json:"align"`
	Code  int      `json:"code"`
	Lang  string   `json:"lang"`
	Text  []string `json:"text"`
}

func Translate(lang, text string) (TranslateAPIResponse, error) {
	req, err := http.NewRequest("GET", "https://just-translated.p.rapidapi.com?lang=" + url.QueryEscape(lang) + "&text=How?", nil)
	if err != nil {
		return TranslateAPIResponse{}, err
	}
	req.Header["x-rapidapi-key"] = []string{"561a41f76amsha7c0323d47335aep1986ecjsn4e41c7c2b518"}
	req.Header["x-rapidapi-host"] = []string{"just-translated.p.rapidapi.com"}
	req.Header["useQueryString-type"] = []string{"true"}
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return TranslateAPIResponse{}, err
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TranslateAPIResponse{}, err
	}
	var out TranslateAPIResponse
	err = json.Unmarshal(response, &out)
	if err != nil {
		return TranslateAPIResponse{}, err
	}
	return out, err
}

func DetectLanguage(text string) (string, error){
	state := detectlanguage.New("c71fb63df8bdc8ea8bc4c0b20771aa5f")
	detections, err := state.Detect(text)
	if err != nil {
		return "", err
	}
	return detections[0].Language, nil
}

func main() {
	// Initializing PostgreSQL DB
	var err error
	db, err = gorm.Open(postgres.Open("host=ec2-54-247-158-179.eu-west-1.compute.amazonaws.com user=ixjrwpdkbmlzas password=bc6cf01fc60625a154f3398bf506891016818fa5183154abc2b2c8b1f09a2b7f dbname=dc6ad698nlfv2m port=5432 TimeZone=Europe/Moscow"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Initializing bot
	const botToken string = "1757580922:AAHc0hz1sdqG216mFnqMZXTaDJA2cKuC1H0"
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false // >:(


}
