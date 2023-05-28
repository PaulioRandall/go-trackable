package trackerr

// UntrackedError represents an untracked error in an error stack.
type UntrackedError struct {
	msg   string
	cause error
}

// Untracked returns a new error without a tracking ID.
//
// This is the same as calling errors.New except for the handy fmt.Sprintf
// function signature and the resultant error has a few extra receiving
// functions for attaching causal errors.
func Untracked(msg string, args ...any) *UntrackedError {
	return because(msg, args...)
}

// Because constructs a cause from msg and args.
//
//		wrapper := trackerr.New("wrapper message")
//
//		e := wrapper.Because("cause message")
//
//		// wrapper message
//		// ⤷ cause message
func (e UntrackedError) Because(msg string, args ...any) error {
	e.cause = because(msg, args...)
	return &e
}

// BecauseOf first calls cause.Because with the error msg and args as arguments
// then attaches the resultant error as the cause of the receiving error.
//
// The cause must satisfy the ErrorWrapper interface.
//
// Put another way, the cause (ErrorWrapper) becomes an intermediate error in
// the error stack. This allows a single call to add two errors to the error
// stack at once.
//
//		top := trackerr.New("top level message")
//		mid := trackerr.New("mid level message")
//
//		e := top.BecauseOf(mid, "low level message")
//
//		// top level message
//		// ⤷ mid level message
//		// ⤷ low level message
func (e UntrackedError) BecauseOf(cause ErrorWrapper, msg string, args ...any) error {
	e.cause = cause.Because(msg, args...)
	return &e
}

// CausedBy wraps the passed cause.
//
//		wrapper := trackerr.New("wrapper message")
//		cause := trackerr.Untracked("cause message")
//
//		e := wrapper.CausedBy(cause)
//
//		// wrapper message
//		// ⤷ cause message
func (e UntrackedError) CausedBy(cause error) error {
	e.cause = cause
	return &e
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

// WrapBy sets the cause of the wrapper error as the receiving error.
//
// Put another way, it performs wrapper.CausedBy(receivingError).
//
//		cause := trackerr.Untracked("cause message")
//		wrapper := trackerr.New("wrapper message")
//
//		e := cause.WrapBy(wrapper)
//
//		// wrapper message
//		// ⤷ cause message
func (e UntrackedError) WrapBy(wrapper ErrorWrapper) error {
	return wrapper.CausedBy(e)
}
