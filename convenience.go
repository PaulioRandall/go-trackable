package trackable

import (
	"errors"
	"fmt"
)

// IsTracked returns true if the error has a trackable ID greater than zero.
func IsTracked(e error) bool {
	te, ok := e.(*trackable)
	return ok && te.id > 0
}

// Is is an alias for errors.Is.
func Is(e, target error) bool {
	return errors.Is(e, target)
}

// All returns true only if errors.Is returns true for all targets.
func All(e error, targets ...error) bool {
	for _, t := range targets {
		if !errors.Is(e, t) {
			return false
		}
	}
	return true
}

// Any returns true if errors.Is returns true for any of the targets.
func Any(e error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(e, t) {
			return true
		}
	}
	return false
}

// Debug is a convenience for fmt.Print("ERROR: ", e.Error()).
func Debug(e error) (int, error) {
	if e == nil {
		return fmt.Print("ERROR: nil")
	}
	return fmt.Print("ERROR: ", e.Error())
}

// Debugln is a convenience for fmt.Println("ERROR: ", e.Error()).
func Debugln(e error) (int, error) {
	if e == nil {
		return fmt.Print("ERROR: nil")
	}
	return fmt.Println("ERROR: ", e.Error())
}
