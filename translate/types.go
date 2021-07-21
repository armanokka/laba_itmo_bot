package translate

import "fmt"

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

type Player struct {
    Lang string // Language code into we have to translate (ISO-6391)
}

type TTSError struct {
    Code int
    Description string
}

func (e TTSError) Error() string {
    return fmt.Sprintf("TTSError [%d]: %s", e.Code, e.Description)
}