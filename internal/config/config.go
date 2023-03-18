package config

import (
	"fmt"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
)

var (
	botToken, dashBotAPIKey string
	dsn                     string
	db                      *gorm.DB
	analytics               dashbot.DashBot
	bot                     *botapi.BotAPI
	once                    sync.Once
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

func Load() (err error) {
	once.Do(func() {
		err = load()
	})
	return
}

func mustLoadEnv(name string, v *string) {
	str := strings.TrimSpace(os.Getenv(name))
	if str == "" {
		panic("$" + name + " is empty")
	}
	*v = str
}

func load() (err error) {
	mustLoadEnv("TRANSLOBOT_TOKEN", &botToken)
	mustLoadEnv("TRANSLOBOT_DASHBOT_TOKEN", &dashBotAPIKey)
	mustLoadEnv("TRANSLOBOT_DSN", &dsn)

	// Initializing MySQL DB
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return err
	}

	if err = db.AutoMigrate(&tables.Users{}); err != nil {
		return err
	}

	analytics = dashbot.NewAPI(dashBotAPIKey)

	var api *tgbotapi.BotAPI
	api, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Something with bot")
		return err
	}
	bot = &botapi.BotAPI{api}
	bot.Debug = false
	bot.Buffer = 30

	localAPIPort := os.Getenv("TRANSLOBOT_LOCAL_API_PORT")
	if localAPIPort == "" {
		localAPIPort = "8081"
	}
	// Setting local bot api endpoint if it exists
	bot.SetAPIEndpoint("http://localhost:" + localAPIPort + "/bot%s/%s")
	if _, err = bot.GetMe(); err == nil {
		fmt.Println("Detected Local Bot API. Sending requests there")
		return nil
	}
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)
	if _, err = bot.GetMe(); err != nil {
		return err
	}
	fmt.Println("Local Bot API not detected. Using api.telegram.org")
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
