package translate

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)


func GoogleTranslate(from, to, text string) (*TranslateGoogleAPIResponse, error) {
	buf := new(bytes.Buffer)
	buf.WriteString("async=translate,sl:" + url.QueryEscape(from) + ",tl:" + url.QueryEscape(to) + ",st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc")
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
		Text:     doc.Find("span[id=tw-answ-target-text]").Text(),
		FromLang: doc.Find("span[id=tw-answ-detected-sl]").Text(),
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

func DetectLanguageGoogle(text string) (string, error) {
	buf := new(bytes.Buffer)
	buf.WriteString("async=translate,sl:auto,tl:en,st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc")
	req, err := http.NewRequest("POST", "https://www.google.com/async/translate?vet=12ahUKEwjFh8rkyaHxAhXqs4sKHYvmAqAQqDgwAHoECAIQJg..i&ei=SMbMYMXDKernrgSLzYuACg&yv=3", buf)
	if err != nil {
		return "", err
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
		return "", err
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

// MustTranslate translate your text into language you set in Player struct
func (p *Player) MustTranslate(text string) string {
	l, err := DetectLanguageGoogle(text)
	if err != nil {
		return text
	}
	tr, err := GoogleTranslate(l, p.Lang, text)
	if err != nil {
		return text
	}
	return tr.Text
}


func TTS(lang, text string) ([]byte, error) {

	tts := func(text string) ([]byte, error) {
		params := url.Values{}
		params.Add("ei", "-4vsYPWwArKWjgahzpfACw")
		params.Add("yv", "3")
		params.Add("ttsp", "tl:" + url.QueryEscape(lang) + ",txt:" + url.QueryEscape(text) + ",spd:1")
		params.Add("async", "_fmt:jspb")

		req, err := http.NewRequest("GET", "https://www.google.com/async/translate_tts?" + params.Encode(), nil)
		if err != nil {
			return nil, err
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
			return nil, err
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			return nil, TTSError{
				Code:        res.StatusCode,
				Description: "Non 200 HTTP code",
			}
		}
		idx := strings.IndexByte(string(body), '\n') + 1
		var out = struct {
			TranslateTTS []string `json:"translate_tts"`
		}{}
		err = json.Unmarshal(body[idx:], &out)
		if err != nil {
			return nil, err
		}
		if len(out.TranslateTTS) == 0 {
			return nil, errors.New("translateTTS js object not found")
		}
		sDec, err := base64.StdEncoding.DecodeString(out.TranslateTTS[0])
		return sDec, err
	}

	var inputs = splitIntoChunks(text, 200)

	var out = make([]byte, 0)
	for _, input := range inputs {
		b, err := tts(input)
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
	}
	return out, nil
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
	params := url.Values{}
	params.Set("q", text)
	params.Set("sl", from)
	params.Set("tl", to)
	params.Set("client", "gtx")
	params.Set("format", "html")
	req, err := http.NewRequest("POST", "https://translate.googleapis.com/translate_a/t?anno=3&client=te_lib&format=html&v=1.0&key=AIzaSyBOti4mM-6x9WDnZIjIeyEU21OpBXqWBgw&logld=vTE_20210503_00&sl=en&tl=ru&tc=3&sr=1&tk=488339.105044&mode=1", bytes.NewBufferString(params.Encode()))
	if err != nil {
		return GoogleHTMLTranslation{}, err
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
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GoogleHTMLTranslation{}, err
	}
	if from == "auto" { // придёт два значения
		var out []string
		if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
			return GoogleHTMLTranslation{}, err
		}
		if len(out) != 2 {
			return GoogleHTMLTranslation{}, errors.New("пришло не два значения от переводчика, разделитель \"|\":" + strings.Join(out, "|"))
		}
		return GoogleHTMLTranslation{
			Text: out[0],
			From: out[1],
		}, nil
	}
	var out string
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return GoogleHTMLTranslation{}, err
	}
	return GoogleHTMLTranslation{
		Text: out,
		From: from,
	}, nil
}