package translate

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/armanokka/translobot/pkg/helpers"
	"github.com/go-resty/resty/v2"
	"github.com/k0kubun/pp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/unicode/norm"
	"html"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

var proxies = []string{
	"http://armangokka7796:Ob3mMsp26e0hUVop1YE4N0PVvxo1nFyG@172.96.166.123:10455",
	"http://armangokka7796:77813e@162.19.173.183:10799",
}

func TTS(ctx context.Context, lang, text string) ([]byte, error) {
	parts := SplitIntoChunks(text, 200)

	type result struct {
		idx int
		s   string
	}
	var results = make([]result, len(parts))

	g, _ := errgroup.WithContext(context.Background())
	for i, part := range parts {
		i := i
		part := part
		g.Go(func() error {
			data, err := ttsRequest(ctx, lang, part)
			if err != nil {
				return err
			}
			results = append(results, result{
				idx: i,
				s:   data,
			})
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var out = new(bytes.Buffer)
	for _, res := range results {
		if _, err := out.WriteString(res.s); err != nil {
			return nil, err
		}
	}

	return base64.StdEncoding.DecodeString(out.String())
}

var client *http.Client

//go :embed certs/*.pem
//var embedFS embed.FS

// request requests with proxy on 3 error responses
func request(req *http.Request, times int) (resp *http.Response, err error) {
	if client == nil {
		//certs := x509.NewCertPool()
		//if err = fs.WalkDir(embedFS, "certs", func(path string, d fs.DirEntry, err error) error {
		//	if d.IsDir() {
		//		return nil
		//	}
		//	pemData, err := os.ReadFile(path)
		//	if err != nil {
		//		return err
		//	}
		//	if ok := certs.AppendCertsFromPEM(pemData); !ok {
		//		return fmt.Errorf("request: can't parse %s", path)
		//	}
		//	return nil
		//}); err != nil {
		//	return nil, err
		//}

		//client = &http.Client{
		//Transport: &http.Transport{
		//	TLSClientConfig: &tls.Config{
		//RootCAs: certs,
		//	},
		//},
		//}
		client = http.DefaultClient
	}
	log, ok := req.Context().Value("log").(*zap.Logger)
	if !ok {
		log, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		log.Error("request: context.Value(\"log\") is not *zap.Logger")
	}

	for i := 0; i < times; i++ {
		ctx, _ := context.WithTimeout(req.Context(), time.Second*150)
		req := req.WithContext(ctx)
		resp, err = client.Do(req)
		if err != nil {
			if e, ok := err.(*url.Error); ok {
				switch e.Err.(type) {
				case x509.CertificateInvalidError, x509.HostnameError, x509.SystemRootsError, x509.InsecureAlgorithmError, x509.UnknownAuthorityError, x509.ConstraintViolationError:
					log.Error("hacking detected. this incident will be reported")
					os.Exit(1)
					return nil, nil
				}
			}

			requestDump, _ := httputil.DumpRequest(req, true)
			log.Error("", zap.Error(err), zap.String("request", string(requestDump)))
		}
		if resp != nil && err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			//fmt.Println("request repeated", i+1, "times")
			return
		}
	}
	for i := 0; i < len(proxies); i++ {
		u, err := url.Parse(proxies[i])
		if err != nil {
			return nil, fmt.Errorf("proxies[%d] - invalid proxy url", i)
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(u),
		}
		resp, err = client.Do(req)
		if err != nil {
			requestDump, _ := httputil.DumpRequest(req, true)
			log.Error("", zap.Error(err), zap.String("request", string(requestDump)))
		}
		if resp != nil && err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return resp, err
		}
	}
	return nil, fmt.Errorf("translate.request(): couldn't send request by ourself and with proxy")
}

func requestLocally(req *http.Request, times int) (resp *http.Response, err error) {
	for i := 0; i < times; i++ {
		resp, err = http.DefaultClient.Do(req)
		if resp != nil && err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return
		}
	}
	return
}

func requestProxy(req *http.Request) (resp *http.Response, err error) {
	client := http.DefaultClient
	for i := 0; i < len(proxies); i++ {
		var u *url.URL
		u, err = url.Parse(proxies[i])
		if err != nil {
			return
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(u),
		}
		resp, err = client.Do(req)
		if resp != nil && err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return
		}
	}
	return
}

// request requests with proxy on 3 error responses
//func request(req *http.Request, times int) (*http.Response, error) {
//	// Отправляем одновременно два запроса: один через нас, второй через прокси
//	// Второй сразу отправляем, если размер запроса меньше 5кб. Иначе отправляем через прокси, только если через нас не смогли
//	// Цель - вернуть любой успешный запрос. Неважно откуда
//	ctx, cancel := context.WithCancel(req.Context())
//	g, ctx := errgroup.WithContext(ctx)
//	callProxy := make(chan bool)
//	var out *http.Response
//	g.Go(func() error {
//		resp, err := requestLocally(req, times)
//		if err != nil && !errors.Is(err, context.Canceled) {
//			return err
//		} else if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
//			return err
//		}
//		cancel()
//		pp.Println(resp.StatusCode)
//		out = resp // TODO fix this
//		// TODO control cache from server side. Cloudflare mustn't cache error responses
//		return nil
//	})
//	g.Go(func() error {
//		if req.ContentLength > 1024*5 {
//			select {
//			case <-ctx.Done():
//				switch ctx.Err() {
//				case context.Canceled, context.DeadlineExceeded:
//					return nil
//				default:
//					return ctx.Err()
//				}
//			case <-callProxy:
//			}
//		}
//		resp, err := requestProxy(req)
//		if err != nil && !errors.Is(err, context.Canceled) {
//			return err
//		} else if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
//			return err
//		}
//		cancel()
//		pp.Println(resp.StatusCode)
//		out = resp
//		return nil
//	})
//
//	err := g.Wait()
//	return out, err
//}

func ttsRequest(ctx context.Context, lang, text string) (string, error) {
	params := url.Values{}
	params.Add("ei", "-4vsYPWwArKWjgahzpfACw")
	params.Add("yv", "3")
	params.Add("ttsp", "tl:"+url.QueryEscape(lang)+",txt:"+url.QueryEscape(text)+",spd:1")
	params.Add("async", "_fmt:jspb")

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com/async/translate_tts?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}
	req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-8"}
	req.Header["cookie"] = []string{"NID=217=mKKVUv88-BW4Vouxnh-qItLKFt7zm0Gj3yDLC8oDKb_PuLIb-p6fcPVcsXZWeNwkjDSFfypZ8BKqy27dcJH-vFliM4dKaiKdFrm7CherEXVt-u_DPr9Yecyv_tZRSDU7E52n5PWwOkaN2I0-naa85Tb9-uTjaKjO0gmdbShqba5MqKxuTLY; 1P_JAR=2021-06-18-16; DV=A3qPWv6ELckmsH4dFRGdR1fe4Gj-oRcZWqaFSPtAjwAAAAA"}
	req.Header["origin"] = []string{"https://www.google.com"}
	req.Header["referer"] = []string{"https://www.google.com/"}
	req.Header["sec-fetch-site"] = []string{"cross-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
	req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"}

	res, err := request(req, 3)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", HTTPError{
			Code:        res.StatusCode,
			Description: "Non 200 HTTP code",
		}
	}
	idx := strings.IndexByte(string(body), '\'') + 1
	if idx == 0 {
		return "", fmt.Errorf("couldn't find \"'\"")
	}
	var out = struct {
		TranslateTTS []string `json:"translate_tts"`
	}{}

	err = json.Unmarshal(body[idx:], &out)
	if err != nil {
		return "", err
	}

	if len(out.TranslateTTS) == 0 {
		return "", fmt.Errorf("translateTTS js object not found")
	}
	if out.TranslateTTS[0] == "" {
		return "", ErrTTSLanguageNotSupported
	}
	return out.TranslateTTS[0], err
}

func SplitIntoChunks(s string, chunkLength int) []string {
	runes := []rune(s)
	length := len(runes)

	chunksCount := length / chunkLength
	if length%chunkLength != 0 {
		chunksCount += 1
	}

	chunks := make([]string, chunksCount)

	for i := range chunks {
		from := i * chunkLength
		var to int
		if length < from+chunkLength {
			to = length
		} else {
			to = from + chunkLength
		}
		chunks[i] = string(runes[from:to])
	}

	return chunks
}

type GoogleHTMLTranslation struct {
	Text string
	From string
}

//func generateTkk(needNew bool) (string, error) {
//	f, err := os.OpenFile("tk.cache", os.O_RDWR|os.O_CREATE, 0644)
//	if err != nil {
//		if errors.Is(err, os.ErrNotExist) {
//			f, err = os.Create("tk.cache")
//			if err != nil {
//				return "", err
//			}
//		} else {
//			return "", err
//		}
//	}
//	defer f.Close()
//
//	if !needNew {
//		stat, err := os.Stat(f.Name())
//		if err != nil {
//			return "", err
//		}
//		if stat.ModTime().Hour() == time.Now().Hour() && time.Since(stat.ModTime()) < time.Hour {
//			tkk, err := ioutil.ReadAll(f)
//			if err != nil {
//				return "", err
//			}
//			return string(tkk), nil
//		}
//	}
//
//	res, err := http.def.Get("https://translate.googleapis.com/translate_a/element.js")
//	if err != nil {
//		return "", err
//	}
//	out, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return "", err
//	}
//	if res.StatusCode != 200 {
//		return "", fmt.Errorf("generateTkk: not 200 code, html: %s", out)
//	}
//	text := string(out)
//
//	i1 := strings.Index(text, "c._ctkk='") + len("c._ctkk='")
//	i2 := strings.Index(text[i1:], "'")
//	tkk := text[i1 : i1+i2]
//	if _, err = f.WriteString(tkk); err != nil {
//		return "", err
//	}
//	return tkk, nil
//}

func getTk(tkk string, text string) (string, error) {
	resp, err := resty.New().R().SetFormData(map[string]string{
		"tkk":  tkk,
		"text": text,
	}).Post(fmt.Sprintf("http://ancient-springs-54230.herokuapp.com/tkk.php"))
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

//func GoogleHTMLTranslate(ctx context.Context, from, to, text string) (GoogleHTMLTranslation, error) {
//	text = strings.ReplaceAll(text, "\n", "<br>")
//	tkk, err := generateTkk(true)
//	if err != nil {
//		return GoogleHTMLTranslation{}, err
//	}
//
//	tk, err := getTk(tkk, text)
//	if err != nil {
//		return GoogleHTMLTranslation{}, err
//	}
//
//	formData := url.Values{}
//	formData.Set("q", text)
//	formData.Set("client", "gtx")
//	formData.Set("sl", from)
//	formData.Set("tl", to)
//	formData.Set("format", "html")
//
//	uri := fmt.Sprintf("https://translate.googleapis.com/translate_a/t?anno=3&client=te_lib&format=html&v=1.0&key=AIzaSyBOti4mM-6x9WDnZIjIeyEU21OpBXqWBgw&logld=vTE_20220201&sl=" + from + "&tl=" + to + "&tc=1&sr=1&tk=" + tk + "&mode=1")
//	resp, err := resty.New().R().SetHeaders(map[string]string{
//		"authority":    "translate.googleapis.com",
//		"origin":       "https://stackoverflow.com/",
//		"referrer":     "https://stackoverflow.com/",
//		"content-type": "application/x-www-form-urlencoded; charset=UTF-16",
//	}).SetFormDataFromValues(formData).Post(uri)
//
//	if err != nil {
//		return GoogleHTMLTranslation{}, err
//	}
//
//	ret := resp.String()
//
//	if arr := gjson.Get(ret, "@this").Array(); len(arr) > 0 {
//
//		var from string
//		if len(arr) == 2 {
//			from = arr[1].String()
//		} else {
//			from, err = DetectLanguageGoogle(ctx, cutString(text, 200))
//			if err != nil {
//				return GoogleHTMLTranslation{}, err
//			}
//		}
//
//		var out string
//
//		for _, v := range arr {
//			var s string
//			if err = json.Unmarshal([]byte(v.Raw), &s); err != nil {
//				return GoogleHTMLTranslation{}, err
//			}
//			doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
//			if err != nil {
//				return GoogleHTMLTranslation{}, err
//			}
//			doc.Find("body > b").Each(func(i int, selection *goquery.Selection) {
//				//if !strings.HasPrefix(text[selection.Index():], "<b>") {
//				//	return
//				//}
//				h, err := selection.Html()
//				if err != nil {
//					return
//				}
//				h = html.UnescapeString(h)
//				out += h
//				//if i % 2 != 0 {
//				//	out += h
//				//}
//			})
//		}
//
//		return GoogleHTMLTranslation{
//			Text: html.UnescapeString(strings.ReplaceAll(strings.TrimSpace(out), "<br>", "\n")),
//			From: from,
//		}, nil
//	}
//
//	var tr string
//	if err = json.Unmarshal([]byte(ret), &tr); err != nil {
//		return GoogleHTMLTranslation{}, err
//	}
//
//	if from == "" || from == "auto" {
//		from, err = DetectLanguageGoogle(ctx, cutString(text, 200))
//		if err != nil {
//			return GoogleHTMLTranslation{}, err
//		}
//	}
//
//	return GoogleHTMLTranslation{
//		Text: strings.ReplaceAll(tr, "<br>", "\n"),
//		From: from,
//	}, nil
//}

var englishAlphabet = map[rune]bool{
	69:  true,
	84:  true,
	85:  true,
	99:  true,
	103: true,
	109: true,
	67:  true,
	87:  true,
	105: true,
	112: true,
	117: true,
	81:  true,
	86:  true,
	97:  true,
	68:  true,
	75:  true,
	80:  true,
	66:  true,
	76:  true,
	79:  true,
	82:  true,
	98:  true,
	106: true,
	111: true,
	122: true,
	89:  true,
	118: true,
	119: true,
	71:  true,
	72:  true,
	102: true,
	107: true,
	110: true,
	115: true,
	88:  true,
	100: true,
	113: true,
	114: true,
	77:  true,
	65:  true,
	73:  true,
	74:  true,
	83:  true,
	101: true,
	104: true,
	108: true,
	120: true,
	90:  true,
	116: true,
	121: true,
	70:  true,
	78:  true,
}

func validHTMLTagNotNotranslate(tag string) bool {
	tag = strings.Trim(tag, "<>/ ")
	spaceIdx := strings.IndexAny(tag, "\r\n\t\f\v ")
	if spaceIdx != -1 {
		tag = tag[:spaceIdx]
	}
	tag = strings.TrimSpace(tag)
	if tag == "notranslate" {
		return false
	}
	for _, r := range tag {
		if _, ok := englishAlphabet[r]; !ok && !unicode.IsDigit(r) && r != '_' && r != '-' && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// plainTag: <a href=""> --> a, <b/> --> b, <code> --> code
func plainTag(s string) string {
	tag := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, "<"), ">"))
	var i int
	for idx, ch := range tag {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			continue
		}
		i = idx
		break
	}
	if i == 0 {
		i = len(tag)
	}
	return strings.TrimPrefix(strings.TrimSuffix(tag[:i], "/"), "/")
}

func nextNotranslate(s string) (start int, end int) {
	start, end, skip := -1, -1, 0
	for {
		start, end = IndexesOfFirstHTMLTag(s[skip:])
		if start == -1 || end == -1 {
			return -1, -1
		}
		//pp.Println(plainTag(s[start:end]), start, end, skip, s[skip:])
		if plainTag(s[start+skip:end+skip]) == "notranslate" {
			return start + skip, end + skip
		}
		skip += end
	}
}

func IndexesOfFirstHTMLTag(s string) (start int, end int) {
	start = strings.Index(s, "<")
	if start == -1 {
		return -1, -1
	}

	end = closingBrace(s[start:])
	if end == -1 {
		return -1, -1
	}
	end++
	return start, start + end
}

func closingBrace(s string) int {
	var quote uint16
	for i, ch := range utf16.Encode([]rune(s)) {
		if ch == '\'' || ch == '"' {
			if quote == ch {
				quote = uint16(0)
			} else if quote == uint16(0) {
				quote = ch
			}
		}
		if ch == '>' && quote == uint16(0) {
			return i
		}
	}
	return -1
}

func insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func replace(s, old, new string, count, skip int) string {
	var startIdx int

	for i := 0; i < skip; i++ {
		idx := strings.Index(s[startIdx:], old)
		if idx == -1 {
			return s
		}
		startIdx += idx + len(old)
	}

	return s[:startIdx] + strings.Replace(s[startIdx:], old, new, count)
}

// plainTag: <a href=""> --> a, <b/> --> b, <code> --> code
func PlainTag(s string) string {
	tag := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(s, "<"), ">"))
	var i int
	for idx, ch := range tag {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '/' {
			continue
		}
		i = idx
		break
	}
	if i == 0 {
		i = len(tag)
	}
	return strings.TrimPrefix(strings.TrimSuffix(tag[:i], "/"), "/")
}

func NextNotranslate(s string) (start int, end int) {
	start, end, skip := -1, -1, 0
	for {
		start, end = IndexesOfFirstHTMLTag(s[skip:])
		if start == -1 || end == -1 {
			return -1, -1
		}
		//pp.Println(plainTag(s[start:end]), start, end, skip, s[skip:])
		if PlainTag(s[start+skip:end+skip]) == "notranslate" {
			return start + skip, end + skip
		}
		skip += end
	}
}

// NotranslateRanges: <notranslate>hey</notranslate> -> [][2]int{{0, 30}}
func NotranslateRanges(s string) [][2]int {
	out := make([][2]int, 0, 2)
	skip := 0
	for {
		openingStart, openingEnd := NextNotranslate(s)
		if openingStart == -1 || openingEnd == -1 {
			return out
		}
		_, endingEnd := NextNotranslate(s[openingEnd:])
		if endingEnd == -1 {
			return out
		}
		out = append(out, [2]int{openingStart + skip, endingEnd + openingEnd + skip})
		skip += endingEnd + openingEnd
		s = s[endingEnd+openingEnd:]
	}
}

func GoogleTranslate(ctx context.Context, from, to, text string) (out TranslateGoogleAPIResponse, err error) {
	// use log in ctx
	if !helpers.In(SupportedLanguageCodes, from, to) {
		return TranslateGoogleAPIResponse{}, fmt.Errorf("GoogleTranslate: unsupported language code %s %s", from, to)
	}
	if strings.TrimSpace(text) == "" {
		return TranslateGoogleAPIResponse{
			Text:     text,
			FromLang: from,
		}, nil
	}
	// Получаем тексты между тегами
	texts := make([]string, 0, 3)
	start, end := -1, -1
	s := text
	for {
		start, end = IndexesOfFirstHTMLTag(s)
		if start == -1 || end == -1 {
			texts = append(texts, s)
			break
		}
		texts = append(texts, s[:start])
		if PlainTag(s[start:end]) == "notranslate" {
			start1, end1 := NextNotranslate(s[end:])
			if start1 != -1 && end1 != -1 {
				s = s[end+end1:] // идем в конец
				continue
			}
		}
		s = s[end:]
	}
	for i, text := range texts { // TODO: реплейсить тексты в строке
		chunks := SplitIntoChunksBySentences(text, 400) // TODO чанки багаются при тексте
		if len(chunks) > 1 {
			texts[i] = chunks[0]
			texts = append(texts, chunks[1:]...)
		}
	}

	// Отправляем в переводчики
	notranslateRanges := NotranslateRanges(text)
	g, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	for _, chunk := range texts {
		chunk := chunk
		if chunk == "" {
			continue
		}
		g.Go(func() error {
			// Нужно сохранить пробелы в началах чанков
			// заменяем пробелы
			prefixSpaces, postfixSpaces := "", ""
			if to != "ar" && !strings.HasPrefix(to, "zh") && to != "th" && to != "ja" && to != "ko" {
				prefixSpaces, postfixSpaces = helpers.GetPrefixSpaces(chunk), helpers.GetPostfixSpaces(chunk)
			}
			out, err = googleTranslate(ctx, from, to, chunk)
			if err != nil {
				return err
			}
			if prefixSpaces != "" {
				out.Text = prefixSpaces + out.Text
			}
			if postfixSpaces != "" {
				out.Text += postfixSpaces
			}

			if len(notranslateRanges) > 0 {
				i, skip := strings.Index(text, chunk), 0 // idx - это индекс, с которого надо искать следующий strings.Index. skip - это сколько совпадений мы пропустили
				for ; i != -1; i, skip = strings.Index(text[i:], chunk), skip+1 {
					for _, block := range notranslateRanges {
						notranslateStart, notranslateEnd := block[0], block[1]
						if i >= notranslateStart && i < notranslateEnd {
							continue
						}
						mu.Lock()
						defer mu.Unlock()
						if out.FromLang == "" {
							out.FromLang = out.FromLang
						}
						text = replace(text, chunk, out.Text, 1, skip)
						return nil
					}
					if i == 0 {
						i++
					}
				}
			}
			mu.Lock()
			defer mu.Unlock()
			if out.FromLang == "" {
				out.FromLang = out.FromLang
			}
			text = replace(text, chunk, out.Text, 1, 0)
			return nil
		})
	}
	err = g.Wait()
	out.Text = html.UnescapeString(text)
	return out, err
}

func googleTranslate(ctx context.Context, from, to, text string) (result TranslateGoogleAPIResponse, err error) {
	text = strings.ReplaceAll(text, "\n", "<br>")
	result, err = googleTranslateRequest(ctx, from, to, text)
	result.Text = strings.NewReplacer(` \ n`, "\n", `\ n`, "\n", "<br>", "\n").Replace(result.Text)
	return result, err
}

func googleTranslateRequest(ctx context.Context, from, to, text string) (result TranslateGoogleAPIResponse, err error) {
	buf := new(bytes.Buffer)
	buf.WriteString("async=translate,sl:" + url.QueryEscape(from) + ",tl:" + url.QueryEscape(to) + ",st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc,format:html")
	req, err := http.NewRequestWithContext(ctx, "POST", "https://www.google.com/async/translate?vet=12ahUKEwjFh8rkyaHxAhXqs4sKHYvmAqAQqDgwAHoECAIQJg..i&ei=SMbMYMXDKernrgSLzYuACg&yv=3&cs=0&rlz=1C1GCEA_enUZ1012UZ1012", buf)
	if err != nil {
		return TranslateGoogleAPIResponse{}, err
	}
	req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-16"}
	req.Header["accept"] = []string{"*/*"}
	req.Header["accept-ch"] = []string{"Sec-CH-UA-Platform", "Sec-CH-UA-Platform", "Sec-CH-UA-Full-Version", "Sec-CH-UA-Arch", "Sec-CH-UA-Model", "Sec-CH-UA-Bitness", "Sec-CH-UA-Full-Version-List", "Sec-CH-UA-WoW64"}
	req.Header["cookie"] = []string{"NID=217=mKKVUv88-BW4Vouxnh-qItLKFt7zm0Gj3yDLC8oDKb_PuLIb-p6fcPVcsXZWeNwkjDSFfypZ8BKqy27dcJH-vFliM4dKaiKdFrm7CherEXVt-u_DPr9Yecyv_tZRSDU7E52n5PWwOkaN2I0-naa85Tb9-uTjaKjO0gmdbShqba5MqKxuTLY; 1P_JAR=2021-06-18-16; DV=A3qPWv6ELckmsH4dFRGdR1fe4Gj-oRcZWqaFSPtAjwAAAAA"}
	req.Header["origin"] = []string{"https://www.google.com"}
	req.Header["referer"] = []string{"https://www.google.com/"}
	req.Header["sec-fetch-site"] = []string{"cross-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
	req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"}

	resp, err := request(req, 3)
	if err != nil {
		return TranslateGoogleAPIResponse{}, err
	}
	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return TranslateGoogleAPIResponse{}, err
		}
		//fmt.Println(string(body))
		doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
		if err != nil {
			return TranslateGoogleAPIResponse{}, err
		}
		if result.FromLang == "" {
			result.FromLang = doc.Find("span[id=tw-answ-detected-sl]").Text()
		}
		if result.FromLang == "" {
			return TranslateGoogleAPIResponse{}, fmt.Errorf("googleTranslateRequest: didn't find \"tw-answ-detected-sl\" in:\n%s", string(body))
		}
		result.Text = doc.Find("span[id=tw-answ-target-text]").Text()
		if result.Text == "" {
			if doc.Find("span[id=tw-answ-id]").Text() != "" {
				return TranslateGoogleAPIResponse{
					Text:     "",
					FromLang: to,
				}, nil
			}
			return TranslateGoogleAPIResponse{}, fmt.Errorf("googleTranslateRequest: didn't find \"tw-answ-target-text\" in:\n%s", string(body))
		}
		doc.Find("div.tw-bilingual-entry").Each(func(_ int, selection *goquery.Selection) {
			translations := strings.Fields(selection.Find("div").Text())
			for i, item := range translations {
				translations[i] = strings.Trim(strings.TrimSpace(item), ",")
			}
			if result.ReverseTranslations == nil {
				result.ReverseTranslations = make(map[string][]string, 1)
			}
			result.ReverseTranslations[selection.Find("span").Children().Text()] = translations
		})
		verified, err := strconv.ParseBool(doc.Find("span.tw-answ-community-verified").Text())
		if err == nil {
			result.CommunityVerified = verified
		}

	case 413:
		parts := SplitIntoChunksBySentences(text, utf8.RuneCountInString(text)/2)
		g, ctx := errgroup.WithContext(ctx)
		var mu sync.Mutex
		from := ""
		for i, part := range parts {
			part := part
			g.Go(func() error {
				tr, err := googleTranslateRequest(ctx, from, to, part)
				if err != nil {
					return err
				}
				mu.Lock()
				defer mu.Unlock()
				if from == "" {
					from = tr.FromLang
				}
				parts[i] = norm.NFKC.String(tr.Text)
				return nil
			})
		}
		if err = g.Wait(); err != nil {
			return TranslateGoogleAPIResponse{}, err
		}
		result.Text = strings.Join(parts, "")
		result.FromLang = from
	default:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return TranslateGoogleAPIResponse{}, err
		}
		return TranslateGoogleAPIResponse{}, HTTPError{
			Code:        resp.StatusCode,
			Description: fmt.Sprintf("googleTranslateRequest: got not 200 http code [%d]: %s", resp.StatusCode, string(body)),
		}
	}
	r, err := regexp.Compile("\\<\\s*[bB][rR]\\s*\\>")
	if err != nil {
		return TranslateGoogleAPIResponse{}, err
	}
	result.Text = r.ReplaceAllString(result.Text, "\n")

	//result := &TranslateGoogleAPIResponse{
	//	Text:     doc.Find("span[id=tw-answ-target-text]").Text(),
	//	FromLang: doc.Find("span[id=tw-answ-detected-sl]").Text(),
	//FromLangNativeName: doc.Find("span[id=tw-answ-detected-sl-name]").Text(),
	//SourceRomanization: doc.Find("span[id=tw-answ-source-romanization]").Text(),
	//}

	//doc.Find(`div[class~=tw-bilingual-entry]`).Each(func(i int, s *goquery.Selection) {
	//	result.Variants = append(result.Variants, &Variant{
	//		Word:    s.Find("span > span").Text(),
	//		Meaning: s.Find("div").Text(),
	//	})
	//})
	//doc.Find("img[data-src]").Each(func(i int, selection *goquery.Selection) {
	//	link, _ := selection.Attr("data-src")
	//	result.Images = append(result.Images, link)
	//})
	return result, nil
}

// cutString cut string using runes by limit
func cutString(text string, limit int) string {
	runes := []rune(text)
	if len(runes) > limit {
		return string(runes[:limit])
	}
	return text
}

func ReversoTranslate(ctx context.Context, from, to, text string) (ReversoTranslation, error) {
	if _, ok := ReversoSupportedLangs[from]; !ok {
		return ReversoTranslation{}, ErrLangNotSupported
	}
	if _, ok := ReversoSupportedLangs[to]; !ok {
		return ReversoTranslation{}, ErrLangNotSupported
	}
	if from == to {
		return ReversoTranslation{}, SameLangsWerePassed
	}

	j, err := json.Marshal(ReversoRequestTranslate{
		Input:  text,
		From:   from,
		To:     to,
		Format: "text",
		Options: ReversoRequestTranslateOptions{
			Origin:            "translation.web",
			SentenceSplitter:  true,
			ContextResults:    true,
			LanguageDetection: true,
		},
	})
	if err != nil {
		return ReversoTranslation{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.reverso.net/translate/v1/translation", bytes.NewBuffer(j))
	if err != nil {
		return ReversoTranslation{}, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("accept-language", "en-US,en;q=0.8")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")
	req.Header.Add("X-Reverso-Origin", "translation.web")
	req.Header.Add("Referrer", "https://www.reverso.net/")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Length", strconv.Itoa(len(text)))

	res, err := request(req, 3)
	if err != nil {
		return ReversoTranslation{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReversoTranslation{}, err
	}
	if res.StatusCode != 200 {
		//if res.StatusCode == 403 {
		//	return ReversoTranslation{}, nil
		//}
		pp.Println("error", string(body))
		return ReversoTranslation{}, fmt.Errorf("Not 200 CODE from Reverso [%d]\n%s->%s,text:%s\nresponse:%s", res.StatusCode, from, to, text, string(body))
	}

	var ret ReversoTranslation
	if err = json.Unmarshal(body, &ret); err != nil {
		return ReversoTranslation{}, fmt.Errorf("ReversoTranslate: unmarshal error %s\nresponse:%s", err.Error(), string(body))
	}
	return ret, nil
}

func ReversoQueryService(ctx context.Context, sourceText, sourceLang, targetText, targetLang string) (ReversoQueryResponse, error) {
	j, err := json.Marshal(ReversoQueryRequest{
		SourceText: sourceText,
		TargetText: targetText,
		SourceLang: sourceLang,
		TargetLang: targetLang,
		Npage:      2,
		Mode:       0,
	})
	if err != nil {
		return ReversoQueryResponse{}, err
	}

	buf := bytes.NewBuffer(j)

	request := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", "https://context.reverso.net/bst-query-service", buf)
		if err != nil {
			return nil, err
		}
		//req.Header.Add("Content-Type", "application/json; charset=UTF-8")
		//req.Header.Add("Accept-Language", "en-US,en;q=0.8")
		//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")
		//
		req.Header = http.Header{
			"Content-Type":    []string{"application/json; charset=UTF-8"},
			"Accept-Language": []string{"en-US,en;q=0.8"},
			"User-Agent":      []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36"},
		}
		return request(req, 3)
	}

	res, err := request()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReversoQueryResponse{}, err
	}

	var ret ReversoQueryResponse

	if err = json.Unmarshal(body, &ret); err != nil {
		return ReversoQueryResponse{}, err
	}
	return ret, nil
}

func GoogleDictionary(ctx context.Context, lang, text string) (dict []string, phonetics string, err error) {
	for i := 0; i < 3; i++ {
		dict, phonetics, err = googleDictionary(ctx, lang, text)
		if err == nil {
			break
		}
	}
	return
}

func googleDictionary(ctx context.Context, lang, text string) (dict []string, phonetics string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://content-dictionaryextension-pa.googleapis.com/v1/dictionaryExtensionData?term="+url.PathEscape(text)+"&corpus="+url.PathEscape(lang)+"&key=AIzaSyA6EEtrDCfBkHV8uU2lgGY-N383ZgAOo7Y", nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("X-Origin", "chrome-extension://mgijmajocgfcbeboacabfgobmjgjcoja")

	resp, err := request(req, 3)
	if err != nil {
		return nil, "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if (resp.StatusCode != 200 && resp.StatusCode != 404) || (gjson.GetBytes(body, "status").Int() != 200 && gjson.GetBytes(body, "status").Int() != 404) {
		return nil, "", fmt.Errorf("googleDictionary: not 200 code [%d]\nlang:%s\ntext:%s\nresponse:%s", resp.StatusCode, lang, text, string(body))
	}
	// debug
	//f, err := os.Create("response.json")
	//if err != nil {
	//	return nil, err
	//}
	//defer f.Close()
	//f.Write(body)
	dict = make([]string, 0, 20)
	for _, v := range gjson.GetBytes(body, "dictionaryData").Array() {
		for _, entry := range v.Get("entries").Array() {
			for _, phonetic := range entry.Get("phonetics").Array() {
				if p := phonetic.Get("text").String(); p != "" {
					phonetics = p
					break
				}
			}
			if s := entry.Get("etymology.etymology.text").String(); s != "" {
				dict = append(dict, s)
			}
			if s := entry.Get("note.text").String(); s != "" {
				dict = append(dict, s)
			}
			for _, senseFamily := range entry.Get("senseFamilies").Array() {

				for _, sense := range senseFamily.Get("senses").Array() {
					dict = append(dict, sense.Get("conciseDefinition").String())
					//for _, subSence := range sense.Get("subsenses").Array() {
					//	result = append(result, subSence.Get("conciseDefinition").String())
					//}
				}
			}
		}
	}

	return dict, phonetics, err
}

func YandexTranscription(ctx context.Context, from, to, text string) (YandexTranscriptionResponse, error) {
	fromto := url.PathEscape(from) + "-" + url.PathEscape(to)
	req, err := http.NewRequest("GET", "https://dictionary.yandex.net/dicservice.json/lookupMultiple?ui=en&srv=tr-text&text="+url.PathEscape(text)+"&type=regular,syn,ant,deriv&lang="+fromto+"&flags=7591&dict="+fromto, nil)
	if err != nil {
		return YandexTranscriptionResponse{}, err
	}
	req.Header["authority"] = []string{"mc.yandex.ru"}
	req.Header["cookie"] = []string{""}
	req.Header["origin"] = []string{"https://translate.yandex.ru"}
	req.Header["referrer"] = []string{"https://translate.yandex.ru/?lang=ru-en"}
	req.Header["content-type"] = []string{"application/x-www-form-urlencoded; charset=UTF-8"}
	req.Header["sec-fetch-site"] = []string{"cross-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
	req.Header["accept"] = []string{"*/*"}
	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
	req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"}

	resp, err := request(req, 3)
	if err != nil {
		return YandexTranscriptionResponse{}, err
	}
	var result = make(map[string]interface{})
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return YandexTranscriptionResponse{}, err
	}

	if code, ok := result["code"]; ok {
		return YandexTranscriptionResponse{
			StatusCode: code.(float64),
		}, nil
	}

	if _, ok := result[fromto]; !ok {
		return YandexTranscriptionResponse{
			StatusCode: -2,
		}, nil
	}
	data, _ := json.Marshal(result[fromto])

	var out yandexTranscriptionResponse
	if err = json.Unmarshal(data, &out); err != nil {
		return YandexTranscriptionResponse{}, err
	}
	if len(out.Regular) == 0 {
		return YandexTranscriptionResponse{
			StatusCode: -1,
		}, nil
	}

	return YandexTranscriptionResponse{
		StatusCode:    200,
		Transcription: out.Regular[0].Ts,
		Pos:           out.Regular[0].Pos.Tooltip,
	}, nil
}

func ReversoSuggestions(ctx context.Context, from, to, text string) (ReversoSuggestionsResponse, error) {
	data, err := json.Marshal(reversoSuggestionRequest{
		Search:     text,
		SourceLang: from,
		TargetLang: to,
	})
	if err != nil {
		return ReversoSuggestionsResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://context.reverso.net/bst-suggest-service", bytes.NewBuffer(data))
	if err != nil {
		return ReversoSuggestionsResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	resp, err := request(req, 3)
	if err != nil {
		return ReversoSuggestionsResponse{}, err
	}
	var result ReversoSuggestionsResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ReversoSuggestionsResponse{}, err
	}
	return result, err
}

type MicrosoftTranslation struct {
	From           string
	TranslatedText string
}

func MicrosoftTranslate(ctx context.Context, from, to, text string) (MicrosoftTranslation, error) { // с помощью расширения Mate Translate
	if helpers.In(MicrosoftUnsupportedLanguages, from, to) {
		return MicrosoftTranslation{}, ErrLangNotSupported
	}
	ctx, _ = context.WithTimeout(ctx, time.Second*150)
	params := url.Values{}
	params.Set("appId", `"000000000A9F426B41914349A3EC94D7073FF941"`)
	texts, err := json.Marshal(strings.Split(text, "."))
	if err != nil {
		return MicrosoftTranslation{}, err
	}
	params.Set("texts", string(texts))
	params.Set("to", `"`+to+`"`)
	params.Set("loc", "en")
	params.Set("ctr", "")
	params.Set("ref", "WidgetV3")
	rand.Seed(time.Now().UnixNano())
	params.Set("rgp", strconv.FormatInt(int64(math.Floor(1e9*rand.Float64())), 16))

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.microsofttranslator.com/v2/ajax.svc/TranslateArray?"+params.Encode(), nil)
	if err != nil {
		return MicrosoftTranslation{}, err
	}
	req.Header["Content-Type"] = []string{"application/json; charset=UTF-8"}
	req.Header["Accept-Language"] = []string{"ru-RU,ru;q=0.9"}
	req.Header["Accept"] = []string{"application/json, text/javascript, */*; q=0.01"}
	req.Header["User-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
	req.Header["origin"] = []string{"https://stackoverflow.com"}
	req.Header["referrer"] = []string{"https://stackoverflow.com/"}
	req.Header["sec-fetch-site"] = []string{"cross-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?1"}
	req.Header["sec-ch-ua-platform"] = []string{`"Android"`}
	req.Header["sec-ch-ua"] = []string{`" Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"`}

	resp, err := request(req, 3)
	if err != nil {
		return MicrosoftTranslation{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MicrosoftTranslation{}, err
	}

	i := bytes.Index(body, []byte("["))
	if i == -1 {
		response := string(body)
		//if strings.HasPrefix(strings.TrimSpace(response), `"ArgumentOutOfRangeException: 'from' must be a valid language`) {
		//	return MicrosoftTranslation{}, nil
		//}
		//fmt.Println(string(body))
		return MicrosoftTranslation{}, fmt.Errorf("MicrosoftTranslate [%d]: bytes.Index(body, \"[\" not found:%s->%s\ntext:%s\nresponse:%s", resp.StatusCode, from, to, text, response)
	}
	body = body[i:]

	out := ""
	from = ""
	for i, elem := range gjson.ParseBytes(body).Array() {
		if from == "" {
			from = elem.Get("From").String()
		}
		if i > 0 {
			out += "."
		}
		out += norm.NFKC.String(elem.Get("TranslatedText").String())
	}

	return MicrosoftTranslation{
		From:           from,
		TranslatedText: out,
	}, nil
}

func ReversoParaphrase(ctx context.Context, lang, text string) (resp []string, err error) {
	if l := len(strings.Fields(text)); l < 3 || l > 30 {
		return nil, err
	}
	for i := 0; i < 3; i++ {
		resp, err = reversoParaphrase(ctx, lang, text)
		if err == nil {
			break
		}
	}
	return resp, err
}

func reversoParaphrase(ctx context.Context, lang, text string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://rephraser-api.reverso.net/v1/rephrase?language="+url.PathEscape(lang)+"&sentence="+url.PathEscape(text)+"&candidates=6", nil)
	if err != nil {
		return nil, err
	}
	req.Header["Content-Type"] = []string{"application/json; charset=UTF-16"}
	req.Header["Accept"] = []string{"*/*"}
	req.Header["Accept-Language"] = []string{"ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7"}
	//req.Header["access-control-request-headers"] = []string{"x-reverso-origin"}
	//req.Header["access-control-request-method"] = []string{"GET"}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
	req.Header["origin"] = []string{"https://www.reverso.net"}
	req.Header["referrer"] = []string{"https://www.reverso.net/"}
	req.Header["sec-fetch-site"] = []string{"same-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
	req.Header["x-reverso-origin"] = []string{`translation.web`}

	resp, err := request(req, 3)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode >= 400 && resp.StatusCode <= 500 && resp.StatusCode != 403 { // lang is not supported or other shit
			return []string{}, nil
		}
		if gjson.GetBytes(body, "message").String() == "No message available" {
			return nil, nil
		}
		return nil, fmt.Errorf("reversoParaphrase: not 200 http [%d], lang %s, text %s\nresponse:%s", resp.StatusCode, lang, text, string(body))
	}
	arr := gjson.GetBytes(body, "candidates").Array()
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		out = append(out, v.Get("candidate").String())
	}
	return out, nil
}
