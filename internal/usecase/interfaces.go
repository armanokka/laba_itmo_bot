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
	GetUserPassedLaboratoriesIDs(userID int64, userThread int) ([]int, error)

	// Queues
	GetQueueBySubject(threadID int, labID int) ([]entity.QueueUser, error)
	AddUserToQueue(userID int64, userThreadID int, labID int) error
	RemoveUserFromQueue(userID int64, threadID int, labID int) error
	MarkLabAsNotPassed(studentID int64, lab int) error
	MarkLabAsPassed(studentID int64, lab int) error
	UserPassedLab(studentID int64, lab int) (exists bool, err error)
	UserInAnyQueue(studentID int64, subject entity.Subject) (in bool, err error) // whether user is in any queue of the subject

	// Threads
	GetThreadsBySubject(subject entity.Subject) ([]entity.Thread, error)
	AddThread(n string, subject entity.Subject) error
	DeleteThread(threadID int) error
	GetThreadNameByID(threadID int) (string, error)
	RenameThread(threadID int, newName string) error

	// Laboratories
	GetLaboratoriesBySubject(subject entity.Subject) ([]entity.Laboratory, error)
	GetLaboratoryNameByID(labID int) (string, error)
}
