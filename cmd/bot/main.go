package bot

import (
	"context"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"sync"
)

func Run(ctx context.Context, bot *botapi.BotAPI, db *gorm.DB, analytics dashbot.DashBot, logger *zap.SugaredLogger) error {
	if _, err := os.Stat("logo.jpg"); err != nil {
		return err
	}
	app := app{
		bot:       bot,
		db:        db,
		analytics: analytics,
		log:       logger,
		logs: make(chan tables.UsersLogs, 100),
		messageState: new(sync.Map),
	}
	return app.Run(ctx)
}