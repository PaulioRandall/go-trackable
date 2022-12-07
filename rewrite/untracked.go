package track

// UntrackedError represents an untrackable node in an error stack trace.
//
// This interface is primarily for documentation.
type UntrackedError interface {
	ErrorWrap

	// Because returns a copy of the receiving error constructing a cause from
	// msg and args.
	Because(msg string, args ...any) error

	// Because returns a copy of the receiving error constructing a cause by
	// wrapping the passed cause with the error msg and args.
	BecauseOf(cause error, msg string, args ...any) error

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

func (e untrackedError) Copy() error {
	return e
}

func (e untrackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

func (e untrackedError) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = becauseOf(cause, msg, args...)
	return &e
}

func (e untrackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}
