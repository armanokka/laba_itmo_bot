// t.me/translobot source code
package main

import (
	"encoding/json"
	"fmt"
	"github.com/armanokka/translobot/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)


// botRun is main handler of bot
func botRun(update *tgbotapi.Update) {
	if update.Message != nil {
		handleMessage(update)
	} else if update.CallbackQuery != nil {
		handleCallback(update)
	} else if update.InlineQuery != nil {
		handleInline(update)
	}
	// f, err := os.Create("mem.out")
	// if err != nil {
	// 	panic(err)// !
	// }
	// if err := pprof.WriteHeapProfile(f); err != nil {
	// 	panic(err) // !
	// }
}


func main() {
	// Initializing PostgreSQL DB
	var err error
	db, err = gorm.Open(mysql.Open("f0568401_user:NlEbEgda@tcp(141.8.193.236:3306)/f0568401_user?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{SkipDefaultTransaction: true, PrepareStmt: false})
	if err != nil {
		panic(err)
	}
	
	analytics = dashbot.NewAPI(DashBotAPIKey, WarnErrorAdmin)
	
	// Initializing bot
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot = &BotAPI{api}
	bot.Debug = false // >:(

	// Ports for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	
	if _, err = os.Stat("ad.jpg"); err != nil {
		panic(err)
	}
	
	// bot.Buffer = 0
	// updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	// for update := range updates {
	// 	go botRun(&update)
	// }
	
	
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


