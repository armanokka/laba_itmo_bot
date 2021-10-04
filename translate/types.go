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

type reversoRequestTranslate struct {
    Format  string `json:"format"`
    From    string `json:"from"`
    Input   string `json:"input"`
    To      string `json:"to"`
    Options reversoRequestTranslateOptions `json:"options"`
}

type reversoRequestTranslateOptions struct {
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
    languageDetection  `json:"languageDetection"`
    contextResults  `json:"contextResults"`
    Truncated bool `json:"truncated"`
    TimeTaken int  `json:"timeTaken"`
}

type languageDetection struct {
    DetectedLanguage                string `json:"detectedLanguage"`
    IsDirectionChanged              bool   `json:"isDirectionChanged"`
    OriginalDirection               string `json:"originalDirection"`
    OriginalDirectionContextMatches int    `json:"originalDirectionContextMatches"`
    ChangedDirectionContextMatches  int    `json:"changedDirectionContextMatches"`
    TimeTaken                       int    `json:"timeTaken"`
}

type contextResults struct {
    RudeWords      bool `json:"rudeWords"`
    Colloquialisms bool `json:"colloquialisms"`
    RiskyWords     bool `json:"riskyWords"`
    Results        []results `json:"results"`
    TotalContextCallsMade int `json:"totalContextCallsMade"`
    TimeTakenContext      int `json:"timeTakenContext"`
}

type results struct {
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