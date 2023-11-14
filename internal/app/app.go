package app

import (
	"context"
	"fmt"
	"github.com/armanokka/laba_itmo_bot/config"
	"github.com/armanokka/laba_itmo_bot/internal/controller/bot"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/repo"
	"github.com/armanokka/laba_itmo_bot/pkg/botapi"
	"github.com/armanokka/laba_itmo_bot/pkg/logger"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) error {
	log := logger.New(cfg.Environment)
	log.Info("displaying config")
	pp.Println(cfg)

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

	log.Info("connecting to database...")
	var dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s TimeZone=Europe/Moscow", cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
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
	log.Info("connected to database successfully")

	api, err := botapi.NewBotAPIWithEndpoint(cfg.BotToken, cfg.BotAPIEndpoint)
	if err != nil {
		return err
	}
	log.Debug("", zap.Any("config", cfg))

	return bot.Run(ctx, api, translationRepo, log, cfg.AdminID)
}
