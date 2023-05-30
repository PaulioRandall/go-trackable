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

// BecauseOf creates a new error using the msg, args, and cause as arguments
// then attaches the result as the cause of the receiving error.
//
// Put another way, the cause argument becomes the root cause in
// the error stack.
//
//		top := trackerr.New("top message")
//		cause := trackerr.New("root cause message")
//
//		e := top.BecauseOf(cause, "middle message")
//
//		// top message
//		// ⤷ middle message
//		// ⤷ root cause message
func (e UntrackedError) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = Untracked(msg, args...).CausedBy(cause)
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

// ContextFor wraps the cause with the intention that it provides context for
// it. The rootCause is first wrapped by the cause if it's not nil.
//
//		context := trackerr.New("context message")
//		cause := trackerr.New("cause message")
//		rootCause := trackerr.Untracked("root cause message")
//
//		e := context.ContextFor(cause, rootCause)
//
//		// context message
//		// ⤷ cause message
//		// ⤷ root cause message
func (e UntrackedError) ContextFor(cause ErrorThatWraps, rootCause error) error {
	c := cause.CausedBy(rootCause)
	return e.CausedBy(c)
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
