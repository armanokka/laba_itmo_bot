package translate

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func generateSid() string {
	i := strconv.FormatInt(time.Now().UnixMilli(), 16)
	for t, e := 0, 16 - len(i); t < e; t++ {
		rand.Seed(time.Now().UnixNano() + int64(t))
		i += strconv.FormatInt(int64(math.Floor(16 * rand.Float64())), 16)
	}
	return i
}

func YandexTranslate(from, to, text string) (string, error) {
	if _, ok := YandexSupportedLanguages[from]; !ok {
		return "", ErrLangNotSupported
	}
	if _, ok := YandexSupportedLanguages[to]; !ok {
		return "", ErrLangNotSupported
	}
	parts := splitIntoChunksBySentences(text, 2000)
	out := ""
	for _, part := range parts {
		params := url.Values{}
		params.Add("text", part)
		params.Add("translateMode", "balloon")
		params.Add("context_title",  url.PathEscape(cutString(part, 16)))
		params.Add("id", generateSid() + `-0-0`)
		params.Add("srv", "yabrowser")
		params.Add("lang", from+`-`+to)
		params.Add("format", "html")
		params.Add("options", "0")
		req, err := http.NewRequest("GET", `https://browser.translate.yandex.net/api/v1/tr.json/translate?`+ params.Encode(), nil)
		if err != nil {
			return "", err
		}
		req.Header["Content-Type"] = []string{"application/json; charset=UTF-16"}
		req.Header["Accept-Language"] = []string{"ru,en;q=0.9"}
		req.Header["User-aAent"] =  []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
		req.Header["origin"] = []string{"https://stackoverflow.com"}
		req.Header["referrer"] = []string{"https://stackoverflow.com/"}
		req.Header["sec-fetch-site"] = []string{"cross-site"}
		req.Header["sec-fetch-mode"] = []string{"cors"}
		req.Header["sec-fetch-dest"] = []string{"empty"}
		req.Header["sec-ch-ua-mobile"] = []string{"?0"}
		req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if gjson.GetBytes(body, "code").Int() != 200 {
			return "", fmt.Errorf("YandexTranslate:" + gjson.GetBytes(body, "message").String())
		}
		for _, result := range gjson.GetBytes(body, "text").Array() {
			out += result.String()
		}
	}
	return out, nil
}
