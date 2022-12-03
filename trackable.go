package trackable

// TODO: Think about how to integrate file names and line numbers.
// TODO: - How, where, and when to collect them?
// TODO: - How to optimise print outs with them?
// TODO: - May have to redesign the Debug function?

// TODO: Create a proper example with real test data files

// trackable represents an error that uses an integer as an identifier.
//
// If an ID is less than or equal to zero then the error is categorised as
// untracked according to errors.Is.
type trackable struct {
	id    int
	msg   string
	cause error
	iface string
}

func (e trackable) Error() string {
	if e.cause == nil {
		return e.msg
	}

	return e.msg + ": " + e.cause.Error()
}

func (e trackable) String() string {
	return e.msg
}

func (e trackable) Unwrap() error {
	return e.cause
}

func (e trackable) AsError() error {
	return e
}

func (e trackable) Is(target error) bool {
	if e.id <= 0 {
		return false
	}

	if it, ok := target.(*trackable); ok {
		return e.id == it.id
	}

	return false
}

func (e trackable) Wrap(cause error) error {
	e.cause = cause
	return &e
}

func (e trackable) Because(msg string, args ...any) error {
	e.cause = Untracked(msg, args...)
	return &e
}

func (e trackable) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = Wrap(cause, msg, args...)
	return &e
}

func (e trackable) Interface(cause error, name string) error {
	e.cause = cause
	e.iface = name
	return &e
}

func (e trackable) InterfaceName() string {
	return e.iface
}
