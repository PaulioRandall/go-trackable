package track

type trackedError struct {
	*untrackedError // TODO: might have to implement some errors specifically
	id              int
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

type checkpointError trackedError

func (e *checkpointError) checkpoint() {}
