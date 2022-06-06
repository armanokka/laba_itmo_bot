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
	"github.com/armanokka/translobot/pkg/helpers"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	"github.com/armanokka/translobot/repos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/qiniu/iconv"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/unicode/norm"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"sync"
)

type App struct {
	htmlTagsRe          *regexp.Regexp
	reSpecialCharacters *regexp.Regexp
	deepl               translate2.Deepl
	bot                 *botapi.BotAPI
	log                 *zap.Logger
	db                  repos.BotDB
	analytics           dashbot.DashBot
	bc                  *bitcask.Bitcask
}

func New(bot *botapi.BotAPI, db repos.BotDB, analytics dashbot.DashBot, log *zap.Logger, bc *bitcask.Bitcask) (App, error) {
	iconv.Open("utf-8", "")
	app := App{
		htmlTagsRe:          regexp.MustCompile("<\\s*[^>]+>(.*?)"),
		reSpecialCharacters: regexp.MustCompile(`[[:punct:]]`),
		//deepl:      translate2.Deepl{},
		bot:       bot,
		log:       log,
		db:        db,
		analytics: analytics,
		bc:        bc,
	}
	return app, nil
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
					if update.MyChatMember != nil {
						if update.MyChatMember.From.LanguageCode == "" || !in(config.BotLocalizedLangs, update.MyChatMember.From.LanguageCode) {
							update.MyChatMember.From.LanguageCode = "en"
						}
						app.onMyChatMember(*update.MyChatMember)
					} else if update.Message != nil {
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

func (app App) SuperTranslate(ctx context.Context, user tables.Users, chatID int64, from, to, text string, userMessage tgbotapi.Message) error {
	entities := userMessage.Entities
	if len(userMessage.CaptionEntities) > 0 {
		entities = userMessage.CaptionEntities
	}
	text = norm.NFKC.String(text)
	text = applyEntitiesHtml(text, entities)

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
	}

	tr, from, err := app.translate(ctx, from, to, text) // examples мы сохраняем, чтобы соединить с keyboard.Examples и положить в кэш
	if err != nil {
		return errors.Unwrap(err)
	}
	//if !validHtml(tr) {
	//	tr =
	//	pp.Println("escaping")
	//	tr = html.EscapeString(tr)
	//}

	//app.bot.Send(tgbotapi.NewDeleteMessage(chatID, userMessage.MessageID))
	chunks := translate2.SplitIntoChunksBySentences(tr, 4096)
	for _, chunk := range chunks {
		chunk = closeUnclosedTags(chunk)
		switch {
		case userMessage.Poll != nil:
			options := make([]string, 0, len(userMessage.Poll.Options))
			for _, opt := range userMessage.Poll.Options {
				options = append(options, opt.Text)
			}
			_, err = app.bot.Send(tgbotapi.SendPollConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           chatID,
					ReplyToMessageID: 0,
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
						ReplyToMessageID: 0,
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
			chunk = helpers.CutStringUTF16(chunk, 1024) // MEDIA_CAPTION_TOO_LONG
			maxResolutionPhoto := userMessage.Photo[len(userMessage.Photo)-1]
			_, err = app.bot.Send(tgbotapi.PhotoConfig{
				BaseFile: tgbotapi.BaseFile{
					BaseChat: tgbotapi.BaseChat{
						ChatID:      chatID,
						ReplyMarkup: userMessage.ReplyMarkup,
					},
					File: tgbotapi.FileID(maxResolutionPhoto.FileID),
				},
				Caption:   chunk,
				ParseMode: tgbotapi.ModeHTML,
			})
		default:
			var keyboard interface{}
			if userMessage.ReplyMarkup != nil {
				keyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton(langs[user.Lang][user.MyLang]+" "+flags[user.MyLang].Emoji),
						tgbotapi.NewKeyboardButton("↔️"),
						tgbotapi.NewKeyboardButton(langs[user.Lang][user.ToLang]+" "+flags[user.ToLang].Emoji)))
			} else {
				keyboard = userMessage.ReplyMarkup
			}
			_, err = app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: keyboard,
				},
				Text:                  chunk,
				ParseMode:             tgbotapi.ModeHTML,
				Entities:              nil,
				DisableWebPagePreview: false,
			})
		}
		if err != nil {
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   config.AdminID,
					ChannelUsername:          "",
					ProtectContent:           false,
					ReplyToMessageID:         0,
					ReplyMarkup:              nil,
					DisableNotification:      false,
					AllowSendingWithoutReply: false,
				},
				Text:                  fmt.Sprintf("Error: %s\nUser's text:%s\nTranslation:%s", err.Error(), text, tr),
				ParseMode:             "",
				Entities:              nil,
				DisableWebPagePreview: false,
			})
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
