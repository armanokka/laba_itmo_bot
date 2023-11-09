package usecase

import (
	"github.com/armanokka/laba_itmo_bot/internal/usecase/entity"
)

type Logger interface {
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
	With(args ...interface{}) Logger
}

type Repo interface {
	// Users
	GetUserByID(id int64) (entity.User, error)
	CreateUser(user entity.User) error
	UpdateUserByID(id int64, columnValue ...interface{}) error
	GetAllUsersIDs() ([]int64, error)
	GetPassedLabsByID(userID int64) ([]int, error)

	// Queues
	GetQueueByThreadID(threadID int) ([]entity.QueueUser, error)
	AddUserToQueue(userID int64, threadID int, labID int) error
	RemoveUserFromQueue(userID int64, threadID int) error
	GradeLab(studentID int64, threadID int, passed bool) (labID int, err error)
	UserRetakesLab(studentID int64, threadID int) (in bool, err error)
	UserInQueue(studentID int64, threadID int) (in bool, err error) // whether user is in any queue of the subject

	// Threads
	GetThreadByID(threadID int) (result entity.Thread, err error)
	GetThreadsBySubject(subject entity.Subject) ([]entity.Thread, error)
	AddThread(n string, subject entity.Subject) error
	DeleteThread(threadID int) error
	RenameThread(threadID int, newName string) error

	// Laboratories
	GetNextLab(userID int64, subject entity.Subject) (entity.Laboratory, error)
	GetLaboratoriesBySubject(subject entity.Subject) ([]entity.Laboratory, error)
	GetLaboratoryNameByID(labID int) (string, error)
}
