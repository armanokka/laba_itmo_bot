package repo

import (
	"errors"
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
	"gorm.io/gorm"
	"time"
)

var ErrNotFound = errors.New("record not found")

type TranslationRepo struct {
	client *gorm.DB
}

type Users struct {
	ID         int64   `gorm:"primaryKey; index; unique"`
	FirstName  *string `gorm:"default:null"`
	LastName   *string `gorm:"default:null"`
	Patronymic *string `gorm:"default:null"`

	OPDThreadID *int `gorm:"default:null"` // 3.6 -> 36

	ProgrammingThreadID *int `gorm:"default:null"` // 2.5 -> 25

	ITThreadID *int `gorm:"default:null"`

	Action         *string         `gorm:"default:null"`
	TeacherSubject *entity.Subject `gorm:"default:null"`
}

// TODO configure foreign keys in gorm

type Queues struct {
	ID           int       `gorm:"primaryKey; index; autoIncrement; not null"`
	UserID       int64     `gorm:"not null"`
	ThreadID     int       `gorm:"not null"`
	LaboratoryID int       `gorm:"not null"`
	Checked      bool      `gorm:"not null;default:false"`
	Retake       bool      `gorm:"not null;default:false"`
	Passed       bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"` // timestampz
}

type Threads struct {
	ID      int            `gorm:"<-:false; primaryKey; index; unique"`
	Subject entity.Subject `gorm:"not null; uniqueIndex:subject_thread_uindex"`
	Name    string         `gorm:"not null; uniqueIndex:subject_thread_uindex"`
}

type Laboratories struct {
	ID      int            `gorm:"primaryKey; index; unique"`
	Subject entity.Subject `gorm:"not null; uniqueIndex:subject_name_uindex"`
	Name    string         `gorm:"not null; uniqueIndex:subject_name_uindex"`
}

func New(client *gorm.DB) (TranslationRepo, error) {
	if err := client.AutoMigrate(&Threads{}, &Users{}, &Queues{}, &Laboratories{}); err != nil {
		return TranslationRepo{}, err
	}
	return TranslationRepo{client: client}, nil
}
