// t.me/translobot source code
package main

import (
	"context"
	"github.com/armanokka/translobot/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)


// botRun is main handler of bot
func botRun(update *tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			WarnAdmin("panic:", err)
		}
	}()
	if update.Message != nil {
		handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		handleCallback(update.CallbackQuery)
	} else if update.InlineQuery != nil {
		handleInline(update.InlineQuery)
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

	if _, err := os.Stat("logo.jpg"); err != nil {
		panic(err)
	}

	loc, err = time.LoadLocation(TimeLocation)
	if err != nil {
		panic(err) // проверяем на валидность константы TimeLocation
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		select {
		case <-c:
			cancel()
		}
	}()

	var wg sync.WaitGroup

	bot.MakeRequest("deleteWebhook", tgbotapi.Params{})
	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	for {
		select {
		case update := <-updates:
			wg.Add(1)
			go func() {
				defer wg.Done()
				botRun(&update)
			}()
		case <-ctx.Done():
			wg.Wait()
			break
		}
	}
}


