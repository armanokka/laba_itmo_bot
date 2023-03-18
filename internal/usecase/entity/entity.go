package entity

import "time"

// Translations is a DB table
type Translations struct {
	ID        int64     `gorm:"primaryKey; index; <-:create; not null"`
	MyLang    string    `gorm:"default:es; size:5"`
	ToLang    string    `gorm:"default:zh; size:5"`
	Act       *string   `gorm:"default:null"`
	Usings    int       `gorm:"default:0; not null"`
	Blocked   bool      `gorm:"default:false; not null"`
	Lang      *string   `gorm:"default:null"`
	TTS       bool      `gorm:"default:true; not null"`
	Deeplink  *string   `gorm:"default:null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:true; not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:true; not null"`
}
