package main

import (
    "database/sql"
    "github.com/armanokka/translobot/dashbot"
    cache "github.com/patrickmn/go-cache"
    "gorm.io/gorm"
    "time"
)

const (
    DashBotAPIKey = "cjVjdWDRijXDk5kl9yGi5TTS9XImME7HbZMOg09F"
    AdminID       = 579515224
    botToken string = "1737819626:AAEJyD8fnSHdkh6VP3ePdwFkpEnrirLMHp4" //
    LanguagesPaginationLimit int = 20
)



var (
    db  *gorm.DB
    bot *BotAPI
    analytics dashbot.DashBot
    InlineCacheTime int = 864000
    c = cache.New(6 * time.Hour, 12 * time.Hour)
    logs = make(chan UsersLogs, 10)
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
    Blocked bool `gorm:"default:false"`
}

type UsersLogs struct {
    ID int64 // fk users.id
    Intent sql.NullString // varchar(25)
    Text string // varchar(518)
    FromBot bool
    Date time.Time
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
