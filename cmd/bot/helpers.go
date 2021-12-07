/*
Helper functions
 */
package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	translate2 "github.com/armanokka/translobot/pkg/translate"
	iso6391 "github.com/emvi/iso-639-1"
	"github.com/go-errors/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k0kubun/pp"
	"golang.org/x/sync/errgroup"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"
)

func format(i int64) string {
    return strconv.FormatInt(i, 10)
}
// cutString cut string using runes by limit
func cutStringUTF16 (text string, limit int) string {
    points := utf16.Encode([]rune(text))
    if len(points) > limit {
        return string(utf16.Decode(points[:limit]))
    }
    return text
}

func in(arr []string, k string) bool {
    for _, v := range arr {
        if k == v {
             return true
        }
    }
    return false
}


func GetTickedCallbacks(keyboard tgbotapi.InlineKeyboardMarkup) []string {
    callbacks := make([]string, 0)
    for _, row := range keyboard.InlineKeyboard {
        for _, button := range row {
            if strings.HasPrefix(button.Text, "✅") {
                callbacks = append(callbacks, *button.CallbackData)
            }
        }
    }
    return callbacks
}

func TickByCallback(uniqueCallbackData string, keyboard *tgbotapi.InlineKeyboardMarkup) {
    var done bool
    for i1, row := range keyboard.InlineKeyboard {
        if done {
            break
        }
        for i2, button := range row {
            if *button.CallbackData == uniqueCallbackData && !strings.HasPrefix(*button.CallbackData, "✅ ") {
                keyboard.InlineKeyboard[i1][i2].Text = "✅ " + button.Text
                done = true
                break
            }
        }
    }
}

func UnTickByCallback(uniqueCallbackData string, keyboard *tgbotapi.InlineKeyboardMarkup) {
    var done bool
    for i1, row := range keyboard.InlineKeyboard {
        if done {
            break
        }
        for i2, button := range row {
            if *button.CallbackData == uniqueCallbackData {
                keyboard.InlineKeyboard[i1][i2].Text = strings.TrimPrefix(button.Text, "✅ ")
                done = true
                break
            }
        }
    }
}

func IsTicked(callback string, keyboard *tgbotapi.InlineKeyboardMarkup) bool {
    for _, row := range keyboard.InlineKeyboard {
        for _, button := range row {
            if *button.CallbackData != callback {
                continue
            }
            if strings.HasPrefix(button.Text, "✅") {
                return true
            }
        }
    }
    return false
}

func applyEntitiesHtml(text string, entities []tgbotapi.MessageEntity) string {
    if len(entities) == 0 {
        return text
    }

    encoded := utf16.Encode([]rune(text))
    pointers := make(map[int]string)

    for _, entity := range entities {
        var startTag string
        switch entity.Type {
        case "code", "pre":
            startTag = `<label class="notranslate"><code>`
        case "mention", "hashtag", "cashtag", "bot_command", "url", "email", "phone_number":
            startTag = `<label class="notranslate">` // very important to keep '<label class="notranslate">' strongly correct, without any spaces or another
        case "bold":
            startTag = `<b>`
        case "italic":
            startTag = `<i>`
        case "underline":
            startTag = `<u>`
        case "strikethrough":
            startTag = `<s>`
        case "text_link":
            startTag = `<a href="` + entity.URL + `">`
        case "text_mention":
            startTag = `<a href="tg://user?id=` + strconv.FormatInt(entity.User.ID, 10) + `">`
        }

        pointers[entity.Offset] += startTag


        //startTag = strings.TrimPrefix(startTag, "<")
        var endTag string
        switch entity.Type {
        case "code", "pre":
            endTag = "</code></label>" // very important to keep '</label>' strongly correct, without any spaces or another
        case "mention", "hashtag", "cashtag", "bot_command", "url", "email", "phone_number":
            endTag = `</label>`
        case "bold":
            endTag = `</b>`
        case "italic":
            endTag = `</i>`
        case "underline":
            endTag = `</u>`
        case "strikethrough":
            endTag = `</s>`
        case "text_link", "text_mention":
            endTag = `</a>`
        }
        pointers[entity.Offset+entity.Length] += endTag
    }

    var out = make([]uint16, 0, len(encoded))

    for i, ch := range encoded {
        if m, ok := pointers[i]; ok {
            out = append(out, utf16.Encode([]rune(m))...)
        }
        out = append(out, ch)

        if i == len(encoded) - 1 {
            if m, ok := pointers[i+1]; ok {
                out = append(out, utf16.Encode([]rune(m))...)
            }
        }
    }
    ret := string(utf16.Decode(out))
    ret = strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "").Replace(ret)
    ret = strings.ReplaceAll(ret, `<br>`, "\n")
    return ret
}



func makeArticle(id string, title, description, messageText string) tgbotapi.InlineQueryResultArticle {
    //keyboard := tgbotapi.NewInlineKeyboardMarkup(
    //    tgbotapi.NewInlineKeyboardRow(
    //        tgbotapi.InlineKeyboardButton{
    //            Text:                         "translate",
    //            SwitchInlineQueryCurrentChat: &description,
    //        }))
    return tgbotapi.InlineQueryResultArticle{
        Type:                "article",
        ID:                  id,
        Title:               title,
        InputMessageContent: map[string]interface{}{
            "message_text": messageText,
            "disable_web_page_preview":false,
        },
        //ReplyMarkup: &keyboard,
        //URL:                 "https://t.me/TransloBot?start=from_inline",
        //HideURL:             true,
        Description:         description,
    }
}

func inMapValues(m map[string]string, values ...string) bool {
    for _, v := range values {
        var ok bool
        for _, v1 := range m {
            if v  == v1 {
                ok = true
                break
            }
        }
        if !ok {
            return false
        }
    }
    return true
}

func (app app) sendSpeech(lang, text string, callbackID string, user User) error {
    sdec, err := translate2.TTS(lang, text)
    if err != nil {
        if err == translate2.ErrTTSLanguageNotSupported {
            call := tgbotapi.NewCallback(callbackID, user.Localize("%s language is not supported", iso6391.Name(lang)))
            call.ShowAlert = true
            app.bot.AnswerCallbackQuery(call)
            return nil
        }
        if e, ok := err.(translate2.HTTPError); ok {
            if e.Code == 500 || e.Code == 414 {
                call := tgbotapi.NewCallback(callbackID, "Too big text")
                call.ShowAlert = true
                app.bot.AnswerCallbackQuery(call)
                return nil
            }
        }
        app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "Iternal error"))
        pp.Println(err)
        return err
    }
    app.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, "⏳"))

    f, err := os.CreateTemp("", "")
    if err != nil {
        return err
    }
    defer func() {
        if err = f.Close(); err != nil {
            app.notifyAdmin(err)
        }
    }()
    _, err = f.Write(sdec)
    if err != nil {
        return err
    }
    audio := tgbotapi.NewAudio(user.ID, f.Name())
    audio.Title = text
    kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("❌", "delete")))
    audio.ReplyMarkup = kb
    app.bot.Send(audio)
    return nil
}


func WitAiSpeech(wav io.Reader, lang string, bits int) (string, error) {
    var key, ok = config.WitAPIKeys[lang]
    if !ok {
        return "", errors.New("no wit.ai key for lang " + lang)
    }
    req, err := http.NewRequest("POST", "https://api.wit.ai/speech?v=20210928&bits=" + strconv.Itoa(bits), wav)
    if err != nil {
        panic(err)
    }
    req.Header.Set("Authorization", "Bearer " + key)
    req.Header.Set("Content-Type", "audio/wave")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    parts := strings.Split(string(body), "\r\n")
    if len(parts) == 0 {
        return "", errors.New("empty parts: " + string(body))
    }
    var result struct{
        Text string `json:"text"`
        Error string `json:"error"`
        Code string `json:"code"`
    }
    if err = json.Unmarshal([]byte(parts[len(parts)-1]), &result); err != nil {
        return "", errors.WrapPrefix(err, string(body), 0)
    }
    if result.Error != "" {
        return "", errors.New(result.Error)
    }
    return result.Text, nil
}

func BuildSupportedLanguagesKeyboard() (tgbotapi.InlineKeyboardMarkup, error) {
    keyboard := tgbotapi.NewInlineKeyboardMarkup()
    for i, code := range config.BotLocalizedLangs {
        lang, ok := langs[code]
        if !ok {
            return tgbotapi.InlineKeyboardMarkup{}, errors.New("no such code "+ code + " in langs")
        }

        if i % 2 == 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name,  "set_bot_lang_and_register:"  + code)))
        } else {
            l := len(keyboard.InlineKeyboard)-1
            keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(lang.Emoji + " " + lang.Name, "set_bot_lang_and_register:"  + code))
        }
    }
    return keyboard, nil
}

type SuperTranslation struct {
    From string
    TranslatedText string
    Examples, Translations, Dictionary bool
}

func (app app) SuperTranslate(from, to, text string, entities []tgbotapi.MessageEntity) (ret SuperTranslation, err error) {
    text = html.EscapeString(text)
    text = applyEntitiesHtml(text, entities)

    var (
        rev = translate2.ReversoTranslation{}
        dict = translate2.GoogleDictionaryResponse{}
    )

    l := len(text)

    g, _ := errgroup.WithContext(context.Background())

    g.Go(func() error {
        if l > 100 {
            return nil
        }
        dict, err = translate2.GoogleDictionary(from, text)
        return err
    })

    g.Go(func() error {
        if l > 100 {
            return nil
        }
        if inMapValues(translate2.ReversoSupportedLangs(), from, to) && from != to {
            rev, err = translate2.ReversoTranslate(translate2.ReversoIso6392(from), translate2.ReversoIso6392(to), strings.ToLower(text))
        }
        return err
    })

    g.Go(func() error {
        tr, err := translate2.GoogleHTMLTranslate(from, to, strings.ReplaceAll(text, "\n", "<br>"))
        if err != nil {
            return err
        }
        if tr.Text == "" && text != "" {
            return err
        }
        ret.From = tr.From
        ret.TranslatedText = tr.Text
        ret.TranslatedText = strings.NewReplacer("<br> ", "<br>").Replace(ret.TranslatedText)
        ret.TranslatedText = strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "",  `<br>`, "\n").Replace(ret.TranslatedText)
        return nil
    })

    if err = g.Wait(); err != nil {
        app.notifyAdmin(err)
        return SuperTranslation{}, err
    }

    if len(rev.ContextResults.Results) > 0 {
        if len(rev.ContextResults.Results[0].SourceExamples) > 0 {
            ret.Examples = true
        }
        if rev.ContextResults.Results[0].Translation != "" {
            ret.Translations = true
        }
    }

    if dict.Status == 200 && dict.DictionaryData != nil {
        ret.Dictionary = true
    }

    return ret, nil
}

type Message struct {
    Text string
    Keyboard tgbotapi.ReplyKeyboardMarkup
}

// buildLettersKeyboard создает клаву 5x5 с заданными коллбэками. В коллбэк напишите %s и туда вставится letter
func buildLettersKeyboard(callbackData string) tgbotapi.InlineKeyboardMarkup {
    keyboard := tgbotapi.NewInlineKeyboardMarkup()
    for i, letter := range lettersSlice {
        btn := tgbotapi.NewInlineKeyboardButtonData(letter, fmt.Sprintf(callbackData, letter))
        if i % 5 == 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
            continue
        }
        last := len(keyboard.InlineKeyboard) - 1
        keyboard.InlineKeyboard[last] = append(keyboard.InlineKeyboard[last], btn)
    }
    return keyboard
}

// buildLettersKeyboard создает клаву с языками на данную букву. В коллбэк напишите %s и туда вставится код языка
func buildOneLetterKeyboard(letter, callbackData, backCallbackData string) tgbotapi.InlineKeyboardMarkup {
    keyboard := tgbotapi.NewInlineKeyboardMarkup()
    for i, lang := range lettersLangs[letter] {
        btn := tgbotapi.NewInlineKeyboardButtonData(lang.Name + lang.Emoji, fmt.Sprintf(callbackData, iso6391.CodeForName(lang.Name)))
        if i % 3 == 0 {
            keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
        } else {
            last := len(keyboard.InlineKeyboard) - 1
            keyboard.InlineKeyboard[last] = append(keyboard.InlineKeyboard[last], btn)
        }
    }
    keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
        tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("↩", backCallbackData)))
    return keyboard
}

// buildLangsPagination создает пагинацию как говорил F d
// в калбак передайте что-то типа set_my_lang:%s, где %s станет код выбранного языка
func buildLangsPagination(offset int, callback, back, next string) (tgbotapi.InlineKeyboardMarkup, error) {
    if offset < 0 || offset > len(codes) - 1 {
        return tgbotapi.InlineKeyboardMarkup{}, nil
    }
    out := tgbotapi.NewInlineKeyboardMarkup()
    for i, code := range codes[offset:offset+18] {
        lang, ok := langs[code]
        if !ok {
            return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("не нашел %s в langs", code)
        }
        btn := tgbotapi.NewInlineKeyboardButtonData(lang.Name + " " + lang.Emoji, fmt.Sprintf(callback, code))
        if i % 3 == 0 {
            out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
        } else {
            l := len(out.InlineKeyboard) - 1
            if l < 0 {
                l = 0
            }
            out.InlineKeyboard[l] = append(out.InlineKeyboard[l], btn)
        }
    }

    out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("<--- Back", back),
        tgbotapi.NewInlineKeyboardButtonData("Next --->", next)))
    return out, nil
}