package lingvo

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"net/url"
)



func TutorCards(from, to, text string) (*[]TutorCard, error) {
	srcLang, ok := Lingvo[from]
	if !ok {
		return nil, fmt.Errorf("LingvoTutorCards: no such code %s", from)
	}
	dstLang, ok := Lingvo[to]
	if !ok {
		return nil, fmt.Errorf("LingvoTutorCards: no such code %s", from)
	}

	dst := fmt.Sprintf("https://api.lingvolive.com/Translation/tutor-cards?text=%s&srcLang=%d&dstLang=%d", url.PathEscape(text), srcLang, dstLang)

	res, err := resty.R().SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		"Origin": "https://www.lingvolive.com",
		"Referrer": "https://www.lingvolive.com/",
		"Host": "api.lingvolive.com",
		"Accept": "application/json, text/plain, */*",
		"Accept-Language": "ru-ru",
		//"LL-GA-ClientId": "1865379537.1639590360",
		"Connection": "keep-alive",
	}).Get(dst)
	if err != nil {
		return nil, err
	}

	var ret = make([]TutorCard, 0, 5)
	if err := json.Unmarshal(res.Body(), &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetDictionary(from, to, text string) ([]Dictionary, error) {
	srcLang, ok := Lingvo[from]
	if !ok {
		return nil, fmt.Errorf("LingvoDictionary: no such code %s", from)
	}
	dstLang, ok := Lingvo[to]
	if !ok {
		return nil, fmt.Errorf("LingvoDictionary: no such code %s", to)
	}

	dst := fmt.Sprintf("https://api.lingvolive.com/Translation/tutor-cards?text=%s&srcLang=%d&dstLang=%d", url.PathEscape(text), srcLang, dstLang)

	res, err := resty.R().SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		"Origin": "https://www.lingvolive.com",
		"Referrer": "https://www.lingvolive.com/",
		"Host": "api.lingvolive.com",
		"Accept": "application/json, text/plain, */*",
		"Accept-Language": "ru-ru",
		//"LL-GA-ClientId": "1865379537.1639590360",
		"Connection": "keep-alive",
	}).Get(dst)
	if err != nil {
		return nil, err
	}

	var out = make([]Dictionary, 0, 1)
	if err := json.Unmarshal(res.Body(), &out); err != nil {
		return nil, err
	}
	return out, err
}

type SuggestionResult struct {
	Items []SuggestionsItem `json:"items"`
	StartPos string `json:"startPos"`
	HasNextPage bool `json:"hasNextPage"`
	Prefix string `json:"prefix"`
	TargetLanguageID int `json:"targetLanguageId"`
	SourceLanguageID int `json:"sourceLanguageId"`
}
type SuggestionsItem struct {
	Heading string `json:"heading"`
	LingvoTranslations string `json:"lingvoTranslations"`
	LingvoSoundFileName string `json:"lingvoSoundFileName"`
	SocialTranslations interface{} `json:"socialTranslations"`
	LingvoDictionaryName string `json:"lingvoDictionaryName"`
	Type string `json:"type"`
	SrcLangID int `json:"srcLangId"`
	DstLangID int `json:"dstLangId"`
	Source string `json:"source"`
}

func Suggestions(from, to, text string, count, offset int) (*SuggestionResult, error) {
	srcLang, ok := Lingvo[from]
	if !ok {
		return nil, fmt.Errorf("lingvo.Suggestions: no such code %s", from)
	}
	dstLang, ok := Lingvo[to]
	if !ok {
		return nil, fmt.Errorf("lingvo.Suggestions: no such code %s", to)
	}

	dst := fmt.Sprintf("https://api.lingvolive.com/Translation/WordListPart?prefix=%s&srcLang=%d&dstLang=%d&pageSize=%d&startIndex=%d", url.PathEscape(text), srcLang, dstLang, count, offset)
	res, err := resty.R().SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		"Origin": "https://www.lingvolive.com",
		"Referrer": "https://www.lingvolive.com/",
		"Host": "api.lingvolive.com",
		"Accept": "application/json, text/plain, */*",
		"Accept-Language": "ru-ru",
		//"LL-GA-ClientId": "1865379537.1639590360",
		"Connection": "keep-alive",
	}).Get(dst)
	if err != nil {
		return nil, err
	}

	var out SuggestionResult
	if err := json.Unmarshal(res.Body(), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

