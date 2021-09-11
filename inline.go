package main

import (
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "strconv"
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

    user := NewUser(update.From.ID, warn)
    user.Fill()
    languages := make([]string, 0, len(codes)+2)

    if user.Exists() {
        myLang := user.MyLang + " (your)"
        toLang := user.ToLang + " (your)"
        languages = []string{myLang, toLang}
    }

    languages = append(languages, codes...)
    
    var offset int
    var err error
    if update.Offset != "" { // Ищем смещение
        offset, err = strconv.Atoi(update.Offset)
        if err != nil {
            warn(err)
            return
        }
    }
    languagesLen := len(languages)
    
    if offset >= languagesLen { // Слишком большое смещение
        warn(err)
        return
    }
    
    end := offset + 10
    if end > languagesLen - 1 {
        end = languagesLen - 1
    }
    results := make([]interface{}, 0, 10)
    
    //var sponsorship Sponsorships
    //err = db.Model(&Sponsorships{}).Where("start <= current_timestamp AND finish >= current_timestamp").Limit(1).Find(&sponsorship).Error
    //if err != nil {
    //   WarnAdmin(err)
    //} else { // no error
    //   langs := strings.Split(sponsorship.ToLangs, ",")
    //   if inArray(update.From.LanguageCode, langs) {
    //       inputMessageContent := map[string]interface{}{
    //           "message_text": sponsorship.Text,
    //           "parse_mode": tgbotapi.ModeHTML,
    //           "disable_web_page_preview":false,
    //       }
    //
    //       results = append(results, tgbotapi.InlineQueryResultArticle{
    //           Type:                "article",
    //           ID:                  "advertise",
    //           Title:               "Ad",
    //           InputMessageContent: inputMessageContent,
    //           URL:                 "https://t.me/TransloBot?start=from_inline",
    //           HideURL:             true,
    //           Description:         sponsorship.Text,
    //       })
    //   }
    //}
    
    
    for ;offset < end; offset++ {
        to := languages[offset] // language code to translate
        tr, err := translate.GoogleHTMLTranslate("auto", to, update.Query)
        if err != nil {
            warn(err)
            return
        }
        if tr.Text == "" {
            continue // ну не вышло, так не вышло, че бубнить-то
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
            ID:                  strconv.Itoa(offset),
            Title:               iso6391.Name(to),
            InputMessageContent: inputMessageContent,
            ReplyMarkup: &keyboard,
            URL:                 "https://t.me/TransloBot?start=from_inline",
            HideURL:             true,
            Description:         tr.Text,
        })
    }
    
    var nextOffset int
    if end < len(languages) {
        nextOffset = end
    }
    pmtext := "Translo"
    if update.Query == "" {
        pmtext = "Enter text"
    }
    
    if _, err = bot.AnswerInlineQuery(tgbotapi.InlineConfig{
        InlineQueryID:     update.ID,
        Results:           results,
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
