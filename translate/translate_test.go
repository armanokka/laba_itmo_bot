package translate_test

import (
	"github.com/armanokka/translobot/translate"
	"strings"
	"testing"
)



func TestDetectLanguageYandex(t *testing.T) {
	var tests = map[string]string{
		".":"",
		"hello":"en",
		"привет":"ru",
		":)":"",
		"":"",
		"тетрагидропиранилциклопентилтетрагидропиридопиридиновые":"ru",
		"pneumonoultramicroscopicsilicovolcanoconiosis":"en",
		"❗️❗":"emj",
		"/v,x]'!p3(_+":"en",
		"++_!!":"",
	}
	for input, result := range tests {
		out, err := translate.DetectLanguageYandex(input)
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
		if out.Lang != result {
			t.Error("waited:", input, "got:", out)
		}
	}
}


func TestTranslateGoogle(t *testing.T) {
	var tests = map[string]string{
		"":"",
		"!+-*(":"! + - * (",
		"английский":"English",
		"гаргулия":"gargulia",
		"😂😂😂":"😂😂😂",
		"123":"one hundred twenty-three",
		"🤡😈👍👌018^$@#&()??>":"🤡😈👍👌018 ^ $ @ # & () ??>",
	}
	for input, waited := range tests {
		got, err := translate.TranslateGoogle("auto", "en", input)
		if err != nil {
			t.Error(err)
			continue
		}
		if got != waited {
			t.Error("waited:", waited, "got:", got)
		}
	}
}

func TestTranslateYandex(t *testing.T) {
	tests := make([]map[string]string, 10, 50)
	tests = []map[string]string{
		{
			"from":"en",
			"to":"ru",
			"text":"hello",
			"result":"привет",
		},
		{
			"from":"en",
			"to":"ru",
			"text":"",
			"result":"",
		},
		{
			"from":"en",
			"to":"ru",
			"text":"!()_+%^$",
			"result":"!()_+%^$",
		},
		{
			"from":"en",
			"to":"ru",
			"text":"💋😻Ё!\"№;%:💄💄",
			"result":"💋😻E!\"№;%:💄💄",
		},
		{
			"from":"en",
			"to":"ru",
			"text":"   ",
			"result":"  ",
		},
		{
			"from":"en",
			"to":"ru",
			"text":"😂",
			"result":"Tears of joy",
		},
	}
	for _, arr := range tests {
		if _, ok := arr["from"]; !ok {
			panic(`"from" key does not exists in test`)
		}
		if _, ok := arr["to"]; !ok {
			panic(`"to" key does not exists in test`)
		}
		if _, ok := arr["text"]; !ok {
			panic(`"text" key does not exists in test`)
		}
		if _, ok := arr["result"]; !ok {
			panic(`"result" key does not exists in test`)
		}
		got, err := translate.TranslateYandex(arr["from"], arr["to"], arr["text"])
		if err != nil {
			if e, ok := err.(translate.YandexTranslateAPIError); ok {
				t.Error(e)
			} else {
				t.Error(err)
			}
			continue
		}
		if len(got.Text) < 1 {
			if strings.Join(strings.Fields(arr["text"]), "") != "" { // хотя бы один символ в отправленном тексте был
				t.Error(err)
			}
		}
	}
}