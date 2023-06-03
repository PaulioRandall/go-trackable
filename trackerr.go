// Package trackerr aims to facilitate creation of referenceable errors and
// elegant stack traces.
//
// It was crafted in frustration trying to navigate Go's printed error stacks
// and the challenge of reliably asserting specific error types while testing.
package trackerr

import (
	"errors"
	"fmt"
)

var (
	// ErrTodo is a convenience trackable error for specifying a TODO.
	//
	// This can be useful if you're taking a stepwise refinement or test driven
	// approach to writing code.
	ErrTodo = New("TODO: Implementation needed")

	// ErrBug is a convenience trackable error for use at the site of known bugs.
	ErrBug = New("BUG: Fix needed")

	// ErrInsane is a convenience trackable error for sanity checking.
	ErrInsane = New("Sanity check failed!!")
)

// ErrorThatWraps represents an error that wraps new untracked errors.
type ErrorThatWraps interface {
	error

	// CausedBy wraps the rootCause within the first item in causes. Then the
	// second item in causes wraps the first. Then the third item wraps the
	// second and so on. Finally, the receiving error wraps the result before
	// returning.
	CausedBy(rootCause error, causes ...ErrorThatWraps) error

	// Unwrap returns the error's underlying cause or nil if none exists.
	Unwrap() error
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

// AllOrdered returns true only if errors.Is returns true for all targets and
// the order of the targets matches the descending order of errors.
//
// This does not mean the depth of the error stack must be the same as the
// number of targets. If three targets are supplied then true is returned if
// during descent of the error stack:
//		THE three targets exist in the error stack
//		AND first target is found first
//		AND the second target is found second
//		And the third target is found last
func AllOrdered(e error, targets ...error) bool {
	for _, t := range targets {
		if !errors.Is(e, t) {
			return false
		}

		for e = Unwrap(e); errors.Is(e, t); e = Unwrap(e) {
		}
	}

	return true
}

// Any returns true if errors.Is returns true for at least one target.
func Any(e error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(e, t) {
			return true
		}
	}

	return false
}

// HasTracked returns true if the error or one of the underlying causes are
// tracked, i.e. those created via the New or Track functions.
func HasTracked(e error) bool {

	type wrapper interface {
		Unwrap() error
	}

	for e != nil {
		if IsTracked(e) {
			return true
		}

		if w, ok := e.(wrapper); ok {
			e = w.Unwrap()
		} else {
			e = nil
		}
	}

	return false
}

// Is is a proxy for errors.Is.
func Is(e, target error) bool {
	return errors.Is(e, target)
}

// IsTracked returns true if the error is being tracked, i.e. those created via
// the New or Track functions.
func IsTracked(e error) bool {
	_, ok := e.(*TrackedError)
	return ok
}

// IsTrackerr returns true if the error is either an UntrackedError or
// TrackedError from this package. That is, if it's an error defined
// by go-trackerr.
func IsTrackerr(e error) bool {
	_, ok := e.(TrackedError)

	if ok {
		return ok
	}

	_, ok = e.(UntrackedError)
	return ok
}

// Unwrap is a proxy for errors.Unwrap.
func Unwrap(e error) error {
	return errors.Unwrap(e)
}

// Debug pretty prints the error stack trace to terminal for debugging
// purposes.
//
// If e is nil then a message will be printed indicating so. This function is
// not designed for logging, just day to day manual debugging.
func Debug(e error) (int, error) {
	s := ErrorStack(e)

	if s == "" {
		return fmt.Print("[DEBUG ERROR] nil error")
	}

	return fmt.Print("[DEBUG ERROR]\n  ", s)
}

// DebugPanic recovers from a panic, prints out the error using the Debug
// function, and finally sets it as the catch error's pointer value.
//
// If nil is passed as the catch then the panic continues after printing.
//
// If the panic value is not an error the panic will continue!
//
// This function is not designed for logging, just day to day manual debugging.
func DebugPanic(catch *error) {
	v := recover()

	if v == nil {
		return
	}

	e, ok := v.(error)
	if !ok {
		panic(v)
	}

	Debug(e)

	if catch == nil {
		panic(e)
	}
	*catch = e
}
