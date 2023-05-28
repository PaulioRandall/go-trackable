package trackerr

import (
	"fmt"
)

// ErrorWrapper represents an error that wraps new untracked errors.
type ErrorWrapper interface {
	error

	// Because constructing a causal error from the msg and args.
	//
	// A call to errors.Unwrap the resultant error will yeild an underlying error
	// with the supplied msg and args.
	Because(msg string, args ...any) error

	// CausedBy wraps the passed causal error.
	CausedBy(cause error) error
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
