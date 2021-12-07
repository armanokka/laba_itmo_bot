package bot

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	"github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"runtime/debug"
	"sync"
	"time"
)

type app struct {
	bot *botapi.BotAPI
	db        *gorm.DB
	analytics dashbot.DashBot
	log       *zap.SugaredLogger
	logs chan tables.UsersLogs
	messageState *sync.Map
	mailer mailer
}

type mailer struct {
	stop context.CancelFunc
}

func (app app) Run(ctx context.Context) error {
	app.bot.MakeRequest("deleteWebhook", map[string]string{})
	updates := app.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Bot have started"))
	g, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	g.Go(func() error {
		for {
			defer func() {
				if err := recover(); err != nil {
					app.log.Error(err)
					app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:" + fmt.Sprint(err)))
				}
			}()
			select {
			case <-ctx.Done():
				wg.Wait()
				return ctx.Err()
			case update := <-updates:
				wg.Add(1)
				go func() {
					defer wg.Done()
					start := time.Now()
					if update.Message != nil {
						if app.messageState == nil {
							app.messageState = new(sync.Map)
						}

						if f, ok := app.messageState.Load(update.Message.From.ID); ok {
							fu, ok := f.(func(message tgbotapi.Message))
							if !ok {
								pp.Printf("Error when casting type in messageState: %v", fu)
							}
							fu(*update.Message)
							return
						}
						app.onMessage(ctx, *update.Message)
					} else if update.CallbackQuery != nil {
						app.onCallbackQuery(*update.CallbackQuery)
					} else if update.InlineQuery != nil {
						app.onInlineQuery(*update.InlineQuery)
					} else if update.MyChatMember != nil {
						app.onMyChatMember(*update.MyChatMember)
					}
					pp.Println("Time spent:", start.String())
				}()

			}
		}
	})
	g.Go(func() error {
		return app.runLogsListener(ctx, time.Minute)
	})
	return g.Wait()
}

func (app app) runLogsListener(ctx context.Context, pushDbInterval time.Duration) error {
	logrus.Info("Message logger has been started")
	if app.logs == nil {
		app.logs = make(chan tables.UsersLogs, 100)
	}
	inserts := make([]tables.UsersLogs, 0, cap(app.logs))
	ticker := time.NewTicker(pushDbInterval)
	for {
		select {
		case <-ctx.Done():
			close(app.logs)
			for log := range app.logs {
				inserts = append(inserts, log)
			}
			if err := app.db.Create(inserts).Error; err != nil {
				app.notifyAdmin(err)
			}
			logrus.Info("Message logger was stopped.")
			return ctx.Err()
		case <-ticker.C:
			for i := 0; i < len(app.logs); i++ {
				inserts = append(inserts, <-app.logs)
			}
			if err := app.db.Create(inserts).Error; err != nil {
				app.notifyAdmin(err)

				for _, insert := range inserts {
					user := app.loadUser(insert.ID, func(err error) {
						app.notifyAdmin(err)
					})
					if !user.Exists() {
						if err = app.db.Create(&tables.Users{
							ID:         insert.ID,
							MyLang:     "en",
							ToLang:     "es",
							Act:        sql.NullString{},
							Mailing:    true,
							Usings:     0,
							Lang:       "en",
							ReferrerID: 0,
							Blocked:    false,
						}).Error; err != nil {
							app.notifyAdmin(err)
						}
					}
				}
			}
			inserts = nil
			logrus.Info("message logs were saved")
		}
	}
}

func (app app) notifyAdmin(args ...interface{}) {
	text := ""
	for _, arg := range args {
		switch v := arg.(type) {
		case error:
			text += "\n" + v.Error() + "\n\n<code>" + string(debug.Stack()) + "</code>"
		default:
			text += "\n\n" + fmt.Sprint(arg)
		}
	}
	if _, err := app.bot.Send(tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{
			ChatID:                   config.AdminID,
		},
		Text:                  text,
		ParseMode:             tgbotapi.ModeHTML,
		Entities:              nil,
		DisableWebPagePreview: false,
	}); err != nil {
		app.log.Error(err)
	}
}

func (app app) setMyCommands(langs []string, commands []tgbotapi.BotCommand) error {
	newCommands := make(map[string][]tgbotapi.BotCommand)
	for _, lang := range langs {
		newCommands[lang] = []tgbotapi.BotCommand{}
		for _, command := range commands {
			tr, err := translate.GoogleHTMLTranslate("en", lang, command.Description)
			if err != nil {
				return err
			}
			newCommands[lang]= append(newCommands[lang], tgbotapi.BotCommand{
				Command:     command.Command,
				Description: tr.Text,
			})
		}
	}

	for lang, command := range newCommands {
		data, err := json.Marshal(command)
		if err != nil {
			return err
		}
		params := tgbotapi.Params{}
		params.AddNonEmpty("commands", string(data))
		params.AddNonEmpty("language_code", lang)
		if _, err = app.bot.MakeRequest("setMyCommands", params); err != nil {
			return err
		}
	}
	return nil
}

func (app app) loadUser(id int64, warn func(err error)) User {
	return User{
		Users: tables.Users{ID: id},
		error: warn,
		bot: app.bot,
		db:    app.db,
	}
}

func (app app) writeBotLog(id int64, intent string, text string) {
	app.logs <- tables.UsersLogs{
		ID:      id,
		Intent:  sql.NullString{
			String: intent,
			Valid:  true,
		},
		Text:    text,
		FromBot: true,
		Date:    time.Now(),
	}
}

func (app app) writeUserLog(id int64, text string) {
	app.logs <- tables.UsersLogs{
		ID:      id,
		Intent:  sql.NullString{
			String: "",
			Valid:  false,
		},
		Text:    text,
		FromBot: false,
		Date:    time.Now(),
	}
}

func (app *app) onNextUserMessage(id int64, f func(msg tgbotapi.Message)) {
	if app.messageState == nil {
		app.messageState = new(sync.Map)
	}
	app.messageState.Store(id, f)
}

func (app *app) stopUserConversation(id int64) {
	app.messageState.Delete(id)
}