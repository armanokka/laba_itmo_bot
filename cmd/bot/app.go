package bot

import (
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/pkg/botapi"
	"github.com/armanokka/translobot/pkg/dashbot"
	"github.com/armanokka/translobot/pkg/lingvo"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	"github.com/armanokka/translobot/repos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"html"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type App struct {
	bot *botapi.BotAPI
	bundle *i18n.Bundle
	log *zap.Logger
	db        repos.BotDB
	analytics dashbot.DashBot
}

func New(bot *botapi.BotAPI, db repos.BotDB, analytics dashbot.DashBot, logger *zap.Logger, bundle *i18n.Bundle) App {
	return App{
		bot:       bot,
		db:        db,
		analytics: analytics,
		log:       logger,
		bundle: bundle,
	}
}


func (app App) Run(ctx context.Context) error {
	app.bot.MakeRequest("deleteWebhook", map[string]string{})
	updates := app.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Bot have started"))
	g, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	g.Go(func() error {
		for {
			defer func() {
				if err := recover(); err != nil {
					app.log.Error("%w", zap.Any("error", err))
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
					defer func() {
						if err := recover(); err != nil {
							pp.Println("Panic:", err)
							app.notifyAdmin("Panic:", err)
						}
					}()
					start := time.Now()
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
					pp.Println("Time spent:", time.Since(start).String())
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
		app.log.Error("%w", zap.Error(err))
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


func (app App) SuperTranslate(from, to, text string, entities []tgbotapi.MessageEntity) (ret SuperTranslation, err error) {
	text = html.EscapeString(text)
	text = applyEntitiesHtml(text, entities)

	var (
		rev = translate2.ReversoTranslation{}
		dict = translate2.GoogleDictionaryResponse{}
		suggestions *lingvo.SuggestionResult
		lingv []lingvo.Dictionary
	)

	l := len(text)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		if l > 100 {
			return nil
		}
		dict, err = translate2.GoogleDictionary(from, strings.ToLower(text))
		return err
	})

	g.Go(func() error {
		_, ok1 := lingvo.Lingvo[from]
		_, ok2 := lingvo.Lingvo[to]
		if !ok1 || !ok2 {
			return nil
		}
		lingv, err = lingvo.GetDictionary(from, to, text)
		return err
	})

	g.Go(func() error {
		if l > 100 {
			return nil
		}
		if inMapValues(translate2.ReversoSupportedLangs(), from, to) && from != to {
			rev, err = translate2.ReversoTranslate(translate2.ReversoIso6392(from), translate2.ReversoIso6392(to), strings.ToLower(text))
		}
		return err
	})

	g.Go(func() error {
		_, ok1 := lingvo.Lingvo[from]
		_, ok2 := lingvo.Lingvo[to]
		if !ok1 || !ok2 {
			return nil
		}
		suggestions, err = lingvo.Suggestions(from, to, strings.ToLower(text), 1, 0)
		return err
	})

	g.Go(func() error {
		tr, err := translate2.GoogleHTMLTranslate(from, to, text)
		if err != nil {
			return err
		}
		if tr.Text == "" && text != "" {
			return err
		}
		ret.From = tr.From
		ret.TranslatedText = tr.Text
		ret.TranslatedText = strings.NewReplacer("<br> ", "<br>").Replace(ret.TranslatedText)
		ret.TranslatedText = strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "",  `<br>`, "\n").Replace(ret.TranslatedText)
		return nil
	})

	if err = g.Wait(); err != nil {
		app.notifyAdmin(err)
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

	if dict.Status == 200 && dict.DictionaryData != nil || len(lingv) > 0 {
		ret.Dictionary = true
	}

	if suggestions != nil && len(suggestions.Items) > 0 {
		ret.Suggestions = true
	}

	return ret, nil
}

func (app App) sendSpeech(id int64, lang, text string, callbackID string, localizer *i18n.Localizer) error {
	sdec, err := translate2.TTS(lang, text)
	if err != nil {
		if err == translate2.ErrTTSLanguageNotSupported {
			locale, err := localizer.LocalizeMessage(&i18n.Message{ID: "%s language is not supported"})
			if err != nil {
				return err
			}
			call := tgbotapi.NewCallback(callbackID, locale)
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
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "Iternal error"))
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
	audio := tgbotapi.NewAudio(id, f.Name())
	audio.Title = text
	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
	audio.ReplyMarkup = kb
	app.bot.Send(audio)
	return nil
}
