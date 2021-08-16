package main

var langs = []Lang{
    {
        Code:  "en",
        Name:  "English",
        Emoji: "",
    },
    {
        Code:  "ru",
        Name:  "Russian",
        Emoji: "🇷🇺",
    },
    {
        Code:  "la",
        Name:  "Latin",
        Emoji: "🇱🇦",
    },
    {
        Code:  "ja",
        Name:  "Japanese",
        Emoji: "",
    },
    {
        Code:  "ar",
        Name:  "Arabic",
        Emoji: "🇦🇷",
    },
    {
        Code:  "fr",
        Name:  "French",
        Emoji: "🇫🇷",
    },
    {
        Code:  "de",
        Name:  "German",
        Emoji: "🇩🇪",
    },
    {
        Code:  "af",
        Name:  "Afrikaans",
        Emoji: "🇦🇫",
    },
    {
        Code:  "uk",
        Name:  "Ukrainian",
        Emoji: "",
    },
    {
        Code:  "uz",
        Name:  "Uzbek",
        Emoji: "🇺🇿",
    },
    {
        Code:  "es",
        Name:  "Spanish",
        Emoji: "🇪🇸",
    },
    {
        Code:  "ko",
        Name:  "Korean",
        Emoji: "",
    },
    {
        Code:  "zh",
        Name:  "Chinese",
        Emoji: "",
    },
    {
        Code:  "hi",
        Name:  "Hindi",
        Emoji: "",
    },
    {
        Code:  "bn",
        Name:  "Bengali",
        Emoji: "🇧🇳",
    },
    {
        Code:  "pt",
        Name:  "Portuguese",
        Emoji: "🇵🇹",
    },
    {
        Code:  "mr",
        Name:  "Marathi",
        Emoji: "🇲🇷",
    },
    {
        Code:  "te",
        Name:  "Telugu",
        Emoji: "",
    },
    {
        Code:  "ms",
        Name:  "Malay",
        Emoji: "🇲🇸",
    },
    {
        Code:  "tr",
        Name:  "Turkish",
        Emoji: "🇹🇷",
    },
    {
        Code:  "vi",
        Name:  "Vietnamese",
        Emoji: "🇻🇮",
    },
    {
        Code:  "ta",
        Name:  "Tamil",
        Emoji: "",
    },
    {
        Code:  "ur",
        Name:  "Urdu",
        Emoji: "",
    },
    {
        Code:  "jv",
        Name:  "Javanese",
        Emoji: "",
    },
    {
        Code:  "it",
        Name:  "Italian",
        Emoji: "🇮🇹",
    },
    {
        Code:  "fa",
        Name:  "Persian",
        Emoji: "",
    },
    {
        Code:  "gu",
        Name:  "Gujarati",
        Emoji: "🇬🇺",
    },
    {
        Code:  "ab",
        Name:  "Abkhaz",
        Emoji: "",
    },
    {
        Code:  "aa",
        Name:  "Afar",
        Emoji: "",
    },
    {
        Code:  "ak",
        Name:  "Akan",
        Emoji: "",
    },
    {
        Code:  "sq",
        Name:  "Albanian",
        Emoji: "",
    },
    {
        Code:  "am",
        Name:  "Amharic",
        Emoji: "🇦🇲",
    },
    {
        Code:  "an",
        Name:  "Aragonese",
        Emoji: "🇦🇳",
    },
    {
        Code:  "hy",
        Name:  "Armenian",
        Emoji: "",
    },
    {
        Code:  "as",
        Name:  "Assamese",
        Emoji: "🇦🇸",
    },
    {
        Code:  "av",
        Name:  "Avaric",
        Emoji: "",
    },
    {
        Code:  "ae",
        Name:  "Avestan",
        Emoji: "🇦🇪",
    },
    {
        Code:  "ay",
        Name:  "Aymara",
        Emoji: "",
    },
    {
        Code:  "az",
        Name:  "Azerbaijani",
        Emoji: "🇦🇿",
    },
    {
        Code:  "bm",
        Name:  "Bambara",
        Emoji: "🇧🇲",
    },
    {
        Code:  "ba",
        Name:  "Bashkir",
        Emoji: "🇧🇦",
    },
    {
        Code:  "eu",
        Name:  "Basque",
        Emoji: "",
    },
    {
        Code:  "be",
        Name:  "Belarusian",
        Emoji: "🇧🇪",
    },
    {
        Code:  "bh",
        Name:  "Bihari",
        Emoji: "🇧🇭",
    },
    {
        Code:  "bi",
        Name:  "Bislama",
        Emoji: "🇧🇮",
    },
    {
        Code:  "bs",
        Name:  "Bosnian",
        Emoji: "🇧🇸",
    },
    {
        Code:  "br",
        Name:  "Breton",
        Emoji: "🇧🇷",
    },
    {
        Code:  "bg",
        Name:  "Bulgarian",
        Emoji: "🇧🇬",
    },
    {
        Code:  "my",
        Name:  "Burmese",
        Emoji: "🇲🇾",
    },
    {
        Code:  "ca",
        Name:  "Catalan",
        Emoji: "🇨🇦",
    },
    {
        Code:  "ch",
        Name:  "Chamorro",
        Emoji: "🇨🇭",
    },
    {
        Code:  "ce",
        Name:  "Chechen",
        Emoji: "",
    },
    {
        Code:  "ny",
        Name:  "Chichewa",
        Emoji: "",
    },
    {
        Code:  "cv",
        Name:  "Chuvash",
        Emoji: "🇨🇻",
    },
    {
        Code:  "kw",
        Name:  "Cornish",
        Emoji: "🇰🇼",
    },
    {
        Code:  "co",
        Name:  "Corsican",
        Emoji: "🇨🇴",
    },
    {
        Code:  "cr",
        Name:  "Cree",
        Emoji: "🇨🇷",
    },
    {
        Code:  "hr",
        Name:  "Croatian",
        Emoji: "🇭🇷",
    },
    {
        Code:  "cs",
        Name:  "Czech",
        Emoji: "",
    },
    {
        Code:  "da",
        Name:  "Danish",
        Emoji: "",
    },
    {
        Code:  "dv",
        Name:  "Divehi",
        Emoji: "",
    },
    {
        Code:  "nl",
        Name:  "Dutch",
        Emoji: "🇳🇱",
    },
    {
        Code:  "eo",
        Name:  "Esperanto",
        Emoji: "",
    },
    {
        Code:  "et",
        Name:  "Estonian",
        Emoji: "🇪🇹",
    },
    {
        Code:  "ee",
        Name:  "Ewe",
        Emoji: "🇪🇪",
    },
    {
        Code:  "fo",
        Name:  "Faroese",
        Emoji: "🇫🇴",
    },
    {
        Code:  "fj",
        Name:  "Fijian",
        Emoji: "🇫🇯",
    },
    {
        Code:  "fi",
        Name:  "Finnish",
        Emoji: "🇫🇮",
    },
    {
        Code:  "ff",
        Name:  "Fula",
        Emoji: "",
    },
    {
        Code:  "gl",
        Name:  "Galician",
        Emoji: "🇬🇱",
    },
    {
        Code:  "ka",
        Name:  "Georgian",
        Emoji: "",
    },
    {
        Code:  "el",
        Name:  "Greek",
        Emoji: "",
    },
    {
        Code:  "gn",
        Name:  "Guaraní",
        Emoji: "🇬🇳",
    },
    {
        Code:  "ht",
        Name:  "Haitian",
        Emoji: "🇭🇹",
    },
    {
        Code:  "ha",
        Name:  "Hausa",
        Emoji: "",
    },
    {
        Code:  "he",
        Name:  "Hebrew",
        Emoji: "",
    },
    {
        Code:  "hz",
        Name:  "Herero",
        Emoji: "",
    },
    {
        Code:  "ho",
        Name:  "Hiri Motu",
        Emoji: "",
    },
    {
        Code:  "hu",
        Name:  "Hungarian",
        Emoji: "🇭🇺",
    },
    {
        Code:  "ia",
        Name:  "Interlingua",
        Emoji: "",
    },
    {
        Code:  "id",
        Name:  "Indonesian",
        Emoji: "🇮🇩",
    },
    {
        Code:  "ie",
        Name:  "Interlingue",
        Emoji: "🇮🇪",
    },
    {
        Code:  "ga",
        Name:  "Irish",
        Emoji: "🇬🇦",
    },
    {
        Code:  "ig",
        Name:  "Igbo",
        Emoji: "",
    },
    {
        Code:  "ik",
        Name:  "Inupiaq",
        Emoji: "",
    },
    {
        Code:  "io",
        Name:  "Ido",
        Emoji: "🇮🇴",
    },
    {
        Code:  "is",
        Name:  "Icelandic",
        Emoji: "🇮🇸",
    },
    {
        Code:  "iu",
        Name:  "Inuktitut",
        Emoji: "",
    },
    {
        Code:  "kl",
        Name:  "Kalaallisut",
        Emoji: "",
    },
    {
        Code:  "kn",
        Name:  "Kannada",
        Emoji: "🇰🇳",
    },
    {
        Code:  "kr",
        Name:  "Kanuri",
        Emoji: "🇰🇷",
    },
    {
        Code:  "ks",
        Name:  "Kashmiri",
        Emoji: "",
    },
    {
        Code:  "kk",
        Name:  "Kazakh",
        Emoji: "",
    },
    {
        Code:  "km",
        Name:  "Khmer",
        Emoji: "🇰🇲",
    },
    {
        Code:  "ki",
        Name:  "Kikuyu",
        Emoji: "🇰🇮",
    },
    {
        Code:  "rw",
        Name:  "Kinyarwanda",
        Emoji: "🇷🇼",
    },
    {
        Code:  "ky",
        Name:  "Kyrgyz",
        Emoji: "🇰🇾",
    },
    {
        Code:  "kv",
        Name:  "Komi",
        Emoji: "",
    },
    {
        Code:  "kg",
        Name:  "Kongo",
        Emoji: "🇰🇬",
    },
    {
        Code:  "ku",
        Name:  "Kurdish",
        Emoji: "",
    },
    {
        Code:  "kj",
        Name:  "Kwanyama",
        Emoji: "",
    },
    {
        Code:  "la",
        Name:  "Latin",
        Emoji: "🇱🇦",
    },
    {
        Code:  "lb",
        Name:  "Luxembourgish",
        Emoji: "🇱🇧",
    },
    {
        Code:  "lg",
        Name:  "Ganda",
        Emoji: "",
    },
    {
        Code:  "li",
        Name:  "Limburgish",
        Emoji: "🇱🇮",
    },
    {
        Code:  "ln",
        Name:  "Lingala",
        Emoji: "",
    },
    {
        Code:  "lo",
        Name:  "Lao",
        Emoji: "",
    },
    {
        Code:  "lt",
        Name:  "Lithuanian",
        Emoji: "🇱🇹",
    },
    {
        Code:  "lu",
        Name:  "Luba-Katanga",
        Emoji: "🇱🇺",
    },
    {
        Code:  "lv",
        Name:  "Latvian",
        Emoji: "🇱🇻",
    },
    {
        Code:  "gv",
        Name:  "Manx",
        Emoji: "",
    },
    {
        Code:  "mk",
        Name:  "Macedonian",
        Emoji: "🇲🇰",
    },
    {
        Code:  "mg",
        Name:  "Malagasy",
        Emoji: "🇲🇬",
    },
    {
        Code:  "ml",
        Name:  "Malayalam",
        Emoji: "🇲🇱",
    },
    {
        Code:  "mt",
        Name:  "Maltese",
        Emoji: "🇲🇹",
    },
    {
        Code:  "mi",
        Name:  "Māori",
        Emoji: "",
    },
    {
        Code:  "mh",
        Name:  "Marshallese",
        Emoji: "🇲🇭",
    },
    {
        Code:  "mn",
        Name:  "Mongolian",
        Emoji: "🇲🇳",
    },
    {
        Code:  "na",
        Name:  "Nauru",
        Emoji: "🇳🇦",
    },
    {
        Code:  "nv",
        Name:  "Navajo",
        Emoji: "",
    },
    {
        Code:  "nb",
        Name:  "Norwegian Bokmål",
        Emoji: "",
    },
    {
        Code:  "nd",
        Name:  "Northern Ndebele",
        Emoji: "",
    },
    {
        Code:  "ne",
        Name:  "Nepali",
        Emoji: "🇳🇪",
    },
    {
        Code:  "ng",
        Name:  "Ndonga",
        Emoji: "🇳🇬",
    },
    {
        Code:  "nn",
        Name:  "Norwegian Nynorsk",
        Emoji: "",
    },
    {
        Code:  "no",
        Name:  "Norwegian",
        Emoji: "🇳🇴",
    },
    {
        Code:  "ii",
        Name:  "Nuosu",
        Emoji: "",
    },
    {
        Code:  "nr",
        Name:  "Southern Ndebele",
        Emoji: "🇳🇷",
    },
    {
        Code:  "oc",
        Name:  "Occitan",
        Emoji: "",
    },
    {
        Code:  "oj",
        Name:  "Ojibwe",
        Emoji: "",
    },
    {
        Code:  "cu",
        Name:  "Old Church Slavonic",
        Emoji: "🇨🇺",
    },
    {
        Code:  "om",
        Name:  "Oromo",
        Emoji: "🇴🇲",
    },
    {
        Code:  "or",
        Name:  "Oriya",
        Emoji: "",
    },
    {
        Code:  "os",
        Name:  "Ossetian",
        Emoji: "",
    },
    {
        Code:  "pa",
        Name:  "Panjabi",
        Emoji: "🇵🇦",
    },
    {
        Code:  "pi",
        Name:  "Pāli",
        Emoji: "",
    },
    {
        Code:  "pl",
        Name:  "Polish",
        Emoji: "🇵🇱",
    },
    {
        Code:  "ps",
        Name:  "Pashto",
        Emoji: "🇵🇸",
    },
    {
        Code:  "qu",
        Name:  "Quechua",
        Emoji: "",
    },
    {
        Code:  "rm",
        Name:  "Romansh",
        Emoji: "",
    },
    {
        Code:  "rn",
        Name:  "Kirundi",
        Emoji: "",
    },
    {
        Code:  "ro",
        Name:  "Romanian",
        Emoji: "🇷🇴",
    },
    {
        Code:  "sa",
        Name:  "Sanskrit",
        Emoji: "🇸🇦",
    },
    {
        Code:  "sc",
        Name:  "Sardinian",
        Emoji: "🇸🇨",
    },
    {
        Code:  "sd",
        Name:  "Sindhi",
        Emoji: "🇸🇩",
    },
    {
        Code:  "se",
        Name:  "Northern Sami",
        Emoji: "🇸🇪",
    },
    {
        Code:  "sm",
        Name:  "Samoan",
        Emoji: "🇸🇲",
    },
    {
        Code:  "sg",
        Name:  "Sango",
        Emoji: "🇸🇬",
    },
    {
        Code:  "sr",
        Name:  "Serbian",
        Emoji: "🇸🇷",
    },
    {
        Code:  "gd",
        Name:  "Scottish Gaelic",
        Emoji: "🇬🇩",
    },
    {
        Code:  "sn",
        Name:  "Shona",
        Emoji: "🇸🇳",
    },
    {
        Code:  "si",
        Name:  "Sinhala",
        Emoji: "🇸🇮",
    },
    {
        Code:  "sk",
        Name:  "Slovak",
        Emoji: "🇸🇰",
    },
    {
        Code:  "sl",
        Name:  "Slovene",
        Emoji: "🇸🇱",
    },
    {
        Code:  "so",
        Name:  "Somali",
        Emoji: "🇸🇴",
    },
    {
        Code:  "st",
        Name:  "Southern Sotho",
        Emoji: "🇸🇹",
    },
    {
        Code:  "su",
        Name:  "Sundanese",
        Emoji: "",
    },
    {
        Code:  "sw",
        Name:  "Swahili",
        Emoji: "",
    },
    {
        Code:  "ss",
        Name:  "Swati",
        Emoji: "🇸🇸",
    },
    {
        Code:  "sv",
        Name:  "Swedish",
        Emoji: "🇸🇻",
    },
    {
        Code:  "tg",
        Name:  "Tajik",
        Emoji: "🇹🇬",
    },
    {
        Code:  "th",
        Name:  "Thai",
        Emoji: "🇹🇭",
    },
    {
        Code:  "ti",
        Name:  "Tigrinya",
        Emoji: "",
    },
    {
        Code:  "bo",
        Name:  "Tibetan Standard",
        Emoji: "🇧🇴",
    },
    {
        Code:  "tk",
        Name:  "Turkmen",
        Emoji: "🇹🇰",
    },
    {
        Code:  "tl",
        Name:  "Tagalog",
        Emoji: "🇹🇱",
    },
    {
        Code:  "tn",
        Name:  "Tswana",
        Emoji: "🇹🇳",
    },
    {
        Code:  "to",
        Name:  "Tonga",
        Emoji: "🇹🇴",
    },
    {
        Code:  "ts",
        Name:  "Tsonga",
        Emoji: "",
    },
    {
        Code:  "tt",
        Name:  "Tatar",
        Emoji: "🇹🇹",
    },
    {
        Code:  "tw",
        Name:  "Twi",
        Emoji: "🇹🇼",
    },
    {
        Code:  "ty",
        Name:  "Tahitian",
        Emoji: "",
    },
    {
        Code:  "ug",
        Name:  "Uyghur",
        Emoji: "🇺🇬",
    },
    {
        Code:  "ve",
        Name:  "Venda",
        Emoji: "🇻🇪",
    },
    {
        Code:  "vo",
        Name:  "Volapük",
        Emoji: "",
    },
    {
        Code:  "wa",
        Name:  "Walloon",
        Emoji: "",
    },
    {
        Code:  "cy",
        Name:  "Welsh",
        Emoji: "🇨🇾",
    },
    {
        Code:  "wo",
        Name:  "Wolof",
        Emoji: "",
    },
    {
        Code:  "fy",
        Name:  "Western Frisian",
        Emoji: "",
    },
    {
        Code:  "xh",
        Name:  "Xhosa",
        Emoji: "",
    },
    {
        Code:  "yi",
        Name:  "Yiddish",
        Emoji: "",
    },
    {
        Code:  "yo",
        Name:  "Yoruba",
        Emoji: "",
    },
    {
        Code:  "za",
        Name:  "Zhuang",
        Emoji: "🇿🇦",
    },
}