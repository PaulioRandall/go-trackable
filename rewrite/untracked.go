package track

import (
	"fmt"
)

type untrackedError struct {
	msg   string
	cause error
}

func (e untrackedError) Error() string {
	return e.msg
}

func (e untrackedError) Unwrap() error {
	return e.cause
}

func (e untrackedError) Wrap(cause error) error {
	e.cause = cause
	return &e
}

func (e untrackedError) Copy() error {
	return e
}

func (e untrackedError) Because(msg string, args ...any) error {
	e.cause = &untrackedError{
		msg: fmt.Sprintf(msg, args...),
	}
	return &e
}

func (e untrackedError) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = &untrackedError{
		msg:   fmt.Sprintf(msg, args...),
		cause: cause,
	}
	return &e
}

func (e untrackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = &checkpointError{
		msg:   fmt.Sprintf(msg, args...),
		cause: cause,
	}
	return &e
}
