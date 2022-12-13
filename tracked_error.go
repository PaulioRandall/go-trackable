package trackerr

// TrackedError represents a trackable node in an error stack trace.
//
// A tracked error may also represents a checkpoint in an error stack. The
// primary purpose being to note interfaces in stack traces, that is, denote
// the key boundary between packages, libraries, systems, and other key
// integration points.
//
// Checkpoints allow stack trace partitioning. Thus making them more
// meaningful, readable, and navigable during debugging.
type TrackedError struct {
	id    int
	msg   string
	cause error
	isCp  bool
}

// Error satisfies the error interface.
func (e TrackedError) Error() string {
	return e.msg
}

// Unwrap returns the error's underlying cause or nil if none exists.
//
// It is designed to work with the Is function exposed by the standard
// errors package.
func (e TrackedError) Unwrap() error {
	return e.cause
}

// Wrap returns a copy of the receiving error with the passed error as the
// underlying cause.
//
//		cause := trackerr.Untracked("cause message")
//		wrapper := trackerr.Tracked("wrapper message")
//
//		error := wrapper.Wrap(cause)
//
//		// wrapper message
//		// ⤷ cause message
func (e TrackedError) Wrap(cause error) error {
	e.cause = cause
	return &e
}

// Because returns a copy of the receiving error constructing a cause from
// msg and args.
//
//		wrapper := trackerr.Tracked("wrapper message")
//
//		error := wrapper.Because("cause message")
//
//		// wrapper message
//		// ⤷ cause message
func (e TrackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

// CausedBy returns a copy of the receiving error constructing a cause by
// wrapping the passed cause with the error msg and args.
//
//		rootCause := trackerr.Untracked("root cause message")
//		wrapper := trackerr.Tracked("wrapper message")
//
//		error := wrapper.CausedBy(cause, "caused by message")
//
//		// wrapper message
//		// ⤷ caused by message
//		// ⤷ root cause message
func (e TrackedError) CausedBy(cause error, msg string, args ...any) error {
	e.cause = causedBy(cause, msg, args...)
	return &e
}

// Checkpoint returns a copy of the receiving error with a checkpoint
// error as an intermediate cause.
//
// The msg and args are for the intermediate CheckpointError's message.
//
//		rootCause := trackerr.Untracked("root cause message")
//		cause := trackerr.Wrap(rootCause, "cause message")
//		wrapper := trackerr.Tracked("wrapper message")
//
//		error := wrapper.Checkpoint(cause, "checkpoint message")
//
//		// wrapper message
//		// ——checkpoint message——
//		// ⤷ cause message
//		// ⤷ root cause message
func (e TrackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}

// IsCheckpoint returns true if the error represents a checkpoint in the stack
// trace.
func (e TrackedError) IsCheckpoint() bool {
	return e.isCp
}

// BecauseOf returns a copy of the receiving error calling Because on the
// passed ErrorWrapper wrapping with the error msg and args.
//
// Unlike the CausedBy function the cause here becomes an intermediate cause
// rather than the root. This allows a single call to add two tracked errors
// to the error stack at once.
func (e TrackedError) BecauseOf(cause ErrorWrapper, msg string, args ...any) error {
	e.cause = cause.Because(msg, args...)
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
