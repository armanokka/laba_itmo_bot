package app

import (
	"context"
	"fmt"
	"github.com/armanokka/laba_itmo_bot/config"
	"github.com/armanokka/laba_itmo_bot/internal/controller/bot"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/repo"
	"github.com/armanokka/laba_itmo_bot/pkg/botapi"
	"github.com/armanokka/laba_itmo_bot/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func Run(cfg *config.Config) error {
	logLevel := "debug"
	if strings.ToLower(strings.TrimSpace(cfg.Environment)) == "production" {
		logLevel = "warn"
	}
	log := logger.New(logLevel)

	// Creating context which finishes on os.Interrupt, SIGINT, SIGQUIT
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1) // Handling system signals
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT)
	defer signal.Stop(c)
	go func() {
		<-c
		cancel()
	}()

	// Connecting to PostgreSQL
	log.Debug("connecting to database...")
	var dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s TimeZone=Europe/Moscow", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 dblogger.Default.LogMode(dblogger.Silent),
	})
	// TODO: set connection timeout
	if err != nil {
		return err
	}
	dbRepo, err := repo.New(db)
	if err != nil {
		return err
	}
	log.Debug("connected to database successfully")

	// Creating Telegram bot API instance
	api, err := botapi.NewBotAPIWithEndpoint(cfg.BotToken, cfg.BotAPIEndpoint)
	if err != nil {
		return err
	}

	return bot.Run(ctx, api, dbRepo, log, cfg.AdminID)
}
