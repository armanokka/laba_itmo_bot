package main

import (
    "errors"
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "github.com/sirupsen/logrus"
    "html"
    "sort"
    "strconv"
    "sync"
    "sync/atomic"
)

func handleInline(update tgbotapi.InlineQuery) {
    analytics.User(update.Query, update.From)

    warn := func(err error) {
        bot.AnswerInlineQuery(tgbotapi.InlineConfig{
            InlineQueryID:     update.ID,
            SwitchPMText:      "Error, sorry.",
            SwitchPMParameter: "from_inline",
        })
        WarnAdmin(err)
        logrus.Error(err)
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
    var user = NewUser(update.From.ID, warn)
    userExists := user.Exists()
    if userExists {
        user.Fill()
    }

    tr, err := translate.GoogleHTMLTranslate("auto", "en", update.Query) // определяем язык
    if err != nil {
        warn(err)
    }

    from := tr.From
    if from == "" {
        from = user.MyLang
    }
    //sortOffset := int64(len(popular))
    var sortOffset int64

    if start == 0 && userExists { // Юзер существует

        if from != user.MyLang {
            sortOffset++
            //fromLang := translate.ReversoIso6392(from)
            //toLang := translate.ReversoIso6392(user.MyLang)
            //
            //if fromLang != "" && toLang != "" {
            //
            //    wg.Add(1)
            //    go func() {
            //        defer wg.Done()
            //        tr, err := translate.ReversoTranslate(fromLang, toLang, update.Query)
            //        if err != nil {
            //            warn(err)
            //            return
            //        }
            //        var description string
            //        for i, result := range tr.ContextResults.Results {
            //            if i > 0 {
            //                description += ", "
            //            }
            //            description += result.Translation
            //        }
            //
            //        if len(tr.Translation) == 0 {
            //            return
            //        }
            //
            //        results = append(results, makeArticle("my_lang-" + tr.Translation[0], iso6391.Name(user.MyLang) + " 🔥", html.UnescapeString(description), html.UnescapeString(tr.Translation[0])))
            //    }()
            //} else {
                wg.Add(1)
                go func() {
                    defer wg.Done()
                    myLangTr, err := translate.GoogleHTMLTranslate(from, user.MyLang, update.Query)
                    if err != nil {
                        warn(err)
                        return
                    }
                    html.UnescapeString(myLangTr.Text)
                    results = append(results, makeArticle("my_lang", iso6391.Name(user.MyLang) + " 🔥", myLangTr.Text, myLangTr.Text))
                }()
            //}

        }
        if from != user.ToLang {

            sortOffset++
            //fromLang := translate.ReversoIso6392(from)
            //toLang := translate.ReversoIso6392(user.ToLang)
            //
            //if fromLang != "" && toLang != "" {
            //    wg.Add(1)
            //    go func() {
            //        defer wg.Done()
            //        tr, err := translate.ReversoTranslate(fromLang, toLang, update.Query)
            //        if err != nil {
            //            warn(err)
            //            return
            //        }
            //
            //        var description string
            //        for i, result := range tr.ContextResults.Results {
            //            if i > 0 {
            //                description += ", "
            //            }
            //            description += result.Translation
            //        }
            //
            //        if len(tr.Translation) == 0 {
            //            return
            //        }
            //        results = append(results, makeArticle("to_lang-" + tr.Translation[0], iso6391.Name(user.ToLang) + " 🔥", html.UnescapeString(description), html.UnescapeString(tr.Translation[0])))
            //    }()
            //} else {
                wg.Add(1)
                go func() {
                    defer wg.Done()
                    toLangTr, err := translate.GoogleHTMLTranslate("auto", user.ToLang, update.Query)
                    if err != nil {
                        warn(err)
                        return
                    }
                    toLangTr.Text = html.UnescapeString(toLangTr.Text)
                    results = append(results, makeArticle("to_lang", iso6391.Name(user.ToLang) + " 🔥", toLangTr.Text, toLangTr.Text))
               }()
            //}
        }
    }

    wg.Wait()

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
            tr.Text = html.UnescapeString(tr.Text)

            if tr.Text == "" {
                return // ну не вышло, так не вышло, че бубнить-то
            }
            results = append(results, makeArticle(strconv.Itoa(i), iso6391.Name(lang), tr.Text, tr.Text))
        }()
    }
    wg.Wait()
    sortOffsetLoaded := int(atomic.LoadInt64(&sortOffset))
    sort.Slice(results[sortOffset:], func(i, j int) bool {
       return results[i+sortOffsetLoaded].(tgbotapi.InlineQueryResultArticle).Title < results[j+sortOffsetLoaded].(tgbotapi.InlineQueryResultArticle).Title
    })

    pmtext := "From: " + langs[from].Name
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

    pp.Println(results)

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
        logrus.Error(err)
        pp.Println(results)
    }
    
    analytics.Bot(update.From.ID, "Inline succeeded", "Inline succeeded")

    if user.MyLang != "" {
        user.UpdateLastActivity()
        user.WriteUserLog(update.Query)
        user.WriteBotLog("inline_succeeded", "")
    }
}
