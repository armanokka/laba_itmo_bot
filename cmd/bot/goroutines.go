package bot

import (
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/lingvo"
	"github.com/armanokka/translobot/pkg/translate"
	"github.com/forPelevin/gomoji"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
	"time"
)

var ErrCacheNotFound = fmt.Errorf("seekForCache: cache not found")

func (app App) seekForCache(ctx context.Context, from, to, text string) (RecordTranslationCache, error) {
	cursor, err := app.arangodb.Query(ctx, "FOR doc IN cache FILTER doc.text == @text AND doc.from == @from AND doc.to == @to LIMIT 1 RETURN doc", map[string]interface{}{
		"text": text,
		"from": from,
		"to":   to,
	})
	if err != nil {
		return RecordTranslationCache{}, err
	}
	defer cursor.Close()

	var record RecordTranslationCache
	meta, err := cursor.ReadDocument(ctx, &record)
	if err != nil && !driver.IsNoMoreDocuments(err) {
		return RecordTranslationCache{}, err
	}
	if meta.Key != "" {
		return record, nil
	}
	return record, ErrCacheNotFound
}

// translate returns translation, examples and error
func (app App) translate(ctx context.Context, user tables.Users, from, to, text string) (string, map[string]string, error) {
	log := app.log.With(zap.String("from", from), zap.String("to", to), zap.String("text", text))
	g, ctx := errgroup.WithContext(ctx)
	var (
		MicrosoftTr    string
		GoogleFromToTr string
		GoogleToFromTr string
		LingvoTr       string
		YandexTr       string
		examples       map[string]string
	)
	if app.htmlTagsRe.MatchString(text) { // есть html теги
		g.Go(func() error {
			start := time.Now()
			text = strings.ReplaceAll(text, "\n", "<br>")
			tr, err := translate.MicrosoftTranslate(ctx, from, to, text)
			if err != nil {
				return err
			}
			MicrosoftTr = tr.TranslatedText
			log.Info("microsoft finished", zap.String("spent", time.Since(start).String()))
			return nil
		})
	} else {
		g.Go(func() error {
			tr, err := translate.GoogleTranslate(ctx, from, to, text)
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
			tr.Text = strings.ReplaceAll(tr.Text, ` \ n`, `\n`)
			tr.Text = strings.ReplaceAll(tr.Text, `\ n`, `\n`)
			GoogleFromToTr = tr.Text
			return nil
		})
		g.Go(func() error {
			tr, err := translate.GoogleTranslate(ctx, to, from, text)
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
			tr.Text = strings.ReplaceAll(tr.Text, ` \ n`, `\n`)
			tr.Text = strings.ReplaceAll(tr.Text, `\ n`, `\n`)
			GoogleToFromTr = tr.Text
			return nil
		})
	}

	_, ok1 := translate.YandexSupportedLanguages[from]
	_, ok2 := translate.YandexSupportedLanguages[to]

	if ok1 && ok2 {
		g.Go(func() error {
			var err error
			YandexTr, err = translate.YandexTranslate(ctx, from, to, text)
			return errors.Wrap(err)
		})
	}
	_, ok1 = lingvo.Lingvo[user.MyLang]
	_, ok2 = lingvo.Lingvo[user.ToLang]
	if len(text) < 50 && ok1 && ok2 {
		g.Go(func() (err error) {
			LingvoTr, examples, err = app.lingvo(ctx, user.MyLang, user.ToLang, text)
			return errors.Wrap(err)
		})
	}

	if err := g.Wait(); err != nil {
		return "", nil, err
	}

	if LingvoTr != "" {
		log.Info("translated via lingvo", zap.String("translation", LingvoTr))
		return LingvoTr, examples, nil
	} else if GoogleFromToTr != "" && GoogleToFromTr != "" {
		if diff(text, GoogleFromToTr) > diff(text, GoogleToFromTr) || (diff(GoogleFromToTr, YandexTr) < diff(GoogleFromToTr, GoogleToFromTr)) && YandexTr != "" {
			log.Info("translated via google", zap.String("translation", GoogleFromToTr))
			return GoogleFromToTr, examples, nil
		} else {
			log.Info("translated via google", zap.String("translation", GoogleToFromTr))
			return GoogleToFromTr, examples, nil
		}
	} else if YandexTr != "" {
		log.Info("translated via yandex", zap.String("translation", YandexTr))
		return YandexTr, examples, nil
	} else if MicrosoftTr != "" {
		MicrosoftTr = strings.ReplaceAll(MicrosoftTr, "<br>", "\n")
		log.Info("translated via microsoft", zap.String("translation", MicrosoftTr))
		return MicrosoftTr, examples, nil
	}
	return "", nil, fmt.Errorf("all translators returned empty result\n%s->%s\n%s", from, to, text)
}

func (app App) keyboard(ctx context.Context, from, to, text string) (Keyboard, error) {
	var (
		dict                 []string
		paraphrase           []string
		examples             map[string]string
		err                  error
		reversedTranslations map[string][]string
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { // dictionary
		dict, err = app.dictionary(ctx, from, text)
		return errors.Wrap(err)
	})

	if inMapValues(translate.ReversoSupportedLangs(), from, to) && from != to {
		if len(text) < 100 {
			g.Go(func() error { // examples
				reversedTranslations, examples, err = app.reverseTranslationsExamples(ctx, from, to, text)
				return errors.Wrap(err)
			})
		}
		if l := len(strings.Fields(app.reSpecialCharacters.ReplaceAllString(gomoji.RemoveEmojis(text), ""))); l > 2 && l < 31 {
			g.Go(func() error { // paraphrase
				paraphrase, err = translate.ReversoParaphrase(ctx, from, text)
				return errors.Wrap(err)
			})
		}
	}

	if err = g.Wait(); err != nil {
		return Keyboard{}, err
	}
	return Keyboard{
		Dictionary:          dict,
		Paraphrase:          paraphrase,
		Examples:            examples,
		ReverseTranslations: reversedTranslations,
	}, nil

}

// lingvo returns its translation and examples and error
func (app App) lingvo(ctx context.Context, from, to, text string) (string, map[string]string, error) {
	text = strings.ToLower(text)
	var (
		d1  []lingvo.Dictionary
		d2  []lingvo.Dictionary
		err error
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		d1, err = lingvo.GetDictionary(ctx, from, to, text)
		return errors.Wrap(err)
	})
	g.Go(func() error {
		d2, err = lingvo.GetDictionary(ctx, to, from, text)
		return errors.Wrap(err)
	})

	if err = g.Wait(); err != nil {
		return "", nil, err
	}

	if len(d1) == 0 && len(d2) > 0 { // помещаем d2 в d1, если d1 пустой
		d1, d2 = d2, d1
	}
	examples := make(map[string]string, len(d1)*2)
	for _, d := range d1 {
		if d.Examples == "" {
			continue
		}
		lines := strings.Split(d.Examples, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			parts := strings.Split(line, "—")
			if len(parts) != 2 {
				app.log.Info("", zap.Error(fmt.Errorf("lingvo's dictionary examples: examples aren't splitted by \n and -\n%s", d.Examples)))
				continue
				//return "", nil, fmt.Errorf("lingvo's dictionary examples: examples aren't splitted by \n and -\n%s", d.Examples)
			}
			examples[parts[0]] = parts[1]
		}
	}
	return writeLingvo(d1), examples, nil
}

func (app App) reverseTranslationsExamples(ctx context.Context, from, to, text string) (map[string][]string, map[string]string, error) {
	from = translate.ReversoIso6392(from)
	to = translate.ReversoIso6392(to)
	rev, err := translate.ReversoTranslate(ctx, from, to, text)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}
	reversedTranslations := make(map[string][]string, len(rev.ContextResults.Results))
	var mu sync.Mutex
	if len(rev.ContextResults.Results) > 10 {
		rev.ContextResults.Results = rev.ContextResults.Results[:10]
	}
	examples := make(map[string]string, 3)
	g, ctx := errgroup.WithContext(ctx)
	for _, result := range rev.ContextResults.Results {
		result := result
		if result.Translation == "" || len(result.SourceExamples) == 0 {
			continue
		}
		r := strings.NewReplacer("<em>", "<b>", "</em>", "</b>")
		for i := 0; i < len(result.SourceExamples) && len(examples) <= 3; i++ {
			examples[r.Replace(result.SourceExamples[i])] = r.Replace(result.TargetExamples[i])
		}

		reversedTranslations[result.Translation] = make([]string, 0, 4)
		g.Go(func() error {
			tr, err := translate.ReversoTranslate(ctx, to, from, result.Translation)
			if err != nil {
				return errors.Wrap(err)
			}
			if len(tr.ContextResults.Results) > 5 {
				tr.ContextResults.Results = tr.ContextResults.Results[:5]
			}
			trs := make([]string, 0, 5)
			for _, result := range tr.ContextResults.Results {
				if result.Translation == "" || result.Translation == text {
					continue
				}
				trs = append(trs, result.Translation)
			}
			mu.Lock()
			reversedTranslations[result.Translation] = trs
			mu.Unlock()
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		return nil, nil, err
	}
	return reversedTranslations, examples, nil
}

func (app App) dictionary(ctx context.Context, lang, text string) ([]string, error) {
	meaning, err := translate.GoogleDictionary(ctx, lang, strings.ToLower(text))
	if err != nil {
		return nil, errors.Wrap(err)
	}
	out := make([]string, 0, 7)
	for _, data := range meaning.DictionaryData {
		for _, entry := range data.Entries {
			for _, senseFamily := range entry.SenseFamilies {
				for _, sense := range senseFamily.Senses {
					out = append(out, sense.Definition.Text)
				}
			}
		}
	}
	return out, nil
}