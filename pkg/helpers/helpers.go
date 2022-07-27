package helpers

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"unicode/utf16"
)

func In(array []string, keys ...string) bool {
	for _, v := range array {
		for _, k := range keys {
			if v == k {
				return true
			}
		}
	}
	return false
}

func InMapValues(m map[string]string, values ...string) bool {
	for _, v := range values {
		var ok = false
		for _, k := range m {
			if v == k {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

// cutString cut string using runes by limit
func CutStringUTF16(text string, limit int) string {
	points := utf16.Encode([]rune(text))
	if len(points) > limit {
		return string(utf16.Decode(points[:limit]))
	}
	return text
}

func LenUTF16(text string) int {
	return len(utf16.Encode([]rune(text)))
}

func CutString(text string, limit int) string {
	runes := []rune(text)
	if len(runes) > limit {
		return string(runes[:limit])
	}
	return text
}

func ReversoType(reversoType string) string {
	switch reversoType {
	case "v.", "infl.":
		return "verb"
	case "n.", "nf.", "nn.", "nm.":
		return "noun"
	case "adv.":
		return "adverb"
	case "conj.":
		return "conjunction"
	case "pron.":
		return "pronoun"
	case "adj.":
		return "adjective"
	}
	return ""
}

func ReversoTypes(types []string) string {
	for _, t := range types {
		out := ReversoType(t)
		if out != "" {
			return out
		}
	}
	return ""
}

func ApplyEntitiesHtml(text string, entities []tgbotapi.MessageEntity) string {
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
		case "spoiler":
			startTag = "<span class=\"tg-spoiler\">"
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
		case "spoiler":
			endTag = "</span>"
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
