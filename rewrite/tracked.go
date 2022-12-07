package track

type (

	// TrackedError represents a trackable node in an error stack trace.
	//
	// This interface is primarily for documentation.
	TrackedError interface {
		UntrackedError

		// Is returns true if the passed error is equivalent to the receiving
		// error.
		//
		// This is a shallow comparison so causes are not checked. It is designed
		// to work with the Is function exposed by the standard errors package.
		Is(error) bool
	}

	// CheckpointError represents a noteworthy node in an error stack trace.
	//
	// The aim is to enable easier reading and debugging of by allowing stack
	// trace printing to highlight key information for navigating to issues. This
	// allows stack traces to be partitioned so they are more meaningful,
	// readable, and navigable.
	//
	// The primary intended purpose is to note interfaces in stack traces, that
	// is, denote the key boundary between packages, libraries, systems, and
	// other key integration points.
	//
	// This interface is primarily for documentation.
	CheckpointError interface {
		TrackedError
		checkpointError()
	}
)

type trackedError struct {
	id    int
	msg   string
	cause error
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

func (e trackedError) Is(other error) bool {
	if e2, ok := other.(*trackedError); ok {
		return e.id == e2.id
	}
	return false
}

func (e trackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

func (e trackedError) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = becauseOf(cause, msg, args...)
	return &e
}

func (e trackedError) Checkpoint(cause error, msg string, args ...any) error {
	e.cause = checkpoint(cause, msg, args...)
	return &e
}

type checkpointError = trackedError

func (e checkpointError) checkpoint() {}
