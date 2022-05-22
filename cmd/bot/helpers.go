/*
Helper functions
*/
package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/lingvo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func parseKeyboard(messageText string) tgbotapi.InlineKeyboardMarkup {
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
		return tgbotapi.InlineKeyboardMarkup{}
	}
	return keyboard
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
			if k == v || diff(k, v) == 1 {
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
		pointers[entity.Offset+entity.Length] = endTag + pointers[entity.Offset+entity.Length]
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
	return strings.NewReplacer(`<label class="notranslate">`, "", `</label>`, "", "<br>", "\n").Replace(string(utf16.Decode(out)))
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
		return "", fmt.Errorf("no wit.ai key for lang " + lang)
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
		return "", fmt.Errorf("empty parts: " + string(body))
	}
	var result struct {
		Text  string `json:"text"`
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err = json.Unmarshal([]byte(parts[len(parts)-1]), &result); err != nil {
		return "", errors.Wrap(err)
	}
	if result.Error != "" {
		return "", fmt.Errorf(result.Error)
	}
	return result.Text, nil
}

func index(arr []string, k string) int {
	for i, v := range arr {
		if k == v {
			return i
		}
	}
	return 0
}

func highlightDiffs(s1, s2, start, stop string) string {
	first := strings.Fields(s1)
	highlited := false
	var out bytes.Buffer
	for i, w := range strings.Fields(s2) {
		idx := index(first, w)
		if idx == 0 && first[0] != w {
			if !highlited {
				highlited = true
				w = start + w
			}

		} else if highlited {
			highlited = false
			out.WriteString(stop)
		}
		if i > 0 {
			out.WriteString(" ")
		}
		out.WriteString(w)
	}
	if highlited {
		out.WriteString("</b>")
	}
	return out.String()
}

// buildLangsPagination создает пагинацию как говорил F d
// в калбак передайте что-то типа set_my_lang:%s, где %s станет код выбранного языка
func buildLangsPagination(user tables.Users, offset int, count int, tickLang, buttonSelectLangCallback, buttonBackCallback, buttonNextCallback string) (tgbotapi.InlineKeyboardMarkup, error) {
	if offset < 0 || offset > len(codes[user.Lang])-1 {
		return tgbotapi.InlineKeyboardMarkup{}, nil
	}
	out := tgbotapi.NewInlineKeyboardMarkup()

	for i, code := range codes[user.Lang][offset : offset+count] {
		lang, ok := langs[user.Lang][code]
		//if offset+count <= 18 {
		//	lang += " 📌"
		//}
		if code == tickLang {
			lang += "✅"
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
		tgbotapi.NewInlineKeyboardButtonData("⬅️", buttonBackCallback),
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(offset)+"/"+strconv.Itoa(len(codes[user.Lang])/18*18), buttonBackCallback),
		tgbotapi.NewInlineKeyboardButtonData("➡️", buttonNextCallback)))
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

func concatNonEmpty(separator string, ss ...string) string {
	b := new(bytes.Buffer)
	for i, s := range ss {
		if strings.TrimSpace(s) == "" {
			continue
		}
		if i > 0 {
			b.WriteString(separator)
		}
		b.WriteString(s)
	}
	return b.String()
}

// diff counts difference between s1 and s2 by comparing their characters
func diff(s1, s2 string) (n int) {
	reader := strings.NewReader(s1)
	for _, ch := range s2 {
		r, _, err := reader.ReadRune()
		if err == nil {
			if r != ch {
				n++
			}
			continue
		}
		break
	}
	if l1, l2 := len(s1), len(s2); l1 != l2 {
		if l1 > l2 {
			n += l1 - l2
		} else {
			n += l2 - l1
		}
	}
	return n
}

func buildKeyboard(from, to string, ret Keyboard) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔉", fmt.Sprintf("speech:%s:%s", from, to))))
	if len(ret.Examples) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("💬", fmt.Sprintf("examples:%s:%s", from, to)))
	}
	if len(ret.Dictionary) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("📖", fmt.Sprintf("dictionary:%s:%s", from, to)))
	}
	if len(ret.Paraphrase) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("✨", fmt.Sprintf("paraphrase:%s:%s", from, to)))
	}
	if len(ret.ReverseTranslations) > 0 {
		keyboard.InlineKeyboard[0] = append(keyboard.InlineKeyboard[0], tgbotapi.NewInlineKeyboardButtonData("📚", fmt.Sprintf("reverse_translations:%s:%s", from, to)))
	}
	return keyboard
}

func IsCtxError(err error) bool {
	if e, ok := err.(errors.Error); ok {
		return IsCtxError(e.Err)
	}
	return errors.Is(err, context.Canceled)
}

func writeLingvo(lingvo []lingvo.Dictionary) string {
	out := new(bytes.Buffer)
	usedWords := make([]string, 0, 10)
	//examples := make([]string, 0, 3)
	lastLineLen := 0
	for _, r := range lingvo {
		words := strings.FieldsFunc(r.Translations, func(r rune) bool {
			return r == ',' || r == ';'
		})
		for _, word := range words {
			if word == "" {
				continue
			}
			if len(usedWords) > 11 {
				break
			}
			word = strings.TrimSpace(word)
			if inFuzzy(usedWords, word) {
				continue
			}
			usedWords = append(usedWords, word)
			word += "; "
			if lastLineLen+len(word) > 40 {
				out.WriteString("\n")
				lastLineLen = len(word)
			} else {
				lastLineLen += len(word)
			}
			out.WriteString(word)
		}

		//examplesSlice := strings.FieldsFunc(r.Examples, func(r rune) bool {
		//	return r == '\n' || r == '\r' || r == ','
		//})
		//for _, e := range examplesSlice {
		//	if strings.TrimSpace(e) == "" {
		//		continue
		//	}
		//	if len(examples) > 2 {
		//		break
		//	}
		//	e = strings.TrimSuffix(e, "...")
		//	examples = append(examples, e)
		//}
	}
	//if len(examples) > 0 {
	//	out.WriteString("\n\n")
	//	out.WriteString(strings.Join(examples, "\n"))
	//}
	return out.String()
}
