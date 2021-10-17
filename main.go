// t.me/translobot source code
package main

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"
)


func main() {
	// Logging
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
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
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxOpenConns(24)
	sqlDb.SetMaxIdleConns(24)
	sqlDb.SetConnMaxLifetime(15 * time.Minute)

	
	analytics = dashbot.NewAPI(DashBotAPIKey, WarnErrorAdmin)

	cronjob := cron.New()
	if err = cronjob.AddFunc("@daily", func() {
		if err = db.Model(&UsersLogs{}).Exec("DELETE FROM users_logs WHERE date < (NOW() - INTERVAL 30 DAY)").Error; err != nil {
			WarnAdmin(err)
		}
	}); err != nil {
		panic(err)
	}

	// Initializing bot
	api, err := tgbotapi.NewBotAPI(strings.TrimSpace(botToken))
	if err != nil {
		panic(err)
	}
	bot = &BotAPI{api}
	bot.Debug = false // >:(
	bot.Buffer = 30


	if _, err := os.Stat("logo.jpg"); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP)
	defer signal.Stop(stop)
	go func() {
		<-stop
		// Write profile of mem and cpu
		logrus.Info("Programm was stopped, cancelling context...")

		cancel()
	}()

	go runLogger(logs, stop, 1 * time.Minute)


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
							WarnAdmin("panic:", e.Error(), "\n\n", string(debug.Stack()))
							return
						}
						WarnAdmin("panic:", err)
						logrus.Panicln(err)
					}
				}()
				start := time.Now()
				if update.Message != nil {
					handleMessage(*update.Message)
				} else if update.CallbackQuery != nil {
					handleCallback(*update.CallbackQuery)
				} else if update.InlineQuery != nil {
					handleInline(*update.InlineQuery)
				} else if update.MyChatMember != nil {
					handleMyChatMember(*update.MyChatMember)
				}
				logrus.Print("Time spent ", time.Since(start).String())
			}()
		case <-ctx.Done():
			logrus.Info("Context was stopped, waiting for goroutines...")
			wg.Wait()
			logrus.Info("Bot was disabled.")
			return
		}
	}
}


