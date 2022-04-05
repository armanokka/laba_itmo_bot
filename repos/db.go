package repos

import (
	"database/sql"
	"github.com/armanokka/translobot/internal/tables"
)

type BotDB interface {
	// GetUserByID
	// Errors: gorm.ErrRecordNotFound, unknown
	GetUserByID(id int64) (tables.Users, error)
	// CreateUser
	// Errors: unknown
	CreateUser(user tables.Users) error
	// UpdateUser updates non-default fields in passed struct
	// Errors: gorm.ErrRecordNotFound, unknown
	UpdateUser(id int64, updates tables.Users) error
	// UpdateUserByMap
	// Errors: gorm.ErrRecordNotFound, unknown
	UpdateUserByMap(id int64, updates map[string]interface{}) error
	// UpdateUserMetrics is a union of UpdateUserLastActivity, IncreaseUserUsings and LogUserMessage
	UpdateUserMetrics(id int64, message string) error

	GetRandomUser() (tables.Users, error)
	// GetAllUsers
	// Errors: unknown
	GetAllUsers() ([]tables.Users, error)
	BotLogs
	Mailing
}

type BotLogs interface {
	LogUserMessage(id int64, text string) error
	LogBotMessage(toID int64, intent string, text string) error
	GetUserLogs(id int64, limit int) ([]tables.UsersLogs, error)
}

type Mailing interface { // temp table for mailings
	GetMailersRows() (rows *sql.Rows, err error)
	MailingExists() (bool, error)
	DeleteMailuser(id int64) error
	DropMailings() error
	CreateMailingTable() error
}
