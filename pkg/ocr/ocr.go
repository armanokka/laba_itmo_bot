package ocr

import (
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"os"
	"strings"
)

type YandexOcr struct {
	DetectedLang string
	Text         string
}

func Yandex(filename string) (YandexOcr, error) {
	f, err := os.Open(filename)
	if err != nil {
		return YandexOcr{}, err
	}
	defer f.Close()
	res, err := resty.DefaultClient.R().SetFileReader("file", filename, f).SetHeaders(map[string]string{
		"Content-type":    "multipart/form-data",
		"Accept":          "*/*",
		"Accept-Language": "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
		"Origin":          "https://translate.yandex.ru",
		"Referrer":        "https://translate.yandex.ru/ocr",
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
	}).SetContentLength(true).Post("https://translate.yandex.net/ocr/v1.1/recognize?srv=tr-image&sid=33580fba.6264eb54.2a1326e4.74722d696d616765&lang=*&yu=3182696151650778339&yum=1650748639252550033")
	if res.StatusCode() != 200 {
		return YandexOcr{}, fmt.Errorf("yandexOcrResty: not 200\nresponse:%s", res.String())
	}
	out := make([]string, 0, 15)
	for _, block := range gjson.GetBytes(res.Body(), "data.blocks").Array() {
		for _, box := range block.Get("boxes").Array() {
			out = append(out, box.Get("text").String())
		}
	}
	return YandexOcr{
		DetectedLang: gjson.GetBytes(res.Body(), "data.detected_lang").String(),
		Text:         strings.Join(out, " "),
	}, nil
}
