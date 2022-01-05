package repos

import "github.com/armanokka/translobot/internal/tables"

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
	// UpdateUserLastActivity is a wrapper to UpdateUserByMap
	UpdateUserLastActivity(id int64) error
	// IncreaseUserUsings is a wrapper to UpdateUserByMap
	IncreaseUserUsings(id int64) error

	// GetAllUsers
	// Errors: unknown
	GetAllUsers() ([]tables.Users, error)
	BotLogs
}

type BotLogs interface {
	LogUserMessage(id int64, text string) error
	LogBotMessage(toID int64, intent string, text string) error
}

