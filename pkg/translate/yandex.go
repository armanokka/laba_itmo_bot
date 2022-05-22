package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func generateSid() string {
	i := strconv.FormatInt(time.Now().UnixMilli(), 16)
	for t, e := 0, 16-len(i); t < e; t++ {
		rand.Seed(time.Now().UnixNano() + int64(t))
		i += strconv.FormatInt(int64(math.Floor(16*rand.Float64())), 16)
	}
	return i
}

func YandexTranslate(ctx context.Context, from, to, text string) (string, error) {
	tr, err := yandexTranslate(ctx, from, to, text)
	if err != nil {
		tr, err = yandexTranslate(ctx, from, to, text)
	}
	return tr, err
}

func yandexTranslate(ctx context.Context, from, to, text string) (string, error) {
	if _, ok := YandexSupportedLanguages[from]; !ok {
		return "", ErrLangNotSupported
	}
	if _, ok := YandexSupportedLanguages[to]; !ok {
		return "", ErrLangNotSupported
	}
	parts := SplitIntoChunksBySentences(text, 1800)
	out := new(bytes.Buffer)
	for _, part := range parts {
		params := url.Values{}
		for _, chunk := range SplitIntoChunksBySentences(part, 400) {
			params.Add("text", chunk)
		}

		params.Add("translateMode", "balloon")
		params.Add("context_title", url.PathEscape(cutString(part, 16)))
		params.Add("id", generateSid()+`-0-0`)
		params.Add("srv", "yabrowser")
		params.Add("lang", from+`-`+to)
		params.Add("format", "html")
		params.Add("options", "0")
		req, err := http.NewRequestWithContext(ctx, "GET", `https://browser.translate.yandex.net/api/v1/tr.json/translate?`+params.Encode(), nil)
		if err != nil {
			return "", err
		}
		req.Header["Content-Type"] = []string{"application/json; charset=UTF-16"}
		req.Header["Accept-Language"] = []string{"ru,en;q=0.9"}
		req.Header["User-aAent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
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
			//if errors.Is(err, net.err)
			return "", err
		}

		if gjson.GetBytes(body, "code").Int() != 200 {
			return "", fmt.Errorf("YandexTranslate:" + gjson.GetBytes(body, "message").String())
		}
		for _, result := range gjson.GetBytes(body, "text").Array() {
			out.WriteString(result.String())
		}
	}
	return out.String(), nil
}
func DetectLanguageYandex(ctx context.Context, text string) (lang string, err error) {
	for i := 0; i < 3; i++ {
		lang, err = detectLanguageYandex(ctx, text)
		if err == nil {
			break
		}
	}
	return lang, err
}

func detectLanguageYandex(ctx context.Context, text string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", `https://translate.yandex.net/api/v1/tr.json/detect?sid=`+generateSid()+`-0-0&srv=tr-text&text=`+url.PathEscape(text)+`&options=1&yu=9527670361648278346&yum=1648283397624970429`, nil)
	if err != nil {
		return "", err
	}
	req.Header["Content-Type"] = []string{"application/json; charset=UTF-16"}
	req.Header["Accept-Language"] = []string{"ru,en;q=0.9"}
	req.Header["User-aAent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
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
	out := struct {
		Code    int    `json:"code"`
		Lang    string `json:"lang"`
		Message string `json:"message"`
	}{}
	if err = json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("json.Unmarshal: %s\nresponse:%s", err, string(body))
	}

	if out.Code != 200 {
		return "", fmt.Errorf("detectLanguageYandex: %s\nBody:%s", out.Message, string(body))
	}
	return out.Lang, nil
}

func YandexQueryCorpus(ctx context.Context, from, to, text string) ([]string, error) {
	pair := from + "-" + to
	escaped := url.PathEscape(text)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://dictionary.yandex.net/dicservice.json/queryCorpus?sid="+generateSid()+"&ui=ru&srv=tr-text&text="+escaped+"&type&lang="+pair+"&flags=1063&src="+escaped+"&dst&options=0&chunks=1&maxlen=200&v=2&yu=9456730081652787081&yum=1648283397624970429", nil)
	if err != nil {
		return nil, err
	}
	req.Header["Accept"] = []string{"*/*"}
	req.Header["Accept-Language"] = []string{"ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7"}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"}
	req.Header["Origin"] = []string{"https://translate.yandex.ru"}
	req.Header["Referrer"] = []string{"https://translate.yandex.ru/?lang=" + pair}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("YandexDictionary error [HTTP %d]:%s->%s, text: %s, response:%s", res.StatusCode, from, to, text, string(body))
	}
	results := gjson.GetBytes(body, "result.examples").Array()
	if len(results) == 0 {
		return nil, nil
		//return Dictionary{}, fmt.Errorf("word %s not found\n%s->%s\nresponse:%s", text, from, to, string(body))
	}
	examples := make([]string, 0, len(results))
	for _, result := range results {
		e := result.Get("dst").String()
		e = strings.TrimSpace(strings.Trim(e, "–-"))
		e = strings.NewReplacer("<", "<b>", ">", "</b>").Replace(e)
		examples = append(examples, e)
	}
	sort.Slice(examples, func(i, j int) bool {
		return len(examples[i]) < len(examples[j])
	})
	if len(examples) > 3 {
		examples = examples[:3]
	}
	return examples, nil
}
