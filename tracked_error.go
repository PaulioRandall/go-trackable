package trackerr

// TrackedError represents a trackable node in an error stack.
type TrackedError struct {
	id    int
	msg   string
	cause error
}

// New is an alias for Track.
func New(msg string, args ...any) *TrackedError {
	return Track(msg, args...)
}

// Track returns a new tracked error from this package's global Realm.
//
// This is the recommended way to use to create all trackable errors.
func Track(msg string, args ...any) *TrackedError {
	checkInitState()
	return globalRealm.Track(msg, args...)
}

// Because constructs a cause from msg and args.
//
//		wrapper := trackerr.New("wrapper message")
//
//		e := wrapper.Because("cause message")
//
//		// wrapper message
//		// ⤷ cause message
func (e TrackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

// BecauseOf creates a new error using the msg, args, and cause as arguments
// then attaches the result as the cause of the receiving error.
//
// Put another way, the cause argument becomes the root cause in
// the error stack.
//
//		top := trackerr.New("top message")
//		cause := trackerr.New("root cause message")
//
//		e := top.BecauseOf(cause, "middle message")
//
//		// top message
//		// ⤷ middle message
//		// ⤷ root cause message
func (e TrackedError) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = Untracked(msg, args...).CausedBy(cause)
	return &e
}

// CausedBy wraps the passed cause.
//
//		wrapper := trackerr.New("wrapper message")
//		cause := trackerr.Untracked("cause message")
//
//		e := wrapper.CausedBy(cause)
//
//		// wrapper message
//		// ⤷ cause message
func (e TrackedError) CausedBy(cause error) error {
	e.cause = cause
	return &e
}

// Error satisfies the error interface.
func (e TrackedError) Error() string {
	return e.msg
}

// Is returns true if the passed error is equivalent to the receiving
// error. This is a shallow comparison so causes are not checked.
//
// It satisfies the Is function referenced by errors.Is in the standard errors
// package.
func (e TrackedError) Is(other error) bool {
	if e2, ok := other.(*TrackedError); ok {
		return e.id == e2.id
	}
	return false
}

// Unwrap returns the error's underlying cause or nil if none exists.
//
// It is designed to work with errors.Is exposed by the standard errors
// package.
func (e TrackedError) Unwrap() error {
	return e.cause
}

// WrapBy sets the cause of the wrapper error as the receiving error.
//
// Put another way, it performs wrapper.CausedBy(receivingError).
//
//		cause := trackerr.Untracked("cause message")
//		wrapper := trackerr.New("wrapper message")
//
//		e := cause.WrapBy(wrapper)
//
//		// wrapper message
//		// ⤷ cause message
func (e TrackedError) WrapBy(wrapper ErrorThatWraps) error {
	return wrapper.CausedBy(e)
}
