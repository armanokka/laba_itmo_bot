package translate_test

import (
	"github.com/armanokka/translobot/translate"
	"testing"
)


func TestTranslateGoogle(t *testing.T) {

	tests := []struct{
		Input, Result string
	}{
		{"",""},
		{"!+-*(", "!+-*( en"},
		{"английский", "English"},
		{"гаргулия", "gargulia"},
		{"😂😂😂", "😂😂😂"},
		{"123", "123 en"},
		{"🤡😈👍👌018^$@#&()??>", "🤡😈👍👌018 ^ $ @ # & () ??>"},
	}
	for _, test := range tests {
		got, err := translate.GoogleHTMLTranslate("auto", "en", test.Input)
		if err != nil {
			t.Error(err)
			continue
		}
		if got.Text != test.Result {
			t.Error("waited:", test.Result, "got:", got)
		}
	}
}

func TestReversoTranslate(t *testing.T) {
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
		_, err := translate.ReversoTranslate("eng", "rus", test.Input)
		if err != nil {
			t.Error(err)
		}
	}
}

// надо еще затестить DetectLanguageGoogle, хотя он идентичен TranslateGoogle, просто ищет другой атрибут
