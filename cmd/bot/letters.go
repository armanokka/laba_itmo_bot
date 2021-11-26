package bot

type Language struct {
	Name string
	Emoji string
}

var lettersSlice = []string{
	"A",
	"B",
	"C",
	"D",
	"E",
	"F",
	"G",
	"H",
	"I",
	"J",
	"K",
	"L",
	"M",
	"N",
	"O",
	"P",
	"R",
	"S",
	"T",
	"U",
	"V",
	"W",
	"X",
	"Y",
}

var lettersLangs = map[string][]Language{
	"D": {
		Language{
			Name:  "Danish",
			Emoji: "",
		},
		Language{
			Name:  "Dutch",
			Emoji: "🇳🇱",
		},
	},
	"F": {
		Language{
			Name:  "Finnish",
			Emoji: "🇫🇮",
		},
		Language{
			Name:  "French",
			Emoji: "🇫🇷",
		},
		Language{
			Name:  "Fula",
			Emoji: "",
		},
	},
	"H": {
		Language{
			Name:  "Haitian",
			Emoji: "🇭🇹",
		},
		Language{
			Name:  "Hausa",
			Emoji: "",
		},
		Language{
			Name:  "Hebrew",
			Emoji: "",
		},
		Language{
			Name:  "Herero",
			Emoji: "",
		},
		Language{
			Name:  "Hindi",
			Emoji: "",
		},
		Language{
			Name:  "Hiri Motu",
			Emoji: "",
		},
		Language{
			Name:  "Hungarian",
			Emoji: "🇭🇺",
		},
	},
	"O": {
		Language{
			Name:  "Ojibwe",
			Emoji: "",
		},
		Language{
			Name:  "Old Church Slavonic",
			Emoji: "🇨🇺",
		},
		Language{
			Name:  "Oriya",
			Emoji: "",
		},
	},
	"X": {
		Language{
			Name:  "Xhosa",
			Emoji: "",
		},
	},
	"Y": {
		Language{
			Name:  "Yiddish",
			Emoji: "",
		},
		Language{
			Name:  "Yoruba",
			Emoji: "",
		},
	},
	"W": {
		Language{
			Name:  "Walloon",
			Emoji: "",
		},
		Language{
			Name:  "Welsh",
			Emoji: "🇨🇾",
		},
		Language{
			Name:  "Western Frisian",
			Emoji: "",
		},
	},
	"A": {
		Language{
			Name:  "Afrikaans",
			Emoji: "🇦🇫",
		},
		Language{
			Name:  "Amharic",
			Emoji: "🇦🇲",
		},
		Language{
			Name:  "Arabic",
			Emoji: "🇦🇷",
		},
		Language{
			Name:  "Aragonese",
			Emoji: "🇦🇳",
		},
		Language{
			Name:  "Armenian",
			Emoji: "",
		},
		Language{
			Name:  "Avaric",
			Emoji: "",
		},
		Language{
			Name:  "Avestan",
			Emoji: "🇦🇪",
		},
		Language{
			Name:  "Azerbaijani",
			Emoji: "🇦🇿",
		},
	},
	"C": {
		Language{
			Name:  "Catalan",
			Emoji: "🇨🇦",
		},
		Language{
			Name:  "Chamorro",
			Emoji: "🇨🇭",
		},
		Language{
			Name:  "Chechen",
			Emoji: "",
		},
		Language{
			Name:  "Chichewa",
			Emoji: "",
		},
		Language{
			Name:  "Chinese",
			Emoji: "",
		},
		Language{
			Name:  "Chuvash",
			Emoji: "🇨🇻",
		},
		Language{
			Name:  "Cornish",
			Emoji: "🇰🇼",
		},
		Language{
			Name:  "Corsican",
			Emoji: "🇨🇴",
		},
		Language{
			Name:  "Cree",
			Emoji: "🇨🇷",
		},
		Language{
			Name:  "Croatian",
			Emoji: "🇭🇷",
		},
		Language{
			Name:  "Czech",
			Emoji: "",
		},
	},
	"I": {
		Language{
			Name:  "Icelandic",
			Emoji: "🇮🇸",
		},
		Language{
			Name:  "Ido",
			Emoji: "🇮🇴",
		},
		Language{
			Name:  "Igbo",
			Emoji: "",
		},
		Language{
			Name:  "Indonesian",
			Emoji: "🇮🇩",
		},
		Language{
			Name:  "Irish",
			Emoji: "🇬🇦",
		},
		Language{
			Name:  "Italian",
			Emoji: "🇮🇹",
		},
	},
	"L": {
		Language{
			Name:  "Lao",
			Emoji: "",
		},
		Language{
			Name:  "Latin",
			Emoji: "🇱🇦",
		},
		Language{
			Name:  "Latin",
			Emoji: "🇱🇦",
		},
		Language{
			Name:  "Latvian",
			Emoji: "🇱🇻",
		},
		Language{
			Name:  "Limburgish",
			Emoji: "🇱🇮",
		},
		Language{
			Name:  "Lithuanian",
			Emoji: "🇱🇹",
		},
		Language{
			Name:  "Luba-Katanga",
			Emoji: "🇱🇺",
		},
		Language{
			Name:  "Luxembourgish",
			Emoji: "🇱🇧",
		},
	},
	"P": {
		Language{
			Name:  "Panjabi",
			Emoji: "🇵🇦",
		},
		Language{
			Name:  "Pashto",
			Emoji: "🇵🇸",
		},
		Language{
			Name:  "Persian",
			Emoji: "",
		},
		Language{
			Name:  "Polish",
			Emoji: "🇵🇱",
		},
		Language{
			Name:  "Portuguese",
			Emoji: "🇵🇹",
		},
		Language{
			Name:  "Pāli",
			Emoji: "",
		},
	},
	"R": {
		Language{
			Name:  "Romanian",
			Emoji: "🇷🇴",
		},
		Language{
			Name:  "Russian",
			Emoji: "🇷🇺",
		},
	},
	"U": {
		Language{
			Name:  "Ukrainian",
			Emoji: "",
		},
		Language{
			Name:  "Urdu",
			Emoji: "",
		},
		Language{
			Name:  "Uyghur",
			Emoji: "🇺🇬",
		},
		Language{
			Name:  "Uzbek",
			Emoji: "🇺🇿",
		},
	},
	"K": {
		Language{
			Name:  "Kannada",
			Emoji: "🇰🇳",
		},
		Language{
			Name:  "Kanuri",
			Emoji: "🇰🇷",
		},
		Language{
			Name:  "Kazakh",
			Emoji: "",
		},
		Language{
			Name:  "Khmer",
			Emoji: "🇰🇲",
		},
		Language{
			Name:  "Kikuyu",
			Emoji: "🇰🇮",
		},
		Language{
			Name:  "Kinyarwanda",
			Emoji: "🇷🇼",
		},
		Language{
			Name:  "Komi",
			Emoji: "",
		},
		Language{
			Name:  "Kongo",
			Emoji: "🇰🇬",
		},
		Language{
			Name:  "Korean",
			Emoji: "",
		},
		Language{
			Name:  "Kurdish",
			Emoji: "",
		},
		Language{
			Name:  "Kwanyama",
			Emoji: "",
		},
		Language{
			Name:  "Kyrgyz",
			Emoji: "🇰🇾",
		},
	},
	"S": {
		Language{
			Name:  "Samoan",
			Emoji: "🇸🇲",
		},
		Language{
			Name:  "Sardinian",
			Emoji: "🇸🇨",
		},
		Language{
			Name:  "Scottish Gaelic",
			Emoji: "🇬🇩",
		},
		Language{
			Name:  "Serbian",
			Emoji: "🇸🇷",
		},
		Language{
			Name:  "Shona",
			Emoji: "🇸🇳",
		},
		Language{
			Name:  "Sindhi",
			Emoji: "🇸🇩",
		},
		Language{
			Name:  "Sinhala",
			Emoji: "🇸🇮",
		},
		Language{
			Name:  "Slovak",
			Emoji: "🇸🇰",
		},
		Language{
			Name:  "Slovene",
			Emoji: "🇸🇱",
		},
		Language{
			Name:  "Somali",
			Emoji: "🇸🇴",
		},
		Language{
			Name:  "Southern Ndebele",
			Emoji: "🇳🇷",
		},
		Language{
			Name:  "Southern Sotho",
			Emoji: "🇸🇹",
		},
		Language{
			Name:  "Spanish",
			Emoji: "🇪🇸",
		},
		Language{
			Name:  "Sundanese",
			Emoji: "",
		},
		Language{
			Name:  "Swahili",
			Emoji: "",
		},
		Language{
			Name:  "Swedish",
			Emoji: "🇸🇻",
		},
	},
	"V": {
		Language{
			Name:  "Vietnamese",
			Emoji: "🇻🇮",
		},
	},
	"B": {
		Language{
			Name:  "Basque",
			Emoji: "",
		},
		Language{
			Name:  "Belarusian",
			Emoji: "🇧🇪",
		},
		Language{
			Name:  "Bengali",
			Emoji: "🇧🇳",
		},
		Language{
			Name:  "Bosnian",
			Emoji: "🇧🇸",
		},
		Language{
			Name:  "Bulgarian",
			Emoji: "🇧🇬",
		},
		Language{
			Name:  "Burmese",
			Emoji: "🇲🇾",
		},
	},
	"E": {
		Language{
			Name:  "English",
			Emoji: "🇬🇧",
		},
		Language{
			Name:  "Esperanto",
			Emoji: "",
		},
		Language{
			Name:  "Estonian",
			Emoji: "🇪🇹",
		},
	},
	"G": {
		Language{
			Name:  "Galician",
			Emoji: "🇬🇱",
		},
		Language{
			Name:  "Georgian",
			Emoji: "",
		},
		Language{
			Name:  "German",
			Emoji: "🇩🇪",
		},
		Language{
			Name:  "Greek",
			Emoji: "",
		},
		Language{
			Name:  "Gujarati",
			Emoji: "🇬🇺",
		},
	},
	"J": {
		Language{
			Name:  "Japanese",
			Emoji: "",
		},
		Language{
			Name:  "Javanese",
			Emoji: "",
		},
	},
	"M": {
		Language{
			Name:  "Macedonian",
			Emoji: "🇲🇰",
		},
		Language{
			Name:  "Malagasy",
			Emoji: "🇲🇬",
		},
		Language{
			Name:  "Malay",
			Emoji: "🇲🇸",
		},
		Language{
			Name:  "Malayalam",
			Emoji: "",
		},
		Language{
			Name:  "Maltese",
			Emoji: "🇲🇹",
		},
		Language{
			Name:  "Marathi",
			Emoji: "🇲🇷",
		},
		Language{
			Name:  "Marshallese",
			Emoji: "🇲🇭",
		},
		Language{
			Name:  "Mongolian",
			Emoji: "🇲🇳",
		},
		Language{
			Name:  "Māori",
			Emoji: "",
		},
	},
	"N": {
		Language{
			Name:  "Navajo",
			Emoji: "",
		},
		Language{
			Name:  "Ndonga",
			Emoji: "🇳🇬",
		},
		Language{
			Name:  "Nepali",
			Emoji: "🇳🇪",
		},
		Language{
			Name:  "Northern Ndebele",
			Emoji: "",
		},
		Language{
			Name:  "Northern Sami",
			Emoji: "🇸🇪",
		},
		Language{
			Name:  "Norwegian",
			Emoji: "🇳🇴",
		},
		Language{
			Name:  "Norwegian Bokmål",
			Emoji: "",
		},
		Language{
			Name:  "Norwegian Nynorsk",
			Emoji: "",
		},
		Language{
			Name:  "Nuosu",
			Emoji: "",
		},
	},
	"T": {
		Language{
			Name:  "Tagalog",
			Emoji: "🇹🇱",
		},
		Language{
			Name:  "Tahitian",
			Emoji: "",
		},
		Language{
			Name:  "Tajik",
			Emoji: "🇹🇬",
		},
		Language{
			Name:  "Tamil",
			Emoji: "",
		},
		Language{
			Name:  "Tatar",
			Emoji: "🇹🇹",
		},
		Language{
			Name:  "Telugu",
			Emoji: "",
		},
		Language{
			Name:  "Thai",
			Emoji: "🇹🇭",
		},
		Language{
			Name:  "Turkish",
			Emoji: "🇹🇷",
		},
		Language{
			Name:  "Turkmen",
			Emoji: "🇹🇰",
		},
	},
}
