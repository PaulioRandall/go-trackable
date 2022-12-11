// Package trackerr aims to facilitate creation of referenceable errors and
// elegant stack traces.
//
// The recommended way to create errors is via the Track, Checkpoint,
// Untracked, and Wrap package functions. It is not recommended to create
// trackable errors after initialisation but Realms exist for any potential
// use cases.
//
// It is recommended to call the Initialised function from an init function in
// package main to prevent creation of trackable errors after initialisation.
package trackerr

import (
	"errors"
)

// TODO: Think about how to integrate file names and line numbers.
// TODO: - How, where, and when to collect them?
// TODO: - How to optimise print outs with them?
// TODO: - May have to redesign the Debug function?

var (
	// ErrTodo is a convenience tracked error for specifying a TODO.
	//
	// This can be useful if you're taking a stepwise refinement or test driven
	// approach to writing code.
	ErrTodo = Track("TODO: Implementation needed")

	// ErrBug is a convenience tracked error for use at the site of known bugs.
	ErrBug = Track("BUG: Fix needed")

	// ErrInsane is a convenience tracked error for sanity checking.
	ErrInsane = Track("Sanity check failed!!")
)

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

// IsTrackerr returns true if the error implements either the UntrackedError or
// TrackedError interfaces.
//
// That is, if it's an error defined outside of go-trackerr. However, any error
// that implements the interface will return true.
//
// The primary use case is to distinguish go-trackerr errors from second and
// third party errors.
func IsTrackerr(e error) bool {
	_, ok := e.(UntrackedError)
	return ok
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
