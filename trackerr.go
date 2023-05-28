// Package trackerr aims to facilitate creation of referenceable errors and
// elegant stack traces.
//
// It was crafted in frustration trying to navigate Go's printed error stacks
// and the challenge of reliably asserting specific error types while testing.
//
// It's important to define errors created via New or Track as package scooped
// (global) or you won't be able to reference them. It is not recommended to
// create trackable errors after initialisation, but if you need to then Realms
// exist for such use cases.
//
// It's also recommended to call Initialised from an init function in package
// main to prevent the creation of trackable errors after program
// initialisation.
//
// You can return a tracked error directly but most of the time you would call
// one of the following receiving functions CausedBy, Because, or BecauseOf
// with additional information.
//
// For manual debugging there's Debug and the deferable DebugPanic which will
// print a readable stack trace.
//
// TODO
// 		* Think about how to integrate file names and line numbers
//		  - How, where, and when to collect them? (reflection)
//		  - How to optimise print outs with them?
//		  - May have to redesign the Debug function?
package trackerr

import (
	"errors"
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

func checkInitState() {
	if globalInitialised {
		panic(Untracked("No tracked errors may be created after initialisation."))
	}
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
