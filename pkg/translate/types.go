package translate

import (
    "errors"
    "fmt"
)

type HTTPError struct {
    Code int
    Description string
}

func (c HTTPError) Error() string {
    return fmt.Sprintf("HTTP Error [code:%v]:%s", c.Code, c.Description)
}


type Variant struct {
    Word string
    Meaning string
}

type TranslateGoogleAPIResponse struct {
    Text string `json:"text"`
    FromLang string `json:"translated_from"`
    FromLangNativeName string `json:"from_lang_native_name"`
    Variants []*Variant       `json:"synonyms,omitempty"`
    SourceRomanization string `json:"source_romanization"`
    Images []string `json:"images"`
}


var ErrTTSLanguageNotSupported = errors.New("translateTTS js object not found")
var ErrLangNotSupported = errors.New("language is not supported by Reverso")
var SameLangsWerePassed = errors.New("the same languages were passed to ReversoTranslate()")

type ReversoRequestTranslate struct {
    Format  string `json:"format"`
    From    string `json:"from"`
    Input   string `json:"input"`
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
    ID                string      `json:"id"`
    From              string      `json:"from"`
    To                string      `json:"to"`
    Input             []string    `json:"input"`
    CorrectedText     interface{} `json:"correctedText"`
    Translation       []string    `json:"translation"`
    Engines           []string          `json:"engines"`
    LanguageDetection LanguageDetection `json:"languageDetection"`
    ContextResults    ContextResults    `json:"contextResults"`
    Truncated         bool              `json:"truncated"`
    TimeTaken int  `json:"timeTaken"`
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
    RudeWords      bool `json:"rudeWords"`
    Colloquialisms bool `json:"colloquialisms"`
    RiskyWords     bool       `json:"riskyWords"`
    Results        []Results  `json:"results"`
    TotalContextCallsMade int `json:"totalContextCallsMade"`
    TimeTakenContext      int `json:"timeTakenContext"`
}

type Results struct {
    Translation    string   `json:"translation"`
    SourceExamples []string `json:"sourceExamples"`
    TargetExamples []string `json:"targetExamples"`
    Rude           bool     `json:"rude"`
    Colloquial     bool     `json:"colloquial"`
    PartOfSpeech   string   `json:"partOfSpeech"`
}

var reversoSupportedLangs = map[string]string{
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
    v, ok := reversoSupportedLangs[iso6392]
    if !ok {
        return ""
    }
    return v
}

func ReversoIso6392(iso6391 string) string {
    for k, v := range reversoSupportedLangs {
        if v == iso6391 {
            return k
        }
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
    NrowsExact          int      `json:"nrows_exact"`
    Pagesize            int      `json:"pagesize"`
    Npages              int      `json:"npages"`
    Page                int      `json:"page"`
    RemovedSuperstrings      []string                                  `json:"removed_superstrings"`
    DictionaryEntryList      []ReversoQueryResponseDictionaryEntryList `json:"dictionary_entry_list"`
    DictionaryOtherFrequency int                                       `json:"dictionary_other_frequency"`
    DictionaryNrows          int `json:"dictionary_nrows"`
    TimeMs                   int                                       `json:"time_ms"`
    Request                  ReversoQueryResponseRequest               `json:"request"`
    Suggestions              []ReversoQueryResponseSuggestions         `json:"suggestions"`
    DymCase                  int                                       `json:"dym_case"`
    DymList                 []interface{} `json:"dym_list"`
    DymApplied              interface{}   `json:"dym_applied"`
    DymNonadaptedSearch     interface{}   `json:"dym_nonadapted_search"`
    DymPairApplied          interface{}   `json:"dym_pair_applied"`
    DymNonadaptedSearchPair interface{}   `json:"dym_nonadapted_search_pair"`
    DymPair                 interface{}   `json:"dym_pair"`
    ExtractedPhrases        []interface{} `json:"extracted_phrases"`
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
    IsFromLOCD       bool                                                  `json:"isFromLOCD"`
    InflectedForms []ReversoQueryResponseDictionaryEntryListInflectedForms `json:"inflectedForms"`
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
    Trans string `json:"trans,omitempty"`
    Orig string `json:"orig,omitempty"`
    Backend int `json:"backend,omitempty"`
    Translit string `json:"translit,omitempty"`
    SrcTranslit string `json:"src_translit,omitempty"`
}
type Entry struct {
    Word string `json:"word"`
    ReverseTranslation []string `json:"reverse_translation"`
    Score float64 `json:"score"`
}
type Dict struct {
    Pos string `json:"pos"`
    Terms []string  `json:"terms"`
    Entry []Entry   `json:"entry"`
    BaseForm string `json:"base_form"`
    PosEnum int `json:"pos_enum"`
}
type Alternative struct {
    WordPostproc string `json:"word_postproc"`
    Score int `json:"score"`
    HasPrecedingSpace bool `json:"has_preceding_space"`
    AttachToNextToken bool `json:"attach_to_next_token"`
    Backends []int `json:"backends"`
}
type Srcunicodeoffsets struct {
    Begin int `json:"begin"`
    End int `json:"end"`
}
type AlternativeTranslations struct {
    SrcPhrase string                      `json:"src_phrase"`
    Alternative []Alternative             `json:"alternative"`
    Srcunicodeoffsets []Srcunicodeoffsets `json:"srcunicodeoffsets"`
    RawSrcSegment string                  `json:"raw_src_segment"`
    StartPos int `json:"start_pos"`
    EndPos int `json:"end_pos"`
}
type Spell struct {
}
type LdResult struct {
    Srclangs []string `json:"srclangs"`
    SrclangsConfidences []float64 `json:"srclangs_confidences"`
    ExtendedSrclangs []string `json:"extended_srclangs"`
}
type Example struct {
    Text string `json:"text"`
    SourceType int `json:"source_type"`
    DefinitionID string `json:"definition_id"`
}
type Examples struct {
    Example []Example `json:"example"`
}

type getSamplesRequest struct {
    Direction string `json:"direction"`
    Source    string `json:"source"`
    Translation string `json:"translation"`
    AppID     string `json:"appId"`
}

type GetSamplesResponse struct {
    Success     bool   `json:"success"`
    Error       bool   `json:"error"`
    Source      string `json:"source"`
    Translation string           `json:"translation"`
    Samples []Samples            `json:"samples"`
    ReverseTranslations []string `json:"reverseTranslations"`
    InflectedForms      []string      `json:"inflectedForms"`
    Transliteration     string        `json:"transliteration"`
    Vowels              interface{}   `json:"vowels"`
    Suggestions         []string `json:"suggestions"`
}

type Samples struct {
    Source string `json:"source"`
    Target string `json:"target"`
    FavoriteID interface{} `json:"favoriteId"`
    IsGood bool `json:"isGood"`
}

type GoogleDictionaryResponse struct {
    DictionaryData []DictionaryData `json:"dictionaryData"`
    Status int                      `json:"status"`
}
type Desktop struct {
    Width int `json:"width"`
    Height int `json:"height"`
    URL string `json:"url"`
}
type Mobile struct {
    Width int `json:"width"`
    Height int `json:"height"`
    URL string `json:"url"`
}
type Tablet struct {
    Width int `json:"width"`
    Height int `json:"height"`
    URL string `json:"url"`
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
    Text string           `json:"text"`
}
type Etymology struct {
    Images    Images        `json:"images"`
    Etymology DeepEtymology `json:"etymology"`
}
type Phonetics struct {
    Text string `json:"text"`
    OxfordAudio string `json:"oxfordAudio"`
}
type FormType struct {
    PosTag string `json:"posTag"`
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
    Text string           `json:"text"`
}
type RelevantTopics struct {
    Name string `json:"name"`
    ClusterName string `json:"clusterName"`
}
type Senses struct {
    ConciseDefinition string `json:"conciseDefinition"`
    DomainClasses  []string         `json:"domainClasses"`
    ExampleGroups  []ExampleGroups  `json:"exampleGroups"`
    Definition     Definition       `json:"definition"`
    SourceID       string           `json:"sourceId"`
    RelevantTopics []RelevantTopics `json:"relevantTopics"`
    WsdSenseIds    []string         `json:"wsdSenseIds"`
    AdditionalExamples []string `json:"additionalExamples"`
}
type PartsOfSpeech struct {
    Value string `json:"value"`
}
type SenseFamilies struct {
    MorphUnits []MorphUnits       `json:"morphUnits"`
    Senses []Senses               `json:"senses"`
    PartsOfSpeech []PartsOfSpeech `json:"partsOfSpeech"`
}
type Corpus struct {
    Name string `json:"name"`
    Language string `json:"language"`
}
type Entries struct {
    Headword      string          `json:"headword"`
    Etymology     Etymology       `json:"etymology"`
    Phonetics     []Phonetics     `json:"phonetics"`
    SenseFamilies []SenseFamilies `json:"senseFamilies"`
    Locale        string          `json:"locale"`
    HeadwordMatchesUserQuery bool `json:"headwordMatchesUserQuery"`
    EntrySeqNo int `json:"entrySeqNo"`
    SourceID      string          `json:"sourceId"`
    Corpus        Corpus          `json:"corpus"`
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
    StatusCode float64
    Transcription string
    Pos string // noun, verb...
}

type Pos struct {
    Code string `json:"code"`
    Text string `json:"text"`
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
    Tr []Tr     `json:"tr"`
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
    Headers []string `json:"headers"`
    Rows [][]string `json:"rows"`
}
type Data struct {
    Tables []Tables `json:"tables"`
}
type Prdg struct {
    Irreg bool  `json:"irreg"`
    Data []Data `json:"data"`
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
    Fuzzy1 []ReversoFuzzy           `json:"fuzzy1"`
    Fuzzy2 []ReversoFuzzy           `json:"fuzzy2"`
    Request struct {
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