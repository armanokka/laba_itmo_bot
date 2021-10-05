package translate_test

import (
	"github.com/armanokka/translobot/translate"
	"github.com/k0kubun/pp"
	"testing"
)


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

func TestReversoQueryService(t *testing.T) {
	pp.Println(translate.ReversoQueryService("beautiful", "en", "красивый", "ru"))
}

// надо еще затестить DetectLanguageGoogle, хотя он идентичен TranslateGoogle, просто ищет другой атрибут
