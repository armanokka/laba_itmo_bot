package config

import (
	"fmt"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/robfig/cron"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

const (
	DashBotAPIKey        = "cjVjdWDRijXDk5kl9yGi5TTS9XImME7HbZMOg09F"
	AdminID              = 579515224
	botToken      string = "1737819626:AAEJyD8fnSHdkh6VP3ePdwFkpEnrirLMHp4" //
)

var (
	db        *gorm.DB
	analytics dashbot.DashBot
	bot       *botapi.BotAPI
)

var BotLocalizedLangs = []string{
	"en", "ru", "de", "es", "uk", "uz", "id", "it", "pt", "ar",
}

var WitAPIKeys = map[string]string{
	"en": "6X4I6P3TLPAW7EBOQKOIET7NHJYJ3TQ3",
	"ru": "XARRYZ2OGP7UPXZJG5MJOL2FTJMHFW74",
	"es": "KQIFMTDYRPS3POH3J5QK2AK3E4GDCEHB",
	"uk": "X3YGCUD5TKZMJYLD3JG7SMF2HYIHDQAW",
	"pt": "4ID2IR4RTLFRBUSPGTVHDMDIEBRESLRA",
	"uz": "BBGPN3S6RF6LTK4Y4UW46D4IGISJSKHW", // beta
	"id": "5M22F2VA4Z5HKA336VRT5EUETLWLHETV", // beta
	"it": "PVN465FPYP5BHUFD3DCSUR7EGNBYG57J", // beta
}

func Load() error {

	// Initializing PostgreSQL DB
	var err error
	db, err = gorm.Open(mysql.Open("f0568401_user:NlEbEgda@tcp(141.8.193.236:3306)/f0568401_user?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        1000,
		PrepareStmt:            true,
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

	// Initializing bot
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		token = strings.TrimSpace(botToken)
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot = &botapi.BotAPI{api}
	bot.Debug = false // >:(
	bot.Buffer = 30

	if _, err = bot.GetMe(); err != nil {
		return err
	}

	// Initializing analytics
	analytics = dashbot.NewAPI(DashBotAPIKey, func(err error) {
		pp.Println(err)
		bot.Send(tgbotapi.NewMessage(AdminID, fmt.Sprint(err)))
	})

	// Running cron job
	cronjob := cron.New()
	if err = cronjob.AddFunc("@daily", func() {
		if err = db.Model(&tables.UsersLogs{}).Exec("DELETE FROM users_logs WHERE date < (NOW() - INTERVAL 30 DAY)").Error; err != nil {
			pp.Println(err)
			bot.Send(tgbotapi.NewMessage(AdminID, fmt.Sprint(err)+"\n\nIT'S FROM config.go"))
		}
	}); err != nil {
		panic(err)
	}

	return nil
}

func DB() *gorm.DB {
	return db
}

func BotAPI() *botapi.BotAPI {
	return bot
}

func Analytics() dashbot.DashBot {
	return analytics
}
