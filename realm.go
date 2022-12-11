package trackerr

// Realm represents a space where each trackable error (stack trace node)
// has its own unique ID.
//
// There is a private global Realm in the package that will suffice for most
// purposes. It is used via the package level Track and Checkpoint functions.
//
// The receiving functions are designed to be called during package
// initialisation. This means it should only be used to initialise package
// global variables and within init functions. The exception is where
// multiple Realms are in use. Testing is the only use case currently
// conceivable.
//
// Furthermore, all functions return a shallow copy of any passed or
// receiving errors creating a somewhat immutability based ecosystem.
//
// This interface is primarily for documentation.
type Realm interface {

	// Track returns a new tracked error, that is, one with a tracking ID.
	Track(msg string, args ...any) *TrackedError

	// Checkpoint returns a new tracked checkpoint error, that is, one with a
	// tracking ID and indicates a key node within a stack trace.
	Checkpoint(msg string, args ...any) *TrackedError
}

// IntRealm is a Realm that uses a simple incrementing integer field as the
// pool of unique IDs.
//
// The recommended way to use this package is to ignore this struct and use the
// Error and Checkpoint package functions instead. If this package's API is
// used as intended then it would be impossible to cause an integer overflow
// scenario in any real world use case. However, Realms were conceived for such
// an event and for those who really hate the idea of relying on a global
// Realm they have no control over.
//
// The incrementation happens on each call to Track and Checkpoint receiving
// functions.
type IntRealm struct {
	idPool *int
	locked bool
}

// Track returns a new tracked error.
//
// Calls to HasTracked, IsTracked, and IsTrackerr will all return true when
// the error is passed to them.
func (r *IntRealm) Track(msg string, args ...any) *TrackedError {
	return &TrackedError{
		id:  r.newID(),
		msg: fmtMsg(msg, args...),
	}
}

// Checkpoint returns a new trackable checkpoint error.
//
// Calls to HasTracked, IsTracked, IsTrackerr, and IsCheckpoint will all return
// true when the error is passed to them.
func (r *IntRealm) Checkpoint(msg string, args ...any) *TrackedError {
	return &TrackedError{
		id:           r.newID(),
		isCheckpoint: true,
		msg:          fmtMsg(msg, args...),
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
