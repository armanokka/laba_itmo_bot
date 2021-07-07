// t.me/translobot source code
package main

import (
	"github.com/armanokka/translobot/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
}


func main() {
	// Initializing PostgreSQL DB
	var err error
	db, err = gorm.Open(postgres.Open("host=ec2-63-34-97-163.eu-west-1.compute.amazonaws.com user=wzlryrrgxbgsbw password=b578bdbc77b5394a60f57660487149ca2238e0cbaf1cdbfb8b931f1168af24c7 dbname=d21k8q9pl6acl4 port=5432 TimeZone=Europe/Moscow"), &gorm.Config{SkipDefaultTransaction: true, PrepareStmt: false})
	if err != nil {
		panic(err)
	}
	
	analytics = dashbot.NewAPI(DashBotAPIKey, nil)
	
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
	
	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	for update := range updates {
		go botRun(&update)
	}
	
	
	bot.Send(tgbotapi.NewMessage(579515224, "Bot started."))
	// requestHandler := func(ctx *fasthttp.RequestCtx) {
	// 	switch string(ctx.Path()) {
	// 	case "/" + botToken:
	// 		if isPost := ctx.IsPost(); isPost {
	// 			data := ctx.PostBody()
	// 			var update tgbotapi.Update
	// 			err := json.Unmarshal(data, &update)
	// 			if err != nil {
	// 				fmt.Fprint(ctx, "can't parse")
	// 			} else {
	// 				go botRun(&update)
	// 			}
	// 		} else {
	// 			fmt.Fprint(ctx, "no way")
	// 		}
	// 	default:
	// 		_, err = fmt.Fprintln(ctx, "ok")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }
	// if err = fasthttp.ListenAndServe(":"+port, requestHandler); err != nil {
	// 	log.Fatalf("Error in ListenAndServe: %s", err)
	// }

}


