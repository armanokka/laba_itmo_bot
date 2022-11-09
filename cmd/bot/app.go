package bot

import (
	"context"
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/go-translo"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	"github.com/armanokka/translobot/repos"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

type App struct {
	htmlTagsRe          *regexp.Regexp
	reSpecialCharacters *regexp.Regexp
	translo             translo.API
	//deepl               translate2.Deepl
	limiter   sync.Map
	bot       *botapi.BotAPI
	log       *zap.Logger
	db        repos.BotDB
	analytics dashbot.DashBot
	bc        *bitcask.Bitcask
}

type FloodLimitation struct {
	usedAt     time.Time
	waitingNow *atomic.Bool
}

func New(bot *botapi.BotAPI, db repos.BotDB, analytics dashbot.DashBot, log *zap.Logger, bc *bitcask.Bitcask) (*App, error) {
	//iconv.Open("utf-8", "")
	app := App{
		htmlTagsRe:          regexp.MustCompile("<\\s*[^>]+>(.*?)"),
		reSpecialCharacters: regexp.MustCompile(`[[:punct:]]`),
		//deepl:      translate2.Deepl{},
		bot:       bot,
		translo:   translo.NewAPIWithClient("00aa06372bmsh6f0a2efa1a6cac4p1e73b6jsn33e40d59cbc5", &http.Client{Timeout: 30 * time.Second}),
		log:       log,
		db:        db,
		analytics: analytics,
		limiter:   sync.Map{},
		bc:        bc,
	}
	return &app, nil
}

func (app App) Run(ctx context.Context) error {
	app.bot.MakeRequest("deleteWebhook", map[string]string{})
	updates := app.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	if _, err := app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Bot have started")); err != nil {
		app.log.Error("sending bot started notification to config.AdminID", zap.Error(err))
	}
	g, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	g.Go(func() error { // бот
		for _, code := range config.BotLocalizedLangs {
			user := tables.Users{Lang: &code}
			resp, err := app.bot.Request(tgbotapi.SetMyCommandsConfig{
				Commands: []tgbotapi.BotCommand{
					{
						Command:     "start",
						Description: user.Localize("show keyboard"),
					},
					{
						Command:     "set_bot_lang",
						Description: user.Localize("set language of the bot"),
					},
				},
				Scope: &tgbotapi.BotCommandScope{
					Type: "all_private_chats",
				},
				LanguageCode: code,
			})
			if err != nil {
				app.log.Error("", zap.Error(err))
			}
			if !resp.Ok {
				out := make([]byte, 0, 256)
				if err = resp.Result.UnmarshalJSON(out); err != nil {
					app.log.Error("", zap.Error(err))
				}
				app.log.Error("couldn't set commands for the bot", zap.String("response", string(out)), zap.Int("status_code", resp.ErrorCode))
			}
		}
		for {
			select {
			case <-ctx.Done():
				app.bot.StopReceivingUpdates()
				wg.Wait()
				return ctx.Err()
			case update := <-updates:
				defer func() {
					if err := recover(); err != nil {
						app.log.Error("%w", zap.Any("error", err))
						app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
					}
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					if update.MyChatMember != nil {
						if update.MyChatMember.From.LanguageCode == "" || !in(config.BotLocalizedLangs, update.MyChatMember.From.LanguageCode) {
							update.MyChatMember.From.LanguageCode = "en"
						}
						app.onMyChatMember(*update.MyChatMember)
					} else if update.Message != nil {
						if update.Message.From.LanguageCode == "" || !in(config.BotLocalizedLangs, update.Message.From.LanguageCode) {
							update.Message.From.LanguageCode = "en"
						}
						limit, loaded := app.limiter.LoadOrStore(update.Message.Chat.ID, rate.NewLimiter(0.5, 5))
						if loaded {
							floodLimit := limit.(*rate.Limiter)
							reserve := floodLimit.Reserve()
							if !reserve.OK() || reserve.Delay() != 0 {
								user := tables.Users{ID: update.Message.From.ID, Lang: &update.Message.From.LanguageCode}
								app.bot.Send(tgbotapi.MessageConfig{
									BaseChat: tgbotapi.BaseChat{
										ChatID:              update.Message.From.ID,
										DisableNotification: true,
									},
									Text:      user.Localize("<b>Пожалуйста, не флудите!</b> Подождите 3 секунды после каждого запроса"),
									ParseMode: tgbotapi.ModeHTML,
								})
								return
							}
							app.limiter.Store(update.Message.From.ID, floodLimit)
						}
						app.onMessage(ctx, *update.Message)
					} else if update.CallbackQuery != nil {
						if update.CallbackQuery.From.LanguageCode == "" || !in(config.BotLocalizedLangs, update.CallbackQuery.From.LanguageCode) {
							update.CallbackQuery.From.LanguageCode = "en"
						}
						app.onCallbackQuery(ctx, *update.CallbackQuery)
					} else if update.InlineQuery != nil {
						if update.InlineQuery.From.LanguageCode == "" || !in(config.BotLocalizedLangs, update.InlineQuery.From.LanguageCode) {
							update.InlineQuery.From.LanguageCode = "en"
						}
						app.onInlineQuery(ctx, *update.InlineQuery)
					}
				}()

			}
		}
	})
	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				app.log.Error("", zap.Any("error", err))
				app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
			}
		}()
		keyboard := tgbotapi.InlineKeyboardMarkup{}
		k, err := app.bc.Get([]byte("mailing_keyboard_raw_text"))
		if err != nil {
			if errors.Is(err, bitcask.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		keyboard = parseKeyboard(string(k))

		app.bot.Send(tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   config.AdminID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              tgbotapi.NewRemoveKeyboard(false),
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			Text:                  "рассылка начата",
			ParseMode:             "",
			Entities:              nil,
			DisableWebPagePreview: false,
		})

		mailingMessageId, err := app.bc.Get([]byte("mailing_message_id"))
		if err != nil {
			if errors.Is(err, bitcask.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		mailingMessageIdInt, err := strconv.Atoi(string(mailingMessageId))
		if err != nil {
			return err
		}

		app.bot.Send(tgbotapi.CopyMessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   config.AdminID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              &keyboard,
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			FromChatID:          config.AdminID,
			FromChannelUsername: "",
			MessageID:           mailingMessageIdInt,
			Caption:             "",
			ParseMode:           "",
			CaptionEntities:     nil,
		})

		usersNumber, err := app.db.GetUsersNumber()
		if err != nil {
			return err
		}

		slice := make([]int64, 0, 100)
		var i int64
		for ; usersNumber/100 < i; i++ { // iterate over each 100 users
			offset := i*100 + 100
			if err = app.db.GetUsersSlice(offset, 100, slice); err != nil {
				return err
			}
			for j := 0; j < 100; j++ {
				if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:                   slice[j],
						ChannelUsername:          "",
						ReplyToMessageID:         0,
						ReplyMarkup:              &keyboard,
						DisableNotification:      false,
						AllowSendingWithoutReply: false,
					},
					FromChatID:          config.AdminID,
					FromChannelUsername: "",
					MessageID:           mailingMessageIdInt,
					Caption:             "",
					ParseMode:           "",
					CaptionEntities:     nil,
				}); err != nil {
					pp.Println(err)
				}
				app.log.Info("mailing was sent", zap.Int64("recepient_id", slice[j]), zap.Int64("queue_position", i*100+int64(j)))
			}
		}

		app.bot.Send(tgbotapi.NewMessage(config.AdminID, "рассылка закончена"))
		if err = app.bc.Delete([]byte("mailing_keyboard_raw_text")); err != nil {
			return err
		}
		if err = app.bc.Delete([]byte("mailing_message_id")); err != nil {
			return err
		}

		return nil
	})
	g.Go(func() error {
		ticker := time.NewTicker(time.Hour * 6)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				return func() error {
					f, err := os.CreateTemp("", "")
					if err != nil {
						app.log.Error("", zap.Error(err))
						return nil
					}
					defer f.Close()
					defer os.Remove(f.Name())

					users, err := app.db.GetAllUsers()
					if err != nil {
						app.log.Error("", zap.Error(err))
						return nil
					}
					for _, user := range users {
						if _, err = f.WriteString(strconv.FormatInt(user.ID, 10) + "\r\n"); err != nil {
							app.log.Error("", zap.Error(err))
							return nil
						}
					}

					resp, err := resty.New().R().SetFile("file", f.Name()).Post(fmt.Sprintf("https://api.botstat.io/create/%s/%s", app.bot.Token, config.BotstatToken))
					if err != nil {
						app.log.Error("", zap.Error(err))
						return nil
					}
					if resp.StatusCode() != 200 && resp.StatusCode() != 409 {
						app.log.Error("couldn't send the file of users", zap.Int("http_code", resp.StatusCode()), zap.String("response", resp.String()))
						return nil
					}
					app.log.Debug("sent data to botstat", zap.String("response", resp.String()), zap.Int("http_code", resp.StatusCode()))
					return nil
				}()
			}
		}
	})
	return g.Wait()
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
	text += "\n\n" + string(debug.Stack())
	if _, err := app.bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: config.AdminID,
		},
		Text:                  text,
		ParseMode:             "",
		Entities:              nil,
		DisableWebPagePreview: false,
	}); err != nil {
		app.log.Error("", zap.Error(err), zap.String("stack", string(debug.Stack())))
	}
}

func (app App) sendSpeech(user tables.Users, lang, text string, callbackID string) error {
	if user.Lang == nil {
		en := "en"
		user.Lang = &en
	}
	sdec, err := translate2.TTS(lang, text)
	if err != nil {
		if err == translate2.ErrTTSLanguageNotSupported {
			app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callbackID, user.Localize("%s не поддерживается", langs[*user.Lang][lang])))
			return nil
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "Internal error"))
		return err
	}
	app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "OK"))

	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(sdec); err != nil {
		return err
	}

	app.bot.Send(tgbotapi.AudioConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID:                   user.ID,
				ChannelUsername:          "",
				ReplyToMessageID:         0,
				ReplyMarkup:              nil,
				DisableNotification:      false,
				AllowSendingWithoutReply: false,
			},
			File: tgbotapi.FilePath(f.Name()),
		},
		Thumb:           nil,
		Caption:         "",
		ParseMode:       "",
		CaptionEntities: nil,
		Duration:        0,
		Performer:       "",
		Title:           text,
	})
	return nil
}
