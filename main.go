// t.me/translobot source code
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/detectlanguage/detectlanguage-go"
	iso6391 "github.com/emvi/iso-639-1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	db  *gorm.DB
	bot *tgbotapi.BotAPI
	detectLangAPIKeys []string
)

type Users struct {
	ID     int64 `gorm:"primaryKey;index;not null"`
	MyLang string `gorm:"default:en"`
	ToLang string `gorm:"default:ar;n"`
	Act    string
}

type TranslateAPIResponse struct {
	Error error `json:"err"`
	Result string `json:"result"`
}

// pingAdmin attempts admin your error
func pingAdmin(err error) {
	bot.Send(tgbotapi.NewMessage(579515224, "ERROR: "+err.Error()))
}

func setUserStep(chatID int64, step string) error {
	return db.Model(&Users{ID: chatID}).Where("id = ?", chatID).Limit(1).Update("act", step).Error
}

func getUserStep(chatID int64) (string, error) {
	var user Users
	err := db.Model(&Users{ID: chatID}).Select("act").Where("id", chatID).Limit(1).Find(&user).Error
	return user.Act, err
}
// botRun is main handler of bot
func botRun(update *tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.Text == "" {
			return
		}
		switch update.Message.Text {
		case "/start":
			var user = Users{ID: update.Message.From.ID}
			if update.Message.From.LanguageCode == "" {
				update.Message.From.LanguageCode = "en"
			}
			err := db.FirstOrCreate(&user, &Users{ID: update.Message.From.ID}).Error
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #1002, try again later"))
				pingAdmin(err)
				return
			}
			user.MyLang = iso6391.Name(user.MyLang)
			user.ToLang = iso6391.Name(user.ToLang)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your language is - <b>"+user.MyLang+"</b>, and translate language - <b>"+user.ToLang+"</b>.\n\nChange your lang /my_lang\nChange translate lang /to_lang")
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
		case "/my_lang":
			edit := tgbotapi.NewMessage(update.Message.Chat.ID, "Send few words *in your language*.")
			edit.ParseMode = tgbotapi.ModeMarkdown
			kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
			edit.ReplyMarkup = &kb
			bot.Send(edit)

			err := setUserStep(update.Message.Chat.ID, "set_my_lang")
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #032, try again later"))
				pingAdmin(err)
				return
			}
		case "/to_lang":
			edit := tgbotapi.NewMessage(update.Message.Chat.ID, "Send few words *in your language into you want translate*.")
			edit.ParseMode = tgbotapi.ModeMarkdown
			kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
			edit.ReplyMarkup = &kb
			bot.Send(edit)

			err := setUserStep(update.Message.Chat.ID, "set_translate_lang")
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #032, try again later"))
				pingAdmin(err)
				return
			}
		default: // Сообщение не является командой.
			userStep, err := getUserStep(update.Message.Chat.ID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "error #37289, try again later"))
				pingAdmin(err)
				return
			}
			switch userStep {
			case "set_my_lang", "set_translate_lang":
				languageDetections, err := DetectLanguage(update.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Could not detect language, please send something else again"))
					return
				}

				keyboard := tgbotapi.NewInlineKeyboardMarkup()
				for _, lang := range languageDetections { // Собираем клавиатуру из языков
					row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(iso6391.Name(lang.Language), userStep + ":"+ lang.Language)) // set_my_lang:en or set_translate_lang:en, set_my_lang:ru etc.
					keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
				}

				if len(keyboard.InlineKeyboard) == 0 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, but this language is unsupported."))
					return
				}

				keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please, select one of this languages:\n\nP.s. If yours isn't here, send another message")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
				return
			default: // У пользователя нет шага и сообщение не команда
				var user Users // Contains only MyLang and ToLang
				err = db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&user).Error
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "err #2020, please try again later"))
					pingAdmin(err)
					return
				}
				msg, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Translating..."))
				if err != nil {
					return
				}
				messageLanguages, err := DetectLanguage(update.Message.Text)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "err #2040, please try again later"))
					pingAdmin(err)
					return
				}
				UserMessageLang := messageLanguages[0].Language
				if UserMessageLang == user.MyLang {
					translate, err := Translate(user.MyLang, user.ToLang, update.Message.Text)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "err #2090, please try again later"))
						pingAdmin(err)
						return
					}
					pp.Println(translate)
					bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, translate.Result))
				} else {
					translate, err := Translate(UserMessageLang, user.MyLang, update.Message.Text)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "err #2090, please try again later"))
						pingAdmin(err)
						return
					}
					pp.Println(translate)
					bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, translate.Result))
				}
			}

		}
	}
	if update.CallbackQuery != nil {
		switch update.CallbackQuery.Data {
		case "back":
			defer bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			var user Users
			err := db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.CallbackQuery.From.ID).Limit(1).Find(&user).Error
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "error #034, try again later"))
				pingAdmin(err)
				return
			}

			err = db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Limit(1).Updates(map[string]interface{}{"act": nil}).Error
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.From.ID, "err #333 please try later"))
				pingAdmin(err)
				return
			}
			user.MyLang = iso6391.Name(user.MyLang)
			user.ToLang = iso6391.Name(user.ToLang)
			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Your language is - <b>"+user.MyLang+"</b>, and translate language - <b>"+user.ToLang+"</b>.\n\nChange your lang /my_lang\nChange translate lang /to_lang")
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
		}
		if arr := strings.Split(update.CallbackQuery.Data, ":"); len(arr) == 2 {
			switch arr[0] {
			case "set_my_lang":
				err := db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Updates(map[string]interface{}{"act": nil, "my_lang": arr[1]}).Error
				if err != nil {
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "error #434"))
					pingAdmin(err)
					return
				}
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
				edit := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Now your language is "+iso6391.Name(arr[1]), replyMarkup)
				bot.Send(edit)
			case "set_translate_lang":
				err := db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Updates(map[string]interface{}{"act": nil, "to_lang": arr[1]}).Error
				if err != nil {
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "error #435"))
					pingAdmin(err)
					return
				}
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
				edit := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Now translate language is "+iso6391.Name(arr[1]), replyMarkup)
				bot.Send(edit)
			}
		}
	}
	//if update.InlineQuery != nil {
	//	var user Users
	//	err := db.Model(&Users{}).Where("id = ?", update.InlineQuery.From.ID).Limit(1).Find(&user).Error
	//	if err != nil {
	//
	//	}
	//
	//}
}


func Translate(source, target, text string) (TranslateAPIResponse, error) {
	source = source + "_" + strings.ToUpper(source)
	target = target + "_" + strings.ToUpper(target)
	body := map[string]string{"from": source,
		"to": target,
		"data": text,
		"platform": "api"}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return TranslateAPIResponse{}, err
	}
	req, err := http.NewRequest("POST", "https://lingvanex-translate.p.rapidapi.com/translate", bytes.NewBuffer(bodyJson))
	if err != nil {
		return TranslateAPIResponse{}, err
	}

	req.Header["x-rapidapi-key"] = []string{"561a41f76amsha7c0323d47335aep1986ecjsn4e41c7c2b518"}
	req.Header["x-rapidapi-host"] = []string{"lingvanex-translate.p.rapidapi.com"}
	req.Header["useQueryString-type"] = []string{"true"}
	req.Header["content-type"] = []string{"application/json"}
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return TranslateAPIResponse{}, err
	}
	response, err := ioutil.ReadAll(res.Body)
	var out TranslateAPIResponse
	err = json.Unmarshal(response, &out)
	return out, err
}

func DetectLanguage(text string) ([]*detectlanguage.DetectionResult, error) {
	rand.Seed(time.Now().UnixNano())
	state := detectlanguage.New(detectLangAPIKeys[rand.Intn(len(detectLangAPIKeys)-1)])
	detections, err := state.Detect(text)
	return detections, err
}

func main() {
	// Initializing PostgreSQL DB
	var err error
	db, err = gorm.Open(postgres.Open("host=ec2-63-34-97-163.eu-west-1.compute.amazonaws.com user=wzlryrrgxbgsbw password=b578bdbc77b5394a60f57660487149ca2238e0cbaf1cdbfb8b931f1168af24c7 dbname=d21k8q9pl6acl4 port=5432 TimeZone=Europe/Moscow"), &gorm.Config{SkipDefaultTransaction: true, PrepareStmt: false})
	if err != nil {
		panic(err)
	}

	// Initializing bot
	const botToken string = "1737819626:AAEoc8WyCq_8rFQcY4q0vtkhqCKro8AudfI" // 1737819626:AAEoc8WyCq_8rFQcY4q0vtkhqCKro8AudfI - @translobot && 1878391408:AAH5lvGVaRcNlFYx9sM31mwDYttx5AUR_LA - @translobetabot
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = false // >:(

	// Ports for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	detectLangAPIKeys = []string{"5ad64c8763b9293233bdc9164765037e", "c71fb63df8bdc8ea8bc4c0b20771aa5f", "11ac4f75a5b18eb919618c073c458241"}


	//updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	//for update := range updates {
	//	go botRun(&update)
	//}

	//conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	//if amqpUrl := os.Getenv("CLOUDAMQP_URL"); amqpUrl == "" {
	//	conn, err = amqp.Dial(amqpUrl)
	//}
	//defer conn.Close()

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


