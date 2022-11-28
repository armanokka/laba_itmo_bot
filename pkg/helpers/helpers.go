package helpers

import (
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/unicode/norm"
	"html"
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

func inRune(r rune, variants []rune) bool {
	for _, v := range variants {
		if r == v {
			return true
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

var htmlEscape = map[uint16][]uint16{
	utf16.Encode([]rune(">"))[0]: utf16.Encode([]rune("&gt;")),
	utf16.Encode([]rune("<"))[0]: utf16.Encode([]rune("&lt;")),
	utf16.Encode([]rune("&"))[0]: utf16.Encode([]rune("&amp;")),
	utf16.Encode([]rune("’"))[0]: utf16.Encode([]rune("&#39;")),
	utf16.Encode([]rune(`"`))[0]: utf16.Encode([]rune("&quot;")),
}

// SplitIntoChunksBySentences. You can merge output by ""
func SplitIntoChunksBySentences(text string, limit int) []string {
	if len(text) < limit {
		return []string{text}
	}
	chunks := make([]string, 0, len(text)/limit+1)
	points := utf16.Encode([]rune(text))
	for i := 0; i < len(points); {
		offset := indexDelim(points[i:], limit, ".!?;\r\n\t\f\v*)")
		ch := string(utf16.Decode(points[i : i+offset]))
		chunks = append(chunks, ch)
		i += offset
	}
	return chunks
}

func indexDelim(text []uint16, limit int, delims string) (offset int) {
	if len(text) < limit {
		return len(text)
	} else if len(text) > limit {
		text = text[:limit]
	}
	offset = len(text)
	delimeters := utf16.Encode([]rune(delims))
	for i := len(text) - 1; i >= 0; i-- {
		if in(delimeters, text[i]) {
			return i + 1
		}
	}
	return offset
}
func in(arr []uint16, keys ...uint16) bool {
	for _, v := range arr {
		for _, k := range keys {
			if k == v {
				return true
			}
		}
	}
	return false
}

// TODO объединить GetPrefixSpaces и GetPostfixSpaces в одну GetEdgeSpaces
func GetPrefixSpaces(s string) (prefix string) {
	for _, ch := range s {
		if ch != '\n' && ch != '\r' && ch != '\t' && ch != '\f' && ch != '\v' && ch != ' ' {
			break
		}
		prefix += string(ch)
	}
	return
}

func GetPostfixSpaces(s string) (postfix string) {
	for i := len(s) - 1; i >= 0; i-- {
		ch := s[i]
		if ch != '\n' && ch != '\r' && ch != '\t' && ch != '\f' && ch != '\v' && ch != ' ' {
			break
		}
		postfix += string(ch)
	}
	return
}

// ApplyEntitiesHtml adds <notranslate></notranslate> to some types of entities
func ApplyEntitiesHtml(text string, entities []tgbotapi.MessageEntity, messageLengthLimit int) []string {
	chunks := SplitIntoChunksBySentences(text, messageLengthLimit)
	if len(entities) == 0 {
		for i, chunk := range chunks {
			chunks[i] = html.EscapeString(chunk)
		}
		return chunks
	}

	var chunkOffset int
	for i := 0; i < len(chunks); i += 1 {
		chunk := utf16.Encode([]rune(chunks[i]))
		pointers := make(map[int]string)
		for _, entity := range entities {
			entityStart := entity.Offset
			entityEnd := entityStart + entity.Length
			if entityEnd < chunkOffset || entityStart > chunkOffset+len(chunk) {
				continue
			}
			var before, after string
			switch entity.Type {
			case "code", "pre":
				before, after = `<notranslate><code>`, `</code></notranslate>`
			case "bold":
				before, after = `<b>`, `</b>`
			case "italic":
				before, after = `<i>`, `</i>`
			case "underline":
				before, after = `<u>`, `</u>`
			case "strikethrough":
				before, after = `<s>`, `</s>`
			case "text_link":
				before, after = `<notranslate><a href="`+entity.URL+`">`, `</a></notranslate>`
			case "text_mention":
				before, after = `<notranslate><a href="tg://user?id=`+strconv.FormatInt(entity.User.ID, 10)+`">`, `</a></notranslate>`
			case "spoiler":
				before, after = "<span class=\"tg-spoiler\">", "</span>"
			case "mention", "hashtag", "cashtag", "bot_command", "url", "email", "phone_number", "custom_emoji":
				before, after = "<notranslate>", "</notranslate>"
			}

			//pointers[entity.Offset] += before
			//pointers[entity.Offset+entity.Length] = after + pointers[entity.Offset+entity.Length]

			// All scenarios:
			// 1. tag starts before chunk and ends before/in/after chunk
			// 2. tag starts in chunk and ends in/after chunk
			// 3. tag starts after chunk and ends after chunk
			// Scenarios we have to handle:
			// 1. tag ends after chunk [*]
			// 2. tag starts before chunk and ends after chunk

			if entityStart < chunkOffset {
				pointers[0] += before
			} else {
				pointers[entityStart-chunkOffset] += before
			}
			if entityEnd > chunkOffset+len(chunk) {
				pointers[len(chunk)] += after
			} else {
				pointers[entityEnd-chunkOffset] += after
			}
		}

		var out = make([]uint16, 0, len(chunk))
		for i, ch := range chunk {
			if m, ok := pointers[i]; ok {
				out = append(out, utf16.Encode([]rune(m))...)
			}
			if escaped, ok := htmlEscape[ch]; ok {
				out = append(out, escaped...)
			} else {
				out = append(out, ch)
			}
		}
		if m, ok := pointers[len(chunk)]; ok {
			out = append(out, utf16.Encode([]rune(m))...)
		}
		chunks[i] = norm.NFKC.String(strings.ReplaceAll(string(utf16.Decode(out)), "<br>", "\n"))
		chunkOffset += len(chunk)
	}
	return chunks
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
