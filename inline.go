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
    
    var registered bool
    err := db.Model(&Users{}).Raw("SELECT EXISTS(SELECT usings FROM users WHERE id=?)", update.From.ID).Find(&registered).Error
    if err != nil {
        warn(err)
        return
    }
    languages := make([]string, 0)
    if registered {
        var user Users
        err := db.Model(&Users{}).Select("my_lang", "to_lang").Where("id = ?", update.From.ID).Find(&user).Error
        if err != nil {
            warn(err)
            return
        }
        languages = []string{user.MyLang, user.ToLang}
    }
    
    languages = append(languages, []string{"en","ru","la","ja","ar","fr","de","af","uk","uz","es","ko","zh","hi","bn","pt","mr","te","ms","tr","vi","ta","ur","jv","it","fa","gu","ab","aa","ak","sq","am","an","hy","as","av","ae","ay","az","bm","ba","eu","be","bh","bi","bs","br","bg","my","ca","ch","ce","ny","cv","kw","co","cr","hr","cs","da","dv","nl","eo","et","ee","fo","fj","fi","ff","gl","ka","el","gn","ht","ha","he","hz","ho","hu","ia","id","ie","ga","ig","ik","io","is","iu","kl","kn","kr","ks","kk","km","ki","rw","ky","kv","kg","ku","kj","la","lb","lg","li","ln","lo","lt","lu","lv","gv","mk","mg","ml","mt","mi","mh","mn","na","nv","nb","nd","ne","ng","nn","no","ii","nr","oc","oj","cu","om","or","os","pa","pi","pl","ps","qu","rm","rn","ro","sa","sc","sd","se","sm","sg","sr","gd","sn","si","sk","sl","so","st","su","sw","ss","sv","tg","th","ti","bo","tk","tl","tn","to","ts","tt","tw","ty","ug","ve","vo","wa","cy","wo","fy","xh","yi","yo","za"}...)
    
    //from, err := translate.DetectLanguageGoogle(update.Query)
    //if err != nil {
    //    warn(err)
    //    return
    //}
    
    var offset int
    
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
    //var sponsorText string
    //err = db.Model(&Sponsorships{}).Select("text", "to_langs").Where("start <= current_timestamp AND finish >= current_timestamp").Limit(1).Find(&sponsorship).Error
    //if err != nil {
    //    WarnAdmin(err)
    //} else { // no error
    //    langs := strings.Split(sponsorship.ToLangs, ",")
    //    if inArray(update.From.LanguageCode, langs) {
    //        sponsorText = "\n⚡️" + sponsorship.Text
    //    }
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
            "message_text":tr.Text + "\n\n» @translobot",
            "parse_mode": tgbotapi.ModeHTML,
            "disable_web_page_preview":true,
        }
        query := cutString(tr.Text, 255)
        keyboard := tgbotapi.NewInlineKeyboardMarkup(
            tgbotapi.NewInlineKeyboardRow(
                tgbotapi.InlineKeyboardButton{
                    Text:                         "translate",
                    SwitchInlineQueryCurrentChat: &query,
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
