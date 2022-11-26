package translate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
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

func YandexTranslate(ctx context.Context, from, to, text string) (tr string, err error) {
	if _, ok := YandexSupportedLanguages[from]; !ok {
		return "", fmt.Errorf("YandexTranslate: unsupported \"from\": %s", from)
	}
	if _, ok := YandexSupportedLanguages[to]; !ok {
		return "", fmt.Errorf("YandexTranslate: unsupported \"to\": %s", to)
	}
	for i := 0; i < 3; i++ {
		tr, err = yandexTranslate(ctx, from, to, text)
		if err == nil {
			return
		}
	}
	return
}

func yandexTranslate(ctx context.Context, from, to, text string) (string, error) {
	parts := SplitIntoChunksBySentences(text, 400)
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	for i, part := range parts {
		i := i
		part := part
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, "GET", `https://browser.translate.yandex.net/api/v1/tr.json/translate?translateMode=balloon&srv=yabrowser&format=html&options=0&lang=`+from+"-"+to+"&id="+generateSid()+`-0-0&context_title=`+url.PathEscape(cutString(part, 16))+"&text="+url.PathEscape(part), nil)
			if err != nil {
				return err
			}
			req.Header["Content-Type"] = []string{"application/json; charset=UTF-16"}
			req.Header["Accept-Language"] = []string{"ru,en;q=0.9"}
			req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
			req.Header["Origin"] = []string{"https://stackoverflow.com"}
			req.Header["Referrer"] = []string{"https://stackoverflow.com/"}

			resp, err := request(req, 3)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				//if errors.Is(err, net.err)
				return err
			}

			if gjson.GetBytes(body, "code").Int() != 200 {
				return fmt.Errorf("YandexTranslate:" + gjson.GetBytes(body, "message").String())
			}
			mu.Lock()
			defer mu.Unlock()
			parts[i] = ""
			for _, result := range gjson.GetBytes(body, "text").Array() {
				parts[i] += result.String()
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return "", err
	}
	translation := strings.Join(parts, ".")
	if text == translation {
		return "", nil
	}
	return translation, nil
}

//func yandexTranslateRequest(ctx context.Context, from, to, text string) (string, error) {
//	params := url.Values{}
//	params.Add("text", text)
//	params.Add("translateMode", "balloon")
//	params.Add("context_title", url.PathEscape(cutString(text, 16)))
//	params.Add("id", generateSid()+`-0-0`)
//	params.Add("srv", "yabrowser")
//	params.Add("lang", from+`-`+to)
//	params.Add("format", "html")
//	params.Add("options", "0")
//
//	req, err := http.NewRequestWithContext(ctx, "GET", `https://browser.translate.yandex.net/api/v1/tr.json/translate?`+params.Encode(), nil)
//	if err != nil {
//		return "", err
//	}
//	req.Header["Content-Type"] = []string{"application/json; charset=UTF-16"}
//	req.Header["Accept-Language"] = []string{"ru,en;q=0.9"}
//	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
//	req.Header["Origin"] = []string{"https://stackoverflow.com"}
//	req.Header["Referrer"] = []string{"https://stackoverflow.com/"}
//
//	resp, err := request(req, 3)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		//if errors.Is(err, net.err)
//		return "", err
//	}
//
//	if gjson.GetBytes(body, "code").Int() != 200 {
//		return fmt.Errorf("YandexTranslate:" + gjson.GetBytes(body, "message").String())
//	}
//}

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

	resp, err := request(req, 3)
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
