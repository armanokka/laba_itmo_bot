package translate

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-errors/errors"
	"github.com/k0kubun/pp"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)


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
	for i:=0;i<3;i++ {
		var err error
		res, err = request()
		if err == nil && res.StatusCode == 200 {
			break
		}
	}

	if res.StatusCode != 200 {
		return "", HTTPError{
			Code:        res.StatusCode,
			Description: "got non 200 http code",
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
		return "", errors.New("couldn't find \"'\"")
	}
	var out = struct {
		TranslateTTS []string `json:"translate_tts"`
	}{}

	err = json.Unmarshal(body[idx:], &out)
	if err != nil {
		return "", err
	}

	if len(out.TranslateTTS) == 0 {
		pp.Println("translateTTS js object not found")
		return "", errors.New("translateTTS js object not found")
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
		if length < from + chunkLength {
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

func GoogleHTMLTranslate(from, to, text string) (GoogleHTMLTranslation, error) {
	request := func() (*http.Response, error) {
		params := url.Values{}
		params.Set("q", text)
		params.Set("sl", from)
		params.Set("tl", to)
		params.Set("client", "gtx")
		params.Set("format", "html")
		req, err := http.NewRequest("POST", "https://translate.googleapis.com/translate_a/t?anno=3&client=te_lib&format=html&v=1.0&key=AIzaSyBOti4mM-6x9WDnZIjIeyEU21OpBXqWBgw&logld=vTE_20210503_00&sl=en&tl=ru&tc=3&sr=1&tk=488339.105044&mode=1", bytes.NewBufferString(params.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header["authority"] = []string{"translate.googleapis.com"}
		req.Header["cookie"] = []string{""}
		req.Header["origin"] = []string{"https://lingvanex.com"}
		req.Header["referrer"] = []string{"https://lingvanex.com"}
		req.Header["content-type"] = []string{"application/x-www-form-urlencoded; charset=UTF-8"}
		req.Header["sec-fetch-site"] = []string{"cross-site"}
		req.Header["sec-fetch-mode"] = []string{"cors"}
		req.Header["sec-fetch-dest"] = []string{"empty"}
		req.Header["sec-ch-ua-mobile"] = []string{"?0"}
		req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
		req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"}
		return http.DefaultClient.Do(req)
	}

	var res *http.Response

	for i:=0;i<3;i++ {
		var err error
		res, err = request()
		if err == nil && res.StatusCode == 200 {
			break
		}
	}


	if from == "auto" { // придёт два значения
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return GoogleHTMLTranslation{}, err
		}

		var out []string
		if err = json.Unmarshal(body, &out); err != nil {
			return GoogleHTMLTranslation{}, errors.WrapPrefix(err, string(body), 0)
		}
		if len(out) != 2 {
			return GoogleHTMLTranslation{}, errors.New("пришло не два значения от переводчика, разделитель \"|\":" + strings.Join(out, "|"))
		}
		return GoogleHTMLTranslation{
			Text: out[0],
			From: out[1],
		}, nil
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return GoogleHTMLTranslation{}, err
		}

		var out string
		if err = json.Unmarshal(body, &out); err != nil {
			return GoogleHTMLTranslation{}, errors.WrapPrefix(body, string(body), 0)
		}
		return GoogleHTMLTranslation{
			Text: out,
			From: from,
		}, nil
	}
}

func ReversoTranslate(from, to, text string) (ReversoTranslation, error) {
	if _, ok := reversoSupportedLangs[from]; !ok {
		return ReversoTranslation{}, ErrReversoLangNotSupported
	}
	if _, ok := reversoSupportedLangs[to]; !ok {
		return ReversoTranslation{}, ErrReversoLangNotSupported
	}
	if from == to {
		return ReversoTranslation{}, SameLangsWerePassed
	}


	j, err := json.Marshal(ReversoRequestTranslate{
		Input:   text,
		From:    from,
		To:      to,
		Format:  "text",
		Options: ReversoRequestTranslateOptions{
			Origin:            "reversodesktop",
			SentenceSplitter:  true,
			ContextResults:    true,
			LanguageDetection: true,
		},
	})
	if err != nil {
		return ReversoTranslation{}, err
	}
	buf := bytes.NewBuffer(j)

	request := func() (*http.Response, error) {
		req, err := http.NewRequest("POST", "https://api.reverso.net/translate/v1/translation", buf)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")
		req.Header.Add("accept-language", "en-US,en;q=0.8")
		req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")

		return http.DefaultClient.Do(req)
	}

	var res *http.Response
	for i:=0;i<3;i++ {
		res, err = request()
		if err == nil && res.StatusCode == 200 {
			break
		}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ReversoTranslation{}, err
	}

	var ret ReversoTranslation
	if err = json.Unmarshal(body, &ret); err != nil {
		return ReversoTranslation{}, errors.WrapPrefix(err, string(body), 0)
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
			"Content-Type": []string{"application/json; charset=UTF-8"},
			"Accept-Language": []string{"en-US,en;q=0.8"},
			"User-Agent": []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36"},
		}
		return http.DefaultClient.Do(req)
	}

	var res *http.Response
	for i:=0;i<3;i++ {
		res, err = request()
		//body, err := ioutil.ReadAll(res.Body)
		//pp.Println(string(body))
		if err == nil {
			break
		}
	}

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
	req, err := http.NewRequest("POST", "https://translate.googleapis.com/translate_a/single?dt=t&dt=bd&dt=qc&dt=rm&dt=ex&client=gtx&sl=" + url.PathEscape(from) + "&tl=" + url.PathEscape(to) + "&q=" + url.PathEscape(text) + "&dj=1&dt=at&ie=UTF-16&oe=UTF-16&otf=2&srcrom=1&ssel=0&tsel=0", nil)
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
	for i:=0;i<3;i++ {
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
		return GoogleTranslateSingleResult{}, errors.WrapPrefix(err, string(body), 0)
	}
	return result, err
}

func GetSamples(from, to, source, translation string) (GetSamplesResponse, error) {
	j, err := json.Marshal(getSamplesRequest{
		Direction: from + "-" + to,
		Source:    source,
		Translation: translation,
		AppID:     "26ad41b9-102f-57b8-5cb4-3dcf1dbf7cad",
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
	for i:=0;i<3;i++ {
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
		return GetSamplesResponse{}, errors.WrapPrefix(err, string(body), 0)
	}
	return result, nil
}


func GoogleDictionary(lang, text string) (GoogleDictionaryResponse, error) {
	req, err := http.NewRequest("GET", "https://content-dictionaryextension-pa.googleapis.com/v1/dictionaryExtensionData?term=" + url.PathEscape(text) + "&corpus=" + url.PathEscape(lang) + "&key=AIzaSyA6EEtrDCfBkHV8uU2lgGY-N383ZgAOo7Y", nil)
	if err != nil {
		return GoogleDictionaryResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("X-Origin", "chrome-extension://mgijmajocgfcbeboacabfgobmjgjcoja")

	var res *http.Response
	for i:=0;i<3;i++ {
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
	req, err := http.NewRequest("GET", "https://dictionary.yandex.net/dicservice.json/lookupMultiple?ui=en&srv=tr-text&text=" + url.PathEscape(text) + "&type=regular,syn,ant,deriv&lang=" + fromto + "&flags=7591&dict=" + fromto, nil)
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
	for i:=0;i<3;i++ {
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
		pp.Println(result)
		return YandexTranscriptionResponse{
			StatusCode:    code.(float64),
		}, nil
	}

	if _, ok := result[fromto]; !ok {
		return YandexTranscriptionResponse{
			StatusCode:    -2,
		}, nil
	}
	data, _ := json.Marshal(result[fromto])

	var out yandexTranscriptionResponse
	if err = json.Unmarshal(data, &out); err != nil {
		return YandexTranscriptionResponse{}, err
	}
	if len(out.Regular) == 0 {
		return YandexTranscriptionResponse{
			StatusCode:    -1,
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
	for i:=0;i<3;i++ {
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
