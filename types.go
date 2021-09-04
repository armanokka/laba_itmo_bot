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
    botToken string = "1737819626:AAEoc8WyCq_8rFQcY4q0vtkhqCKro8AudfI"
)

var (
    db  *gorm.DB
    bot *BotAPI
    analytics *dashbot.DashBot
    InlineCacheTime int = 300
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

// Sponsorships is table in DB
type Sponsorships struct {
    ID int64 // ID of user that bought offer
    Text string // Advertise
    ToLangs string // String of languages separated by "," of users that must to receive advertise
    Start time.Time // When offer starts
    Finish time.Time // When offer finish
}

// SponsorshipsOffers is table in DB
type SponsorshipsOffers struct {
    ID int64 // ID of user that bought offer
    Text string // Advertise
    ToLangs string // String of languages separated by "," of users that must to receive advertise
    Days int // Days of sponsorship
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
