package translate

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GoogleTranslationOfficial struct {
	Translation          string
	From                 string
	Dictionary           []Meaning
	Transcription        string
	Synonyms             []string
	SynonymsTranslations map[string][]string
}

type Meaning struct {
	Word    string
	Type    string
	Meaning string
	Example string
}

func GoogleTranslateOfficial(text, from, to string) (GoogleTranslationOfficial, error) {
	if from == "" {
		from = "auto"
	}

	req, err := http.NewRequest("GET", "https://translate.google.com", nil)
	if err != nil {
		return GoogleTranslationOfficial{}, err
	}
	req.Header.Set("user-agent", "got/9.6.0 (https://github.com/sindresorhus/got)")

	res, err := client.Do(req)
	if err != nil {
		return GoogleTranslationOfficial{}, err
	}
	ret, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GoogleTranslationOfficial{}, err
	}
	rets := string(ret)

	fsid := extract("FdrFJe", rets)
	bl := extract("cfb2h", rets)

	queryParams := url.Values{}
	queryParams.Set("rpcids", "MkEWBc")
	queryParams.Set("f.sid", fsid)
	queryParams.Set("bl", bl)
	queryParams.Set("hl", "en-US")
	queryParams.Set("soc-app", "1")
	queryParams.Set("soc-platform", "1")
	queryParams.Set("soc-device", "1")
	queryParams.Set("_reqid", strconv.FormatFloat(math.Floor(float64(1000+rand.Float64()*9000)), 'f', 6, 64))
	queryParams.Set("rt", "c")

	bodyParams := url.Values{}
	bodyParams.Set("f.req", fmt.Sprintf(`[[["MkEWBc","[[\"%s\",\"%s\",\"%s\",true],[null]]",null,"generic"]]]`, url.PathEscape(text), from, to))
	bodyParams.Set("dj", "1")

	req, err = http.NewRequest("POST", "https://translate.google.com/_/TranslateWebserverUi/data/batchexecute?dj=1&"+queryParams.Encode(), strings.NewReader(bodyParams.Encode()))
	if err != nil {
		return GoogleTranslationOfficial{}, err
	}
	req.Header["content-type"] = []string{"application/x-www-form-urlencoded;charset=UTF-8"}

	res, err = client.Do(req)
	if err != nil {
		return GoogleTranslationOfficial{}, err
	}
	ret, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return GoogleTranslationOfficial{}, err
	}

	var out [][]interface{}
	if err = json.Unmarshal([]byte(strings.Split(string(ret), "\n")[3]), &out); err != nil {
		return GoogleTranslationOfficial{}, err
	}

	need := out[0][2].(string)
	fmt.Println("need", need)
	translation := gjson.Get(need, "0.4.0.0").String()
	transcription := gjson.Get(need, "0.0").String()
	fromLang := gjson.Get(need, "1.1").String()

	// Getting synonyms
	synonyms := make([]string, 0, 3)
	for _, value := range gjson.Get(need, "1.0.0.5").Array() {
		for _, value := range value.Get("4").Array() {
			synonym := value.Get("0")
			synonyms = append(synonyms, synonym.String())
		}
	}

	// Getting meanings
	meanings := make([]Meaning, 0, 3)
	word := gjson.Get(need, "3.0").String()
	for _, value := range gjson.Get(need, "3.1.0").Array() {
		//if !value.IsObject() {
		//	continue
		//}
		meanings = append(meanings, Meaning{
			Word:    word,
			Type:    value.Get("0").String(),
			Meaning: value.Get("1.0.0").String(),
			Example: value.Get("1.0.1").String(),
		})
	}

	synonymsTranslations := make(map[string][]string, 3)
	for _, value := range gjson.Get(need, "3.5.0").Array() {
		for _, value := range value.Get("1").Array() {
			term := value.Get("0").String()
			translations := make([]string, 0, 3)
			for _, translation := range value.Get("2").Array() {
				translations = append(translations, translation.String())
			}
			if _, ok := synonymsTranslations[term]; !ok {
				synonymsTranslations[term] = make([]string, 0, len(translations))
			}
			synonymsTranslations[term] = translations
		}
	}
	return GoogleTranslationOfficial{
		Translation:          translation,
		From:                 fromLang,
		Dictionary:           meanings,
		Transcription:        transcription,
		Synonyms:             synonyms,
		SynonymsTranslations: synonymsTranslations,
	}, nil
}

func extract(key, html string) string {
	q := "\"" + key + "\"" + ":" + "\""
	i1 := strings.Index(html, q) + len(q)
	i2 := strings.Index(html[i1:], "\",")
	return html[i1 : i1+i2]
}
