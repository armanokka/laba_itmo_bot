package main

var langs = map[string]Lang{
	"zh": {
		Name:  "Chinese",
		Emoji: "🇨🇳",
	},
	"hi": {
		Name:  "Hindi",
		Emoji: "",
	},
	"ms": {
		Name:  "Malay",
		Emoji: "🇲🇸",
	},
	"ae": {
		Name:  "Avestan",
		Emoji: "🇦🇪",
	},
	"bs": {
		Name:  "Bosnian",
		Emoji: "🇧🇸",
	},
	"sd": {
		Name:  "Sindhi",
		Emoji: "🇸🇩",
	},
	"sg": {
		Name:  "Sango",
		Emoji: "🇸🇬",
	},
	"ab": {
		Name:  "Abkhaz",
		Emoji: "",
	},
	"ig": {
		Name:  "Igbo",
		Emoji: "",
	},
	"lg": {
		Name:  "Ganda",
		Emoji: "",
	},
	"or": {
		Name:  "Oriya",
		Emoji: "",
	},
	"os": {
		Name:  "Ossetian",
		Emoji: "",
	},
	"ro": {
		Name:  "Romanian",
		Emoji: "🇷🇴",
	},
	"ti": {
		Name:  "Tigrinya",
		Emoji: "",
	},
	"ug": {
		Name:  "Uyghur",
		Emoji: "🇺🇬",
	},
	"cy": {
		Name:  "Welsh",
		Emoji: "🇨🇾",
	},
	"vo": {
		Name:  "Volapük",
		Emoji: "",
	},
	"wo": {
		Name:  "Wolof",
		Emoji: "",
	},
	"ru": {
		Name:  "Russian",
		Emoji: "🇷🇺",
	},
	"ko": {
		Name:  "Korean",
		Emoji: "🇰🇵🇰🇷",
	},
	"mg": {
		Name:  "Malagasy",
		Emoji: "🇲🇬",
	},
	"mn": {
		Name:  "Mongolian",
		Emoji: "🇲🇳",
	},
	"oc": {
		Name:  "Occitan",
		Emoji: "",
	},
	"sr": {
		Name:  "Serbian",
		Emoji: "🇸🇷",
	},
	"bn": {
		Name:  "Bengali",
		Emoji: "🇧🇩",
	},
	"kr": {
		Name:  "Kanuri",
		Emoji: "🇰🇷",
	},
	"kk": {
		Name:  "Kazakh",
		Emoji: "",
	},
	"rw": {
		Name:  "Kinyarwanda",
		Emoji: "🇷🇼",
	},
	"tk": {
		Name:  "Turkmen",
		Emoji: "🇹🇰",
	},
	"ay": {
		Name:  "Aymara",
		Emoji: "",
	},
	"pl": {
		Name:  "Polish",
		Emoji: "🇵🇱",
	},
	"qu": {
		Name:  "Quechua",
		Emoji: "",
	},
	"ty": {
		Name:  "Tahitian",
		Emoji: "",
	},
	"vi": {
		Name:  "Vietnamese",
		Emoji: "🇻🇮",
	},
	"my": {
		Name:  "Burmese",
		Emoji: "🇲🇾",
	},
	"fi": {
		Name:  "Finnish",
		Emoji: "🇫🇮",
	},
	"ku": {
		Name:  "Kurdish",
		Emoji: "",
	},
	"ii": {
		Name:  "Nuosu",
		Emoji: "",
	},
	"sw": {
		Name:  "Swahili",
		Emoji: "",
	},
	"en": {
		Name:  "English",
		Emoji: "🇬🇧",
	},
	"tr": {
		Name:  "Turkish",
		Emoji: "🇹🇷",
	},
	"am": {
		Name:  "Amharic",
		Emoji: "🇦🇲",
	},
	"br": {
		Name:  "Breton",
		Emoji: "🇧🇷",
	},
	"et": {
		Name:  "Estonian",
		Emoji: "🇪🇹",
	},
	"li": {
		Name:  "Limburgish",
		Emoji: "🇱🇮",
	},
	"nl": {
		Name:  "Dutch",
		Emoji: "🇳🇱",
	},
	"fj": {
		Name:  "Fijian",
		Emoji: "🇫🇯",
	},
	"gn": {
		Name:  "Guaraní",
		Emoji: "🇬🇳",
	},
	"he": {
		Name:  "Hebrew",
		Emoji: "",
	},
	"mt": {
		Name:  "Maltese",
		Emoji: "🇲🇹",
	},
	"rn": {
		Name:  "Kirundi",
		Emoji: "",
	},
	"tl": {
		Name:  "Tagalog",
		Emoji: "🇹🇱",
	},
	"da": {
		Name:  "Danish",
		Emoji: "",
	},
	"ha": {
		Name:  "Hausa",
		Emoji: "",
	},
	"iu": {
		Name:  "Inuktitut",
		Emoji: "",
	},
	"ks": {
		Name:  "Kashmiri",
		Emoji: "",
	},
	"ki": {
		Name:  "Kikuyu",
		Emoji: "🇰🇮",
	},
	"nn": {
		Name:  "Norwegian Nynorsk",
		Emoji: "",
	},
	"ee": {
		Name:  "Ewe",
		Emoji: "🇪🇪",
	},
	"el": {
		Name:  "Greek",
		Emoji: "",
	},
	"kn": {
		Name:  "Kannada",
		Emoji: "🇰🇳",
	},
	"lt": {
		Name:  "Lithuanian",
		Emoji: "🇱🇹",
	},
	"gv": {
		Name:  "Manx",
		Emoji: "",
	},
	"na": {
		Name:  "Nauru",
		Emoji: "🇳🇦",
	},
	"ce": {
		Name:  "Chechen",
		Emoji: "",
	},
	"ny": {
		Name:  "Chichewa",
		Emoji: "",
	},
	"co": {
		Name:  "Corsican",
		Emoji: "🇨🇴",
	},
	"rm": {
		Name:  "Romansh",
		Emoji: "",
	},
	"sl": {
		Name:  "Slovene",
		Emoji: "🇸🇱",
	},
	"id": {
		Name:  "Indonesian",
		Emoji: "🇮🇩",
	},
	"xh": {
		Name:  "Xhosa",
		Emoji: "",
	},
	"hy": {
		Name:  "Armenian",
		Emoji: "",
	},
	"cr": {
		Name:  "Cree",
		Emoji: "🇨🇷",
	},
	"om": {
		Name:  "Oromo",
		Emoji: "🇴🇲",
	},
	"ss": {
		Name:  "Swati",
		Emoji: "🇸🇸",
	},
	"tn": {
		Name:  "Tswana",
		Emoji: "🇹🇳",
	},
	"bg": {
		Name:  "Bulgarian",
		Emoji: "🇧🇬",
	},
	"nb": {
		Name:  "Norwegian Bokmål",
		Emoji: "",
	},
	"gd": {
		Name:  "Scottish Gaelic",
		Emoji: "🇬🇩",
	},
	"sn": {
		Name:  "Shona",
		Emoji: "🇸🇳",
	},
	"oj": {
		Name:  "Ojibwe",
		Emoji: "",
	},
	"pi": {
		Name:  "Pāli",
		Emoji: "",
	},
	"de": {
		Name:  "German",
		Emoji: "🇩🇪",
	},
	"it": {
		Name:  "Italian",
		Emoji: "🇮🇹",
	},
	"ia": {
		Name:  "Interlingua",
		Emoji: "",
	},
	"nv": {
		Name:  "Navajo",
		Emoji: "",
	},
	"ne": {
		Name:  "Nepali",
		Emoji: "🇳🇪",
	},
	"nr": {
		Name:  "Southern Ndebele",
		Emoji: "🇳🇷",
	},
	"sc": {
		Name:  "Sardinian",
		Emoji: "🇸🇨",
	},
	"th": {
		Name:  "Thai",
		Emoji: "🇹🇭",
	},
	"fr": {
		Name:  "French",
		Emoji: "🇫🇷",
	},
	"ga": {
		Name:  "Irish",
		Emoji: "🇬🇦",
	},
	"su": {
		Name:  "Sundanese",
		Emoji: "",
	},
	"ts": {
		Name:  "Tsonga",
		Emoji: "",
	},
	"lu": {
		Name:  "Luba-Katanga",
		Emoji: "🇱🇺",
	},
	"ja": {
		Name:  "Japanese",
		Emoji: "🇯🇵",
	},
	"es": {
		Name:  "Spanish",
		Emoji: "🇪🇸",
	},
	"ur": {
		Name:  "Urdu",
		Emoji: "",
	},
	"eu": {
		Name:  "Basque",
		Emoji: "",
	},
	"ca": {
		Name:  "Catalan",
		Emoji: "🇨🇦",
	},
	"is": {
		Name:  "Icelandic",
		Emoji: "🇮🇸",
	},
	"uk": {
		Name:  "Ukrainian",
		Emoji: "🇺🇦",
	},
	"fo": {
		Name:  "Faroese",
		Emoji: "🇫🇴",
	},
	"mh": {
		Name:  "Marshallese",
		Emoji: "🇲🇭",
	},
	"sa": {
		Name:  "Sanskrit",
		Emoji: "🇸🇦",
	},
	"sm": {
		Name:  "Samoan",
		Emoji: "🇸🇲",
	},
	"io": {
		Name:  "Ido",
		Emoji: "🇮🇴",
	},
	"sv": {
		Name:  "Swedish",
		Emoji: "🇸🇻",
	},
	"la": {
		Name:  "Latin",
		Emoji: "",
	},
	"pt": {
		Name:  "Portuguese",
		Emoji: "🇵🇹",
	},
	"jv": {
		Name:  "Javanese",
		Emoji: "",
	},
	"av": {
		Name:  "Avaric",
		Emoji: "",
	},
	"ho": {
		Name:  "Hiri Motu",
		Emoji: "",
	},
	"hu": {
		Name:  "Hungarian",
		Emoji: "🇭🇺",
	},
	"af": {
		Name:  "Afrikaans",
		Emoji: "",
	},
	"ba": {
		Name:  "Bashkir",
		Emoji: "🇧🇦",
	},
	"ik": {
		Name:  "Inupiaq",
		Emoji: "",
	},
	"pa": {
		Name:  "Panjabi",
		Emoji: "🇵🇦",
	},
	"st": {
		Name:  "Southern Sotho",
		Emoji: "🇸🇹",
	},
	"ve": {
		Name:  "Venda",
		Emoji: "🇻🇪",
	},
	"uz": {
		Name:  "Uzbek",
		Emoji: "🇺🇿",
	},
	"te": {
		Name:  "Telugu",
		Emoji: "",
	},
	"ta": {
		Name:  "Tamil",
		Emoji: "",
	},
	"ff": {
		Name:  "Fula",
		Emoji: "",
	},
	"ie": {
		Name:  "Interlingue",
		Emoji: "🇮🇪",
	},
	"kv": {
		Name:  "Komi",
		Emoji: "",
	},
	"km": {
		Name:  "Khmer",
		Emoji: "🇰🇲",
	},
	"lo": {
		Name:  "Lao",
		Emoji: "",
	},
	"mi": {
		Name:  "Māori",
		Emoji: "",
	},
	"yo": {
		Name:  "Yoruba",
		Emoji: "",
	},
	"gu": {
		Name:  "Gujarati",
		Emoji: "🇬🇺",
	},
	"an": {
		Name:  "Aragonese",
		Emoji: "🇦🇳",
	},
	"cv": {
		Name:  "Chuvash",
		Emoji: "🇨🇻",
	},
	"ht": {
		Name:  "Haitian",
		Emoji: "🇭🇹",
	},
	"fy": {
		Name:  "Western Frisian",
		Emoji: "",
	},
	"lb": {
		Name:  "Luxembourgish",
		Emoji: "🇱🇧",
	},
	"lv": {
		Name:  "Latvian",
		Emoji: "🇱🇻",
	},
	"si": {
		Name:  "Sinhala",
		Emoji: "🇸🇮",
	},
	"so": {
		Name:  "Somali",
		Emoji: "🇸🇴",
	},
	"se": {
		Name:  "Northern Sami",
		Emoji: "🇸🇪",
	},
	"to": {
		Name:  "Tonga",
		Emoji: "🇹🇴",
	},
	"cs": {
		Name:  "Czech",
		Emoji: "",
	},
	"kl": {
		Name:  "Kalaallisut",
		Emoji: "",
	},
	"ln": {
		Name:  "Lingala",
		Emoji: "",
	},
	"mk": {
		Name:  "Macedonian",
		Emoji: "🇲🇰",
	},
	"ng": {
		Name:  "Ndonga",
		Emoji: "🇳🇬",
	},
	"cu": {
		Name:  "Old Church Slavonic",
		Emoji: "🇨🇺",
	},
	"fa": {
		Name:  "Persian",
		Emoji: "",
	},
	"sq": {
		Name:  "Albanian",
		Emoji: "",
	},
	"ka": {
		Name:  "Georgian",
		Emoji: "",
	},
	"gl": {
		Name:  "Galician",
		Emoji: "🇬🇱",
	},
	"hz": {
		Name:  "Herero",
		Emoji: "",
	},
	"kg": {
		Name:  "Kongo",
		Emoji: "🇰🇬",
	},
	"as": {
		Name:  "Assamese",
		Emoji: "🇦🇸",
	},
	"ch": {
		Name:  "Chamorro",
		Emoji: "🇨🇭",
	},
	"kw": {
		Name:  "Cornish",
		Emoji: "🇰🇼",
	},
	"ps": {
		Name:  "Pashto",
		Emoji: "🇵🇸",
	},
	"tg": {
		Name:  "Tajik",
		Emoji: "🇹🇬",
	},
	"be": {
		Name:  "Belarusian",
		Emoji: "🇧🇪",
	},
	"bh": {
		Name:  "Bihari",
		Emoji: "🇧🇭",
	},
	"dv": {
		Name:  "Divehi",
		Emoji: "",
	},
	"sk": {
		Name:  "Slovak",
		Emoji: "🇸🇰",
	},
	"tt": {
		Name:  "Tatar",
		Emoji: "🇹🇹",
	},
	"bo": {
		Name:  "Tibetan Standard",
		Emoji: "🇧🇴",
	},
	"za": {
		Name:  "Zhuang",
		Emoji: "🇿🇦",
	},
	"mr": {
		Name:  "Marathi",
		Emoji: "🇲🇷",
	},
	"bi": {
		Name:  "Bislama",
		Emoji: "🇧🇮",
	},
	"ky": {
		Name:  "Kyrgyz",
		Emoji: "🇰🇬",
	},
	"kj": {
		Name:  "Kwanyama",
		Emoji: "",
	},
	"nd": {
		Name:  "Northern Ndebele",
		Emoji: "",
	},
	"no": {
		Name:  "Norwegian",
		Emoji: "🇳🇴",
	},
	"ar": {
		Name:  "Arabic",
		Emoji: "🇪🇬",
	},
	"eo": {
		Name:  "Esperanto",
		Emoji: "",
	},
	"wa": {
		Name:  "Walloon",
		Emoji: "",
	},
	"az": {
		Name:  "Azerbaijani",
		Emoji: "🇦🇿",
	},
	"hr": {
		Name:  "Croatian",
		Emoji: "🇭🇷",
	},
	"ml": {
		Name:  "Malayalam",
		Emoji: "",
	},
	"yi": {
		Name:  "Yiddish",
		Emoji: "",
	},
}

//var popular = []string{
//	"ru",
//	"en",
//	"es",
//	"ar",
//	"uz",
//	"pt",
//	"hi",
//	"fr",
//	"la",
//	"ja",
//	"ko",
//	"tr",
//}

var codes = []string{
	"af",
	"am",
	"ar",
	"an",
	"hy",
	"av",
	"ae",
	"az",
	"eu",
	"be",
	"bn",
	"bs",
	"bg",
	"my",
	"ca",
	"ch",
	"ce",
	"ny",
	"zh",
	"cv",
	"kw",
	"co",
	"cr",
	"hr",
	"cs",
	"da",
	"nl",
	"en",
	"eo",
	"et",
	"fi",
	"fr",
	"ff",
	"gl",
	"ka",
	"de",
	"el",
	"gu",
	"ht",
	"ha",
	"he",
	"hz",
	"hi",
	"ho",
	"hu",
	"is",
	"io",
	"ig",
	"id",
	"ga",
	"it",
	"ja",
	"jv",
	"kn",
	"kr",
	"kk",
	"km",
	"ki",
	"rw",
	"kv",
	"kg",
	"ko",
	"ku",
	"kj",
	"ky",
	"lo",
	"la",
	"la",
	"lv",
	"li",
	"lt",
	"lu",
	"lb",
	"mk",
	"mg",
	"ms",
	"ml",
	"mt",
	"mr",
	"mh",
	"mn",
	"mi",
	"nv",
	"ng",
	"ne",
	"nd",
	"se",
	"no",
	"nb",
	"nn",
	"ii",
	"oj",
	"cu",
	"or",
	"pa",
	"ps",
	"fa",
	"pl",
	"pt",
	"pi",
	"ro",
	"ru",
	"sm",
	"sc",
	"gd",
	"sr",
	"sn",
	"sd",
	"si",
	"sk",
	"sl",
	"so",
	"nr",
	"st",
	"es",
	"su",
	"sw",
	"sv",
	"tl",
	"ty",
	"tg",
	"ta",
	"tt",
	"te",
	"th",
	"tr",
	"tk",
	"uk",
	"ur",
	"ug",
	"uz",
	"vi",
	"wa",
	"cy",
	"fy",
	"xh",
	"yi",
	"yo",
}