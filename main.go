package main

import (
	"encoding/json"
	"fmt"
	"github.com/detectlanguage/detectlanguage-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	db *gorm.DB
	bot *tgbotapi.BotAPI
)


type Users struct {
	ID int64
	MyLang string
	ToLang string
	Act string
}

type TranslateAPIResponse struct {
	Align []string `json:"align"`
	Code  int      `json:"code"`
	Lang  string   `json:"lang"`
	Text  []string `json:"text"`
}

func Translate(lang, text string) (TranslateAPIResponse, error) {
	req, err := http.NewRequest("GET", "https://just-translated.p.rapidapi.com?lang=" + url.QueryEscape(lang) + "&text=" + url.QueryEscape(text), nil)
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
	db, err = gorm.Open(postgres.Open("host=ec2-63-34-97-163.eu-west-1.compute.amazonaws.com user=wzlryrrgxbgsbw password=b578bdbc77b5394a60f57660487149ca2238e0cbaf1cdbfb8b931f1168af24c7 dbname=d21k8q9pl6acl4 port=5432 TimeZone=Europe/Moscow"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Initializing bot
	const botToken string = "1737819626:AAEoc8WyCq_8rFQcY4q0vtkhqCKro8AudfI"
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false // >:(


	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/" + botToken:
			if isPost := ctx.IsPost(); isPost {
				data := ctx.PostBody()
				var update tgbotapi.Update
				err := json.Unmarshal(data, &update)
				if err != nil {
					fmt.Fprint(ctx, "can't parse")
				} else {
					go botRun(&update)
				}
			} else {
				fmt.Fprint(ctx, "no way")
			}
		default:
			_, err = fmt.Fprintln(ctx, "ok")
			if err != nil {
				panic(err)
			}
		}
	}
	if err = fasthttp.ListenAndServe(":"+port, requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}

}

func pingAdmin(err error) {
	bot.Send(tgbotapi.NewMessage(579515224, "ERROR: " + err.Error()))
}

func setUserStep(chatID int64, step string) error {
	return db.Model(&Users{ID: chatID}).Update("act", step).Error
}

func getUserStep(chatID int64) (string, error) {
	var user Users
	err := db.Model(&Users{ID: chatID}).Select("act").Limit(1).Find(&user).Error
	return user.Act, err
}



func botRun(update *tgbotapi.Update) {
	if update.Message != nil {
		switch update.Message.Text {
		case "/start":
			var user Users
			err := db.Model(&Users{ID: update.Message.Chat.ID}).Take(&user).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					if update.Message.From.LanguageCode == "" {
						update.Message.From.LanguageCode = "en"
					}
					err = db.Create(&Users{ID: update.Message.Chat.ID, MyLang: update.Message.From.LanguageCode}).Error
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #9273423, try again later"))
						pingAdmin(err)
						return
					}
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #012033, try again later"))
				pingAdmin(err)
				return
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Your language is - *" + user.MyLang + "*, and translate language - *" + user.ToLang + "*.\n\nChange your lang - /my_lang\nChange translate lang - /to_lang"))
		case "/my_lang":
			defer bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			edit := tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Send few words *in your language*.")
			edit.ParseMode = tgbotapi.ModeMarkdown
			kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
			edit.ReplyMarkup = &kb
			bot.Send(edit)

			err := setUserStep(update.CallbackQuery.From.ID, "set_my_lang")
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "error #032, try again later"))
				pingAdmin(err)
				return
			}
		case "/to_lang":
			defer bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			edit := tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Send few words *in your language into you want translate*.\n\nUsage: send message to translate")
			edit.ParseMode = tgbotapi.ModeMarkdown
			kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
			edit.ReplyMarkup = &kb
			bot.Send(edit)

			err := setUserStep(update.CallbackQuery.From.ID, "set_my_lang")
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "error #032, try again later"))
				pingAdmin(err)
				return
			}
		default:
			userStep, err := getUserStep(update.Message.Chat.ID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #37289, try again later"))
				pingAdmin(err)
				return
			}
			switch userStep {
			case "set_my_lang", "set_translate_lang":
				userLang, err := DetectLanguage(update.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Could not detect language, please send something else again"))
					return
				}
				if userStep == "set_my_lang" {
					err = db.Model(&Users{ID: update.Message.Chat.ID}).Update("my_lang", userLang).Error
				}
				if userStep == "set_translate_lang" {
					err = db.Model(&Users{ID: update.Message.Chat.ID}).Update("to_lang", userLang).Error
				}
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #1010, try again later"))
					pingAdmin(err)
					return
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Language set."))
			}


			var user Users
			err = db.Model(&Users{ID: update.Message.From.ID}).Take(&user).Error
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #034, try again later"))
				pingAdmin(err)
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your language is - *" + user.MyLang + "*, and translate language - *" + user.ToLang + "*.\n\nChange your lang - /my_lang\nChange translate lang - /to_lang\n\nUsage: send message to translate")
			msg.ParseMode = tgbotapi.ModeMarkdown
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("My lang", "set_my_lang")),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Translate lang", "set_translate_lang")))
			bot.Send(msg)
		}
	}
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			var user Users
			err := db.Model(&Users{ID: update.CallbackQuery.From.ID}).Take(&user).Error
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "error #034, try again later"))
				pingAdmin(err)
				return
			}

			db.Model(&Users{ID: update.CallbackQuery.From.ID}).Updates(map[string]interface{}{"act": nil})
			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Your language is - *" + user.MyLang + "*, and translate language - *" + user.ToLang + "*. If you want to change them, click the button. Using: Just send message to translate.")
			msg.ParseMode = tgbotapi.ModeMarkdown
			kb := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("My lang", "set_my_lang")),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Translate lang", "set_translate_lang")))
			msg.ReplyMarkup = &kb
			bot.Send(msg)
		}
	}
}
