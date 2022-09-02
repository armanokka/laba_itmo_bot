/*
Helper functions
*/
package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/lingvo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/html"
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

var SupportedFormattingTags = []string{"b", "strong", "i", "em", "u", "ins", "s", "strike", "del", "span", "tg-spoiler", "a", "code", "pre"}

func validHtml(s string) bool {
	_, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return false
	}
	return true
	//d := xml.NewDecoder(strings.NewReader(s))
	//tags := make(map[string]bool, 10)
	//for {
	//	token, err := d.Token()
	//	if err != nil && err != io.EOF {
	//		return false
	//	}
	//	if token == nil {
	//		break
	//	}
	//	switch t := token.(type) {
	//	case xml.StartElement:
	//		if !in(SupportedFormattingTags, t.Name.Local) {
	//			return false
	//		}
	//		tags[t.Name.Local] = false
	//	case xml.EndElement:
	//		if _, ok := tags[t.Name.Local]; !ok || !in(SupportedFormattingTags, t.Name.Local) { // закрытый тег, не имеющий открытого, или неподдерживаемый тег
	//			return false
	//		}
	//		delete(tags, t.Name.Local)
	//	}
	//}
	//return len(tags) == 0
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

// buildLangsPagination создает пагинацию как говорил F d
// в callback передайте что-то типа set_my_lang:%s, где %s станет код выбранного языка
func buildLangsPagination(user tables.Users, offset int, count int, tickLang, buttonSelectLangCallback, buttonBackCallback, buttonNextCallback string, includeAutoDetect bool) (tgbotapi.InlineKeyboardMarkup, error) {
	if offset < 0 || offset > len(codes[user.Lang])-1 {
		return tgbotapi.InlineKeyboardMarkup{}, nil
	}
	out := tgbotapi.NewInlineKeyboardMarkup()
	if includeAutoDetect {
		if offset == 0 {
			count-- // уменьшаем кол-во кнопок, потому что мы пихаем свою
			out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(user.Localize("Detect language"), fmt.Sprintf(buttonSelectLangCallback, "auto"))))
		} else {
			offset-- // на первой странице мы недопоказали одну кнопку
			if offset+count == len(codes[user.Lang])-1 {
				count++
			}
		}
	}

	for i, code := range codes[user.Lang][offset : offset+count] {
		if offset == 0 && includeAutoDetect {
			i++
		}
		lang, ok := langs[user.Lang][code]
		//if offset+count <= 18 {
		//	lang += " 📌"
		//}
		if code == tickLang {
			lang = "✅" + lang
		}
		if code == "emj" {
			lang = "🆕" + lang
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
	if includeAutoDetect && offset > 0 {
		offset++ // для вида
		//offset = len(codes[user.Lang]) / 18 * 18 // для счетчика снизу, а то на 181 строчке мы уменьшили оффсет
	}
	out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", buttonBackCallback),
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(offset)+"/"+strconv.Itoa(len(codes[user.Lang])/18*18), "none"),
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

func closeUnclosedTags(s string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return s
	}
	raw, err := doc.Html()
	if err != nil {
		return s
	}
	i1 := strings.Index(raw, "<body>") + len("<body>")
	if i1 == -1 {
		return ""
	}
	i2 := strings.Index(raw[i1:], "</body>")
	if i2 == -1 {
		return ""
	}
	return raw[i1 : i1+i2]
}
