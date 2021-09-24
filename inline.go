package main

import (
    "encoding/json"
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "html"
    "sort"
    "strconv"
    "sync"
    "time"
)

func handleInline(update *tgbotapi.InlineQuery) {
    started := time.Now()
    defer func() {
        pp.Println(time.Since(started).String())
    }()

    analytics.User(update.Query, update.From)
    
    warn := func(err error) {
        bot.AnswerInlineQuery(tgbotapi.InlineConfig{
            InlineQueryID:     update.ID,
            SwitchPMText:      "Error, sorry.",
            SwitchPMParameter: "from_inline",
        })
        WarnAdmin(err)
    }

    if update.Query == "" {
        bot.AnswerInlineQuery(
            tgbotapi.InlineConfig{
                InlineQueryID:     update.ID,
                SwitchPMText:      "Type text please",
                SwitchPMParameter: "from_inline",
            })
        return
    }

    var start int // смещение для пагинации
    if update.Offset != "" {
        var err error
        start, err = strconv.Atoi(update.Offset)
        if err != nil {
            warn(err)
            return
        }
    }

    l := len(codes)

    if start >= l {
        warn(errors.New("слишком большое смещение"))
        return
    }


    end := start + 50
    if end > l - 1 {
        end = l - 1
    }
    results := make([]interface{}, 0, 52)


    var wg sync.WaitGroup
    wg.Add(1)
    var user User
    go func() {
        defer wg.Done()
        user = NewUser(update.From.ID, warn)
        if user.Exists() {
            user.Fill()
        }
    }()

    tr, err := translate.GoogleHTMLTranslate("auto", "en", update.Query) // определяем язык
    if err != nil {
        warn(err)
    }
    from := tr.From

    if start == 0 {
        if from != user.MyLang {
            myLangTr, err := translate.GoogleHTMLTranslate("auto", user.MyLang, update.Query)
            if err != nil {
                warn(err)
                return
            }
            results = append(results, makeArticle("my_lang", iso6391.Name(user.MyLang) + " 🔥", html.UnescapeString(myLangTr.Text)))
        }
        if from != user.ToLang {
            toLangTr, err := translate.GoogleHTMLTranslate("auto", user.ToLang, update.Query)
            if err != nil {
                warn(err)
                return
            }
            results = append(results, makeArticle("to_lang", iso6391.Name(user.ToLang) + " 🔥", html.UnescapeString(toLangTr.Text)))
        }
    }


    for i, lang := range codes[start:end] {
        if lang == user.MyLang || lang == user.ToLang || lang == from {
            continue
        }
        wg.Add(1)
        lang := lang
        i := i
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

            results = append(results, makeArticle(strconv.Itoa(i), iso6391.Name(lang), html.UnescapeString(tr.Text)))
        }()
    }
    wg.Wait()

    sort.Slice(results[2:], func(i, j int) bool {
       return results[i].(tgbotapi.InlineQueryResultArticle).Title < results[j].(tgbotapi.InlineQueryResultArticle).Title
    })


    pmtext := "From: " + iso6391.Name(from)
    if update.Query == "" {
        pmtext = "Enter text"
    }

    if len(results) > 50 {
       results = results[:50]
    }
    nextoffset := strconv.Itoa(end)
    if end >= l - 1 {
        nextoffset = ""
    }

    if _, err := bot.AnswerInlineQuery(tgbotapi.InlineConfig{
        InlineQueryID:     update.ID,
        Results:           results,
        CacheTime:         InlineCacheTime,
        NextOffset: 	     nextoffset,
        IsPersonal:        true,
        SwitchPMText:      pmtext,
        SwitchPMParameter: "from_inline",
    }); err != nil {
        warn(err)
        WarnAdmin("из ошибки выше результат был таков:", results)
        j, _ := json.Marshal(results)
        WarnAdmin(string(j))
        pp.Println(results)
    }
    
    analytics.Bot(update.From.ID, "Inline succeeded", "Inline succeeded")

}
