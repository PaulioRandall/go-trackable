package track

// UntrackedError represents an untrackable node in an error stack trace.
//
// This interface is primarily for documentation.
type UntrackedError interface {
	error

	// Unwrap returns the error's underlying cause or nil if none exists.
	//
	// It is designed to work with the Is function exposed by the standard
	// errors package.
	Unwrap() error

	// Wrap returns a copy of the receiving error with the passed error as the
	// underlying cause.
	Wrap(error) error

	// Because returns a copy of the receiving error constructing a cause from
	// msg and args.
	Because(msg string, args ...any) error

	// CausedBy returns a copy of the receiving error constructing a cause by
	// wrapping the passed cause with the error msg and args.
	CausedBy(cause error, msg string, args ...any) error

	// Checkpoint returns a copy of the receiving error with a checkpoint
	// error as an intermediate cause.
	//
	// The msg and args are for the intermediate CheckpointError's message.
	Checkpoint(cause error, msg string, args ...any) error
}

type untrackedError struct {
	msg   string
	cause error
}

func (e untrackedError) Error() string {
	return e.msg
}

func (e untrackedError) Unwrap() error {
	return e.cause
}

func (e untrackedError) Wrap(cause error) error {
	e.cause = cause
	return &e
}

func (e untrackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

func (e untrackedError) CausedBy(cause error, msg string, args ...any) error {
	e.cause = causedBy(cause, msg, args...)
	return &e
}

func (e untrackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}
