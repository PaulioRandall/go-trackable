package trackable

var (
	// ErrTodo is a convenience trackable for specifying a TODO.
	//
	// This can be useful if you're taking a stepwise refinement or test driven
	// approach to writing code.
	ErrTodo = Track("TODO: Implementation needed")

	// ErrBug is a convenience trackable for use at the site of known bugs.
	ErrBug = Track("BUG: Fix needed")

	// ErrInsane is a convenience trackable for sanity checks.
	ErrInsane = Track("Sanity check failed!!")
)

// Trackable represents a trackable error. This interface is primarily for
// documentation.
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

	// Interface does the same as BecauseOf except the trackable error is marked
	// as being at the boundary of a key interface.
	//
	// This allows stack traces to be partitioned so they are more meaningful,
	// readable, and navigable.
	Interface(cause error, msg string, args ...any) error

	// IsInterface returns true if the trackable error was created at the site
	// of a key interface.
	IsInterface() bool
}
