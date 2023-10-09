package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/laba_itmo_bot/internal/usecase"
	botapi "github.com/armanokka/laba_itmo_bot/pkg/botapi"
	"github.com/armanokka/laba_itmo_bot/pkg/errors"
	"github.com/armanokka/laba_itmo_bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

type App struct {
	bot     botapi.BotAPI
	log     logger.Logger
	repo    usecase.Repo
	adminID int64
	loc     *time.Location
}

func (app App) now() time.Time {
	return time.Now().In(app.loc).Add(3 * time.Hour)
}

func (app App) notifyAdmin(args ...interface{}) {
	text := ""
	for _, arg := range args {
		switch v := arg.(type) {
		case errors.Error:
			text += "\n" + v.Error() + "\n" + string(v.Stack())
		case error:
			text += "\n" + v.Error()
		default:
			text += "\n\n" + fmt.Sprint(arg)
		}
	}
	// TODO пересдающие должны быть в конце очереди
	_, file1, line1, _ := runtime.Caller(2)
	_, file2, line2, _ := runtime.Caller(3)
	_, file3, line3, _ := runtime.Caller(4)
	_, file4, line4, _ := runtime.Caller(5)
	text += "\n\n" + file1 + ":<b>" + strconv.Itoa(line1) + "</b>"
	text += "\n" + file2 + ":<b>" + strconv.Itoa(line2) + "</b>"
	text += "\n" + file3 + ":<b>" + strconv.Itoa(line3) + "</b>"
	text += "\n" + file4 + ":<b>" + strconv.Itoa(line4) + "</b>"
	if _, err := app.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: app.adminID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}); err != nil {
		app.log.Error("", zap.Error(err), zap.String("stack", string(debug.Stack())))
	}
}

func Run(ctx context.Context, api botapi.BotAPI, repo usecase.Repo, log logger.Logger, adminID int64) (err error) {
	app := App{
		bot:     api,
		log:     log,
		repo:    repo,
		adminID: adminID,
	}
	updatesConfig := tgbotapi.UpdateConfig{}
	updates := api.GetUpdatesChan(updatesConfig)

	app.loc, err = time.LoadLocation("UTC")
	if err != nil {
		return err
	}

	app.bot.Send(tgbotapi.NewMessage(app.adminID, "Бот запущен /start"))
	for update := range updates {
		switch {
		case update.Message != nil:
			app.onMessage(ctx, *update.Message)
		case update.CallbackQuery != nil:
			app.OnCallbackQuery(ctx, *update.CallbackQuery)
		}
	}
	return nil
}
