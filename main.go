// t.me/translobot source code
package main

import (
	"encoding/json"
	"fmt"
	"github.com/armanokka/translobot/translate"
	iso6391 "github.com/emvi/iso-639-1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	db  *gorm.DB
	bot *tgbotapi.BotAPI
)

type Users struct {
	ID     int64 `gorm:"primaryKey;index;not null"`
	MyLang string `gorm:"default:en"`
	ToLang string `gorm:"default:ar;n"`
	Act, Engine    string
}

type Errors struct {
	UserID int64
	Code int
	Time time.Time
}

// botRun is main handler of bot
func botRun(update *tgbotapi.Update) {
	// warn is error handler
	warn := func(code int, err error) {
		errText := "Error #" + strconv.Itoa(code) + ". Please, try again later and PM @armanokka.\n\nI'll fix the bug in near future, please don't block the bot."
		if update == nil {
			fmt.Println("update is nil")
			return
		}
		var userID int64
		if update.Message != nil {
			userID = update.Message.Chat.ID
			bot.Send(tgbotapi.NewMessage(userID, errText))
		} else if update.CallbackQuery != nil {
			userID = update.CallbackQuery.From.ID
			bot.Send(tgbotapi.NewMessage(userID, errText))
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, errText))
		} else if update.InlineQuery != nil {
			userID = update.InlineQuery.From.ID
			bot.AnswerInlineQuery(tgbotapi.InlineConfig{
				InlineQueryID:     update.InlineQuery.ID,
				Results:           nil,
				CacheTime:         0,
				IsPersonal:        true,
				SwitchPMText:      "Error. Try again later",
				SwitchPMParameter: "from_inline",
			})
		}
		bot.Send(tgbotapi.NewMessage(579515224, fmt.Sprintf("Error [%v]: %v", code, err)))
		
		err = db.Create(&Errors{
			UserID: userID,
			Code:   code,
			Time:   time.Now(),
		}).Error
		if err != nil {
			bot.Send(tgbotapi.NewMessage(579515224, fmt.Sprintf("Couldn't create new row in Errors relation, so this user %v generated error with this code %v now.", userID, code)))
		}
	}

	if update.Message != nil {
		switch update.Message.Text {
		case "/start", "/start from_inline":
			var user = Users{ID: update.Message.From.ID}
			if update.Message.From.LanguageCode == "" {
				update.Message.From.LanguageCode = "en"
			}
			err := db.FirstOrCreate(&user, &Users{ID: update.Message.From.ID}).Error
			if err != nil {
				warn(1002, err)
				return
			}
			user.MyLang = iso6391.Name(user.MyLang)
			user.ToLang = iso6391.Name(user.ToLang)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your language is - <b>"+user.MyLang+"</b>, and translate language - <b>"+user.ToLang+"</b>.\n\nNeed help? /help\nChange your lang /my_lang\nChange translate lang /to_lang")
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
				warn(1003, err)
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
				warn(032, err)
				return
			}
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "<b>What can this bot do?</b>\n▫️ Translo allows you to translate your messages into over than 100 languages. (117)\n<b>How to translate message?</b>\n▫️ Firstly, you have to setup your lang (default: English), then setup translate lang (default; Arabic) then send text messages and bot will translate them quickly.\n<b>How to setup my lang?</b>\n▫️ Send /my_lang then send any message <b>IN YOUR LANGUAGE</b>. Bot will detect and suggest you some variants. Select your lang. Done.\n<b>How to setup translate lang?</b>\n▫️ Send /to_lang then send any message <b>IN LANGUAGE YOU WANT TRANSLATE</b>. Bot will detect and suggest you some variants. Select your lang. Done.\n<b>I have a suggestion or I found bug!</b>\n▫️ 👉 Contact me pls - @armanokka")
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
		case "/engine":
			var user Users // Contains only MyLang and ToLang
			err := db.Model(&Users{}).Select("engine").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&user).Error
			if err != nil {
				warn(2020, err)
				return
			}
			n := map[string]string{"google":"", "yandex":""}
			if _, ok := n[user.Engine]; !ok {
				warn(50, err)
				return
			}
			n[user.Engine] = " ✅"

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You can change translate engine")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Google" + n["google"], "switch_engine:google")),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Yandex" + n["yandex"], "switch_engine:yandex")),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Back", "back")))
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			return
		default: // Сообщение не является командой.
			userStep, err := getUserStep(update.Message.Chat.ID)
			if err != nil {
				warn(37289, err)
				return
			}
			switch userStep {
			case "set_my_lang":
				lowerUserMsg := strings.ToLower(update.Message.Text)

				if nameOfLang := iso6391.Name(lowerUserMsg); nameOfLang != "" { // Юзер отправил сразу код языка
					err := db.Model(&Users{}).Where("id", update.Message.Chat.ID).Updates(map[string]interface{}{"act": nil, "my_lang": lowerUserMsg}).Error
					if err != nil {
						warn(300, err)
						return
					}
					replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+nameOfLang)
					msg.ReplyMarkup = replyMarkup
					bot.Send(msg)
					return
				}

				if codeOfLang := iso6391.CodeForName(strings.Title(lowerUserMsg)); codeOfLang != "" { // Юзер полное его название языка
					err := db.Model(&Users{}).Where("id", update.Message.Chat.ID).Updates(map[string]interface{}{"act": nil, "my_lang": codeOfLang}).Error
					if err != nil {
						warn(301, err)
						return
					}
					replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(codeOfLang))
					msg.ReplyMarkup = replyMarkup
					bot.Send(msg)
					return
				}

				langDetects, err := translate.DetectLanguageYandex(update.Message.Text)
				if err != nil || langDetects.Lang == "" {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Could not detect language, please send something else again"))
					return
				}
				err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Updates(map[string]interface{}{"act": nil, "my_lang": langDetects.Lang}).Error
				if err != nil {
					warn(302, err)
					return
				}
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(langDetects.Lang))
				msg.ReplyMarkup = replyMarkup
				bot.Send(msg)
				return
			case "set_translate_lang":
				lowerUserMsg := strings.ToLower(update.Message.Text)

				if nameOfLang := iso6391.Name(lowerUserMsg); nameOfLang != "" { // Юзер отправил сразу код языка
					err := db.Model(&Users{}).Where("id", update.Message.Chat.ID).Updates(map[string]interface{}{"act": nil, "to_lang": lowerUserMsg}).Error
					if err != nil {
						warn(303, err)
						return
					}
					replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+nameOfLang)
					msg.ReplyMarkup = replyMarkup
					bot.Send(msg)
					return
				}

				if codeOfLang := iso6391.CodeForName(strings.Title(lowerUserMsg)); codeOfLang != "" { // Юзер полное его название языка
					err := db.Model(&Users{}).Where("id", update.Message.Chat.ID).Updates(map[string]interface{}{"act": nil, "to_lang": codeOfLang}).Error
					if err != nil {
						warn(304, err)
						return
					}
					replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(codeOfLang))
					msg.ReplyMarkup = replyMarkup
					bot.Send(msg)
					return
				}

				langDetects, err := translate.DetectLanguageYandex(update.Message.Text)
				if err != nil || langDetects.Lang == "" {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Could not detect language, please send something else again"))
					return
				}
				err = db.Model(&Users{}).Where("id", update.Message.Chat.ID).Updates(map[string]interface{}{"act": nil, "to_lang": langDetects.Lang}).Error
				if err != nil {
					warn(305, err)
					return
				}
				replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Now language is "+iso6391.Name(langDetects.Lang))
				msg.ReplyMarkup = replyMarkup
				bot.Send(msg)
				return
			default: // У пользователя нет шага и сообщение не команда
				var user Users // Contains only MyLang and ToLang
				err = db.Model(&Users{}).Select("my_lang", "to_lang", "engine").Where("id = ?", update.Message.Chat.ID).Limit(1).Find(&user).Error
				if err != nil {
					warn(306, err)
					return
				}
				msg, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ Translating..."))
				if err != nil {
					return
				}

				var text = update.Message.Text
				if update.Message.Caption != "" {
					text = update.Message.Caption
				}
				if text == "" {
					bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Please, send text message"))
					return
				}
				
				
				langDetects, err := translate.DetectLanguageYandex(text)
				if err != nil {
					if e, ok := err.(translate.HTTPError); ok {
						if e.Code == 413 {
							bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Too big text"))
							return
						}
					} else {
						warn(308, err)
					}
					return
				}

				if langDetects.Lang == "" {
					bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, text))
					return
				}
				var to string // language into need to translate
				if langDetects.Lang == user.ToLang {
					to = user.MyLang
				} else {
					to = user.ToLang
				}
				var translatedText string
				switch user.Engine {
				case "google":
					tr, err := translate.TranslateGoogle(langDetects.Lang, to, text)
					if err != nil {
						if e, ok := err.(translate.HTTPError); ok {
							if e.Code == 413 {
								bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Too big text"))
								return
							}
						}
						warn(309, err)
						return
					}
					if tr == "" {
						bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Empty result"))
						return
					}
					pp.Println(tr)
				case "yandex":
					tr, err := translate.TranslateYandex(langDetects.Lang, to, text)
					if err != nil {
						if e, ok := err.(translate.HTTPError); ok {
							if e.Code == 413 {
								bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Too big text"))
								return
							}
						}
						warn(310, err)
						return
					}
					if len(tr.Text) == 0 {
						bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Could not translate message"))
						return
					}
					if strings.Join(strings.Fields(tr.Text[0]), "") == "" { // No words
						bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, "Empty result"))
						return
					}
					pp.Println(tr.Text[0])
				default:
					warn(311, err)
					return
				}
				bot.Send(tgbotapi.NewEditMessageText(update.Message.Chat.ID, msg.MessageID, translatedText))
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
				warn(312, err)
				return
			}

			err = db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Limit(1).Updates(map[string]interface{}{"act": nil}).Error
			if err != nil {
				warn(313, err)
				return
			}
			user.MyLang = iso6391.Name(user.MyLang)
			user.ToLang = iso6391.Name(user.ToLang)
			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Your language is - <b>"+user.MyLang+"</b>, and translate language - <b>"+user.ToLang+"</b>.\n\nChange your lang /my_lang\nChange translate lang /to_lang")
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
		case "none":
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
		}

		arr := strings.Split(update.CallbackQuery.Data, ":")
		if len(arr) == 0 {
			return
		}
		switch arr[0] {
		case "set_my_lang": // arr[1] - language code
			err := db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Updates(map[string]interface{}{"act": nil, "my_lang": arr[1]}).Error
			if err != nil {
				warn(314, err)
				return
			}
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
			edit := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Now your language is "+iso6391.Name(arr[1]), replyMarkup)
			bot.Send(edit)
		case "set_translate_lang": // arr[1] - language code
			err := db.Model(&Users{}).Where("id", update.CallbackQuery.From.ID).Updates(map[string]interface{}{"act": nil, "to_lang": arr[1]}).Error
			if err != nil {
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "error #435"))
				warn(315, err)
				return
			}
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			replyMarkup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", "back")))
			edit := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "Now translate language is "+iso6391.Name(arr[1]), replyMarkup)
			bot.Send(edit)
		case "switch_engine":
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			err := db.Model(&Users{}).Where("id = ?", update.CallbackQuery.From.ID).Update("engine", arr[1]).Error
			if err != nil {
				warn(316, err)
				return
			}
			var user Users // Contains only MyLang and ToLang
			err = db.Model(&Users{}).Select("engine").Where("id = ?", update.CallbackQuery.From.ID).Limit(1).Find(&user).Error
			if err != nil {
				warn(317, err)
				return
			}
			n := map[string]string{"google":"", "yandex":""}
			if _, ok := n[user.Engine]; !ok {
				warn(318, err)
				return
			}
			n[user.Engine] = " ✅"
			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, "You can change translate engine")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Google" + n["google"], "switch_engine:google")),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Yandex" + n["yandex"], "switch_engine:yandex")),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Back", "back")))
			msg.ReplyMarkup = &keyboard
			bot.Send(msg)
		}

	}
	//if update.InlineQuery != nil {
	//	if strings.Join(strings.Fields(update.InlineQuery.Query), " ") == "" {
	//		bot.AnswerInlineQuery(tgbotapi.InlineConfig{
	//			InlineQueryID:     update.InlineQuery.ID,
	//			Results:           nil,
	//			CacheTime:         0,
	//			IsPersonal:        false,
	//			NextOffset:        "",
	//			SwitchPMText:      "Type correct query",
	//			SwitchPMParameter: "from_inline",
	//		})
	//		return
	//	}
	//	langDetects, err := DetectLanguage(update.InlineQuery.Query)
	//	if err != nil {
	//		warn("-09", err)
	//		return
	//	}
	//	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔽", "open_translate_langs_from_inline_query")))
	//	if langDetects.Lang == update.InlineQuery.From.LanguageCode { // Юзер отправил текст на своем языке. Языке платформы
	//		result := tgbotapi.InlineQueryResultArticle{
	//			Type:                "article",
	//			ID:                  "1",
	//			Title:               "Buttons",
	//			ReplyMarkup:         &keyboard,
	//			Description:         update.InlineQuery.Query,
	//		}
	//
	//		bot.AnswerInlineQuery(tgbotapi.InlineConfig{
	//			InlineQueryID:     update.CallbackQuery.ID,
	//			Results:           []interface{}{result},
	//			CacheTime:         15,
	//			IsPersonal:        false,
	//		})
	//	} else { // Юзер отправил текст НЕ на своем языке
	//		transl, err := Translate(langDetects.Lang, update.InlineQuery.From.LanguageCode, update.InlineQuery.Query)
	//		if err != nil {
	//			warn("-44", err)
	//			return
	//		}
	//
	//		result := tgbotapi.InlineQueryResultArticle{
	//			Type:                "article",
	//			ID:                  "1",
	//			Title:               "Buttons",
	//			ReplyMarkup:         &keyboard,
	//			Description:         transl.Text[0],
	//		}
	//		bot.AnswerInlineQuery(tgbotapi.InlineConfig{
	//			InlineQueryID:     update.CallbackQuery.ID,
	//			Results:           []interface{}{result},
	//			CacheTime:         15,
	//			IsPersonal:        false,
	//		})
	//	}
	//
	//
	//	bot.AnswerInlineQuery(tgbotapi.InlineConfig{
	//		InlineQueryID:     update.CallbackQuery.ID,
	//		Results:           nil,
	//		CacheTime:         0,
	//		IsPersonal:        false,
	//		NextOffset:        "",
	//		SwitchPMText:      "",
	//		SwitchPMParameter: "",
	//	})
	//}
}


func setUserStep(chatID int64, step string) error {
	return db.Model(&Users{ID: chatID}).Where("id = ?", chatID).Limit(1).Update("act", step).Error
}

func getUserStep(chatID int64) (string, error) {
	var user Users
	err := db.Model(&Users{ID: chatID}).Select("act").Where("id", chatID).Limit(1).Find(&user).Error
	return user.Act, err
}

func main() {
	// Initializing PostgreSQL DB
	var err error
	db, err = gorm.Open(postgres.Open("host=ec2-63-34-97-163.eu-west-1.compute.amazonaws.com user=wzlryrrgxbgsbw password=b578bdbc77b5394a60f57660487149ca2238e0cbaf1cdbfb8b931f1168af24c7 dbname=d21k8q9pl6acl4 port=5432 TimeZone=Europe/Moscow"), &gorm.Config{SkipDefaultTransaction: true, PrepareStmt: false})
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

	// Ports for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	//updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	//for update := range updates {
	//	go botRun(&update)
	//}

	//conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	//if amqpUrl := os.Getenv("CLOUDAMQP_URL"); amqpUrl == "" {
	//	conn, err = amqp.Dial(amqpUrl)
	//}
	//defer conn.Close()
	
	bot.Send(tgbotapi.NewMessage(579515224, "Bot started."))
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


