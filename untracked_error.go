package trackerr

// UntrackedError represents an untrackable node in an error stack trace.
//
// An untracked error may also represent a checkpoint in an error stack. The
// primary purpose being to note interfaces in stack traces, that is, denote
// the key boundary between packages, libraries, systems, and other key
// integration points.
//
// Checkpoints allow stack trace partitioning. Thus making them more
// meaningful, readable, and navigable during debugging.
type UntrackedError struct {
	msg   string
	cause error
	isCp  bool
}

// Error satisfies the error interface.
func (e UntrackedError) Error() string {
	return e.msg
}

// Unwrap returns the error's underlying cause or nil if none exists.
//
// It is designed to work with errors.Is exposed by the standard errors
// package.
func (e UntrackedError) Unwrap() error {
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
func (e UntrackedError) Wrap(cause error) error {
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
func (e UntrackedError) Because(msg string, args ...any) error {
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
func (e UntrackedError) CausedBy(cause error, msg string, args ...any) error {
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
func (e UntrackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}

// IsCheckpoint returns true if the error represents a checkpoint in the stack
// trace.
func (e UntrackedError) IsCheckpoint() bool {
	return e.isCp
}

// BecauseOf returns a copy of the receiving error calling Because on the
// passed ErrorWrapper wrapping with the error msg and args.
//
// Unlike the CausedBy function the cause here becomes an intermediate error
// rather than the root. This allows a single call to add two tracked errors
// to the error stack at once.
func (e UntrackedError) BecauseOf(cause ErrorWrapper, msg string, args ...any) error {
	e.cause = cause.Because(msg, args...)
	return &e
}
