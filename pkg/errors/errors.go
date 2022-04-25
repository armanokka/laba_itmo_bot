package errors

import (
	"errors"
	"runtime/debug"
)

type Error struct {
	err   error
	stack []byte
}

func (e Error) Error() string {
	return e.err.Error()
}
func (e Error) Stack() []byte {
	return e.stack
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return Error{
		err:   err,
		stack: debug.Stack(),
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
