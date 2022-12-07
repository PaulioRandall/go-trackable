package track

import (
	"fmt"
)

type trackedError struct {
	id    int
	msg   string
	cause error
}

func (e trackedError) Error() string {
	return e.msg
}

func (e trackedError) Unwrap() error {
	return e.cause
}

func (e trackedError) Wrap(cause error) error {
	e.cause = cause
	return &e
}

func (e trackedError) Copy() error {
	return e
}

func (e trackedError) Is(other error) bool {
	if e2, ok := other.(*trackedError); ok {
		return e.id == e2.id
	}
	return false
}

func (e trackedError) Because(msg string, args ...any) error {
	e.cause = &untrackedError{
		msg: fmt.Sprintf(msg, args...),
	}
	return &e
}

func (e trackedError) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = &untrackedError{
		msg:   fmt.Sprintf(msg, args...),
		cause: cause,
	}
	return &e
}

func (e trackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = &checkpointError{
		msg:   fmt.Sprintf(msg, args...),
		cause: cause,
	}
	return &e
}

type checkpointError = trackedError

func (e checkpointError) checkpoint() {}
