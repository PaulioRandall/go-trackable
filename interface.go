package trackable

var (
	// ErrTodo is a convenience trackable for specifying a TODO.
	//
	// This can be useful if you're taking a stepwise refinement or test driven
	// approach to writing code.
	ErrTodo Trackable = Track("TODO: Implementation needed")

	// ErrBug is a convenience trackable for use at the site of known bugs.
	ErrBug Trackable = Track("BUG: Fix needed")

	// ErrInsane is a convenience trackable for sanity checks.
	ErrInsane Trackable = Track("Sanity check failed!!")
)

// Trackable represents a trackable error.
//
// Trackable errors are just errors that one can use to precisely compare
// without reference to the error message and easily construct elegant and
// readable stack traces.
type Trackable interface {
	// Unwrap returns the underlying cause or nil if none exists.
	//
	// It is designed to work with the Is function exposed by the standard errors
	// package.
	Unwrap() error

	// Is returns true if the passed error is equivalent to the receiving
	// trackable error.
	//
	// This is a shallow comparison so causes are not checked. It is designed to
	// work with the Is function exposed by the standard errors package.
	Is(error) bool

	// Wrap returns a copy of the receiving error with the passed cause.
	Wrap(cause error) error

	// Because returns a copy of the receiving error constructing a cause from
	// msg and args.
	Because(msg string, args ...any) error

	// Because returns a copy of the receiving error constructing a cause by
	// wrapping the passed cause with the error msg and args.
	BecauseOf(cause error, msg string, args ...any) error
}
