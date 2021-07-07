package main

import (
    "database/sql"
    "github.com/armanokka/translobot/dashbot"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "gorm.io/gorm"
)

const (
    DashBotAPIKey = "cjVjdWDRijXDk5kl9yGi5TTS9XImME7HbZMOg09F"
    AdminID       = 579515224
)

var (
    db  *gorm.DB
    bot *tgbotapi.BotAPI
    analytics *dashbot.DashBot
)

// Users is table in DB
type Users struct {
    ID     int64 `gorm:"primaryKey;index;not null"`
    MyLang string `gorm:"default:en"`
    ToLang string `gorm:"default:fr"`
    Act sql.NullString `gorm:"default:null"`
}

type Referrers struct {
    ID int64 `gorm:"primaryKey;index;not null"`
    Code string
    Users int64
}
