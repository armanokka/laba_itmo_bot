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

    var offset int
    if update.Offset != "" { // Ищем смещение
        var err error
        offset, err = strconv.Atoi(update.Offset)
        if err != nil {
            warn(err)
            return
        }
    }
    l := len(codes)

    if offset >= l { // Слишком большое смещение
        warn(errors.New("слишком большое смещение"))
        return
    }
    
    end := offset + 50
    if end > l - 1 {
        end = l - 1
    }
    results := make([]interface{}, 0, 50)

    var from string
    
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
            if from == "" {
                from = tr.From
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
                ID:                  strconv.Itoa(offs), // +2, потому что могут быть языки юзера
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

    var nextOffset string
    if end < l {
        nextOffset = strconv.Itoa(end)
    }
    pmtext := "From: " + iso6391.Name(from)
    if update.Query == "" {
        pmtext = "Enter text"
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
