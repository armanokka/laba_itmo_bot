package translate

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/armanokka/translobot/pkg/errors"
	"github.com/armanokka/translobot/pkg/lingvo"
	"github.com/go-resty/resty/v2"
	"github.com/k0kubun/pp"
	"github.com/tidwall/gjson"
	"html"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func FlexibleTranslate(from, to, text string) (string, error) {
	if len(text) < 50 {
		if v, err := lingvo.GetDictionary(from, to, text); err == nil && len(v) > 0 {
			out := ""
			for i, r := range v {
				if i > 0 {
					out += "\n"
				}
				out += r.Translations
			}
			return out, nil
		}
		if v, err := lingvo.GetDictionary(to, from, text); err == nil && len(v) > 0 {
			out := ""
			for i, r := range v {
				if i > 0 {
					out += "\n"
				}
				out += r.Translations
			}
			return out, nil
		}
	}

	if html.UnescapeString(text) != html.EscapeString(text) { // есть html теги
		tr, err := MicrosoftTranslate(from, to, text)
		if err != nil {
			return "", err
		}
		return tr.TranslatedText, nil
	}
	tr, err := YandexTranslate(from, to, text)
	if err != nil {
		return "", err
	}
	return tr, nil
}

//func GoogleTranslate(from, to, text string) (TranslateGoogleAPIResponse, error) {
//	buf := new(bytes.Buffer)
//	buf.WriteString("async=translate,sl:" + url.QueryEscape(from) + ",tl:" + url.QueryEscape(to) + ",st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc")
//	req, err := http.NewRequest("POST", "https://www.google.com/async/translate?vet=12ahUKEwjFh8rkyaHxAhXqs4sKHYvmAqAQqDgwAHoECAIQJg..i&ei=SMbMYMXDKernrgSLzYuACg&yv=3", buf)
//	if err != nil {
//		return TranslateGoogleAPIResponse{}, err
//	}
//	req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-8"}
//	// req.Header["accept"] = []string{"*/*"}
//	// req.Header["accept-encoding"] = []string{"gzip, deflate, br"}
//	// req.Header["accept-language"] = []string{"ru-RU,ru;q=0.9"}
//	req.Header["cookie"] = []string{"NID=217=mKKVUv88-BW4Vouxnh-qItLKFt7zm0Gj3yDLC8oDKb_PuLIb-p6fcPVcsXZWeNwkjDSFfypZ8BKqy27dcJH-vFliM4dKaiKdFrm7CherEXVt-u_DPr9Yecyv_tZRSDU7E52n5PWwOkaN2I0-naa85Tb9-uTjaKjO0gmdbShqba5MqKxuTLY; 1P_JAR=2021-06-18-16; DV=A3qPWv6ELckmsH4dFRGdR1fe4Gj-oRcZWqaFSPtAjwAAAAA"}
//	req.Header["origin"] = []string{"https://www.google.com"}
//	req.Header["referer"] = []string{"https://www.google.com/"}
//	req.Header["sec-fetch-site"] = []string{"cross-site"}
//	req.Header["sec-fetch-mode"] = []string{"cors"}
//	req.Header["sec-fetch-dest"] = []string{"empty"}
//	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
//	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
//	req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"}
//
//
//	var client http.Client
//	res, err := client.Do(req)
//	if err != nil {
//		return TranslateGoogleAPIResponse{}, err
//	}
//	if res.StatusCode != 200 {
//		return TranslateGoogleAPIResponse{}, HTTPError{
//			Code:        res.StatusCode,
//			Description: "got non 200 http code",
//		}
//	}
//
//	doc, err := goquery.NewDocumentFromReader(res.Body)
//	if err != nil {
//		return TranslateGoogleAPIResponse{}, err
//	}
//	result := TranslateGoogleAPIResponse{
//		Text:     doc.Find("span[id=tw-answ-target-text]").Text(),
//		FromLang: doc.Find("span[id=tw-answ-detected-sl]").Text(),
//		FromLangNativeName: doc.Find("span[id=tw-answ-detected-sl-name]").Text(),
//		SourceRomanization: doc.Find("span[id=tw-answ-source-romanization]").Text(),
//	}
//
//	doc.Find(`div[class~=tw-bilingual-entry]`).Each(func(i int, s *goquery.Selection) {
//		result.Variants = append(result.Variants, &Variant{
//			Word:    s.Find("span > span").Text(),
//			Meaning: s.Find("div").Text(),
//		})
//	})
//	doc.Find("img[data-src]").Each(func(i int, selection *goquery.Selection) {
//		link, _ := selection.Attr("data-src")
//		result.Images = append(result.Images, link)
//	})
//	return result, err
//}

func DetectLanguageGoogle(text string) (string, error) {
	d, err := detectLanguageGoogle(text)
	if err != nil {
		d, err = detectLanguageGoogle(text)
	}
	return d, err
}

func detectLanguageGoogle(text string) (string, error) {
	buf := new(bytes.Buffer)
	request := func() (*http.Response, error) {
		buf.WriteString("async=translate,sl:auto,tl:en,st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc")
		req, err := http.NewRequest("POST", "https://www.google.com/async/translate?vet=12ahUKEwjFh8rkyaHxAhXqs4sKHYvmAqAQqDgwAHoECAIQJg..i&ei=SMbMYMXDKernrgSLzYuACg&yv=3", buf)
		if err != nil {
			return nil, err
		}
		req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-8"}
		// req.Header["accept"] = []string{"*/*"}
		// req.Header["accept-encoding"] = []string{"gzip, deflate, br"}
		// req.Header["accept-language"] = []string{"ru-RU,ru;q=0.9"}
		req.Header["cookie"] = []string{"NID=217=mKKVUv88-BW4Vouxnh-qItLKFt7zm0Gj3yDLC8oDKb_PuLIb-p6fcPVcsXZWeNwkjDSFfypZ8BKqy27dcJH-vFliM4dKaiKdFrm7CherEXVt-u_DPr9Yecyv_tZRSDU7E52n5PWwOkaN2I0-naa85Tb9-uTjaKjO0gmdbShqba5MqKxuTLY; 1P_JAR=2021-06-18-16; DV=A3qPWv6ELckmsH4dFRGdR1fe4Gj-oRcZWqaFSPtAjwAAAAA"}
		req.Header["origin"] = []string{"https://www.google.com"}
		req.Header["referer"] = []string{"https://www.google.com/"}
		req.Header["sec-fetch-site"] = []string{"cross-site"}
		req.Header["sec-fetch-mode"] = []string{"cors"}
		req.Header["sec-fetch-dest"] = []string{"empty"}
		req.Header["sec-ch-ua-mobile"] = []string{"?0"}
		req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
		req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"}

		return http.DefaultClient.Do(req)
	}

	var res *http.Response
	for i := 0; i < 3; i++ {
		var err error
		res, err = request()
		if err == nil && res.StatusCode == 200 {
			break
		}
	}

	if res.StatusCode != 200 {
		return "", HTTPError{
			Code:        res.StatusCode,
			Description: "detectLanguageGoogle: non 200 http code",
		}
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	return doc.Find("span[id=tw-answ-detected-sl]").Text(), err
}

func TTS(lang, text string) ([]byte, error) {
	parts := splitIntoChunks(text, 200)

	type result struct {
		idx int
		s   string
	}
	var results = make([]result, len(parts))
	var errs = make(chan error, len(parts))

	var wg sync.WaitGroup
	for i, part := range parts {
		i := i
		part := part
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := ttsRequest(lang, part)
			if err != nil {
				errs <- err
				return
			}
			results = append(results, result{
				idx: i,
				s:   data,
			})
		}()
	}
	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return nil, <-errs
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var out string
	for _, res := range results {
		out += res.s
	}

	return base64.StdEncoding.DecodeString(out)
}

func ttsRequest(lang, text string) (string, error) {
	params := url.Values{}
	params.Add("ei", "-4vsYPWwArKWjgahzpfACw")
	params.Add("yv", "3")
	params.Add("ttsp", "tl:"+url.QueryEscape(lang)+",txt:"+url.QueryEscape(text)+",spd:1")
	params.Add("async", "_fmt:jspb")

	req, err := http.NewRequest("GET", "https://www.google.com/async/translate_tts?"+params.Encode(), nil)
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

	res, err := http.DefaultClient.Do(req)
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

func splitIntoChunks(s string, chunkLength int) []string {
	length := len(s)

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
		chunks[i] = s[from:to]
	}

	return chunks
}

type GoogleHTMLTranslation struct {
	Text string
	From string
}

func generateTkk(needNew bool) (string, error) {
	f, err := os.OpenFile("tk.cache", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err = os.Create("tk.cache")
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	defer f.Close()

	if !needNew {
		stat, err := os.Stat(f.Name())
		if err != nil {
			return "", err
		}
		if stat.ModTime().Hour() == time.Now().Hour() && time.Since(stat.ModTime()) < time.Hour {
			tkk, err := ioutil.ReadAll(f)
			if err != nil {
				return "", err
			}
			return string(tkk), nil
		}
	}

	res, err := http.DefaultClient.Get("https://translate.googleapis.com/translate_a/element.js")
	if err != nil {
		return "", err
	}
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("generateTkk: not 200 code, html: %s", out)
	}
	text := string(out)

	i1 := strings.Index(text, "c._ctkk='") + len("c._ctkk='")
	i2 := strings.Index(text[i1:], "'")
	tkk := text[i1 : i1+i2]
	if _, err = f.WriteString(tkk); err != nil {
		return "", err
	}
	return tkk, nil
}

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

func GoogleHTMLTranslate(from, to, text string) (GoogleHTMLTranslation, error) {
	text = strings.ReplaceAll(text, "\n", "<br>")
	tkk, err := generateTkk(true)
	if err != nil {
		return GoogleHTMLTranslation{}, err
	}

	tk, err := getTk(tkk, text)
	if err != nil {
		return GoogleHTMLTranslation{}, err
	}

	formData := url.Values{}
	formData.Set("q", text)
	formData.Set("client", "gtx")
	formData.Set("sl", from)
	formData.Set("tl", to)
	formData.Set("format", "html")

	uri := fmt.Sprintf("https://translate.googleapis.com/translate_a/t?anno=3&client=te_lib&format=html&v=1.0&key=AIzaSyBOti4mM-6x9WDnZIjIeyEU21OpBXqWBgw&logld=vTE_20220201&sl=" + from + "&tl=" + to + "&tc=1&sr=1&tk=" + tk + "&mode=1")
	resp, err := resty.New().R().SetHeaders(map[string]string{
		"authority":    "translate.googleapis.com",
		"origin":       "https://stackoverflow.com/",
		"referrer":     "https://stackoverflow.com/",
		"content-type": "application/x-www-form-urlencoded; charset=UTF-16",
	}).SetFormDataFromValues(formData).Post(uri)

	if err != nil {
		return GoogleHTMLTranslation{}, err
	}

	ret := resp.String()

	if arr := gjson.Get(ret, "@this").Array(); len(arr) > 0 {

		var from string
		if len(arr) == 2 {
			from = arr[1].String()
		} else {
			from, err = DetectLanguageGoogle(cutString(text, 200))
			if err != nil {
				return GoogleHTMLTranslation{}, err
			}
		}

		var out string

		for _, v := range arr {
			var s string
			if err = json.Unmarshal([]byte(v.Raw), &s); err != nil {
				return GoogleHTMLTranslation{}, err
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
			if err != nil {
				return GoogleHTMLTranslation{}, err
			}
			doc.Find("body > b").Each(func(i int, selection *goquery.Selection) {
				//if !strings.HasPrefix(text[selection.Index():], "<b>") {
				//	return
				//}
				h, err := selection.Html()
				if err != nil {
					return
				}
				h = html.UnescapeString(h)
				pp.Println(h)
				out += h
				//if i % 2 != 0 {
				//	out += h
				//}
			})
		}

		return GoogleHTMLTranslation{
			Text: html.UnescapeString(strings.ReplaceAll(strings.TrimSpace(out), "<br>", "\n")),
			From: from,
		}, nil
	}

	var tr string
	if err = json.Unmarshal([]byte(ret), &tr); err != nil {
		return GoogleHTMLTranslation{}, err
	}

	if from == "" || from == "auto" {
		from, err = DetectLanguageGoogle(cutString(text, 200))
		if err != nil {
			return GoogleHTMLTranslation{}, err
		}
	}

	return GoogleHTMLTranslation{
		Text: strings.ReplaceAll(tr, "<br>", "\n"),
		From: from,
	}, nil
}

func GoogleTranslate(from, to, text string) (*TranslateGoogleAPIResponse, error) {
	tr, err := googleTranslate(from, to, text)
	if err != nil {
		tr, err = googleTranslate(from, to, text)
	}
	return tr, err
}

func googleTranslate(from, to, text string) (*TranslateGoogleAPIResponse, error) {
	//text = cutString(text, 100)
	buf := new(bytes.Buffer)
	buf.WriteString("async=translate,sl:" + url.QueryEscape(from) + ",tl:" + url.QueryEscape(to) + ",st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc,format:html")
	req, err := http.NewRequest("POST", "https://www.google.com/async/translate?vet=12ahUKEwjFh8rkyaHxAhXqs4sKHYvmAqAQqDgwAHoECAIQJg..i&ei=SMbMYMXDKernrgSLzYuACg&yv=3", buf)
	if err != nil {
		return &TranslateGoogleAPIResponse{}, err
	}
	req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-8"}
	// req.Header["accept"] = []string{"*/*"}
	// req.Header["accept-encoding"] = []string{"gzip, deflate, br"}
	// req.Header["accept-language"] = []string{"ru-RU,ru;q=0.9"}
	req.Header["cookie"] = []string{"NID=217=mKKVUv88-BW4Vouxnh-qItLKFt7zm0Gj3yDLC8oDKb_PuLIb-p6fcPVcsXZWeNwkjDSFfypZ8BKqy27dcJH-vFliM4dKaiKdFrm7CherEXVt-u_DPr9Yecyv_tZRSDU7E52n5PWwOkaN2I0-naa85Tb9-uTjaKjO0gmdbShqba5MqKxuTLY; 1P_JAR=2021-06-18-16; DV=A3qPWv6ELckmsH4dFRGdR1fe4Gj-oRcZWqaFSPtAjwAAAAA"}
	req.Header["origin"] = []string{"https://www.google.com"}
	req.Header["referer"] = []string{"https://www.google.com/"}
	req.Header["sec-fetch-site"] = []string{"cross-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
	req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"}

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return &TranslateGoogleAPIResponse{}, err
	}
	if res.StatusCode != 200 {
		return &TranslateGoogleAPIResponse{}, HTTPError{
			Code:        res.StatusCode,
			Description: "got non 200 http code",
		}
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return &TranslateGoogleAPIResponse{}, err
	}
	result := &TranslateGoogleAPIResponse{
		Text:               doc.Find("span[id=tw-answ-target-text]").Text(),
		FromLang:           doc.Find("span[id=tw-answ-detected-sl]").Text(),
		FromLangNativeName: doc.Find("span[id=tw-answ-detected-sl-name]").Text(),
		SourceRomanization: doc.Find("span[id=tw-answ-source-romanization]").Text(),
	}

	doc.Find(`div[class~=tw-bilingual-entry]`).Each(func(i int, s *goquery.Selection) {
		result.Variants = append(result.Variants, &Variant{
			Word:    s.Find("span > span").Text(),
			Meaning: s.Find("div").Text(),
		})
	})
	doc.Find("img[data-src]").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Attr("data-src")
		result.Images = append(result.Images, link)
	})
	return result, err
}

// cutString cut string using runes by limit
func cutString(text string, limit int) string {
	runes := []rune(text)
	if len(runes) > limit {
		return string(runes[:limit])
	}
	return text
}

func ReversoTranslate(from, to, text string) (ReversoTranslation, error) {
	if _, ok := reversoSupportedLangs[from]; !ok {
		return ReversoTranslation{}, ErrLangNotSupported
	}
	if _, ok := reversoSupportedLangs[to]; !ok {
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

	req, err := http.NewRequest("POST", "https://api.reverso.net/translate/v1/translation", bytes.NewBuffer(j))
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ReversoTranslation{}, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReversoTranslation{}, err
	}
	if res.StatusCode != 200 {
		pp.Println("Not 200 CODE from Reverso", res.StatusCode, string(body))
	}

	var ret ReversoTranslation
	if err = json.Unmarshal(body, &ret); err != nil {
		return ReversoTranslation{}, errors.Wrap(err)
	}
	return ret, nil
}

// ReversoSupportedLangs returns list of supported languages by Reverso, where first string is code in  ISO639-2/T and second is code in ISO639-1
func ReversoSupportedLangs() map[string]string {
	return reversoSupportedLangs
}

func ReversoQueryService(sourceText, sourceLang, targetText, targetLang string) (ReversoQueryResponse, error) {
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
		req, err := http.NewRequest("POST", "https://context.reverso.net/bst-query-service", buf)
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
		return http.DefaultClient.Do(req)
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

func GoogleTranslateSingle(from, to, text string) (GoogleTranslateSingleResult, error) {
	var res *http.Response
	var err error
	req, err := http.NewRequest("POST", "https://translate.googleapis.com/translate_a/single?dt=t&dt=bd&dt=qc&dt=rm&dt=ex&client=gtx&sl="+url.PathEscape(from)+"&tl="+url.PathEscape(to)+"&q="+url.PathEscape(text)+"&dj=1&dt=at&ie=UTF-16&oe=UTF-16&otf=2&srcrom=1&ssel=0&tsel=0", nil)
	req.Header["host"] = []string{"translate.googleapis.com"}
	req.Header["content-type"] = []string{"application/json; charset=UTF-8"}
	req.Header["sec-fetch-site"] = []string{"cross-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
	req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"}
	req.Header["x-client-data"] = []string{"CI+2yQEIpLbJAQjBtskBCKmdygEIq9HKAQjv8ssBCJ75ywEItP/LAQjnhMwBCLWFzAEI2IXMAQjLicwB"}
	req.Header["connection"] = []string{"keep-alive"}
	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		if err == nil && res.StatusCode == 200 {
			break
		}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GoogleTranslateSingleResult{}, err
	}

	var result GoogleTranslateSingleResult
	if err = json.Unmarshal(body, &result); err != nil {
		return GoogleTranslateSingleResult{}, errors.Wrap(err)
	}
	return result, err
}

func GetSamples(from, to, source, translation string) (GetSamplesResponse, error) {
	j, err := json.Marshal(getSamplesRequest{
		Direction:   from + "-" + to,
		Source:      source,
		Translation: translation,
		AppID:       "26ad41b9-102f-57b8-5cb4-3dcf1dbf7cad",
	})
	if err != nil {
		return GetSamplesResponse{}, err
	}

	buf := bytes.NewBuffer(j)
	req, err := http.NewRequest("POST", "https://cps.reverso.net/api2/GetSamples", buf)
	if err != nil {
		return GetSamplesResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")

	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		//body, err := ioutil.ReadAll(res.Body)
		//pp.Println(string(body))
		if err == nil && res.StatusCode == 200 {
			break
		}
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GetSamplesResponse{}, err
	}

	var result GetSamplesResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return GetSamplesResponse{}, errors.Wrap(err)
	}
	return result, nil
}

func GoogleDictionary(lang, text string) (GoogleDictionaryResponse, error) {
	req, err := http.NewRequest("GET", "https://content-dictionaryextension-pa.googleapis.com/v1/dictionaryExtensionData?term="+url.PathEscape(text)+"&corpus="+url.PathEscape(lang)+"&key=AIzaSyA6EEtrDCfBkHV8uU2lgGY-N383ZgAOo7Y", nil)
	if err != nil {
		return GoogleDictionaryResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("X-Origin", "chrome-extension://mgijmajocgfcbeboacabfgobmjgjcoja")

	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		//body, err := ioutil.ReadAll(res.Body)
		//pp.Println(string(body))
		if err == nil {
			break
		}
	}
	var result GoogleDictionaryResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return GoogleDictionaryResponse{}, err
	}
	return result, err
}

func YandexTranscription(from, to, text string) (YandexTranscriptionResponse, error) {
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

	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		if err == nil {
			break
		}
	}
	var result = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
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

func ReversoSuggestions(from, to, text string) (ReversoSuggestionsResponse, error) {
	data, err := json.Marshal(reversoSuggestionRequest{
		Search:     text,
		SourceLang: from,
		TargetLang: to,
	})
	if err != nil {
		return ReversoSuggestionsResponse{}, err
	}
	req, err := http.NewRequest("POST", "https://context.reverso.net/bst-suggest-service", bytes.NewBuffer(data))
	if err != nil {
		return ReversoSuggestionsResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		if err == nil {
			break
		}
	}
	var result ReversoSuggestionsResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return ReversoSuggestionsResponse{}, err
	}
	return result, err
}

type MicrosoftTranslation struct {
	From           string
	TranslatedText string
}

func MicrosoftTranslate(from, to, text string) (MicrosoftTranslation, error) { // с помощью расширения Mate Translate
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

	req, err := http.NewRequest("GET", "https://api.microsofttranslator.com/v2/ajax.svc/TranslateArray?"+params.Encode(), nil)
	if err != nil {
		return MicrosoftTranslation{}, err
	}
	req.Header["Content-Type"] = []string{"application/json; charset=UTF-8"}
	req.Header["Accept-Language"] = []string{"ru-RU,ru;q=0.9"}
	req.Header["Accept"] = []string{"application/json, text/javascript, */*; q=0.01"}
	req.Header["User-aAent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.174 YaBrowser/22.1.5.810 Yowser/2.5 Safari/537.36"}
	req.Header["origin"] = []string{"https://stackoverflow.com"}
	req.Header["referrer"] = []string{"https://stackoverflow.com/"}
	req.Header["sec-fetch-site"] = []string{"cross-site"}
	req.Header["sec-fetch-mode"] = []string{"cors"}
	req.Header["sec-fetch-dest"] = []string{"empty"}
	req.Header["sec-ch-ua-mobile"] = []string{"?1"}
	req.Header["sec-ch-ua-platform"] = []string{`"Android"`}
	req.Header["sec-ch-ua"] = []string{`" Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"`}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return MicrosoftTranslation{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MicrosoftTranslation{}, err
	}

	i := bytes.Index(body, []byte("["))
	if i == -1 {
		time.Sleep(time.Second)
		return MicrosoftTranslate(from, to, text)
		//return MicrosoftTranslation{}, fmt.Errorf("%s", body)
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
		out += elem.Get("TranslatedText").String()
	}

	return MicrosoftTranslation{
		From:           from,
		TranslatedText: out,
	}, nil
}
