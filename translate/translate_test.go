package translate_test

import (
	"github.com/armanokka/translobot/translate"
	"github.com/k0kubun/pp"
	"strings"
	"testing"
)



func TestDetectLanguageYandex(t *testing.T) {
	tests := []struct{
		Input, Result string
	}{
		{".", ""},
		{"привет", "ru"},
		{"hello", "en"},
		{":)", ""},
		{"тетрагидропиранилциклопентилтетрагидропиридопиридиновые", "ru"},
		{"pneumonoultramicroscopicsilicovolcanoconiosis", "en"},
		{"❗️❗", "emj"},
		{"/v,x]'!p3(_+", "en"},
		{"++_!!", ""},
	}
	for _, test := range tests {
		out, err := translate.DetectLanguageYandex(test.Input)
		if err != nil {
			if e, ok := err.(translate.YandexDetectAPIError); ok {
				if e.Code != 502 {
					t.Error(e)
				}
			} else {
				t.Error(err)
			}
			continue
		}
		if out.Lang != test.Result {
			t.Error("waited:", test.Result, "got:", out)
		}
	}
}


func TestTranslateGoogle(t *testing.T) {

	tests := []struct{
		Input, Result string
	}{
		{"",""},
		{"!+-*(", "! + - * ("},
		{"английский", "English"},
		{"гаргулия", "gargulia"},
		{"😂😂😂", "😂😂😂"},
		{"123", "one hundred twenty-three"},
		{"🤡😈👍👌018^$@#&()??>", "🤡😈👍👌018 ^ $ @ # & () ??>"},
	}
	for _, test := range tests {
		got, err := translate.TranslateGoogle("auto", "en", test.Input)
		if err != nil {
			t.Error(err)
			continue
		}
		if got != test.Result {
			t.Error("waited:", test.Result, "got:", got)
		}
	}
}

func TestTranslateYandex(t *testing.T) {
	tests := []struct{
		FromLang, ToLang, Text, Result string
	}{
		{"en", "ru", "hello", "привет"},
		{"en", "ru", "", ""},
		{"en", "ru", "!()_+%^$", "!()_+%^$"},
		{"en", "ru", "💋😻Ё!\"№;%:💄💄", "💋😻Ё!\"№;%:💄💄"},
		{"en", "ru", "   ", "   "},
		{"en", "ru", "😂", "Tears of joy"},
	}
	for _, test := range tests {
		got, err := translate.TranslateYandex(test.FromLang, test.ToLang, test.Text)
		if err != nil {
			if e, ok := err.(translate.YandexTranslateAPIError); ok {
				t.Error(e)
			} else {
				t.Error(err)
			}
			continue
		}
		if len(got.Text) < 1 {
			if strings.Join(strings.Fields(test.Text), "") != "" { // хотя бы один символ в отправленном тексте был
				t.Error(err)
			}
		}
		pp.Println("passed", test.Text)
		
	}
}