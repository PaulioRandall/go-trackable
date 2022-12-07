package track

type untrackedError struct {
	msg   string
	cause error
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
