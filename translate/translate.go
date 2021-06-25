package translate

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

type HTTPError struct {
	Code int
	Description string
}

func (c HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error [code:%v]:%s", c.Code, c.Description)
}


type TranslateGoogleAPIResponse struct {
	Text, FromLang string
}

func TranslateGoogle(from, to, text string) (*TranslateGoogleAPIResponse, error) {
	buf := new(bytes.Buffer)
	buf.WriteString("async=translate,sl:" + url.QueryEscape(from) + ",tl:" + url.QueryEscape(to) + ",st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc")
	req, err := http.NewRequest("POST", "https://www.google.com/async/translate?vet=12ahUKEwjFh8rkyaHxAhXqs4sKHYvmAqAQqDgwAHoECAIQJg..i&ei=SMbMYMXDKernrgSLzYuACg&yv=3", buf)
	if err != nil {
		return &TranslateGoogleAPIResponse{}, err
	}
	req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-8"}
	//req.Header["accept"] = []string{"*/*"}
	//req.Header["accept-encoding"] = []string{"gzip, deflate, br"}
	//req.Header["accept-language"] = []string{"ru-RU,ru;q=0.9"}
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
	return &TranslateGoogleAPIResponse{
		Text:     doc.Find("span[id=tw-answ-target-text]").Text(),
		FromLang: doc.Find("span[id=tw-answ-detected-sl]").Text(),
	}, err
}

func DetectLanguageGoogle(text string) (string, error) {
	buf := new(bytes.Buffer)
	buf.WriteString("async=translate,sl:auto,tl:en,st:" + url.QueryEscape(text) + ",id:1624032860465,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc")
	req, err := http.NewRequest("POST", "https://www.google.com/async/translate?vet=12ahUKEwjFh8rkyaHxAhXqs4sKHYvmAqAQqDgwAHoECAIQJg..i&ei=SMbMYMXDKernrgSLzYuACg&yv=3", buf)
	if err != nil {
		return "", err
	}
	req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-8"}
	//req.Header["accept"] = []string{"*/*"}
	//req.Header["accept-encoding"] = []string{"gzip, deflate, br"}
	//req.Header["accept-language"] = []string{"ru-RU,ru;q=0.9"}
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


//type TranslateLingVanexResponse struct {
//	Error string `json:"err"`
//	Result string
//}
//
//func TranslateLingVanex(from, to, text string) (*TranslateLingVanexResponse, error) {
//	params := url.Values{}
//	params.Set("from", from + "_" + strings.ToUpper(from))
//	params.Set("to", to + "_" + strings.ToUpper(to))
//	params.Set("text", text)
//	params.Set("platform", "dp")
//	buf := new(bytes.Buffer)
//	buf.WriteString(params.Encode())
//	req, err := http.NewRequest("POST", "https://api-b2b.backenster.com/b1/api/v3/translate/", buf)
//	if err != nil {
//		return nil, err
//	}
//	req.Header["content-type"] = []string{"application/json; charset=utf-8"}
//	req.Header["user-agent"] = []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"}
//	req.Header["accept"] = []string{"application/json, text/javascript, */*; q=0.01"}
//	req.Header["referer"] = []string{`https://translate.yandex.ru/?lang=ru-en&text=Как дела?`}
//	req.Header["accept-encoding"] = []string{"gzip, deflate, br"}
//	req.Header["accept-language"] = []string{"ru-RU,ru;q=0.9"}
//	req.Header["authorization"] = []string{"Bearer a_25rccaCYcBC9ARqMODx2BV2M0wNZgDCEl3jryYSgYZtF1a702PVi4sxqi2AmZWyCcw4x209VXnCYwesx"}
//	req.Header["cf-request-id"] = []string{"0ac158cd8700003aa7ad1ac000000001"}
//	req.Header["origin"] = []string{"https://lingvanex.com"}
//	req.Header["referer"] = []string{"https://lingvanex.com/"}
//	req.Header["sec-fetch-site"] = []string{"cross-site"}
//	req.Header["sec-fetch-mode"] = []string{"cors"}
//	req.Header["sec-fetch-dest"] = []string{"empty"}
//	req.Header["sec-ch-ua-mobile"] = []string{"?0"}
//	req.Header["sec-ch-ua"] = []string{`" Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`}
//
//	var client http.Client
//	res, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return nil, err
//	}
//	var response TranslateLingVanexResponse
//	err = json.Unmarshal(body, &response)
//	return &response, err
//}
