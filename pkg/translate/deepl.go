package translate

import (
	"encoding/json"
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var DeeplSupportedLangs = []string{
	"es",
	"pl",
	"fi",
	"hu",
	"lt",
	"sl",
	"de",
	"pt",
	"ru",
	"ja",
	"da",
	"el",
	"ro",
	"sk",
	"en",
	"fr",
	"it",
	"nl",
	"bg",
	"cs",
	"sv",
	"zh",
	"et",
	"lv",
}

type Deepl struct {
	client *resty.Client
	sessionId, userId, userAgent string
	headers                              map[string]string
	weights map[string]float64
}

func (d Deepl) getId() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(1000000000 + rand.Intn(8999999999))
}

func (d Deepl) isReplace(id int64) bool {
	return (id+3)%13 == 0 || (id+5)%29 == 0
}



func NewDeepl() (Deepl, error) {
	d := Deepl{
		client: resty.New(),
		sessionId: "821b2aeb-1074-44b7-9285-0cadc8999715",
		userId:    "6346616a-cb7f-4e8a-a0c7-d4c4dafd7d01",
		userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36",
		headers: nil,
		weights:  nil,
	}
	d.headers = map[string]string{
		"accept":          "*/*",
		"accept-encoding": "gzip, deflate, br",
		"accept-language": "uk-UA,uk;q=0.9,en-US;q=0.8,en;q=0.7",
		"Authorization":   "None",
		"Cookie": "", // добавится в loadWeightsAndCookies
		"Content-Type":    "application/json; charset=utf-8",
		"origin":          "chrome-extension://cofdbpoegempjloogbagkncekinflcnj",
		"referer":         "https://www.deepl.com/",
		"sec-fetch-dest":  "empty",
		"sec-fetch-mode":  "cors",
		"sec-fetch-site":  "none",
		"user-agent":      d.userAgent,
	}

	if err := d.loadWeightsAndCookies(); err != nil {
		return Deepl{}, err
	}

	return d, nil
}

func (d *Deepl) loadWeightsAndCookies() error {
	id := d.getId()
	data, err := json.Marshal(DeeplRequest{
		Jsonrpc: "2.0",
		Method:  "LMT_handle_texts",
		Params: Params{
			Texts:     []Texts{{
				Text: "Пример",
			}},
			Html: "enabled",
			Splitting: "newlines",
			Lang: Lang{
				TargetLang:             "en",
				SourceLangUserSelected: "auto",
				Preference:             Preference{Weight: map[string]float64{}},
			},
			Timestamp: time.Now().Unix(),
		},
		ID: id,
	})
	pp.Println(string(data))

	if err != nil {
		return err
	}

	if d.isReplace(id) {
		data = []byte(strings.ReplaceAll(string(data), `"method":"`, `"method" : "`))
	} else {
		data = []byte(strings.ReplaceAll(string(data), `"method":"`, `"method": "`))
	}


	resp, err := d.client.R().
		SetHeaders(d.headers).
		SetBody(data).
		SetContentLength(true).
		Post("https://www2.deepl.com/jsonrpc?client=chrome-extension,0.14.0")
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("Deepl loadWeightsAndCookies: not 200 http code. Response:"+resp.String())
	}

	for _, cookie := range resp.Cookies() {
		d.headers["Cookie"] += cookie.Raw + " "
	}
	pp.Println("Deepl: updated headers:", d.headers["Cookie"])

	weights := make(map[string]float64, 12)
	for k, v  := range gjson.GetBytes(resp.Body(), "result.detectedLanguages").Map() {
		if k == "unsupported" {
			continue
		}
		weights[k] = Le(0, v.Float())
	}
	d.weights = weights
	pp.Println("Deepl: updated weights", d.weights)
	return nil
}

func getCookie(cookies []*http.Cookie, key string) string {
	key = strings.ToLower(key)
	for _, cookie := range cookies {
		if strings.ToLower(cookie.Name) == key {
			return cookie.Value
		}
	}
	return ""
}

type DeeplRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int64  `json:"id"`
}
type Texts struct {
	Text string `json:"text"`
}



type Weight struct {
	BG float64 `json:"BG"`
	CS float64 `json:"CS"`
	DA float64 `json:"DA"`
	DE float64 `json:"DE"`
	EL float64 `json:"EL"`
	EN float64 `json:"EN"`
	ES float64 `json:"ES"`
	ET float64 `json:"ET"`
	FI float64 `json:"FI"`
	FR float64 `json:"FR"`
	HU float64 `json:"HU"`
	IT float64 `json:"IT"`
	JA float64 `json:"JA"`
	LT float64 `json:"LT"`
	LV float64 `json:"LV"`
	NL float64 `json:"NL"`
	PL float64 `json:"PL"`
	PT float64 `json:"PT"`
	RO float64 `json:"RO"`
	RU float64 `json:"RU"`
	SK float64 `json:"SK"`
	SL float64 `json:"SL"`
	SV float64 `json:"SV"`
	ZH float64 `json:"ZH"`
}
type Preference struct {
	Weight map[string]float64 `json:"weight"`
}
type Lang struct {
	TargetLang             string     `json:"target_lang"`
	SourceLangUserSelected string     `json:"source_lang_user_selected"`
	Preference             Preference `json:"preference"`
}
type Params struct {
	Html string `json:"html,omitempty"` //enabled
	Texts     []Texts `json:"texts"`
	Splitting string  `json:"splitting"`
	Lang      Lang    `json:"lang"`
	Timestamp int64   `json:"timestamp"`
}
type DeeplStatisticsReq struct {
	EventID                int                    `json:"eventId"`
	SessionID              string                 `json:"sessionId"`
	UserInfos              UserInfos              `json:"userInfos"`
	TranslationRequestData TranslationRequestData `json:"translationRequestData"`
}
type UserInfos struct {
	UserType int `json:"userType"`
}
type TranslationData struct {
	SourceLang   string `json:"sourceLang"`
	TargetLang   string `json:"targetLang"`
	SourceLength int    `json:"sourceLength"`
}
type TranslationRequestData struct {
	ExtensionVersion   string          `json:"extensionVersion"`
	DomainName         string          `json:"domainName"`
	UserID             string          `json:"userId"`
	TranslationTrigger int             `json:"translationTrigger"`
	TranslationData    TranslationData `json:"translationData"`
}

func Le(t, e float64) float64 {
	const n = 1e-5;
	return n * math.Floor(float64((0.97*t+e)/n))
}


func (d *Deepl) Translate(from, to, text string, withWeights bool) (string, error) {
	text = strings.ReplaceAll(text, "\n", "<br>")
	from, to = strings.ToUpper(from), strings.ToUpper(to)
	parts := splitIntoChunksBySentences(text, 4999)
	texts := make([]Texts, 0, len(parts))
	for i, part := range parts {
		texts = append(texts, Texts{
			Text: part,
		})
		parts[i] = ""
	}

	id := d.getId()
	data := make([]byte, 0, 1024)
	var err error

	if !withWeights {
		data, err = json.Marshal(DeeplRequest{
			Jsonrpc: "2.0",
			Method:  "LMT_handle_texts",
			Params: Params{
				Texts:     texts,
				Html: "enabled",
				Splitting: "newlines",
				Lang: Lang{
					TargetLang:             to,
					SourceLangUserSelected: "auto",
					Preference:             Preference{Weight: map[string]float64{}},
				},
				Timestamp: time.Now().Unix(),
			},
			ID: id,
		})
	} else {
		data, err = json.Marshal(DeeplRequest{
			Jsonrpc: "2.0",
			Method:  "LMT_handle_texts",
			Params: Params{
				Texts:     texts,
				Splitting: "newlines",
				Html: "enabled",
				Lang: Lang{
					TargetLang:             to,
					SourceLangUserSelected: from,
					Preference:             Preference{Weight: d.weights},
				},
				Timestamp: time.Now().Unix(),
			},
			ID: id,
		})
	}

	if err != nil {
		return "", err
	}

	if d.isReplace(id) {
		data = []byte(strings.ReplaceAll(string(data), `"method":"`, `"method" : "`))
	} else {
		data = []byte(strings.ReplaceAll(string(data), `"method":"`, `"method": "`))
	}

	resp, err := d.client.R().
		SetHeaders(d.headers).
		SetBody(data).
		SetContentLength(true).
		Post("https://www2.deepl.com/jsonrpc?client=chrome-extension,0.14.0")
	fmt.Println(resp.String())

	if err != nil {
		return "", err
	}
	if gjson.GetBytes(resp.Body(), "error").Exists() {
		return "", fmt.Errorf("Deepl error\nInput text[%s-%s]:%s\n\n%s", from, to, text, gjson.GetBytes(resp.Body(), "error.message").String())
	}


	out := ""
	for _, v := range gjson.GetBytes(resp.Body(), "result.texts").Array() {
		out += v.Get("text").String()
	}
	out = strings.ReplaceAll(text, "<br>", "\n")


	for lang, weight := range gjson.GetBytes(resp.Body(), "params.lang.preference.weight").Map() {
		d.weights[lang] = Le(d.weights[lang], weight.Float())
	}


	resp, err = d.client.R().SetHeaders(d.headers).SetBody(DeeplStatisticsReq{
		EventID:   60001,
		SessionID: d.sessionId,
		UserInfos: UserInfos{UserType: 1},
		TranslationRequestData: TranslationRequestData{
			ExtensionVersion:   "0.14.0",
			DomainName:         "2ip.ru",
			UserID:             d.userId,
			TranslationTrigger: 3,
			TranslationData: TranslationData{
				SourceLang:   strings.ToLower(gjson.GetBytes(resp.Body(), "result.lang").String()),
				TargetLang:   to,
				SourceLength: len([]rune(text)),
			},
		},
	}).Post("https://s.deepl.com/chrome/statistics")
	if err != nil {
		return "", err
	}


	return out, err
}
