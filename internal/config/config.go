package config

import (
	"database/sql"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	AdminID = 579515224
)

var (
	botToken, dashBotAPIKey                              string
	arangoHost, arangoUser, arangoPassword, arangoDBName string

	db        *gorm.DB
	arangodb  driver.Database
	analytics dashbot.DashBot
	bot       *botapi.BotAPI
	botID     int64
	once      sync.Once
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
	mustLoadEnv("TRANSLOBOT_ARANGODB_HOST", &arangoHost)
	mustLoadEnv("TRANSLOBOT_ARANGODB_USER", &arangoUser)
	mustLoadEnv("TRANSLOBOT_ARANGODB_PASSWORD", &arangoPassword)
	mustLoadEnv("TRANSLOBOT_ARANGODB_DBNAME", &arangoDBName)

	// Initializing MySQL DB
	db, err = gorm.Open(mysql.Open("translo:oEr|ea5uiKS@tcp(94.228.112.221:3306)/translo?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        5000,
		PrepareStmt:            true,
	})
	if err != nil {
		return err
	}
	var sqlDb *sql.DB
	sqlDb, err = db.DB()
	if err != nil {
		return err
	}
	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetMaxIdleConns(20)
	sqlDb.SetConnMaxLifetime(time.Hour)

	var api *tgbotapi.BotAPI
	api, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Something with bot")
		return err
	}
	bot = &botapi.BotAPI{api}
	bot.Debug = false // >:(
	bot.Buffer = 30

	me, err := bot.GetMe()
	if err != nil {
		return err
	}
	botID = me.ID

	// Initializing analytics
	analytics = dashbot.NewAPI(dashBotAPIKey, func(err error) {
		pp.Println(err)
		bot.Send(tgbotapi.NewMessage(AdminID, fmt.Sprint(err)))
	})

	var conn driver.Connection
	conn, err = http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{arangoHost},
	})
	if err != nil {
		return err
	}
	// Client object
	var client driver.Client
	client, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(arangoUser, arangoPassword),
	})
	if err != nil {
		return err
	}
	// Open "examples_books" database
	var exists bool
	exists, err = client.DatabaseExists(nil, arangoDBName)
	if err != nil {
		return err
	}
	if !exists {
		arangodb, err = client.CreateDatabase(nil, arangoDBName, &driver.CreateDatabaseOptions{})
	} else {
		arangodb, err = client.Database(nil, arangoDBName)
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

func ArangoDB() driver.Database {
	return arangodb
}

func BotID() int64 {
	return botID
}
