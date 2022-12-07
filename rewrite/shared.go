package track

import (
	"fmt"
)

func fmtMsg(msg string, args ...any) string {
	return fmt.Sprintf(msg, args...)
}

func because(msg string, args ...any) *untrackedError {
	return &untrackedError{
		msg: fmtMsg(msg, args...),
	}
}

func becauseOf(cause error, msg string, args ...any) *untrackedError {
	return &untrackedError{
		msg:   fmtMsg(msg, args...),
		cause: cause,
	}
}

func checkpoint(cause error, msg string, args ...any) *checkpointError {
	return &checkpointError{
		msg:   fmtMsg(msg, args...),
		cause: cause,
	}
}
