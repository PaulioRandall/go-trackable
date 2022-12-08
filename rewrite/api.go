// Package track aims to facilitate creation of referenceable errors and
// elegant stack traces.
package track

import (
	"errors"
	"fmt"
	"strings"
)

// TODO 1: Write up a realistic example for this interface using test data
// TODO 2: Implement package interface

// TODO: Think about how to integrate file names and line numbers.
// TODO: - How, where, and when to collect them?
// TODO: - How to optimise print outs with them?
// TODO: - May have to redesign the Debug function?

var (
	globalRealm IntRealm

	// ErrTodo is a convenience tracked error for specifying a TODO.
	//
	// This can be useful if you're taking a stepwise refinement or test driven
	// approach to writing code.
	ErrTodo = Error("TODO: Implementation needed")

	// ErrBug is a convenience tracked error for use at the site of known bugs.
	ErrBug = Error("BUG: Fix needed")

	// ErrInsane is a convenience tracked error for sanity checking.
	ErrInsane = Error("Sanity check!!")
)

// Untracked returns a new error without a tracking ID.
//
// It is no different than using errors.New except it has a handy fmt.Sprintf
// signiture and a few extra receiving functions for any niche use cases one
// may encounter.
func Untracked(msg string, args ...any) *untrackedError {
	return globalRealm.Untracked(msg, args...)
}

// Error is an alias for Track. More readable when combined with the package
// name, i.e. 'track.Error(...)' as opposed to 'track.Track(...)'.
//
// The Track function is kept because it maintains alignment with the Realm
// interface.
func Error(msg string, args ...any) *trackedError {
	return globalRealm.Track(msg, args...)
}

// Track returns a new tracked error from this package's singleton Realm.
//
// This is recommended way to use to create all trackable errors outside of
// testing.
func Track(msg string, args ...any) *trackedError {
	return globalRealm.Track(msg, args...)
}

// Checkpoint returns a new trackable checkpoint error from this package's
// singleton Realm.
//
// This is recommended way to use to create all checkpoint errors outside of
// testing.
func Checkpoint(msg string, args ...any) *trackedError {
	return globalRealm.Checkpoint(msg, args...)
}

// Debug pretty prints the error stack trace to terminal for debugging
// purposes.
//
// If e is nil then a message will be printed indicating so. While this
// function can be used for logging it's not designed for it.
func Debug(e error) (int, error) {
	s := ErrorStack(e)

	if s == "" {
		return fmt.Print("[Debugging error] nil error")
	}

	return fmt.Print("[Debugging error]\n", s)
}

// HasTracked returns true if the error or one of the underlying causes are
// tracked, i.e. those created via the Error and Checkpoint functions.
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

// IsTracked returns true if the error is being tracked, i.e. those created via
// the Error and Checkpoint functions.
func IsTracked(e error) bool {
	_, ok := e.(*trackedError)
	return ok
}

// IsCheckpoint returns true if the error is a trackable checkpoint.
func IsCheckpoint(e error) bool {
	if te, ok := e.(*trackedError); ok {
		return te.IsCheckpoint()
	}
	return false
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

// Any returns true if errors.Is returns true for any target.
func Any(e error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(e, t) {
			return true
		}
	}
	return false
}

// ErrorStack returns a human readable stack trace for the error.
func ErrorStack(e error) string {
	sb := strings.Builder{}
	sb.WriteString("  ")

	for i, cause := range AsStack(e) {

		var prefix, suffix string
		if IsCheckpoint(cause) {
			prefix = "\n——"
			suffix = "——"
		} else {
			prefix = "\n⤷ "
		}

		if i > 0 {
			sb.WriteString(prefix)
		}

		s := ErrorWithoutCause(cause)
		sb.WriteString(s)

		if i > 0 {
			sb.WriteString(suffix)
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

// AsStack recursively unwraps the error returning a slice of errors.
//
// The passed error e will be first and root cause last.
func AsStack(e error) []error {
	var stack []error

	for e != nil {
		stack = append(stack, e)
		e = errors.Unwrap(e)
	}

	return stack
}

// ErrorWithoutCause removes the cause from error messages that use the format
// '%s: %w'. Where s is the error message and w is the cause's message.
func ErrorWithoutCause(e error) string {
	s := e.Error()

	if _, ok := e.(*untrackedError); ok {
		return s
	}

	if _, ok := e.(*trackedError); ok {
		return s
	}

	cause := errors.Unwrap(e)

	if cause == nil {
		return s
	}

	s = strings.TrimSuffix(s, cause.Error())
	s = strings.TrimSpace(s)
	return strings.TrimSuffix(s, ":")
}
