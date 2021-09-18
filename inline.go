package main

import (
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
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

    var offset int
    var err error
    if update.Offset != "" { // Ищем смещение
        offset, err = strconv.Atoi(update.Offset)
        if err != nil {
            warn(err)
            return
        }
    }
    languagesLen := len(codes)
    
    if offset >= languagesLen { // Слишком большое смещение
        warn(err)
        return
    }
    
    end := offset + 50
    if end > languagesLen - 1 {
        end = languagesLen - 1
    }
    results := make([]interface{}, 0, 50)

    user := NewUser(update.From.ID, warn)
    if user.Exists() {
        user.Fill()

        tr, err := translate.GoogleHTMLTranslate("auto", user.MyLang, update.Query)
        if err != nil {
            warn(err)
            return
        }

        inputMessageContent := map[string]interface{}{
            "message_text": tr.Text,
            "parse_mode": tgbotapi.ModeHTML,
            "disable_web_page_preview":false,
        }

        keyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.InlineKeyboardButton{
                    Text:                         "translate",
                    SwitchInlineQueryCurrentChat: &tr.Text,
                }))
        results = append(results, tgbotapi.InlineQueryResultArticle{
            Type:                "article",
            ID:                  "0", // надо для рекламы
            Title:               iso6391.Name(user.MyLang),
            InputMessageContent: inputMessageContent,
            ReplyMarkup: &keyboard,
            URL:                 "https://t.me/TransloBot?start=from_inline",
            HideURL:             true,
            Description:         tr.Text,
        })



        tr, err = translate.GoogleHTMLTranslate("auto", user.ToLang, update.Query)
        if err != nil {
            warn(err)
            return
        }

        inputMessageContent = map[string]interface{}{
            "message_text": tr.Text,
            "parse_mode": tgbotapi.ModeHTML,
            "disable_web_page_preview":false,
        }

        keyboard = tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.InlineKeyboardButton{
                    Text:                         "translate",
                    SwitchInlineQueryCurrentChat: &tr.Text,
                }))
        results = append(results, tgbotapi.InlineQueryResultArticle{
            Type:                "article",
            ID:                  "1",
            Title:               iso6391.Name(user.ToLang),
            InputMessageContent: inputMessageContent,
            ReplyMarkup: &keyboard,
            URL:                 "https://t.me/TransloBot?start=from_inline",
            HideURL:             true,
            Description:         tr.Text,
        })
    }

    //var l int
    //if err = db.Model(&Ads{}).Raw("SELECT COUNT(*) FROM ads LIMIT 1").Find(&l).Error; err == nil {
    //    // no error
    //
    //    if l > 0 {
    //        var offer Ads
    //        if err = db.Model(&Ads{}).Where("start_date <= current_timestamp AND finish_date >= current_timestamp").Limit(1).Find(&offer).Error; err == nil {
    //            // here no error
    //            inputMessageContent := map[string]interface{}{
    //                "message_text": offer.Content,
    //                "parse_mode": tgbotapi.ModeHTML,
    //                "disable_web_page_preview":false,
    //            }
    //            results = append(results, tgbotapi.InlineQueryResultArticle{
    //                Type:                "article",
    //                ID:                  "0",
    //                Title:               "Advertise",
    //                InputMessageContent: inputMessageContent,
    //                ReplyMarkup: nil,
    //                URL:                 "https://t.me/TransloBot?start=from_inline",
    //                HideURL:             true,
    //                Description:         offer.Content,
    //            })
    //        }
    //    }
    //}


    
    var wg sync.WaitGroup
    for ;offset < end; offset++ {
        offs := offset
        wg.Add(1)
        go func() {
            defer wg.Done()
            to := codes[offs] // language code to translate
            tr, err := translate.GoogleHTMLTranslate("auto", to, update.Query)
            if err != nil {
                warn(err)
                return
            }
            if tr.Text == "" {
                return // ну не вышло, так не вышло, че бубнить-то
            }

            inputMessageContent := map[string]interface{}{
                "message_text": tr.Text,
                "parse_mode": tgbotapi.ModeHTML,
                "disable_web_page_preview":false,
            }

            keyboard := tgbotapi.NewInlineKeyboardMarkup(
                tgbotapi.NewInlineKeyboardRow(
                    tgbotapi.InlineKeyboardButton{
                        Text:                         "translate",
                        SwitchInlineQueryCurrentChat: &tr.Text,
                    }))
            results = append(results, tgbotapi.InlineQueryResultArticle{
                Type:                "article",
                ID:                  strconv.Itoa(offs+2), // +2, потому что могут быть языки юзера
                Title:               iso6391.Name(to),
                InputMessageContent: inputMessageContent,
                ReplyMarkup: &keyboard,
                URL:                 "https://t.me/TransloBot?start=from_inline",
                HideURL:             true,
                Description:         tr.Text,
            })
        }()
    }
    wg.Wait()
    sort.Slice(results, func(i, j int) bool {
        a, err := strconv.Atoi(results[i].(tgbotapi.InlineQueryResultArticle).ID)
        if err != nil {
            warn(err)
        }
        b, err := strconv.Atoi(results[j].(tgbotapi.InlineQueryResultArticle).ID)
        if err != nil {
            warn(err)
        }
        return a < b
    })

    var nextOffset int
    if end < len(codes) {
        nextOffset = end
    }
    pmtext := "Translo"
    if update.Query == "" {
        pmtext = "Enter text"
    }
    
    if _, err = bot.AnswerInlineQuery(tgbotapi.InlineConfig{
        InlineQueryID:     update.ID,
        Results:           results[:50],
        CacheTime:         InlineCacheTime,
        NextOffset: 	   strconv.Itoa(nextOffset),
        IsPersonal:        false,
        SwitchPMText:      pmtext,
        SwitchPMParameter: "from_inline",
    }); err != nil {
        pp.Println(err)
        warn(err)
    }
    
    analytics.Bot(update.From.ID, "Inline succeeded", "Inline succeeded")
}
