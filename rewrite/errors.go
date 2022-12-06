package track

type intRealm struct {
}

type untrackedError struct {
}

func (e *untrackedError) Unwrap() error {
	panic("TODO errors.untrackedError.Unwrap")
}

func (e *untrackedError) Is(error) bool {
	panic("TODO errors.untrackedError.Is")
}

func (e *untrackedError) Wrap(error) error {
	panic("TODO errors.untrackedError.Wrap")
}

func (e *untrackedError) Copy() error {
	panic("TODO errors.untrackedError.Copy")
}

type trackedError struct {
	untrackedError // TODO: might have to implement some errors specifically
}

func (e *trackedError) Because(msg string, args ...any) error {
	panic("TODO errors.trackedError.Because")
}

func (e *trackedError) BecauseOf(cause error, msg string, args ...any) error {
	panic("TODO errors.trackedError.BecauseOf")
}

func (e *trackedError) Checkpoint(cause error, msg string, args ...any) error {
	panic("TODO errors.trackedError.Checkpoint")
}

type checkpointError untrackedError

func (e *checkpointError) checkpoint() {}
