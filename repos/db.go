package repos

import (
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

	SwapLangs(userID int64) error

	GetUsersNumber() (int64, error)
	GetUsersSlice(offset, count int64, slice []int64) (err error)

	GetAllUsers() ([]tables.Users, error)
}
