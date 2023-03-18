package bot

import (
	"github.com/armanokka/translobot/internal/usecase/repo"
	"github.com/armanokka/translobot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api  *tgbotapi.BotAPI
	log  logger.Logger
	repo repo.TranslationRepo
}

func Run(api *tgbotapi.BotAPI, repo repo.TranslationRepo, log logger.Logger) error {
	app := Bot{
		api:  api,
		log:  log,
		repo: repo,
	}
	updatesConfig := tgbotapi.UpdateConfig{}
	updates := api.GetUpdatesChan(updatesConfig)

	for update := range updates {
		switch {
		case update.Message != nil:
			app.OnMessage(*update.Message)
		case update.CallbackQuery != nil:
			app.OnCallbackQuery(*update.CallbackQuery)
		}
	}
	return nil
}
