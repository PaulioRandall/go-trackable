// Package trackerr aims to facilitate creation of referenceable errors and
// elegant stack traces.
//
// It was crafted in frustration trying to navigate Go's printed error stacks
// and the challenge of reliably asserting specific error types while testing.
//
// It's important to define errors created via Track and Checkpoint as
// package scooped (global) or you won't be able to reference them. It is not
// recommended to create trackable errors after initialisation but Realms do
// exist if such use cases appear.
//
// It is also recommended to call Initialised from an init function in
// package main to prevent the creation of trackable errors after program
// initialisation.
//
// You can return a tracked error directly but it's recommended to call one of
// the receiving functions Wrap, CausedBy, Because, BecauseOf, or Checkpoint
// with additional information.
//
// For manual debugging there's Debug and the deferable DebugPanic which will
// print a readable stack trace.
//
// TODO
// 		* Think about how to integrate file names and line numbers
//		  - How, where, and when to collect them? (reflection)
//		  - How to optimise print outs with them?
//		  - May have to redesign the Debug function?
package trackerr

// ErrorWrapper represents an error that wraps new untracked errors.
type ErrorWrapper interface {
	error

	// Because returns a copy of the receiving error constructing a cause from
	// msg and args essentially creating and wrapping a causal error.
	//
	// Put another way, a call to errors.Unwrap using this functions returned
	// error as input should yeild an underlying error with the supplied error
	// msg and args.
	Because(msg string, args ...any) error
}

var (
	globalRealm       IntRealm
	globalInitialised bool
)

// Initialised causes all future calls to Track and Checkpoint to panic.
//
// When called within an init function in the main package, it prevents
// creation of trackable errors after initialisation. Trackable errors created
// in after initialisation are almost useless.
//
//		package main
//
//		import "github.com/PaulioRandall/go-trackerr"
//
//		func init() {
//			trackerr.Initialised()
//		}
func Initialised() {
	globalInitialised = true
}

// Untracked returns a new error without a tracking ID.
//
// This is no different from errors.New except for the handy fmt.Sprintf
// signiture and a few extra receiving functions for various use cases.
func Untracked(msg string, args ...any) *UntrackedError {
	return because(msg, args...)
}

// Wrap returns a new untracked error that wraps a cause.
//
//		cause := trackerr.Untracked("cause message")
//
//		error := trackerr.Wrap(cause, "wrapper message")
//
//		// wrapper message
//		// ⤷ cause message
func Wrap(cause error, msg string, args ...any) *UntrackedError {
	return because(msg, args...).Wrap(cause).(*UntrackedError)
}

// Track returns a new tracked error from this package's global Realm.
//
// This is recommended way to use to create all trackable errors outside of
// testing.
func Track(msg string, args ...any) *TrackedError {
	checkInitState()
	return globalRealm.Track(msg, args...)
}

// Checkpoint returns a new trackable checkpoint error from this package's
// global Realm.
//
// This is recommended way to use to create all checkpoint errors outside of
// testing.
func Checkpoint(msg string, args ...any) *TrackedError {
	checkInitState()
	return globalRealm.Checkpoint(msg, args...)
}

func checkInitState() {
	if globalInitialised {
		panic(Untracked("No tracked errors may be created after initialisation. Initialise a package variable using this function instead."))
	}
}
