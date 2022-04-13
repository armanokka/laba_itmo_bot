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
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
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
	deepl     translate2.Deepl
	bot       *botapi.BotAPI
	log       *zap.Logger
	db        repos.BotDB
	analytics dashbot.DashBot
	bc        *bitcask.Bitcask
}

func New(bot *botapi.BotAPI, db repos.BotDB, analytics dashbot.DashBot, log *zap.Logger, bc *bitcask.Bitcask /* deepl translate2.Deepl*/) App {
	return App{
		//deepl:     deepl,
		bot:       bot,
		db:        db,
		analytics: analytics,
		log:       log,
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
				pp.Println("Stopping receiving updates...")
				app.bot.StopReceivingUpdates()
				wg.Wait()
				pp.Println("Waiting for goroutines finishing...")
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
						app.onInlineQuery(*update.InlineQuery)
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
				return err
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
		rev            = translate2.ReversoTranslation{}
		dict           = translate2.GoogleDictionaryResponse{}
		suggestions    *lingvo.SuggestionResult
		LingvoTr       string
		YandexTr       string
		DeeplTr        string
		MicrosoftTr    string
		GoogleFromToTr string
		GoogleToFromTr string
		//lingv []lingvo.Dictionary
	)

	l := len([]rune(text))
	lower := strings.ToLower(text)

	if from == "auto" {
		tr, err := translate2.GoogleTranslate(from, to, cutStringUTF16(text, 100))
		if err != nil {
			return SuperTranslation{}, errors.WrapPrefix(err, "g.Go: translate2.GoogleTranslate", 1)
		}
		from = tr.FromLang
	}

	g, _ := errgroup.WithContext(context.Background())
	log := app.log.With(zap.String("from", from), zap.String("to", to), zap.String("text", text))
	if l < 100 {
		g.Go(func() error {
			dict, err = translate2.GoogleDictionary(from, lower)
			if err != nil {
				log.Error("translate2.GoogleDictionary", zap.Error(err))
				return errors.WrapPrefix(err, "g.Go: translate2.GoogleDictionary:", 1)
			}
			definitions := 0
			for _, data := range dict.DictionaryData {
				for _, entry := range data.Entries {
					for _, senseFamily := range entry.SenseFamilies {
						definitions += len(senseFamily.Senses)
					}
				}
			}
			if definitions < 2 {
				dict = translate2.GoogleDictionaryResponse{}
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
				log.Error("translate2.ReversoSupportedLangs", zap.Error(err))
				return errors.WrapPrefix(err, "g.Go: translate2.ReversoTranslate:", 1)
			}
			return nil
		})
	}

	_, ok1 := lingvo.Lingvo[from]
	_, ok2 := lingvo.Lingvo[to]
	if ok1 && ok2 && l < 50 {
		g.Go(func() error {
			suggestions, err = lingvo.Suggestions(from, to, lower, 1, 0)
			if err != nil {
				log.Error("lingvo.Suggestions", zap.Error(err))
			}
			return err
		})
	}

	if l < 50 && in(lingvo.LingvoDictionaryLangs, user.MyLang, user.ToLang) {
		g.Go(func() error {
			v, err := lingvo.GetDictionary(user.MyLang, user.ToLang, lower)
			if err != nil {
				log.Error("lingvo.GetDictionary", zap.Error(err), zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang))
				return err
			}
			if len(v) == 0 {
				v, err = lingvo.GetDictionary(user.ToLang, user.MyLang, lower)
				if err != nil {
					log.Error("lingvo.GetDictionary", zap.Error(err), zap.String("my_lang", user.MyLang), zap.String("to_lang", user.ToLang))
					return err
				}
			}
			if len(v) > 8 {
				v = v[:8]
			}

			out := ""
			usedWords := make([]string, 0, 10)
			for _, r := range v {
				words := strings.Split(r.Translations, ";")
				if len(words) == 0 {
					words[0] = r.Translations
				}
				lines := strings.Split(out, "\n")
				last := lines[len(lines)-1]

				for _, word := range words {
					if word == "" {
						continue
					}
					word = strings.TrimSpace(word)
					if inFuzzy(usedWords, word) {
						continue
					}
					if len(last)+len(words) > 25 {
						out += "\n" + word + ";"
					} else {
						out += " " + word + ";"
					}
					usedWords = append(usedWords, word)
				}
			}
			LingvoTr = out
			return nil
		})
	}

	_, ok1 = translate2.YandexSupportedLanguages[from]
	_, ok2 = translate2.YandexSupportedLanguages[to]

	if ok1 && ok2 {
		g.Go(func() error {
			tr, err := translate2.YandexTranslate(from, to, text)
			if err != nil {
				app.log.Error("translate2.YandexTranslate", zap.Error(err))
				return err
			}
			YandexTr = tr
			return nil
		})
	} else {
		if html.UnescapeString(text) != html.EscapeString(text) { // есть html теги
			g.Go(func() error {
				tr, err := translate2.MicrosoftTranslate(from, to, text)
				if err != nil {
					app.log.Error("translate2.MicrosoftTranslate", zap.Error(err))
					return err
				}
				MicrosoftTr = tr.TranslatedText
				return nil
			})
		}
	}

	g.Go(func() error {
		tr, err := translate2.GoogleTranslate(from, to, text)
		if err != nil {
			app.log.Error("translate2.GoogleTranslate", zap.Error(err))
			return err
		}
		GoogleFromToTr = tr.Text
		return nil
	})
	g.Go(func() error {
		tr, err := translate2.GoogleTranslate(to, from, text)
		if err != nil {
			app.log.Error("translate2.GoogleTranslate", zap.Error(err))
			return err
		}
		GoogleToFromTr = tr.Text
		return nil
	})

	if err = g.Wait(); err != nil {
		app.log.Error("g.Wait()", zap.Error(err))
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

	switch {
	case LingvoTr != "":
		ret.TranslatedText = LingvoTr
		pp.Println("translated via lingvo")
		break
	case DeeplTr != "":
		ret.TranslatedText = DeeplTr
		pp.Println("translated via deepl")
		break
	case YandexTr != "":
		if fuzzy.EditDistance(text, GoogleFromToTr) < fuzzy.EditDistance(text, YandexTr) {
			ret.TranslatedText = YandexTr
			pp.Println("translated via yandex")
			break
		}
		fallthrough
	case GoogleToFromTr != "":
		fallthrough
	case GoogleFromToTr != "":
		if fuzzy.EditDistance(text, GoogleFromToTr) > fuzzy.EditDistance(text, GoogleToFromTr) {
			ret.TranslatedText = GoogleFromToTr
		} else {
			ret.TranslatedText = GoogleToFromTr
		}
		pp.Println("translated via google")

		break
	case MicrosoftTr != "":
		ret.TranslatedText = MicrosoftTr
		pp.Println("translated via microsoft")

	}

	return ret, nil
}

func (app App) sendSpeech(user tables.Users, lang, text string, callbackID string, replyToMessageID int) error {
	sdec, err := translate2.TTS(lang, text)
	if err != nil {
		if err == translate2.ErrTTSLanguageNotSupported {
			call := tgbotapi.NewCallback(callbackID, user.Localize("%s не поддерживается", langs[user.Lang][lang]))
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
	audio.ReplyToMessageID = replyToMessageID
	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
	audio.ReplyMarkup = kb
	app.bot.Send(audio)
	return nil
}
