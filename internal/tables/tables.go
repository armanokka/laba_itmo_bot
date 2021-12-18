package tables

import (
	"database/sql"
	"time"
)

// Users is table in DB
type Users struct {
	ID     int64 `gorm:"primaryKey;index;not null"`
	Usings int `gorm:"default:0"`
	Lang string `gorm:"default:en"`
	Blocked bool `gorm:"default:false"`
	LastLangs string `gorm:"default:''"`
}

type UsersLogs struct {
	ID int64 // fk users.id
	Intent sql.NullString // varchar(25)
	Text string // varchar(518)
	FromBot bool
	Date time.Time
}