/*
Helper functions
*/
package bot

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/armanokka/translobot/internal/config"
	"github.com/armanokka/translobot/internal/tables"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/lingvo"
	"github.com/dlclark/regexp2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"unicode"
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

// inI is in but returning index
func inSlice(arr []string, k string) (int, bool) {
	for i, v := range arr {
		if k == v {
			return i, true
		}
	}
	return 0, false
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

// validHtml doesn't deal with order of tags, e.g. "</b>hey<b>" is valid
func validHtml(s string) bool {
	d := xml.NewDecoder(strings.NewReader(s))
	tags := make(map[string]int, 2)
	for {
		token, err := d.Token()
		if err != nil && err != io.EOF {
			return false
		}
		if token == nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			tags[t.Name.Local]++
		case xml.EndElement:
			tags[t.Name.Local]--
		}
	}
	for _, count := range tags {
		if count != 0 {
			return false
		}
	}
	return true
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
	if offset < 0 || offset > len(codes[*user.Lang])-1 {
		return tgbotapi.InlineKeyboardMarkup{}, nil
	}
	out := tgbotapi.NewInlineKeyboardMarkup()
	if includeAutoDetect {
		if offset == 0 {
			count-- // уменьшаем кол-во кнопок, потому что мы пихаем свою
			out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(user.Localize("Auto"), fmt.Sprintf(buttonSelectLangCallback, "auto"))))
		} else {
			offset-- // на первой странице мы недопоказали одну кнопку
			if offset+count == len(codes[*user.Lang])-1 {
				count++
			}
		}
	}

	for i, code := range codes[*user.Lang][offset : offset+count] {
		if offset == 0 && includeAutoDetect {
			i++
		}
		lang, ok := langs[*user.Lang][code]
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
		//offset = len(codes[*user.Lang]) / 18 * 18 // для счетчика снизу, а то на 181 строчке мы уменьшили оффсет
	}
	query := "hey"
	keyboardForMyLang := strings.HasPrefix(buttonSelectLangCallback, "set_my_lang")
	typeLanguageCallback := "type_my_lang_name"
	if !keyboardForMyLang {
		typeLanguageCallback = "type_to_lang_name"
	}
	out.InlineKeyboard = append(out.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️", buttonBackCallback),
		tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(offset)+"/"+strconv.Itoa(len(codes[*user.Lang])/18*18), "none"),
		tgbotapi.NewInlineKeyboardButtonData("➡️", buttonNextCallback)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(user.Localize(`искать языки 🔎`), typeLanguageCallback)),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:              user.Localize("inline mode"),
				SwitchInlineQuery: &query,
			}))
	return out, nil
}

func hasPrefix(s, prefix string, maxCharsDifference int) bool {
	if maxCharsDifference < 1 {
		maxCharsDifference = 0
	}
	runes := []rune(s)
	prefixRunes := []rune(prefix)
	if len(runes) > len(prefixRunes) {
		runes = runes[:len(prefixRunes)]
	}
	for i, r := range runes {
		if r != prefixRunes[i] {
			maxCharsDifference--
			if maxCharsDifference < 0 {
				return false
			}
		}
	}
	return maxCharsDifference > -1
}

func maxDiff(source string, arrs [][]string) []string {
	i, maxDifference := 0, 0
	for idx, arr := range arrs {
		difference := 0
		for _, v := range arr {
			if v == "" {
				continue
			}
			difference += diff(source, v)
		}
		if difference > maxDifference {
			i, maxDifference = idx, difference
		}
	}
	return arrs[i]
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

func parseKeyboard(text string) (keyboard interface{}, valid bool) {
	kb := tgbotapi.NewInlineKeyboardMarkup()
	if text != "/empty" {
		scanner := bufio.NewScanner(strings.NewReader(text))
		for scanner.Scan() {
			if scanner.Err() != nil {
				return
			}
			btns := strings.Fields(scanner.Text())
			row := tgbotapi.NewInlineKeyboardRow()
			for _, btn := range btns {
				parts := strings.Split(btn, "|") // parts[0] - text on button, parts[1] - link for button
				if len(parts) != 2 {
					return
				}
				row = append(row, tgbotapi.NewInlineKeyboardButtonURL(parts[0], parts[1]))
			}
			kb.InlineKeyboard = append(kb.InlineKeyboard, row)
		}
	}
	if len(kb.InlineKeyboard) > 0 {
		return kb, true
	}
	return
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

func clearGoqueryShit(s string) string {
	i1 := strings.Index(s, "<body>") + len("<body>")
	if i1 == -1 {
		return ""
	}
	i2 := strings.Index(s, "</body>")
	if i2 == -1 {
		return ""
	}
	return s[i1:i2]
}

// closeUnclosedTagsAndClearUnsupported closes all html tags even those that shouldn't be
func closeUnclosedTagsAndClearUnsupported(s string) string {
	r, err := regexp2.Compile("<[^>]*>", regexp2.RE2)
	if err != nil {
		panic(err)
	}
	tags := make([]string, 0, 4) // b /a code /p, but not b/ or /a
	m, _ := r.FindStringMatch(s)
	for m != nil {
		tag := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(m.String(), "<"), ">"))
		var i int
		for idx, ch := range tag {
			if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
				continue
			}
			i = idx
			break
		}
		if i == 0 {
			i = len(tag)
		}
		if strings.HasSuffix(tag, "/") && !strings.HasPrefix(tag, "/") {
			tag = "/" + strings.TrimSuffix(tag, "/")
		}
		if !in([]string{"b", "i", "u", "s", "span", "a", "code", "pre"}, strings.TrimPrefix(strings.TrimSuffix(tag, "/"), "/")) {
			s = strings.Replace(s, m.String(), "", 1)
		} else {
			tags = append(tags, tag)
		}
		m, _ = r.FindNextMatch(m)
	}
	if validHtml(s) {
		return s
	}
	if len(tags) > 0 {
		usedIndexes := make([]int, 0, len(tags)/2+1)
		for i := 0; i < len(tags); i++ {
			tag := tags[i]
			opening := !strings.HasPrefix(tag, "/") && !strings.HasSuffix(tag, "/")
			if opening {
				idx, ok := inSliceNotUsed(tags[i:], usedIndexes, "/"+tag)
				usedIndexes = append(usedIndexes, i)
				if !ok { // все заняты
					s += "</" + tag + ">"
					continue
				}
				usedIndexes = append(usedIndexes, idx)
				continue
			}
			idx, ok := inSliceNotUsed(tags[:i], usedIndexes, tag[1:])
			usedIndexes = append(usedIndexes, i)
			if !ok { // все заняты
				s = "<" + tag[1:] + ">" + s
				continue
			}
			usedIndexes = append(usedIndexes, idx)
		}
	}
	return s
}

func inSliceNotUsed(arr []string, usedIndexes []int, k string) (int, bool) {
	for i, v := range arr {
		if k == v {
			for _, usedIdx := range usedIndexes {
				if i != usedIdx {
					return i, true
				}
			}
		}
	}
	return 0, false
}

func remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}
func removeHtml(s string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return "", err
	}
	return doc.Text(), nil
}

func Tick(callbackData string, inlineKeyboard [][]tgbotapi.InlineKeyboardButton) {
	for i1, row := range inlineKeyboard {
		for i2, button := range row {
			if button.CallbackData == nil {
				continue
			}
			if *button.CallbackData == callbackData && !strings.HasPrefix(*button.CallbackData, "✅ ") {
				inlineKeyboard[i1][i2].Text = "✅ " + button.Text
				break
			}
		}
	}
}

func UntickAll(inlineKeyboard [][]tgbotapi.InlineKeyboardButton) {
	for i1, row := range inlineKeyboard {
		for i2, button := range row {
			inlineKeyboard[i1][i2].Text = strings.TrimPrefix(button.Text, "✅ ")
		}
	}
}
