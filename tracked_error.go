package trackerr

// Because represents errors that can have a cause created and attached to
// them.
type Because interface {
	// Because returns a copy of the receiving error constructing a cause from
	// msg and args.
	Because(msg string, args ...any) error
}

// TrackedError represents a trackable node in an error stack trace.
type TrackedError struct {
	id           int
	msg          string
	cause        error
	isCheckpoint bool
}

func (e TrackedError) Error() string {
	return e.msg
}

func (e TrackedError) Unwrap() error {
	return e.cause
}

func (e TrackedError) Wrap(cause error) error {
	e.cause = cause
	return &e
}

func (e TrackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

func (e TrackedError) CausedBy(cause error, msg string, args ...any) error {
	e.cause = causedBy(cause, msg, args...)
	return &e
}

func (e TrackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}

// Is returns true if the passed error is equivalent to the receiving
// error.
//
// This is a shallow comparison so causes are not checked. It is designed
// to work with the Is function exposed by the standard errors package.
func (e TrackedError) Is(other error) bool {
	if e2, ok := other.(*TrackedError); ok {
		return e.id == e2.id
	}
	return false
}

// IsCheckpoint returns true if the trackable error represents a checkpoint
// in the stack trace.
func (e TrackedError) IsCheckpoint() bool {
	return e.isCheckpoint
}

// BecauseOf returns a copy of the receiving error calling Because on the
// passed cause wrapping with the error msg and args.
//
// Unlike the CausedBy function the cause here becomes an intermediate cause
// rather than the root. This allows a single call to add two tracked errors
// to the error stack at once.
func (e TrackedError) BecauseOf(cause Because, msg string, args ...any) error {
	e.cause = cause.Because(msg, args...)
	return &e
}
