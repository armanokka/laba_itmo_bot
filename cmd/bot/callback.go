package bot

import (
	"bytes"
	"context"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/pkg/errors"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/unicode/norm"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func (app *App) onCallbackQuery(ctx context.Context, callback tgbotapi.CallbackQuery) {
	log := app.log.With(zap.Int64("id", callback.From.ID))
	defer func() {
		if err := recover(); err != nil {
			log.Error("%w", zap.Any("error", err), zap.String("stack_trace", string(debug.Stack())))
			app.bot.Send(tgbotapi.NewMessage(config.AdminID, "Panic:"+fmt.Sprint(err)))
		}
	}()
	warn := func(err error) {
		app.bot.Send(tgbotapi.NewCallback(callback.ID, "Error, sorry"))
		app.notifyAdmin(err)
	}

	var err error
	user, err := app.db.GetUserByID(callback.From.ID)
	if err != nil {
		warn(err)
	}
	user.SetLang(callback.From.LanguageCode)

	arr := strings.Split(callback.Data, ":")

	defer func() {
		if err := app.db.UpdateUserMetrics(callback.From.ID, "callback:"+callback.Data); err != nil {
			app.notifyAdmin(fmt.Errorf("%w", err))
		}
	}()

	switch arr[0] {
	case "delete_cache": // arr[1] - document's ._key
		if _, err = app.cache.RemoveDocument(ctx, arr[1]); err != nil {
			_, err = app.bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(callback.ID, err.Error()))
		} else {
			_, err = app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "OK"))
		}
		if err != nil {
			warn(err)
			log.Error("", zap.Error(err))
			return
		}
		log.Info("cache was deleted", zap.String("_key", arr[1]))
	case "cancel_mailing_act":
		if err := app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{"act": ""}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))
		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.ReplyToMessage.MessageID))
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "OK"))
	case "none":
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		app.db.LogBotMessage(callback.From.ID, "cb_none", "")
	case "delete": // arr[1] - text for callback query
		text := ""
		if len(arr) > 1 {
			text = strings.Join(arr[1:], ":")
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, text))
		app.bot.Send(tgbotapi.DeleteMessageConfig{
			ChatID:    callback.From.ID,
			MessageID: callback.Message.MessageID,
		})
		app.db.LogBotMessage(callback.From.ID, "cb_delete", "")
		return
	case "translate": // arr[1] - from, arr[2] - to
		entites := callback.Message.ReplyToMessage.Entities
		if len(callback.Message.ReplyToMessage.CaptionEntities) > 0 {
			entites = callback.Message.ReplyToMessage.CaptionEntities
		}
		callback.Message.ReplyToMessage.Text = applyEntitiesHtml(norm.NFKC.String(callback.Message.ReplyToMessage.Text), entites)
		pp.Println(arr[1], arr[2], callback.Message.ReplyToMessage.Text)

		if err := app.SuperTranslate(ctx, user, callback.From.ID, callback.Message.MessageID, arr[1], arr[2], callback.Message.ReplyToMessage.Text, true); err != nil {
			warn(err)
			return
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "speech": // arr[1] - from, arr[2] - to
		userMsgText := callback.Message.ReplyToMessage.Text

		botMsgText := callback.Message.Text
		if callback.Message.Caption != "" {
			botMsgText = callback.Message.Caption
		}
		user, err := app.db.GetUserByID(callback.From.ID)
		if err != nil {
			warn(err)
			return
		}
		go func() {
			if err = app.sendSpeech(user, arr[2], botMsgText, callback.ID); err != nil { // озвучиваем непереведенное сообщение
				warn(err)
				return
			}
		}()
		if err = app.sendSpeech(user, arr[1], userMsgText, callback.ID); err != nil { // озвучиваем переведенное сообщение
			warn(err)
			return
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	case "dictionary": // arr[1] - from, arr[2] - to, text in replied message
		callback.Message.ReplyToMessage.Text = strings.ToLower(callback.Message.ReplyToMessage.Text)

		var dictionary []string
		ctx2, cancel := context.WithCancel(ctx)
		g, ctx2 := errgroup.WithContext(ctx2)
		g.Go(func() (err error) {
			r, err := app.dictionary(ctx2, arr[1], callback.Message.ReplyToMessage.Text)
			if err != nil {
				if IsCtxError(err) {
					return nil
				}
				return err
			}
			cancel()
			dictionary = r
			log.Info("got dictionary from scratch")
			return nil
		})
		g.Go(func() error {
			r, err := app.seekForCache(ctx2, arr[1], arr[2], callback.Message.ReplyToMessage.Text)
			if err != nil {
				if IsCtxError(err) || errors.Is(err, ErrCacheNotFound) {
					return nil
				}
				return err
			}
			cancel()
			dictionary = r.Dictionary
			log.Info("got dictionary from cache")
			return nil
		})
		if err := g.Wait(); err != nil {
			warn(err)
			return
		}
		if len(dictionary) == 0 {
			app.bot.Send(tgbotapi.CallbackConfig{
				CallbackQueryID: callback.ID,
				Text:            user.Localize("Something went wrong. Try again later or contact @armanokka if it's no trouble"),
				ShowAlert:       true,
				URL:             "",
				CacheTime:       0,
			})
			return
		}

		text := strings.Join(dictionary, "\n- ")

		if callback.From.LanguageCode != "en" {
			tr, err := translate2.GoogleTranslate(ctx, "en", callback.From.LanguageCode, text)
			if err != nil {
				warn(err)
				return
			}
			text = tr.Text
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("translate:%s:%s", arr[1], arr[2]))))
		if _, err := app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &keyboard,
			},
			Text:                  text,
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
	case "examples": // arr[1] - from, arr[2] - to. Target text in replied message
		var examples map[string]string
		ctx, cancel := context.WithCancel(ctx)
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() (err error) {
			_, r, err := app.reverseTranslationsExamples(ctx, arr[1], arr[2], callback.Message.ReplyToMessage.Text)
			if err != nil {
				if IsCtxError(err) {
					return nil
				}
				return err
			}
			cancel()
			examples = r
			log.Info("got examples anew")
			return nil
		})
		g.Go(func() error {
			r, err := app.seekForCache(ctx, arr[1], arr[2], callback.Message.ReplyToMessage.Text)
			if err != nil {
				if IsCtxError(err) || errors.Is(err, ErrCacheNotFound) {
					return nil
				}
				return err
			}
			cancel()
			examples = r.Examples
			log.Info("got examples from cache")

			return nil
		})
		if err := g.Wait(); err != nil {
			warn(err)
			return
		}

		text := new(bytes.Buffer)

		count := 1
		r := strings.NewReplacer("<em>", "<b>", "</em>", "</b>")
		for source, target := range examples {
			source = r.Replace(source)
			target = r.Replace(target)
			if count > 3 {
				break
			}
			text.WriteString("\n\n<b>")
			text.WriteString(strconv.Itoa(count))
			text.WriteString(".</b> ")
			text.WriteString(source)
			text.WriteString("\n<b>└</b>")
			text.WriteString(target)
			count++
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("translate:%s:%s", arr[1], arr[2]))))
		if _, err := app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &keyboard,
			},
			Text:                  text.String(),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			pp.Println(err)
		}

		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
	case "reverse_translations": // arr[1], arr[2] = from, to (in iso6391)
		callback.Message.ReplyToMessage.Text = strings.ToLower(callback.Message.ReplyToMessage.Text)

		var rev map[string][]string
		ctx, cancel := context.WithCancel(ctx)
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() (err error) {
			r, _, err := app.reverseTranslationsExamples(ctx, translate2.ReversoIso6392(arr[1]), translate2.ReversoIso6392(arr[2]), callback.Message.ReplyToMessage.Text)
			if err != nil {
				if IsCtxError(err) {
					return nil
				}
				return err
			}
			cancel()
			rev = r
			log.Info("got reverse translations anew")
			return nil
		})
		g.Go(func() error {
			r, err := app.seekForCache(ctx, arr[1], arr[2], callback.Message.ReplyToMessage.Text)
			if err != nil {
				if IsCtxError(err) || errors.Is(err, ErrCacheNotFound) {
					return nil
				}
				return err
			}
			cancel()
			rev = r.ReverseTranslations
			log.Info("got reverse translations from cache")
			return nil
		})
		if err := g.Wait(); err != nil {
			warn(err)
			return
		}

		out := bytes.NewBufferString(callback.Message.ReplyToMessage.Text)
		out.WriteString("\n")
		count := 0
		for key, trs := range rev {
			if count >= 4 {
				break
			}
			out.WriteString("\n<b>" + key + "\n└</b>")
			if len(trs) > 4 {
				trs = trs[:4]
			}
			for i, tr := range trs {
				if i > 0 {
					out.WriteString(",")
				}
				out.WriteString("<code>" + tr + "</code>")
			}
			count++
		}

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("translate:%s:%s", arr[1], arr[2]))))
		if _, err := app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &keyboard,
			},
			Text:                  out.String(),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			warn(err)
			return
		}
		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
	case "setup_langs": // arr[1] - source language, arr[2] - direction to translate, in replied message there is source text
		app.bot.Send(tgbotapi.NewChatAction(callback.From.ID, "typing"))

		app.bot.Send(tgbotapi.NewDeleteMessage(callback.From.ID, callback.Message.MessageID))

		from := arr[1]
		to := arr[2]
		if err := app.db.UpdateUserByMap(callback.From.ID, map[string]interface{}{
			"my_lang":       to,
			"to_lang":       from,
			"last_activity": time.Now(),
			"act":           "",
		}); err != nil {
			warn(err)
			return
		}

		app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		app.bot.Send(tgbotapi.NewMessage(callback.From.ID, user.Localize("start_info", langs[callback.From.LanguageCode][from], langs[callback.From.LanguageCode][to])))

		if _, err := app.bot.Send(tgbotapi.DocumentConfig{
			BaseFile: tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   callback.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              nil,
					DisableNotification:      false,
					AllowSendingWithoutReply: false,
				},
				File: tgbotapi.FilePath("inline.gif"),
			},
			Thumb:                       nil,
			Caption:                     user.Localize("Как переводить еще удобнее"),
			ParseMode:                   "",
			CaptionEntities:             nil,
			DisableContentTypeDetection: false,
		}); err != nil {
			pp.Println(err)
		}

		if _, err := app.bot.Send(tgbotapi.NewMessage(callback.From.ID, user.Localize("format_translation_tip"))); err != nil {
			pp.Println(err)
		}
		// user.Localize("")
	case "setup_langs_pagination": // arr[1] - source language of the text, arr[2] - offset
		from := arr[1]
		offset, err := strconv.Atoi(arr[2])
		if err != nil {
			warn(err)
			return
		}
		if offset < 0 || offset > len(codes[user.Lang])-1 {
			warn(fmt.Errorf("offset is too big, len(codes[user.Lang]) is %d, offset ois %d", len(codes[user.Lang]), offset))
			return
		}

		count := 18
		if offset+count > len(codes[user.Lang])-1 {
			count = len(codes[user.Lang]) - 1 - offset
		}

		if count == 0 {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
			return
		}

		back := offset - 18
		if back < 0 {
			back = 0
		}
		kb, err := buildLangsPagination(user, offset, count, arr[1],
			fmt.Sprintf("setup_langs:%s:%s", from, "%s"),
			fmt.Sprintf("setup_langs_pagination:%s:%d", from, back),
			fmt.Sprintf("setup_langs_pagination:%s:%d", from, offset+count))
		if err != nil {
			warn(err)
		}

		if reflect.DeepEqual(*callback.Message.ReplyMarkup, kb) {
			app.bot.Send(tgbotapi.NewCallback(callback.ID, ""))
			return
		}

		_, err = app.bot.Send(tgbotapi.EditMessageReplyMarkupConfig{tgbotapi.BaseEdit{
			ChatID:          callback.From.ID,
			ChannelUsername: "",
			MessageID:       callback.Message.MessageID,
			InlineMessageID: "",
			ReplyMarkup:     &kb,
		}})
		if err != nil {
			app.bot.Send(tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:                   callback.From.ID,
					ChannelUsername:          "",
					ReplyToMessageID:         0,
					ReplyMarkup:              &kb,
					DisableNotification:      true,
					AllowSendingWithoutReply: false,
				},
				Text:                  callback.Message.Text,
				ParseMode:             "",
				Entities:              nil,
				DisableWebPagePreview: false,
			})
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
		//case "inline_translate": //arr[1] - lang
		//	pp.Println(callback)
		//	tr, err := translate2.MicrosoftTranslate("", arr[1], callback.Message.Text)
		//	if err != nil {
		//		warn(err)
		//		return
		//	}
		//	app.bot.Send(tgbotapi.NewEditMessageText(callback.From.ID, callback.Message.MessageID, tr.TranslatedText))
		//	default:
		//	app.bot.Send(tgbotapi.NewMessage(callback.From.ID, "Action is expired. Send /start"))
	case "paraphrase": // Paraphrasing text. arr[1] - from, arr[2] -to, text in replied message
		paraphrase := make([]string, 0, 6)
		ctx, cancel := context.WithCancel(ctx)
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() (err error) {
			p, err := translate2.ReversoParaphrase(ctx, arr[2], callback.Message.Text)
			if err != nil {
				if IsCtxError(err) {
					return nil
				}
				return errors.Wrap(err)
			}
			cancel()
			paraphrase = p
			log.Info("paraphrased from scratch")
			return nil
		})
		g.Go(func() error {
			cache, err := app.seekForCache(ctx, arr[1], arr[2], callback.Message.ReplyToMessage.Text)
			if err != nil {
				if IsCtxError(err) || errors.Is(err, ErrCacheNotFound) {
					return nil
				}
				return errors.Wrap(err)
			}
			cancel()
			log.Info("paraphrased from cache")
			paraphrase = cache.Paraphrase
			return nil
		})
		if err := g.Wait(); err != nil {
			warn(err)
			cancel()
			log.Error("", zap.Error(err))
			return
		}

		var text = bytes.NewBufferString(">" + callback.Message.Text + "\n\n")
		for i, v := range paraphrase {
			if i > 0 {
				text.WriteString("\n\n")
			}
			text.WriteString(highlightDiffs(callback.Message.Text, v, "<b>", "</b>"))
		}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("translate:%s:%s", arr[1], arr[2]))))
		if _, err := app.bot.Send(tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:          callback.From.ID,
				ChannelUsername: "",
				MessageID:       callback.Message.MessageID,
				InlineMessageID: "",
				ReplyMarkup:     &keyboard,
			},
			Text:                  text.String(),
			ParseMode:             tgbotapi.ModeHTML,
			Entities:              nil,
			DisableWebPagePreview: false,
		}); err != nil {
			warn(err)
			return
		}
		app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, ""))
	}
}
