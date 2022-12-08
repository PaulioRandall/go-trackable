package track

// TrackedError represents a trackable node in an error stack trace.
//
// A tracked error may also represents a checkpoint in an error stack. The
// primary purpose being to note interfaces in stack traces, that is, denote
// the key boundary between packages, libraries, systems, and other key
// integration points.
//
// The aim of checkpoints is to enable stack trace partitioning so they are
// more meaningful, readable, navigable. Thus aiding debugging. Key
// information can then be highlighted in stack trace print outs.
//
// This interface is primarily for documentation.
type TrackedError interface {
	UntrackedError

	// Is returns true if the passed error is equivalent to the receiving
	// error.
	//
	// This is a shallow comparison so causes are not checked. It is designed
	// to work with the Is function exposed by the standard errors package.
	Is(error) bool

	// IsCheckpoint returns true if the trackable error represents a checkpoint
	// in the stack trace.
	IsCheckpoint() bool

	// BecauseOf returns a copy of the receiving error calling Because on the
	// passed cause wrapping with the error msg and args.
	//
	// Unlike the CausedBy function the cause here becomes an intermediate cause
	// rather than the root. This allows a single call to add two tracked errors
	// to the error stack at once.
	BecauseOf(cause Because, msg string, args ...any) error
}

// Because represents errors that can have a cause created and attached to
// them.
//
// This interface is primarily for documentation.
type Because interface {
	// Because returns a copy of the receiving error constructing a cause from
	// msg and args.
	Because(msg string, args ...any) error
}

type trackedError struct {
	id           int
	isCheckpoint bool
	msg          string
	cause        error
}

func (e trackedError) Error() string {
	return e.msg
}

func (e trackedError) Unwrap() error {
	return e.cause
}

func (e trackedError) Wrap(cause error) error {
	e.cause = cause
	return &e
}

func (e trackedError) Copy() error {
	return e
}

func (e trackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

func (e trackedError) CausedBy(cause error, msg string, args ...any) error {
	e.cause = causedBy(cause, msg, args...)
	return &e
}

func (e trackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}

func (e trackedError) Is(other error) bool {
	if e2, ok := other.(*trackedError); ok {
		return e.id == e2.id
	}
	return false
}

func (e trackedError) IsCheckpoint() bool {
	return e.isCheckpoint
}

func (e trackedError) BecauseOf(cause Because, msg string, args ...any) error {
	e.cause = cause.Because(msg, args...)
	return &e
}
