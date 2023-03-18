package app

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/config"
	"github.com/armanokka/translobot/internal/controller/bot"
	"github.com/armanokka/translobot/internal/usecase/repo"
	"github.com/armanokka/translobot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) error {
	// Creating context
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1) // Handling system signals
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT)
	defer signal.Stop(c)
	go func() {
		<-c
		cancel()
	}()

	// Creating connection with PostgreSQL
	var dsn = fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s TimeZone=Europe/Moscow", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dblogger.Default.LogMode(dblogger.Silent),
	})
	if err != nil {
		return err
	}
	translationRepo, err := repo.New(db)
	if err != nil {
		return err
	}

	log := logger.New(cfg.Environment)

	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return err
	}
	
	return bot.Run(api)
}
