// t.me/translobot source code
package main

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)


func main() {
	// Logging
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	// Profiling
	go func() {
		r := mux.NewRouter()
		r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		r.PathPrefix("/debug/").Handler(http.DefaultServeMux)
		// Ports for Heroku
		port := os.Getenv("PORT")
		if port == "" {
			port = "80"
		}
		r.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "ok")
		}))
		if err := http.ListenAndServe(":" + port, r); err != nil {
			panic(err)
		}
	}()

	// Initializing PostgreSQL DB
	var err error
	db, err = gorm.Open(mysql.Open("f0568401_user:NlEbEgda@tcp(141.8.193.236:3306)/f0568401_user?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxOpenConns(25)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(6 * time.Hour)

	
	analytics = dashbot.NewAPI(DashBotAPIKey, WarnErrorAdmin)
	
	// Initializing bot
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot = &BotAPI{api}
	bot.Debug = false // >:(
	bot.Buffer = 1


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
		<-c
		cancel()
	}()

	var wg sync.WaitGroup
	go func() {
		bot.Send(tgbotapi.NewMessage(AdminID, "Bot was started."))
	}()
	_, _ = bot.MakeRequest("deleteWebhook", tgbotapi.Params{})
	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	for {
		select {
		case update, more := <-updates:
			if !more {
				wg.Wait()
				break
			}
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Защита от паники
				defer func() {
					if err := recover(); err != nil {
						if e, ok := err.(error); ok {
							WarnAdmin("panic:", e.Error())
							return
						}
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
			}()
		case <-ctx.Done():
			wg.Wait()
			break
		}
	}
}


