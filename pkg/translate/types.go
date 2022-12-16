package translate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"
)

type HTTPError struct {
	Code        int
	Description string
}

func (c HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error [code:%v]:%s", c.Code, c.Description)
}

type Variant struct {
	Word    string
	Meaning string
}

type TranslateGoogleAPIResponse struct {
	Text     string `json:"text"`
	FromLang string `json:"translated_from"`
	//FromLangNativeName string `json:"from_lang_native_name"`
	//Variants           []*Variant `json:"synonyms,omitempty"`
	//SourceRomanization string     `json:"source_romanization"`
	ReverseTranslations map[string][]string `json:"reverse_translations"`
	CommunityVerified   bool                `json:"community_verified"`
	//Images              []string            `json:"images"`
}

var ErrTTSLanguageNotSupported = errors.New("translateTTS js object not found")
var ErrLangNotSupported = errors.New("language is not supported by Reverso")
var SameLangsWerePassed = errors.New("the same languages were passed to ReversoTranslate()")

type ReversoRequestTranslate struct {
	Format  string                         `json:"format"`
	From    string                         `json:"from"`
	Input   string                         `json:"input"`
	To      string                         `json:"to"`
	Options ReversoRequestTranslateOptions `json:"options"`
}

type ReversoRequestTranslateOptions struct {
	Origin            string `json:"origin"`
	SentenceSplitter  bool   `json:"sentenceSplitter"`
	ContextResults    bool   `json:"contextResults"`
	LanguageDetection bool   `json:"languageDetection"`
}

type ReversoTranslation struct {
	ID                string            `json:"id"`
	From              string            `json:"from"`
	To                string            `json:"to"`
	Input             []string          `json:"input"`
	CorrectedText     interface{}       `json:"correctedText"`
	Translation       []string          `json:"translation"`
	Engines           []string          `json:"engines"`
	LanguageDetection LanguageDetection `json:"languageDetection"`
	ContextResults    ContextResults    `json:"contextResults"`
	Truncated         bool              `json:"truncated"`
	TimeTaken         int               `json:"timeTaken"`
}

type LanguageDetection struct {
	DetectedLanguage                string `json:"detectedLanguage"`
	IsDirectionChanged              bool   `json:"isDirectionChanged"`
	OriginalDirection               string `json:"originalDirection"`
	OriginalDirectionContextMatches int    `json:"originalDirectionContextMatches"`
	ChangedDirectionContextMatches  int    `json:"changedDirectionContextMatches"`
	TimeTaken                       int    `json:"timeTaken"`
}

type ContextResults struct {
	RudeWords             bool      `json:"rudeWords"`
	Colloquialisms        bool      `json:"colloquialisms"`
	RiskyWords            bool      `json:"riskyWords"`
	Results               []Results `json:"results"`
	TotalContextCallsMade int       `json:"totalContextCallsMade"`
	TimeTakenContext      int       `json:"timeTakenContext"`
}

type Results struct {
	Translation    string   `json:"translation"`
	SourceExamples []string `json:"sourceExamples"`
	TargetExamples []string `json:"targetExamples"`
	Rude           bool     `json:"rude"`
	Colloquial     bool     `json:"colloquial"`
	PartOfSpeech   string   `json:"partOfSpeech"`
}

var ReversoSupportedLangs = map[string]string{
	"dut": "nl",
	"ita": "it",
	"ger": "de",
	"rus": "ru",
	"ara": "ar",
	"jpn": "ja",
	"chi": "zh",
	"rum": "ro",
	"heb": "he",
	"por": "pt",
	"tur": "tr",
	"spa": "es",
	"pol": "pl",
	"fra": "fr",
	"eng": "en",
}

func ReversoIso6391(iso6392 string) string {
	v, ok := ReversoSupportedLangs[iso6392]
	if !ok {
		return ""
	}
	return v
}

func ReversoIso6392(iso6391 string) string {
	for k, v := range ReversoSupportedLangs {
		if v == iso6391 {
			return k
		}
	}
	if _, ok := ReversoSupportedLangs[iso6391]; ok { // нам и так дали iso6392
		return iso6391
	}
	return ""
}

type ReversoQueryRequest struct {
	Npage      int    `json:"npage"`
	Mode       int    `json:"mode"`
	SourceText string `json:"source_text"`
	TargetText string `json:"target_text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type ReversoQueryResponse struct {
	List                     []ReversoQueryResponseList                `json:"list"`
	Nrows                    int                                       `json:"nrows"`
	NrowsExact               int                                       `json:"nrows_exact"`
	Pagesize                 int                                       `json:"pagesize"`
	Npages                   int                                       `json:"npages"`
	Page                     int                                       `json:"page"`
	RemovedSuperstrings      []string                                  `json:"removed_superstrings"`
	DictionaryEntryList      []ReversoQueryResponseDictionaryEntryList `json:"dictionary_entry_list"`
	DictionaryOtherFrequency int                                       `json:"dictionary_other_frequency"`
	DictionaryNrows          int                                       `json:"dictionary_nrows"`
	TimeMs                   int                                       `json:"time_ms"`
	Request                  ReversoQueryResponseRequest               `json:"request"`
	Suggestions              []ReversoQueryResponseSuggestions         `json:"suggestions"`
	DymCase                  int                                       `json:"dym_case"`
	DymList                  []interface{}                             `json:"dym_list"`
	DymApplied               interface{}                               `json:"dym_applied"`
	DymNonadaptedSearch      interface{}                               `json:"dym_nonadapted_search"`
	DymPairApplied           interface{}                               `json:"dym_pair_applied"`
	DymNonadaptedSearchPair  interface{}                               `json:"dym_nonadapted_search_pair"`
	DymPair                  interface{}                               `json:"dym_pair"`
	ExtractedPhrases         []interface{}                             `json:"extracted_phrases"`
}

type ReversoQueryResponseList struct {
	SText string `json:"s_text"`
	TText string `json:"t_text"`
	Ref   string `json:"ref"`
	Cname string `json:"cname"`
	URL   string `json:"url"`
	Ctags string `json:"ctags"`
	Pba   bool   `json:"pba"`
}

type ReversoQueryResponseDictionaryEntryList struct {
	Frequency        int                                                     `json:"frequency"`
	Term             string                                                  `json:"term"`
	IsFromDict       bool                                                    `json:"isFromDict"`
	IsPrecomputed    bool                                                    `json:"isPrecomputed"`
	Stags            []interface{}                                           `json:"stags"`
	Pos              string                                                  `json:"pos"`
	Sourcepos        []string                                                `json:"sourcepos"`
	Variant          interface{}                                             `json:"variant"`
	Domain           interface{}                                             `json:"domain"`
	Definition       interface{}                                             `json:"definition"`
	Vowels2          interface{}                                             `json:"vowels2"`
	Transliteration2 string                                                  `json:"transliteration2"`
	AlignFreq        int                                                     `json:"alignFreq"`
	ReverseValidated bool                                                    `json:"reverseValidated"`
	PosGroup         int                                                     `json:"pos_group"`
	IsTranslation    bool                                                    `json:"isTranslation"`
	IsFromLOCD       bool                                                    `json:"isFromLOCD"`
	InflectedForms   []ReversoQueryResponseDictionaryEntryListInflectedForms `json:"inflectedForms"`
}

type ReversoQueryResponseDictionaryEntryListInflectedForms struct {
	Frequency        int           `json:"frequency"`
	Term             string        `json:"term"`
	IsFromDict       bool          `json:"isFromDict"`
	IsPrecomputed    bool          `json:"isPrecomputed"`
	Stags            []interface{} `json:"stags"`
	Pos              string        `json:"pos"`
	Sourcepos        []string      `json:"sourcepos"`
	Variant          interface{}   `json:"variant"`
	Domain           interface{}   `json:"domain"`
	Definition       interface{}   `json:"definition"`
	Vowels2          interface{}   `json:"vowels2"`
	Transliteration2 string        `json:"transliteration2"`
	AlignFreq        int           `json:"alignFreq"`
	ReverseValidated bool          `json:"reverseValidated"`
	PosGroup         int           `json:"pos_group"`
	IsTranslation    bool          `json:"isTranslation"`
	IsFromLOCD       bool          `json:"isFromLOCD"`
	InflectedForms   []interface{} `json:"inflectedForms"`
}

type ReversoQueryResponseRequest struct {
	SourceText      string      `json:"source_text"`
	TargetText      string      `json:"target_text"`
	SourceLang      string      `json:"source_lang"`
	TargetLang      string      `json:"target_lang"`
	Npage           int         `json:"npage"`
	Corpus          interface{} `json:"corpus"`
	Nrows           int         `json:"nrows"`
	Adapted         bool        `json:"adapted"`
	NonadaptedText  string      `json:"nonadapted_text"`
	RudeWords       bool        `json:"rude_words"`
	Colloquialisms  bool        `json:"colloquialisms"`
	RiskyWords      bool        `json:"risky_words"`
	Mode            int         `json:"mode"`
	ExprSug         int         `json:"expr_sug"`
	DymApply        bool        `json:"dym_apply"`
	PosReorder      int         `json:"pos_reorder"`
	Device          int         `json:"device"`
	SplitLong       bool        `json:"split_long"`
	HasLocd         bool        `json:"has_locd"`
	AvoidLOCD       bool        `json:"avoidLOCD"`
	SourcePos       interface{} `json:"source_pos"`
	ToolwordRequest bool        `json:"toolwordRequest"`
}

type ReversoQueryResponseSuggestions struct {
	Lang       string `json:"lang"`
	Suggestion string `json:"suggestion"`
	Weight     int    `json:"weight"`
	IsFromDict bool   `json:"isFromDict"`
}

type GoogleTranslateSingleResult struct {
	Sentences               []Sentences               `json:"sentences"`
	Dict                    []Dict                    `json:"dict"`
	Src                     string                    `json:"src"`
	AlternativeTranslations []AlternativeTranslations `json:"alternative_translations"`
	Confidence              float64                   `json:"confidence"`
	Spell                   Spell                     `json:"spell"`
	LdResult                LdResult                  `json:"ld_result"`
	Examples                Examples                  `json:"examples"`
}
type Sentences struct {
	Trans       string `json:"trans,omitempty"`
	Orig        string `json:"orig,omitempty"`
	Backend     int    `json:"backend,omitempty"`
	Translit    string `json:"translit,omitempty"`
	SrcTranslit string `json:"src_translit,omitempty"`
}
type Entry struct {
	Word               string   `json:"word"`
	ReverseTranslation []string `json:"reverse_translation"`
	Score              float64  `json:"score"`
}
type Dict struct {
	Pos      string   `json:"pos"`
	Terms    []string `json:"terms"`
	Entry    []Entry  `json:"entry"`
	BaseForm string   `json:"base_form"`
	PosEnum  int      `json:"pos_enum"`
}
type Alternative struct {
	WordPostproc      string `json:"word_postproc"`
	Score             int    `json:"score"`
	HasPrecedingSpace bool   `json:"has_preceding_space"`
	AttachToNextToken bool   `json:"attach_to_next_token"`
	Backends          []int  `json:"backends"`
}
type Srcunicodeoffsets struct {
	Begin int `json:"begin"`
	End   int `json:"end"`
}
type AlternativeTranslations struct {
	SrcPhrase         string              `json:"src_phrase"`
	Alternative       []Alternative       `json:"alternative"`
	Srcunicodeoffsets []Srcunicodeoffsets `json:"srcunicodeoffsets"`
	RawSrcSegment     string              `json:"raw_src_segment"`
	StartPos          int                 `json:"start_pos"`
	EndPos            int                 `json:"end_pos"`
}
type Spell struct {
}
type LdResult struct {
	Srclangs            []string  `json:"srclangs"`
	SrclangsConfidences []float64 `json:"srclangs_confidences"`
	ExtendedSrclangs    []string  `json:"extended_srclangs"`
}
type Example struct {
	Text         string `json:"text"`
	SourceType   int    `json:"source_type"`
	DefinitionID string `json:"definition_id"`
}
type Examples struct {
	Example []Example `json:"example"`
}

type getSamplesRequest struct {
	Direction   string `json:"direction"`
	Source      string `json:"source"`
	Translation string `json:"translation"`
	AppID       string `json:"appId"`
}

type GetSamplesResponse struct {
	Success             bool        `json:"success"`
	Error               bool        `json:"error"`
	Source              string      `json:"source"`
	Translation         string      `json:"translation"`
	Samples             []Samples   `json:"samples"`
	ReverseTranslations []string    `json:"reverseTranslations"`
	InflectedForms      []string    `json:"inflectedForms"`
	Transliteration     string      `json:"transliteration"`
	Vowels              interface{} `json:"vowels"`
	Suggestions         []string    `json:"suggestions"`
}

type Samples struct {
	Source     string      `json:"source"`
	Target     string      `json:"target"`
	FavoriteID interface{} `json:"favoriteId"`
	IsGood     bool        `json:"isGood"`
}

type GoogleDictionaryResponse struct {
	DictionaryData []DictionaryData `json:"dictionaryData"`
	Status         int              `json:"status"`
}
type Desktop struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}
type Mobile struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}
type Tablet struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}
type Images struct {
	Desktop Desktop `json:"desktop"`
	Mobile  Mobile  `json:"mobile"`
	Tablet  Tablet  `json:"tablet"`
}
type Fragments struct {
	Text string `json:"text"`
}
type DeepEtymology struct {
	Fragments []Fragments `json:"fragments"`
	Text      string      `json:"text"`
}
type Etymology struct {
	Images    Images        `json:"images"`
	Etymology DeepEtymology `json:"etymology"`
}
type Phonetics struct {
	Text        string `json:"text"`
	OxfordAudio string `json:"oxfordAudio"`
}
type FormType struct {
	PosTag      string `json:"posTag"`
	Description string `json:"description"`
}
type MorphUnits struct {
	FormType FormType `json:"formType"`
	WordForm string   `json:"wordForm"`
}
type ExampleGroups struct {
	Examples []string `json:"examples"`
}
type Definition struct {
	Fragments []Fragments `json:"fragments"`
	Text      string      `json:"text"`
}
type RelevantTopics struct {
	Name        string `json:"name"`
	ClusterName string `json:"clusterName"`
}
type Senses struct {
	ConciseDefinition  string           `json:"conciseDefinition"`
	DomainClasses      []string         `json:"domainClasses"`
	ExampleGroups      []ExampleGroups  `json:"exampleGroups"`
	Definition         Definition       `json:"definition"`
	SourceID           string           `json:"sourceId"`
	RelevantTopics     []RelevantTopics `json:"relevantTopics"`
	WsdSenseIds        []string         `json:"wsdSenseIds"`
	AdditionalExamples []string         `json:"additionalExamples"`
}
type PartsOfSpeech struct {
	Value string `json:"value"`
}
type SenseFamilies struct {
	MorphUnits    []MorphUnits    `json:"morphUnits"`
	Senses        []Senses        `json:"senses"`
	PartsOfSpeech []PartsOfSpeech `json:"partsOfSpeech"`
}
type Corpus struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}
type Entries struct {
	Headword                 string          `json:"headword"`
	Etymology                Etymology       `json:"etymology"`
	Phonetics                []Phonetics     `json:"phonetics"`
	SenseFamilies            []SenseFamilies `json:"senseFamilies"`
	Locale                   string          `json:"locale"`
	HeadwordMatchesUserQuery bool            `json:"headwordMatchesUserQuery"`
	EntrySeqNo               int             `json:"entrySeqNo"`
	SourceID                 string          `json:"sourceId"`
	Corpus                   Corpus          `json:"corpus"`
}
type UsageOverTimeImage struct {
	Desktop Desktop `json:"desktop"`
	Mobile  Mobile  `json:"mobile"`
	Tablet  Tablet  `json:"tablet"`
}
type DictionaryData struct {
	Entries            []Entries          `json:"entries"`
	QueryTerm          string             `json:"queryTerm"`
	UsageOverTimeImage UsageOverTimeImage `json:"usageOverTimeImage"`
}

type YandexTranscriptionResponse struct {
	StatusCode    float64
	Transcription string
	Pos           string // noun, verb...
}

type Pos struct {
	Code    string `json:"code"`
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}
type Syn struct {
	Text string `json:"text"`
	Pos  Pos    `json:"pos"`
	Fr   int    `json:"fr"`
}
type Mean struct {
	Text string `json:"text"`
}

type Ex struct {
	Text string `json:"text"`
	Tr   []Tr   `json:"tr"`
}
type Tr struct {
	Text string `json:"text"`
	Pos  Pos    `json:"pos"`
	Fr   int    `json:"fr"`
	Syn  []Syn  `json:"syn,omitempty"`
	Mean []Mean `json:"mean"`
	Ex   []Ex   `json:"ex,omitempty"`
}
type Tables struct {
	Headers []string   `json:"headers"`
	Rows    [][]string `json:"rows"`
}
type Data struct {
	Tables []Tables `json:"tables"`
}
type Prdg struct {
	Irreg bool   `json:"irreg"`
	Data  []Data `json:"data"`
}
type Regular struct {
	Text string `json:"text"`
	Pos  Pos    `json:"pos"`
	Ts   string `json:"ts"`
	Tr   []Tr   `json:"tr"`
	Prdg Prdg   `json:"prdg,omitempty"`
}
type yandexTranscriptionResponse struct {
	Regular []Regular `json:"regular"`
}

type ReversoSuggestionsResponse struct {
	Suggestions []ReversoSuggestion `json:"suggestions"`
	Fuzzy1      []ReversoFuzzy      `json:"fuzzy1"`
	Fuzzy2      []ReversoFuzzy      `json:"fuzzy2"`
	Request     struct {
		Search     string `json:"search"`
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
		MaxResults int    `json:"max_results"`
		Mode       int    `json:"mode"`
	} `json:"request"`
	TimeMs int `json:"time_ms"`
}

type ReversoSuggestion struct {
	Lang       string `json:"lang"`
	Suggestion string `json:"suggestion"`
	Weight     int    `json:"weight"`
	IsFromDict bool   `json:"isFromDict"`
}

type ReversoFuzzy struct {
	Lang       string `json:"lang"`
	Suggestion string `json:"suggestion"`
	Weight     int    `json:"weight"`
	IsFromDict bool   `json:"isFromDict"`
}

type reversoSuggestionRequest struct {
	Search     string `json:"search"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

var YandexSupportedLanguages = map[string]string{
	"af":     "Африкаанс",
	"am":     "Амхарский",
	"ar":     "Арабский",
	"az":     "Азербайджанский",
	"ba":     "Башкирский",
	"be":     "Белорусский",
	"bg":     "Болгарский",
	"bn":     "Бенгальский",
	"bs":     "Боснийский",
	"ca":     "Каталанский",
	"ceb":    "Себуанский",
	"cv":     "Чувашский",
	"cy":     "Валлийский",
	"da":     "Датский",
	"de":     "Немецкий",
	"el":     "Греческий",
	"emj":    "Эмодзи",
	"en":     "Английский",
	"eo":     "Эсперанто",
	"es":     "Испанский",
	"et":     "Эстонский",
	"eu":     "Баскский",
	"fa":     "Персидский",
	"fi":     "Финский",
	"fr":     "Французский",
	"ga":     "Ирландский",
	"gd":     "Шотландский (гэльский)",
	"gl":     "Галисийский",
	"gu":     "Гуджарати",
	"he":     "Иврит",
	"hi":     "Хинди",
	"hr":     "Хорватский",
	"ht":     "Гаитянский",
	"hu":     "Венгерский",
	"hy":     "Армянский",
	"id":     "Индонезийский",
	"is":     "Исландский",
	"it":     "Итальянский",
	"ja":     "Японский",
	"jv":     "Яванский",
	"ka":     "Грузинский",
	"kazlat": "Казахский (латиница)",
	"kk":     "Казахский",
	"km":     "Кхмерский",
	"kn":     "Каннада",
	"ko":     "Корейский",
	"ky":     "Киргизский",
	"la":     "Латынь",
	"lb":     "Люксембургский",
	"lo":     "Лаосский",
	"lt":     "Литовский",
	"lv":     "Латышский",
	"mg":     "Малагасийский",
	"mhr":    "Марийский",
	"mi":     "Маори",
	"mk":     "Македонский",
	"ml":     "Малаялам",
	"mn":     "Монгольский",
	"mr":     "Маратхи",
	"mrj":    "Горномарийский",
	"ms":     "Малайский",
	"mt":     "Мальтийский",
	"my":     "Бирманский",
	"ne":     "Непальский",
	"nl":     "Нидерландский",
	"no":     "Норвежский",
	"pa":     "Панджаби",
	"pap":    "Папьяменто",
	"pl":     "Польский",
	"pt":     "Португальский",
	"ro":     "Румынский",
	"ru":     "Русский",
	"sah":    "Якутский",
	"si":     "Сингальский",
	"sjn":    "Эльфийский (синдарин)",
	"sl":     "Словенский",
	"sq":     "Албанский",
	"su":     "Сунданский",
	"sv":     "Шведский",
	"sw":     "Суахили",
	"ta":     "Тамильский",
	"te":     "Телугу",
	"tg":     "Таджикский",
	"th":     "Тайский",
	"tl":     "Тагальский",
	"tr":     "Турецкий",
	"tt":     "Татарский",
	"udm":    "Удмуртский",
	"uk":     "Украинский",
	"ur":     "Урду",
	"uz":     "Узбекский",
	"uzbcyr": "Узбекский (кириллица)",
	"vi":     "Вьетнамский",
	"xh":     "Коса",
	"yi":     "Идиш",
	"zh":     "Китайский",
	"zu":     "Зулу",
}

var MicrosoftUnsupportedLanguages = []string{
	"ia",
	"lb",
	"pi",
	"la",
	"cu",
	"ng",
	"ln",
	"bi",
	"kr",
	"br",
	"ch",
	"hz",
	"sd",
	"gn",
	"ig",
	"yi",
	"vo",
	"os",
	"ts",
	"xh",
	"co",
	"oc",
	"qu",
	"se",
	"si",
	"su",
	"wo",
	"sc",
	"lu",
	"sa",
	"ki",
	"yo",
	"ho",
	"nv",
	"kv",
	"an",
	"cr",
	"rw",
	"rm",
	"ee",
	"ff",
	"sg",
	"ii",
	"ae",
	"ks",
	"ik",
	"kw",
	"li",
	"av",
	"fy",
	"ve",
	"lg",
	"io",
	"ie",
	"rn",
	"nr",
	"kj",
	"ay",
	"nd",
	"jv",
	"st",
	"ny",
	"tn",
	"eo",
	"na",
	"kl",
	"mh",
	"gd",
	"ce",
	"om",
	"wa",
	"ss",
	"za",
	"bh",
	"ha",
	"ab",
	"nn",
	"kg",
	"sn",
	"gv",
	"oj",
	"be",
	"tl",
}

// SplitIntoChunksBySentences. You can merge output by ""
func SplitIntoChunksBySentences(text string, limit int) []string {
	if len(text) < limit {
		return []string{text}
	}
	chunks := make([]string, 0, len(text)/limit+1)
	points := utf16.Encode([]rune(text))
	for i := 0; i < len(points); {
		offset := indexDelim(points[i:], limit, ".!?;\r\n\t\f\v*)")
		ch := string(utf16.Decode(points[i : i+offset]))
		chunks = append(chunks, ch)
		i += offset
	}
	return chunks
}

func CheckHtmlTags(in, out string) string {
	r, err := regexp.Compile("<\\s*[^>]+>(.*?)")
	if err != nil {
		return "regexp.Compile err"
	}
	inTags := r.FindAllString(in, -1)
	outTags := r.FindAllString(out, -1)
	//pp.Println(inTags, len(inTags), outTags, len(outTags))
	//if len(inTags) != len(outTags) {
	//	return "checkHtmlTags: tags number in "in" and "out" don"t match"
	//}
	inTagsLen := len(inTags)
	outTagsLen := len(outTags)
	for i := 0; i < len(inTags); i++ {
		if i > inTagsLen-1 {
			out += strings.Join(outTags[i:], "")
			break
		} else if i > outTagsLen-1 {
			out += strings.Join(outTags[i:], "")
			break
		}
		if inTags[i] != outTags[i] {
			out = strings.Replace(out, outTags[i], inTags[i], 1)
		}
	}
	return out
}

func indexDelim(text []uint16, limit int, delims string) (offset int) {
	if len(text) < limit {
		return len(text)
	} else if len(text) > limit {
		text = text[:limit]
	}
	offset = len(text)
	delimeters := utf16.Encode([]rune(delims))
	for i := len(text) - 1; i >= 0; i-- {
		if in(delimeters, text[i]) {
			return i + 1
		}
	}
	return offset
}

func in(arr []uint16, keys ...uint16) bool {
	for _, v := range arr {
		for _, k := range keys {
			if k == v {
				return true
			}
		}
	}
	return false
}

//func lastIndexAny(text string, chars string) (idx int, delim string) {
//	idx--
//	for _, ch := range chars {
//		c := string(ch)
//		i := strings.LastIndex(text, c)
//		if i > idx {
//			idx, delim = i, c
//		}
//	}
//	return
//}

var SupportedLanguageCodes = []string{
	"lu",
	"hu",
	"ko",
	"zh-tw",
	"ik",
	"os",
	"th",
	"id",
	"ja",
	"mt",
	"nb",
	"fa",
	"ve",
	"uk",
	"kv",
	"is",
	"be",
	"om",
	"ta",
	"ro",
	"ky",
	"hi",
	"lo",
	"kw",
	"sq",
	"na",
	"rm",
	"tr",
	"av",
	"he",
	"or",
	"as",
	"ay",
	"tk",
	"kj",
	"kg",
	"zh-hk",
	"gn",
	"fo",
	"si",
	"sw",
	"nn",
	"fr",
	"bh",
	"bn",
	"bs",
	"mh",
	"el",
	"to",
	"bo",
	"gv",
	"ka",
	"sk",
	"cs",
	"lb",
	"ho",
	"yi",
	"hz",
	"st",
	"zh",
	"cy",
	"gl",
	"uz",
	"ne",
	"hy",
	"ru",
	"ch",
	"rn",
	"wa",
	"gd",
	"mg",
	"nr",
	"km",
	"no",
	"sd",
	"jv",
	"zh-sg",
	"sa",
	"io",
	"sn",
	"cv",
	"ny",
	"ng",
	"ms",
	"fy",
	"kr",
	"zh-mo",
	"ab",
	"ee",
	"su",
	"ff",
	"sm",
	"ps",
	"an",
	"sg",
	"ie",
	"tn",
	"ks",
	"ss",
	"ha",
	"nv",
	"co",
	"my",
	"et",
	"ga",
	"ae",
	"ia",
	"eo",
	"tl",
	"bg",
	"ki",
	"iu",
	"za",
	"es",
	"pa",
	"cr",
	"ug",
	"gu",
	"sv",
	"ii",
	"ht",
	"so",
	"ty",
	"am",
	"tg",
	"mr",
	"te",
	"ce",
	"zh-cn",
	"qu",
	"hr",
	"it",
	"fi",
	"da",
	"oj",
	"li",
	"lt",
	"de",
	"ig",
	"az",
	"ln",
	"vo",
	"dv",
	"mn",
	"kn",
	"sl",
	"en",
	"af",
	"nd",
	"la",
	"sr",
	"fj",
	"yo",
	"mi",
	"ml",
	"se",
	"kk",
	"pl",
	"vi",
	"sc",
	"oc",
	"bi",
	"br",
	"pi",
	"ku",
	"pt",
	"ti",
	"kl",
	"nl",
	"rw",
	"cu",
	"wo",
	"ar",
	"eu",
	"ba",
	"ur",
	"lv",
	"ca",
	"tt",
	"lg",
	"ts",
	"xh",
	"mk",
}

var DEFAULT_GOOGLE_URLS = []string{"translate.google.ac", "translate.google.ad", "translate.google.ae",
	"translate.google.al", "translate.google.am", "translate.google.as",
	"translate.google.at", "translate.google.az", "translate.google.ba",
	"translate.google.be", "translate.google.bf", "translate.google.bg",
	"translate.google.bi", "translate.google.bj", "translate.google.bs",
	"translate.google.bt", "translate.google.by", "translate.google.ca",
	"translate.google.cat", "translate.google.cc", "translate.google.cd",
	"translate.google.cf", "translate.google.cg", "translate.google.ch",
	"translate.google.ci", "translate.google.cl", "translate.google.cm",
	"translate.google.cn", "translate.google.co.ao", "translate.google.co.bw",
	"translate.google.co.ck", "translate.google.co.cr", "translate.google.co.id",
	"translate.google.co.il", "translate.google.co.in", "translate.google.co.jp",
	"translate.google.co.ke", "translate.google.co.kr", "translate.google.co.ls",
	"translate.google.co.ma", "translate.google.co.mz", "translate.google.co.nz",
	"translate.google.co.th", "translate.google.co.tz", "translate.google.co.ug",
	"translate.google.co.uk", "translate.google.co.uz", "translate.google.co.ve",
	"translate.google.co.vi", "translate.google.co.za", "translate.google.co.zm",
	"translate.google.co.zw", "translate.google.com.af", "translate.google.com.ag",
	"translate.google.com.ai", "translate.google.com.ar", "translate.google.com.au",
	"translate.google.com.bd", "translate.google.com.bh", "translate.google.com.bn",
	"translate.google.com.bo", "translate.google.com.br", "translate.google.com.bz",
	"translate.google.com.co", "translate.google.com.cu", "translate.google.com.cy",
	"translate.google.com.do", "translate.google.com.ec", "translate.google.com.eg",
	"translate.google.com.et", "translate.google.com.fj", "translate.google.com.gh",
	"translate.google.com.gi", "translate.google.com.gt", "translate.google.com.hk",
	"translate.google.com.jm", "translate.google.com.kh", "translate.google.com.kw",
	"translate.google.com.lb", "translate.google.com.ly", "translate.google.com.mm",
	"translate.google.com.mt", "translate.google.com.mx", "translate.google.com.my",
	"translate.google.com.na", "translate.google.com.ng", "translate.google.com.ni",
	"translate.google.com.np", "translate.google.com.om", "translate.google.com.pa",
	"translate.google.com.pe", "translate.google.com.pg", "translate.google.com.ph",
	"translate.google.com.pk", "translate.google.com.pr", "translate.google.com.py",
	"translate.google.com.qa", "translate.google.com.sa", "translate.google.com.sb",
	"translate.google.com.sg", "translate.google.com.sl", "translate.google.com.sv",
	"translate.google.com.tj", "translate.google.com.tr", "translate.google.com.tw",
	"translate.google.com.ua", "translate.google.com.uy", "translate.google.com.vc",
	"translate.google.com.vn", "translate.google.com", "translate.google.cv",
	"translate.google.cz", "translate.google.de", "translate.google.dj",
	"translate.google.dk", "translate.google.dm", "translate.google.dz",
	"translate.google.ee", "translate.google.es", "translate.google.eu",
	"translate.google.fi", "translate.google.fm", "translate.google.fr",
	"translate.google.ga", "translate.google.ge", "translate.google.gf",
	"translate.google.gg", "translate.google.gl", "translate.google.gm",
	"translate.google.gp", "translate.google.gr", "translate.google.gy",
	"translate.google.hn", "translate.google.hr", "translate.google.ht",
	"translate.google.hu", "translate.google.ie", "translate.google.im",
	"translate.google.io", "translate.google.iq", "translate.google.is",
	"translate.google.it", "translate.google.je", "translate.google.jo",
	"translate.google.kg", "translate.google.ki", "translate.google.kz",
	"translate.google.la", "translate.google.li", "translate.google.lk",
	"translate.google.lt", "translate.google.lu", "translate.google.lv",
	"translate.google.md", "translate.google.me", "translate.google.mg",
	"translate.google.mk", "translate.google.ml", "translate.google.mn",
	"translate.google.ms", "translate.google.mu", "translate.google.mv",
	"translate.google.mw", "translate.google.ne", "translate.google.nf",
	"translate.google.nl", "translate.google.no", "translate.google.nr",
	"translate.google.nu", "translate.google.pl", "translate.google.pn",
	"translate.google.ps", "translate.google.pt", "translate.google.ro",
	"translate.google.rs", "translate.google.ru", "translate.google.rw",
	"translate.google.sc", "translate.google.se", "translate.google.sh",
	"translate.google.si", "translate.google.sk", "translate.google.sm",
	"translate.google.sn", "translate.google.so", "translate.google.sr",
	"translate.google.st", "translate.google.td", "translate.google.tg",
	"translate.google.tk", "translate.google.tl", "translate.google.tm",
	"translate.google.tn", "translate.google.to", "translate.google.tt",
	"translate.google.us", "translate.google.vg", "translate.google.vu",
	"translate.google.ws"}
