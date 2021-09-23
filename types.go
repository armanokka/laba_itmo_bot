package main

import (
    "database/sql"
    "github.com/armanokka/translobot/dashbot"
    "gorm.io/gorm"
    "time"
)

const (
    DashBotAPIKey = "cjVjdWDRijXDk5kl9yGi5TTS9XImME7HbZMOg09F"
    AdminID       = 579515224
    botToken string = "1737819626:AAEJyD8fnSHdkh6VP3ePdwFkpEnrirLMHp4 "
    layout string = "2006/01/02 15:04" // для парсинга времени в /ad, год-месяц-день час:минута
    TimeLocation string = "Europe/Moscow"
)

var (
    now time.Time
    loc *time.Location
)



var (
    db  *gorm.DB
    bot *BotAPI
    analytics *dashbot.DashBot
    InlineCacheTime int = 864000
)

// Users is table in DB
type Users struct {
    ID     int64 `gorm:"primaryKey;index;not null"`
    MyLang string `gorm:"default:en"`
    ToLang string `gorm:"default:fr"`
    Act sql.NullString `gorm:"default:null"`
    Mailing bool `gorm:"default:true"`
    Usings int `gorm:"default:0"`
    Lang string `gorm:"default:en"`
    ReferrerID int64 `gorm:"default:null"`
}

type Ads struct {
    ID int64 // user id of admin that made it
    Content string // text content of an ad
    StartDate time.Time
    FinishDate time.Time
    IDWhoseAd int64 // user id who's this ad
    Views int `gorm:"default:0"`
    ToLangs string // en,ru,ja,es - languages of users that must see an ad
}

type AdsOffers struct {
    ID int64 // user id of admin that made it
    Content string // text content of an ad
    StartDate time.Time
    FinishDate time.Time
    IDWhoseAd int64 // user id who's this ad
    ToLangs string // en,ru,ja,es - languages of users that must see an ad
}


type Localization struct {
    LanguageCode string
    Text string
}

type Lang struct {
    Name string
    Emoji string
}

// type Referrers struct {
//     ID int64 `gorm:"primaryKey;index;not null"`
//     Code string
//     Users int64
// }
