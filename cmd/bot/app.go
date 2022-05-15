package bot

import (
	"context"
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/arangodb/go-driver"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	"github.com/armanokka/translobot/pkg/errors"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	"github.com/armanokka/translobot/repos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"sync"
)

//TODO: заменить fuzzywuzzy и отпрофилировать бота, чтобы убрать утечки памяти

type App struct {
	htmlTagsRe          *regexp.Regexp
	reSpecialCharacters *regexp.Regexp
	deepl               translate2.Deepl
	bot                 *botapi.BotAPI
	log                 *zap.Logger
	db                  repos.BotDB
	analytics           dashbot.DashBot
	bc                  *bitcask.Bitcask

	arangodb driver.Database
	cache    driver.Collection
} // TODO: сократить объем lingvo

func New(bot *botapi.BotAPI, db repos.BotDB, analytics dashbot.DashBot, log *zap.Logger, bc *bitcask.Bitcask /* deepl translate2.Deepl*/, arangodb driver.Database) (App, error) {
	app := App{
		arangodb:            arangodb,
		htmlTagsRe:          regexp.MustCompile("<\\s*[^>]+>(.*?)"),
		reSpecialCharacters: regexp.MustCompile(`[[:punct:]]`),
		//deepl:      translate2.Deepl{},
		bot:       bot,
		log:       log,
		db:        db,
		analytics: analytics,
		bc:        bc,
	}
	if err := app.prepareCollections(); err != nil {
		return App{}, err
	}
	return app, nil
}

func (app *App) prepareCollections() error {
	cacheEnabled := true
	options := &driver.CreateCollectionOptions{CacheEnabled: &cacheEnabled, WaitForSync: true}

	cache, err := app.createCollectionIfNotExists("cache", options) // IT MUST HAVE UNIQUE INDEX ON from, to and original_text!
	if err != nil {
		return err
	}
	app.cache = cache
	//if _, _, err = app.cache.EnsureTTLIndex(nil, "createdAt", 60*60*24, nil); err != nil {
	//	return err
	//}
	return nil
}

func (app App) createCollectionIfNotExists(collection string, options *driver.CreateCollectionOptions) (driver.Collection, error) {
	col, err := app.arangodb.Collection(nil, collection)
	if err != nil {
		if driver.IsNotFound(err) {
			col, err = app.arangodb.CreateCollection(nil, collection, options)
			if err != nil {
				return nil, err
			}
		}
	}
	return col, nil
}

func (app App) Run(ctx context.Context) error {
	app.bot.MakeRequest("deleteWebhook", map[string]string{})
	updates := app.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Bot have started"))
	g, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	g.Go(func() error { // бот
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
					if update.Message != nil {
						if update.Message.From.LanguageCode == "" || !in(config.BotLocalizedLangs, update.Message.From.LanguageCode) {
							update.Message.From.LanguageCode = "en"
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
					} else if update.MyChatMember != nil {
						if update.MyChatMember.From.LanguageCode == "" || !in(config.BotLocalizedLangs, update.MyChatMember.From.LanguageCode) {
							update.MyChatMember.From.LanguageCode = "en"
						}
						app.onMyChatMember(*update.MyChatMember)
					}
				}()

			}
		}
	})
	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				app.log.Error("%w", zap.Any("error", err))
				app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
			}
		}()
		fmt.Println("чекаем не завалялась ли рассылка")
		exists, err := app.db.MailingExists()
		if err != nil {
			return err
		}
		if exists {
			pp.Println("рассылка есть, продолжаю")

			mailingMessageIDBytes, err := app.bc.Get([]byte("mailing_message_id"))
			if err != nil {
				return err
			}
			mailingMessageID, err := strconv.Atoi(string(mailingMessageIDBytes))
			if err != nil {
				return err
			}
			mailing_keyboard_raw_text, err := app.bc.Get([]byte("mailing_keyboard_raw_text"))
			if err != nil && !errors.Is(err, bitcask.ErrKeyNotFound) {
				return errors.Wrap(err)
			}
			keyboard := parseKeyboard(string(mailing_keyboard_raw_text))
			withKeyboard := false
			if len(keyboard.InlineKeyboard) > 0 {
				withKeyboard = true
			}

			rows, err := app.db.GetMailersRows()
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var id int64
				if rows.Err() != nil {
					return rows.Err()
				}
				if err = rows.Scan(&id); err != nil {
					return err
				}
				if withKeyboard {
					if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID:                   id,
							ChannelUsername:          "",
							ReplyToMessageID:         0,
							ReplyMarkup:              keyboard,
							DisableNotification:      false,
							AllowSendingWithoutReply: false,
						},
						FromChatID:          config.AdminID,
						FromChannelUsername: "",
						MessageID:           mailingMessageID,
						Caption:             "",
						ParseMode:           "",
						CaptionEntities:     nil,
					}); err != nil {
						return err
					}
				} else {
					if _, err = app.bot.Send(tgbotapi.CopyMessageConfig{
						BaseChat: tgbotapi.BaseChat{
							ChatID:                   id,
							ChannelUsername:          "",
							ReplyToMessageID:         0,
							ReplyMarkup:              nil,
							DisableNotification:      false,
							AllowSendingWithoutReply: false,
						},
						FromChatID:          config.AdminID,
						FromChannelUsername: "",
						MessageID:           mailingMessageID,
						Caption:             "",
						ParseMode:           "",
						CaptionEntities:     nil,
					}); err != nil {
						return err
					}
				}

				if err = app.db.DeleteMailuser(id); err != nil {
					return err
				}
			}

			if err = app.db.DropMailings(); err != nil {
				return err
			}
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "рассылка закончена"))
		} else {
			fmt.Println("рассылок нет")
		}
		return nil
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
		app.log.Error(fmt.Sprintf("%w", err), zap.Error(err))
	}
}

func (app App) SuperTranslate(ctx context.Context, user tables.Users, chatID int64, from, to, text string, replyToMessageID int, userMessage tgbotapi.Message) error {
	log := app.log.With(zap.String("from", from), zap.String("to", to), zap.String("text", text), zap.Int64("chat_id", chatID))
	if user.ID == 0 {
		user.ID = chatID
	}

	g, ctx := errgroup.WithContext(ctx)
	if userMessage.ReplyMarkup != nil {
		for i1, row := range userMessage.ReplyMarkup.InlineKeyboard {
			i1 := i1
			row := row
			for i2, btn := range row {
				i2 := i2
				btn := btn
				g.Go(func() error {
					tr, err := translate2.GoogleTranslate(ctx, from, to, btn.Text)
					if err != nil {
						return errors.Wrap(err)
					}
					userMessage.ReplyMarkup.InlineKeyboard[i1][i2].Text = tr.Text
					return nil
				})
			}
		}
	}
	if userMessage.Poll != nil {
		g.Go(func() error {
			tr, err := translate2.GoogleTranslate(ctx, from, to, userMessage.Poll.Question)
			userMessage.Poll.Question = tr.Text
			return errors.Wrap(err)
		})
		g.Go(func() error {
			tr, err := translate2.GoogleTranslate(ctx, from, to, applyEntitiesHtml(userMessage.Poll.Explanation, userMessage.Poll.ExplanationEntities))
			userMessage.Poll.Explanation = tr.Text
			return errors.Wrap(err)
		})
		for i, q := range userMessage.Poll.Options {
			i := i
			q := q
			g.Go(func() error {
				tr, err := translate2.GoogleTranslate(ctx, from, to, q.Text)
				userMessage.Poll.Options[i].Text = tr.Text
				return errors.Wrap(err)
			})
		}
		// TODO: apply entities to the answer and translate it
	}
	var (
		tr  string
		err error
	)
	g.Go(func() error {
		tr, err = app.translate(ctx, user, from, to, text) // examples мы сохраняем, чтобы соединить с keyboard.Examples и положить в кэш
		return errors.Wrap(err)
	})

	if err = g.Wait(); err != nil {
		return err
	}

	chunks := translate2.SplitIntoChunksBySentences(tr, 4000)
	for i, chunk := range chunks {
		if i > 0 {
			replyToMessageID = 0
		}
		switch {
		case userMessage.Poll != nil:
			options := make([]string, 0, len(userMessage.Poll.Options))
			for _, opt := range userMessage.Poll.Options {
				options = append(options, opt.Text)
			}
			_, err = app.bot.Send(tgbotapi.SendPollConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           chatID,
					ReplyToMessageID: replyToMessageID,
					ReplyMarkup:      userMessage.ReplyMarkup,
				},
				Question:              userMessage.Poll.Question,
				Options:               options,
				IsAnonymous:           userMessage.Poll.IsAnonymous,
				Type:                  userMessage.Poll.Type,
				AllowsMultipleAnswers: userMessage.Poll.AllowsMultipleAnswers,
				CorrectOptionID:       int64(userMessage.Poll.CorrectOptionID),
				Explanation:           userMessage.Poll.Explanation,
				ExplanationParseMode:  tgbotapi.ModeHTML,
				ExplanationEntities:   nil,
				OpenPeriod:            userMessage.Poll.OpenPeriod,
				CloseDate:             userMessage.Poll.CloseDate,
				IsClosed:              userMessage.Poll.IsClosed,
			})
		case userMessage.Audio != nil:
			thumbnail := ""
			if userMessage.Audio.Thumbnail != nil {
				thumbnail = userMessage.Audio.Thumbnail.FileID
			}
			_, err = app.bot.Send(tgbotapi.AudioConfig{
				BaseFile: tgbotapi.BaseFile{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatID,
						ReplyToMessageID: replyToMessageID,
						ReplyMarkup:      userMessage.ReplyMarkup,
					},
					File: tgbotapi.FileID(userMessage.Audio.FileID),
				},
				Thumb:     tgbotapi.FileID(thumbnail),
				Caption:   tr,
				ParseMode: tgbotapi.ModeHTML,
				Duration:  userMessage.Audio.Duration,
				Performer: userMessage.Audio.Performer,
				Title:     userMessage.Audio.Title,
			})
		case len(userMessage.Photo) > 0:
			maxResolutionPhoto := userMessage.Photo[len(userMessage.Photo)-1]
			_, err = app.bot.Send(tgbotapi.PhotoConfig{
				BaseFile: tgbotapi.BaseFile{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           chatID,
						ReplyToMessageID: replyToMessageID,
						ReplyMarkup:      userMessage.ReplyMarkup,
					},
					File: tgbotapi.FileID(maxResolutionPhoto.FileID),
				},
				Caption:   chunk,
				ParseMode: tgbotapi.ModeHTML,
			})
		default:
			_, err = app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           chatID,
					ReplyToMessageID: replyToMessageID,
					ReplyMarkup:      userMessage.ReplyMarkup,
				},
				Text:                  chunk,
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			})
		}
		if err != nil {
			log.Error("", zap.Error(err))
			return err
		}

	}

	return nil
}

func (app App) sendSpeech(user tables.Users, lang, text string, callbackID string) error {
	sdec, err := translate2.TTS(lang, text)
	if err != nil {
		if err == translate2.ErrTTSLanguageNotSupported {
			app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callbackID, user.Localize("%s не поддерживается", langs[user.Lang][lang])))
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
