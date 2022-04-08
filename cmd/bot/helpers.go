/*
Helper functions
*/
package bot

import (
	"encoding/json"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/go-errors/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/paul-mannino/go-fuzzywuzzy"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf16"
)

// cutString cut string using runes by limit
func cutStringUTF16(text string, limit int) string {
	points := utf16.Encode([]rune(text))
	if len(points) > limit {
		return string(utf16.Decode(points[:limit]))
	}
	return text
}

func parseKeyboard(messageText string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	if messageText != "Empty" {
		lines := strings.Split(messageText, "\n")
		for _, line := range lines {
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				continue
			}
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(parts[0], parts[1])))
		}
	}
	if reflect.DeepEqual(keyboard, tgbotapi.NewInlineKeyboardMarkup()) {
		return nil
	}
	return &keyboard
}

func in(arr []string, keys ...string) bool {
	for _, k := range keys {
		exists := false
		for _, v := range arr {
			if k == v {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}
	return true
}

func inFuzzy(arr []string, keys ...string) bool {
	for _, k := range keys {
		exists := false
		for _, v := range arr {
			if k == v || fuzzy.EditDistance(k, v) == 1 {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}
	return true
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

/*
func applyEntitiesHtml(text string, entities []tgbotapi.MessageEntity) string {
    if len(entities) == 0 {
        return text
    }

    //text = strings.ReplaceAll(text, "\n", "<br>")

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
           pp.Println("adding", m)
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
*/

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

		if i == len(encoded)-1 {
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

func inMapValues(m map[string]string, values ...string) bool {
	for _, v := range values {
		var ok bool
		for _, v1 := range m {
			if v == v1 {
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

func WitAiSpeech(wav io.Reader, lang string, bits int) (string, error) {
	var key, ok = config.WitAPIKeys[lang]
	if !ok {
		return "", errors.New("no wit.ai key for lang " + lang)
	}
	req, err := http.NewRequest("POST", "https://api.wit.ai/speech?v=20210928&bits="+strconv.Itoa(bits), wav)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "audio/wave")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
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
	var result struct {
		Text  string `json:"text"`
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err = json.Unmarshal([]byte(parts[len(parts)-1]), &result); err != nil {
		return "", errors.WrapPrefix(err, string(body), 0)
	}
	if result.Error != "" {
		return "", errors.New(result.Error)
	}
	return result.Text, nil
}

func BuildSupportedLanguagesKeyboard(user tables.Users) (tgbotapi.InlineKeyboardMarkup, error) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()
	for i, code := range config.BotLocalizedLangs {
		lang, ok := langs[user.Lang][code]
		if !ok {
			return tgbotapi.InlineKeyboardMarkup{}, errors.New("no such code " + code + " in langs")
		}

		if i%2 == 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lang, "set_bot_lang_and_register:"+code)))
		} else {
			l := len(keyboard.InlineKeyboard) - 1
			keyboard.InlineKeyboard[l] = append(keyboard.InlineKeyboard[l], tgbotapi.NewInlineKeyboardButtonData(lang, "set_bot_lang_and_register:"+code))
		}
	}
	return keyboard, nil
}

type SuperTranslation struct {
	From                                            string
	TranslatedText                                  string
	Examples, Translations, Dictionary, Suggestions bool
}

type Message struct {
	Text     string
	Keyboard tgbotapi.ReplyKeyboardMarkup
}

func reverse(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// buildLangsPagination создает пагинацию как говорил F d
// в калбак передайте что-то типа set_my_lang:%s, где %s станет код выбранного языка
func buildLangsPagination(user tables.Users, offset int, count int, exceptLang, buttonSelectLangCallback, buttonBackCallback, buttonNextCallback string) (tgbotapi.InlineKeyboardMarkup, error) {
	if offset < 0 || offset > len(codes[user.Lang])-1 {
		return tgbotapi.InlineKeyboardMarkup{}, nil
	}
	out := tgbotapi.NewInlineKeyboardMarkup()

	if count == 0 {
		offset -= 19
		count += 19
	}
	for i, code := range codes[user.Lang][offset : offset+count] {
		if code == exceptLang {
			continue
		}
		lang, ok := langs[user.Lang][code]
		if i+offset < 19 {
			lang += " 📌"
		}
		if !ok {
			return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("не нашел %s в langs", code)
		}

		callback := fmt.Sprintf(buttonSelectLangCallback, code)

		btn := tgbotapi.NewInlineKeyboardButtonData(lang, callback)
		if i%3 == 0 {
			out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
		} else {
			l := len(out.InlineKeyboard) - 1
			if l < 0 {
				l = 0
			}
			if len(out.InlineKeyboard) == 0 {
				out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(btn))
				continue
			}
			out.InlineKeyboard[l] = append(out.InlineKeyboard[l], btn)
		}
	}

	out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("<--- Back", buttonBackCallback),
		tgbotapi.NewInlineKeyboardButtonData("Next --->", buttonNextCallback)))
	return out, nil
}

func randid(seed int64) string {
	rand.Seed(seed)
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func remove(arr []string, k string) []string {
	ret := make([]string, 0, len(arr))
	for _, v := range arr {
		if v == k {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

func tickUntick(keyboard tgbotapi.InlineKeyboardMarkup, tickCallback, untickCallback, prefix string) tgbotapi.InlineKeyboardMarkup {
	for i1, row := range keyboard.InlineKeyboard {
		for i2, btn := range row {
			callback := *btn.CallbackData
			if callback == tickCallback {
				btn.Text = prefix + btn.Text
				keyboard.InlineKeyboard[i1][i2] = btn
			}
			if callback == untickCallback {
				btn.Text = strings.TrimPrefix(btn.Text, prefix)
				keyboard.InlineKeyboard[i1][i2] = btn
			}
		}
	}
	return keyboard
}

func min(ints ...float64) float64 {
	if len(ints) == 0 {
		return -1
	}
	min := ints[0]
	for _, v := range ints {
		if v < min {
			min = v
		}
	}
	return min
}

func inlineTranslationKeyboard(lang string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Translate", "none")))
}
