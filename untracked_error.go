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
//		```
//		wrapper message
//		⤷ cause message
//		```
func (e UntrackedError) Because(msg string, args ...any) error {
	return e.BecauseOf(nil, msg, args...)
}

// BecauseOf creates a new error using the msg, args, and cause as arguments
// then attaches the result as the cause of the receiving error.
//
// Put another way, the cause argument becomes the root cause in
// the error stack.
//
//		top := trackerr.New("top message")
//		rootCause := trackerr.New("root cause message")
//
//		e := top.BecauseOf(rootCause, "middle message")
//
//		```
//		top message
//		⤷ middle message
//		⤷ root cause message
//		```
func (e UntrackedError) BecauseOf(rootCause error, msg string, args ...any) error {
	e.cause = Untracked(msg, args...).CausedBy(rootCause)
	return &e
}

// CausedBy wraps the rootCause within the first item in causes. Then the
// second item in causes wraps the first. Then the third item wraps the second
// and so on. Finally, the receiving error wraps the result before returning.
//
//		head := trackerr.New("head message")
//		causeA := trackerr.New("cause message A")
//		causeB := trackerr.New("cause message B")
//		rootCause := trackerr.Untracked("root cause message")
//
//		e := head.CausedBy(rootCause, causeB, causeA)
//
//		```
//		head message
//		⤷ cause message A
//		⤷ cause message B
//		⤷ root cause message
//		```
//
// CausedBy will very often be used to wrap a single error.
//
//		head := trackerr.New("head message")
//		cause := trackerr.Untracked("cause message")
//
//		e := head.CausedBy(cause)
//
//		```
//		head message
//		⤷ cause message
//		```
func (e UntrackedError) CausedBy(rootCause error, causes ...ErrorThatWraps) error {
	c := Stack(rootCause, causes...)
	e.cause = c
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
