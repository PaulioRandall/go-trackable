package trackerr

// UntrackedError represents an untrackable node in an error stack trace.
type UntrackedError struct {
	msg   string
	cause error
}

func (e UntrackedError) Error() string {
	return e.msg
}

// Unwrap returns the error's underlying cause or nil if none exists.
//
// It is designed to work with the Is function exposed by the standard
// errors package.
func (e UntrackedError) Unwrap() error {
	return e.cause
}

// Wrap returns a copy of the receiving error with the passed error as the
// underlying cause.
func (e UntrackedError) Wrap(cause error) error {
	e.cause = cause
	return &e
}

// Because returns a copy of the receiving error constructing a cause from
// msg and args.
func (e UntrackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

// CausedBy returns a copy of the receiving error constructing a cause by
// wrapping the passed cause with the error msg and args.
func (e UntrackedError) CausedBy(cause error, msg string, args ...any) error {
	e.cause = causedBy(cause, msg, args...)
	return &e
}

// Checkpoint returns a copy of the receiving error with a checkpoint
// error as an intermediate cause.
//
// The msg and args are for the intermediate CheckpointError's message.
func (e UntrackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}
