package trackerr

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

func causedBy(cause error, msg string, args ...any) *untrackedError {
	return &untrackedError{
		msg:   fmtMsg(msg, args...),
		cause: cause,
	}
}

func checkpoint(cause error, msg string, args ...any) *trackedError {
	return &trackedError{
		isCheckpoint: true,
		msg:          fmtMsg(msg, args...),
		cause:        cause,
	}
}
