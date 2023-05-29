package trackerr

import (
	"fmt"
)

// ErrorThatWraps represents an error that wraps new untracked errors.
type ErrorThatWraps interface {
	error

	// CausedBy wraps the passed causal error.
	CausedBy(cause error) error

	// Unwrap returns the error's underlying cause or nil if none exists.
	Unwrap() error
}

func fmtMsg(msg string, args ...any) string {
	return fmt.Sprintf(msg, args...)
}

func because(msg string, args ...any) *UntrackedError {
	return &UntrackedError{
		msg: fmtMsg(msg, args...),
	}
}

func causedBy(cause error, msg string, args ...any) *UntrackedError {
	return &UntrackedError{
		msg:   fmtMsg(msg, args...),
		cause: cause,
	}
}
