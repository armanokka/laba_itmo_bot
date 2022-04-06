package lingvo

type Dictionary struct {
	DictionaryName string `json:"dictionaryName"`
	SourceLanguage int    `json:"sourceLanguage"`
	TargetLanguage int    `json:"targetLanguage"`
	Heading        string `json:"heading"`
	ComplexHeading string `json:"complexHeading"`
	Transcription  string `json:"transcription"`
	SoundFileName  string `json:"soundFileName"`
	Translations   string `json:"translations"`
	Examples       string `json:"examples"`
	Comment        string `json:"comment"`
	PartOfSpeech   string `json:"partOfSpeech"`
}

var LingvoDictionaryLangs = []string{"en", "hu", "el", "da", "es", "it", "kk", "zh", "la", "de", "nl", "no", "pl", "pt", "ru", "tt", "tr", "uk", "fi", "fr"}

var Lingvo = map[string]int{
	"zh": 1028,
	"da": 1030,
	"nl": 1043,
	"en": 1033,
	"fi": 1035,
	"fr": 1036,
	"de": 32775,
	"el": 1032,
	"hu": 1038,
	"it": 1040,
	"kk": 1087,
	"la": 1142,
	"no": 1044,
	"pl": 1045,
	"pt": 2070,
	"ru": 1049,
	"es": 1034,
	"tt": 1092,
	"tr": 1055,
	"uk": 1058,
}

type TutorCard struct {
	DictionaryName string `json:"dictionaryName"`
	SourceLanguage int    `json:"sourceLanguage"`
	TargetLanguage int    `json:"targetLanguage"`
	Heading        string `json:"heading"`
	ComplexHeading string `json:"complexHeading"`
	Transcription  string `json:"transcription"`
	SoundFileName  string `json:"soundFileName"`
	Translations   string `json:"translations"`
	Examples       string `json:"examples"`
	Comment        string `json:"comment"`
	PartOfSpeech   string `json:"partOfSpeech"`
}
