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
    Variants []*Variant `json:"synonyms,omitempty"`
    SourceRomanization string `json:"source_romanization"`
    Images []string `json:"images"`
}


var ErrTTSLanguageNotSupported = errors.New("translateTTS js object not found")
var ErrReversoLangNotSupported = errors.New("language is not supported by Reverso")
var SameLangsWerePassed = errors.New("the same languages were passed to ReversoTranslate()")

type ReversoRequestTranslate struct {
    Format  string `json:"format"`
    From    string `json:"from"`
    Input   string `json:"input"`
    To      string `json:"to"`
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
    Engines           []string    `json:"engines"`
    LanguageDetection LanguageDetection  `json:"languageDetection"`
    ContextResults ContextResults  `json:"contextResults"`
    Truncated bool `json:"truncated"`
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
    RiskyWords     bool `json:"riskyWords"`
    Results        []Results `json:"results"`
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

func Iso6391(iso6392 string) string {
    v, ok := reversoSupportedLangs[iso6392]
    if !ok {
        return ""
    }
    return v
}

func Iso6392(iso6391 string) string {
    for k, v := range reversoSupportedLangs {
        if v == iso6391 {
            return k
        }
    }
    return ""
}