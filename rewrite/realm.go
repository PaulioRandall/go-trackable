package track

import (
	"fmt"
)

var globalRealm intRealm

type intRealm struct {
	idPool *int
}

// Untracked returns a new error without a tracking ID.
func (r intRealm) Untracked(msg string, args ...any) *untrackedError {
	return &untrackedError{
		msg: fmt.Sprintf(msg, args...),
	}
}

// Error returns a new tracked error from this package's singleton Realm.
func (r intRealm) Error(msg string, args ...any) *trackedError {
	return &trackedError{
		untrackedError: Untracked(msg, args...),
		id:             r.newID(),
	}
}

// Checkpoint returns a new trackable checkpoint error from this package's
// singleton Realm.
func (r intRealm) Checkpoint(msg string, args ...any) *checkpointError {
	panic("TODO intRealm.Checkpoint")
}

func (r *intRealm) newID() int {
	if r.idPool == nil {
		n := 0
		r.idPool = &n
	}

	(*r.idPool)++
	return *r.idPool
}
