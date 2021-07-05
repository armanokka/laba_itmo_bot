package main

import (
    "github.com/armanokka/translobot/translate"
    iso6391 "github.com/emvi/iso-639-1"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/k0kubun/pp"
    "strconv"
)

func handleInline(update *tgbotapi.Update) {
    warn := func(err error) {
        bot.AnswerInlineQuery(tgbotapi.InlineConfig{
            InlineQueryID:     update.InlineQuery.ID,
            SwitchPMText:      "Error, sorry.",
            SwitchPMParameter: "from_inline",
        })
        WarnAdmin(err)
    }
    
    
    languages := []string{"en","ru","la","ja","ar","fr","de","af","uk","uz","es","ko","zh","hi","bn","pt","mr","te","ms","tr","vi","ta","ur","jv","it","fa","gu","ab","aa","ak","sq","am","an","hy","as","av","ae","ay","az","bm","ba","eu","be","bh","bi","bs","br","bg","my","ca","ch","ce","ny","cv","kw","co","cr","hr","cs","da","dv","nl","eo","et","ee","fo","fj","fi","ff","gl","ka","el","gn","ht","ha","he","hz","ho","hu","ia","id","ie","ga","ig","ik","io","is","iu","kl","kn","kr","ks","kk","km","ki","rw","ky","kv","kg","ku","kj","la","lb","lg","li","ln","lo","lt","lu","lv","gv","mk","mg","ml","mt","mi","mh","mn","na","nv","nb","nd","ne","ng","nn","no","ii","nr","oc","oj","cu","om","or","os","pa","pi","pl","ps","qu","rm","rn","ro","sa","sc","sd","se","sm","sg","sr","gd","sn","si","sk","sl","so","st","su","sw","ss","sv","tg","th","ti","bo","tk","tl","tn","to","ts","tt","tw","ty","ug","ve","vo","wa","cy","wo","fy","xh","yi","yo","za"}
    
    from, err := translate.DetectLanguageGoogle(update.InlineQuery.Query)
    if err != nil {
        warn(err)
        return
    }
    
    var offset int
    
    if update.InlineQuery.Offset != "" { // Ищем смещение
        offset, err = strconv.Atoi(update.InlineQuery.Offset)
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
    for ;offset < end; offset++ {
        to := languages[offset] // language code to translate
        tr, err := translate.GoogleTranslate(from, to, update.InlineQuery.Query)
        if err != nil {
            warn(err)
            return
        }
        if tr.Text == "" {
            continue // ну не вышло, так не вышло, че бубнить-то
        }
        inputMessageContent := map[string]interface{}{
            "message_text":tr.Text,
            "disable_web_page_preview":true,
        }
        results = append(results, tgbotapi.InlineQueryResultArticle{
            Type:                "article",
            ID:                  strconv.Itoa(offset),
            Title:               iso6391.Name(to),
            InputMessageContent: inputMessageContent,
            URL:                 "https://t.me/TransloBot?start=from_inline",
            HideURL:             true,
            Description:         cutString(tr.Text, 40),
        })
    }
    
    var nextOffset int
    if end < len(languages) {
        nextOffset = end
    }
    
    ok, err := bot.AnswerInlineQuery(tgbotapi.InlineConfig{
        InlineQueryID:     update.InlineQuery.ID,
        Results:           results,
        CacheTime:         300,
        NextOffset: 	   strconv.Itoa(nextOffset),
        IsPersonal:        false,
        SwitchPMText:      "Translo",
        SwitchPMParameter: "from_inline",
    })
    if err != nil || !ok {
        pp.Println(err)
    }
}
