package main

import (
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "sort"
    "strconv"
    "sync"
)

func handleInline(update *tgbotapi.InlineQuery) {
    analytics.User(update.Query, update.From)
    
    warn := func(err error) {
        bot.AnswerInlineQuery(tgbotapi.InlineConfig{
            InlineQueryID:     update.ID,
            SwitchPMText:      "Error, sorry.",
            SwitchPMParameter: "from_inline",
        })
        WarnAdmin(err)
    }

    var offset int // смещение для пагинации
    if update.Offset != "" {
        var err error
        offset, err = strconv.Atoi(update.Offset)
        if err != nil {
            warn(err)
            return
        }
    }

    if offset > len(codes) { // Слишком большое смещение
        warn(errors.New("слишком большое смещение"))
        return
    }
    
    end := offset + 50
    if end > len(codes) - 1 {
        end = len(codes) - 1
    }
    results := make([]interface{}, 0, 52)

    user := NewUser(update.From.ID, warn)
    if user.Exists() {
        myLangTr, err := translate.GoogleHTMLTranslate("auto", user.MyLang, update.Query)
        if err != nil {
            warn(err)
            return
        }
        results = append(results, makeArticle(iso6391.Name(user.MyLang) + " 🔥", myLangTr.Text))

        toLangTr, err := translate.GoogleHTMLTranslate("auto", user.MyLang, update.Query)
        if err != nil {
            warn(err)
            return
        }
        results = append(results, makeArticle(iso6391.Name(user.ToLang) + " 🔥", toLangTr.Text))
    }

    var from string

    if len(results) > 0 {
        end -= len(results)
    }

    var wg sync.WaitGroup
    for _, lang := range codes[offset:end]{
        wg.Add(1)
        lang := lang
        go func() {
            defer wg.Done()
            tr, err := translate.GoogleHTMLTranslate("auto", lang, update.Query)
            if err != nil {
                warn(err)
                return
            }
            if tr.Text == "" {
                return // ну не вышло, так не вышло, че бубнить-то
            }
            if from == "" {
                from = tr.From
            }

            results = append(results, makeArticle(iso6391.Name(tr.From), tr.Text))
        }()
    }

    wg.Wait()
    sort.Slice(results, func(i, j int) bool {
        return results[i].(tgbotapi.InlineQueryResultArticle).Title < results[j].(tgbotapi.InlineQueryResultArticle).Title
    })

    var nextOffset string = strconv.Itoa(end)

    pmtext := "From: " + iso6391.Name(from)
    if update.Query == "" {
        pmtext = "Enter text"
    }

    if len(results) > 50 {
        results = results[:50]
    }
    
    if _, err := bot.AnswerInlineQuery(tgbotapi.InlineConfig{
        InlineQueryID:     update.ID,
        Results:           results,
        CacheTime:         InlineCacheTime,
        NextOffset: 	     nextOffset,
        IsPersonal:        false,
        SwitchPMText:      pmtext,
        SwitchPMParameter: "from_inline",
    }); err != nil {
        warn(err)
        WarnAdmin("из ошибки выше результат был таков:", results)
    }
    
    analytics.Bot(update.From.ID, "Inline succeeded", "Inline succeeded")
}
