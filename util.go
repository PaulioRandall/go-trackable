package trackerr

import (
	"errors"
)

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
// tracked, i.e. those created via the Track and Checkpoint functions.
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
// the Track and Checkpoint functions.
func IsTracked(e error) bool {
	_, ok := e.(*TrackedError)
	return ok
}

// IsCheckpoint returns true if the error is a checkpoint.
func IsCheckpoint(e error) bool {
	type checkpointError interface {
		IsCheckpoint() bool
	}

	if ce, ok := e.(checkpointError); ok {
		return ce.IsCheckpoint()
	}
	return false
}

// IsTrackerr returns true if the error is either an UntrackedError or
// TrackedError error from this package. That is, if it's an error defined
// within go-trackerr.
func IsTrackerr(e error) bool {
	_, ok := e.(UntrackedError)
	return ok
}

// Is is a proxy for errors.Is.
func Is(e, target error) bool {
	return errors.Is(e, target)
}

// Unwrap is a proxy for errors.Unwrap.
func Unwrap(e error) error {
	return errors.Unwrap(e)
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
	if len(targets) == 0 {
		return true
	}

	for _, t := range targets {
		if !errors.Is(e, t) {
			return false
		}

		for e = Unwrap(e); errors.Is(e, t); e = Unwrap(e) {
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
