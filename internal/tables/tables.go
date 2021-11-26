package tables

import (
	"database/sql"
	"time"
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
	IsDeveloper sql.NullBool
}

type UsersLogs struct {
	ID int64 // fk users.id
	Intent sql.NullString // varchar(25)
	Text string // varchar(518)
	FromBot bool
	Date time.Time
}