package main

import (
    "database/sql"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "gorm.io/gorm"
)

const AdminID = 579515224

var (
    db  *gorm.DB
    bot *tgbotapi.BotAPI
)

// Users is table in DB
type Users struct {
    ID     int64 `gorm:"primaryKey;index;not null"`
    MyLang string `gorm:"default:en"`
    ToLang string `gorm:"default:fr"`
    Act sql.NullString `gorm:"default:null"`
}
