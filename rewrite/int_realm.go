package track

import (
	"fmt"
)

// IntRealm is a Realm that uses a simple incrementing integer field as the
// pool of unique IDs.
//
// The recommended way to use this package is to ignore this struct and use the
// Error and Checkpoint package functions instead. If this package's API is
// used as intended then it would be impossible to cause an integer overflow
// scenario in any real world use case. However, Realms were conceived for such
// an event and for those who really hate the idea of relying on a singleton
// Realm they have no control over.
//
// The incrementation happens on each call to Error and Checkpoint receiving
// functions.
type IntRealm struct {
	idPool *int
}

// Untracked returns a new error without a tracking ID.
func (r *IntRealm) Untracked(msg string, args ...any) *untrackedError {
	return &untrackedError{
		msg: fmt.Sprintf(msg, args...),
	}
}

// Error returns a new tracked error from this package's singleton Realm.
func (r *IntRealm) Error(msg string, args ...any) *trackedError {
	return &trackedError{
		id:  r.newID(),
		msg: fmt.Sprintf(msg, args...),
	}
}

// Checkpoint returns a new trackable checkpoint error from this package's
// singleton Realm.
func (r *IntRealm) Checkpoint(msg string, args ...any) *checkpointError {
	return &checkpointError{
		id:  r.newID(),
		msg: fmt.Sprintf(msg, args...),
	}
}

func (r *IntRealm) newID() int {
	if r.idPool == nil {
		n := 0
		r.idPool = &n
	}

	(*r.idPool)++
	return *r.idPool
}
