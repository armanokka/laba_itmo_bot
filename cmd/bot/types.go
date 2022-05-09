package bot

type SuperTranslation struct {
	TranslatedText string
	Examples       map[string]string
	Translations   map[string][]string // reverso
	Dictionary     []string
	Paraphrase     []string
}

type Translation struct { // for internal using
	TranslatedText string
	FromCache      bool
	Keyboard       Keyboard
}

type Keyboard struct {
	Dictionary          []string            `json:"dictionary"`
	Paraphrase          []string            `json:"paraphrase"`
	Examples            map[string]string   `json:"examples"`
	ReverseTranslations map[string][]string `json:"reverse_translations"`
}

type RecordTranslationCache struct {
	Key         string `json:"_key,omitempty"`
	From        string `json:"from"`
	To          string `json:"to"`
	Text        string `json:"text"`
	Translation string `json:"translation"`
	//CreatedAt   int64  `json:"createdAt"`
	Keyboard
}
