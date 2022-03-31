package bot

import (
	"context"
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	"github.com/armanokka/translobot/pkg/lingvo"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	"github.com/armanokka/translobot/repos"
	"github.com/go-errors/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"html"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

type App struct {
	bot       *botapi.BotAPI
	log       *zap.Logger
	db        repos.BotDB
	analytics dashbot.DashBot
	bc        *bitcask.Bitcask
}

func New(bot *botapi.BotAPI, db repos.BotDB, analytics dashbot.DashBot, logger *zap.Logger, bc *bitcask.Bitcask) App {
	return App{
		bot:       bot,
		db:        db,
		analytics: analytics,
		log:       logger,
		bc:        bc,
	}
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
						if update.Message.From.LanguageCode == "" {
							update.Message.From.LanguageCode = "en"
						}
						app.onMessage(ctx, *update.Message)
					} else if update.CallbackQuery != nil {
						if update.CallbackQuery.From.LanguageCode == "" {
							update.CallbackQuery.From.LanguageCode = "en"
						}
						app.onCallbackQuery(ctx, *update.CallbackQuery)
					} else if update.InlineQuery != nil {
						if update.InlineQuery.From.LanguageCode == "" {
							update.InlineQuery.From.LanguageCode = "en"
						}
						app.onInlineQuery(*update.InlineQuery)
					} else if update.MyChatMember != nil {
						if update.MyChatMember.From.LanguageCode == "" {
							update.MyChatMember.From.LanguageCode = "en"
						}
						app.onMyChatMember(*update.MyChatMember)
					}
				}()

			}
		}
	})
	g.Go(func() error {
		pp.Println("чекаем не завалялась ли рассылка")
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
			rows, err := app.db.GetMailersRows()
			if err != nil {
				return err
			}

			defer rows.Close()
			for rows.Next() {
				var id int64
				if err = rows.Scan(&id); err != nil {
					return err
				}
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
				if err = app.db.DeleteMailuser(id); err != nil {
					return err
				}
			}

			if err = app.db.DropMailings(); err != nil {
				return err
			}
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "рассылка закончена"))
		} else {
			pp.Println("рассылок нет")
		}
		return nil
	})
	return g.Wait()
}

func (app App) notifyAdmin(args ...interface{}) {
	text := ""
	for _, arg := range args {
		switch v := arg.(type) {
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

//func (app app) setMyCommands(langs []string, commands []tgbotapi.BotCommand) error {
//	newCommands := make(map[string][]tgbotapi.BotCommand)
//	for _, lang := range langs {
//		newCommands[lang] = []tgbotapi.BotCommand{}
//		for _, command := range commands {
//			tr, err := translate.GoogleHTMLTranslate("en", lang, command.Description)
//			if err != nil {
//				return err
//			}
//			newCommands[lang]= append(newCommands[lang], tgbotapi.BotCommand{
//				Command:     command.Command,
//				Description: tr.Text,
//			})
//		}
//	}
//
//	for lang, command := range newCommands {
//		data, err := json.Marshal(command)
//		if err != nil {
//			return err
//		}
//		params := tgbotapi.Params{}
//		params.AddNonEmpty("commands", string(data))
//		params.AddNonEmpty("language_code", lang)
//		if _, err = app.bot.MakeRequest("setMyCommands", params); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (app App) SuperTranslate(user tables.Users, from, to, text string, entities []tgbotapi.MessageEntity) (ret SuperTranslation, err error) {
	text = applyEntitiesHtml(text, entities)
	//text = html.EscapeString(text)

	var (
		rev         = translate2.ReversoTranslation{}
		dict        = translate2.GoogleDictionaryResponse{}
		suggestions *lingvo.SuggestionResult
		//lingv []lingvo.Dictionary
	)

	l := len(text)
	lower := strings.ToLower(text)
	g, _ := errgroup.WithContext(context.Background())

	if l < 100 {
		g.Go(func() error {
			dict, err = translate2.GoogleDictionary(from, lower)
			if err != nil {
				err = errors.WrapPrefix(err, "g.Go: translate2.GoogleDictionary:", 1)
			}
			return err
		})
	}

	if l < 100 {
		g.Go(func() error {
			if inMapValues(translate2.ReversoSupportedLangs(), from, to) && from != to {
				rev, err = translate2.ReversoTranslate(translate2.ReversoIso6392(from), translate2.ReversoIso6392(to), lower)
			}
			if err != nil {
				err = errors.WrapPrefix(err, "g.Go: translate2.ReversoTranslate:", 1)
			}
			return err
		})
	}

	_, ok1 := lingvo.Lingvo[from]
	_, ok2 := lingvo.Lingvo[to]
	if ok1 && ok2 && len([]rune(lower)) < 50 {
		g.Go(func() error {
			suggestions, err = lingvo.Suggestions(from, to, lower, 1, 0)
			if err != nil {
				err = errors.WrapPrefix(err, "g.Go: lingvo.Suggestions:", 1)
			}
			return err
		})
	}

	if from == "auto" {
		tr, err := translate2.GoogleTranslate(from, to, cutStringUTF16(text, 100))
		if err != nil {
			return SuperTranslation{}, errors.WrapPrefix(err, "g.Go: translate2.GoogleTranslate", 1)
		}
		from = tr.FromLang
	}

	g.Go(func() error {
		if len(text) < 50 {
			if v, err := lingvo.GetDictionary(user.MyLang, user.ToLang, lower); err == nil && len(v) > 0 {
				out := ""
				for i, r := range v {
					if i > 0 {
						out += "\n"
					}
					out += r.Translations
				}
				ret.TranslatedText = out
				pp.Println("translated via lingvo")
				return nil
			}
			if v, err := lingvo.GetDictionary(user.ToLang, user.MyLang, lower); err == nil && len(v) > 0 {
				out := ""
				for i, r := range v {
					if i > 0 {
						out += "\n"
					}
					out += r.Translations
				}
				ret.TranslatedText = out
				pp.Println("translated via lingvo")
				return nil
			}
		}

		_, ok1 = translate2.YandexSupportedLanguages[from]
		_, ok2 = translate2.YandexSupportedLanguages[to]
		if ok1 && ok2 {
			tr, err := translate2.YandexTranslate(from, to, text)
			if err != nil {
				return err
			}
			ret.TranslatedText = tr
			pp.Println("translated via yandex")
			return nil
		}

		if html.UnescapeString(text) != html.EscapeString(text) { // есть html теги
			tr, err := translate2.MicrosoftTranslate(from, to, text)
			if err != nil {
				return err
			}
			ret.TranslatedText = tr.TranslatedText
			pp.Println("translated via microsoft")
			return nil
		} else {
			tr, err := translate2.GoogleTranslate(from, to, text)
			if err != nil {
				return err
			}
			ret.TranslatedText = tr.Text
			pp.Println("translated via google")
			return nil
		}
		return nil
	})

	if err = g.Wait(); err != nil {
		return SuperTranslation{}, err
	}

	if len(rev.ContextResults.Results) > 0 {
		if len(rev.ContextResults.Results[0].SourceExamples) > 0 {
			ret.Examples = true
		}
		if rev.ContextResults.Results[0].Translation != "" {
			ret.Translations = true
		}
	}

	if dict.Status == 200 && dict.DictionaryData != nil {
		ret.Dictionary = true
	}

	if suggestions != nil && len(suggestions.Items) > 0 {
		ret.Suggestions = true
	}

	return ret, nil
}

func (app App) sendSpeech(user tables.Users, lang, text string, callbackID string) error {
	sdec, err := translate2.TTS(lang, text)
	if err != nil {
		if err == translate2.ErrTTSLanguageNotSupported {
			call := tgbotapi.NewCallback(callbackID, user.Localize("%s не поддерживается", langs[lang].Name))
			call.ShowAlert = true
			app.bot.AnswerCallbackQuery(call)
			return nil
		}
		if e, ok := err.(translate2.HTTPError); ok {
			if e.Code == 500 || e.Code == 414 {
				call := tgbotapi.NewCallback(callbackID, "Too big text")
				call.ShowAlert = true
				app.bot.AnswerCallbackQuery(call)
				return nil
			}
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "Internal error"))
		pp.Println(err)
		return err
	}
	app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "⏳"))

	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	defer func() {
		if err = f.Close(); err != nil {
			app.notifyAdmin(err)
		}
	}()
	_, err = f.Write(sdec)
	if err != nil {
		return err
	}
	audio := tgbotapi.NewAudio(user.ID, tgbotapi.FilePath(f.Name()))
	audio.Title = text
	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
	audio.ReplyMarkup = kb
	app.bot.Send(audio)
	return nil
}
