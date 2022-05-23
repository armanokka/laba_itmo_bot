package helpers

import (
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
