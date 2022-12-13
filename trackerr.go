// Package trackerr aims to facilitate creation of referenceable errors and
// elegant stack traces.
//
//		I crafted this package in reponse to my frustration in trying to
//		decypher Go's printed error stack traces and the challenge of reliably
//		asserting specific error types while testing.
//
//		Many programmers assert using error messages but I've found this to be
//		unreliable and leave me less than confident. trackerr attempts to rectify
//		this by assigning tracked errors there own unique identifiers which can
//		be checked using errors.Is or one of trackerr's utility functions.
//
//		Paulio
//
// The recommended way to create errors is via the Track, Checkpoint,
// Untracked, and Wrap package functions. It is not recommended to create
// trackable errors after initialisation but Realms exist for such cases.
//
// It is also recommended to call the Initialised function from an init
// function in package main to prevent creation of trackable errors after
// program initialisation.
package trackerr

// ErrorWrapper represents an error that wraps new untracked errors.
type ErrorWrapper interface {
	// Because returns a copy of the receiving error constructing a cause from
	// msg and args.
	Because(msg string, args ...any) error
}

// TODO: Think about how to integrate file names and line numbers.
// TODO: - How, where, and when to collect them? (reflection)
// TODO: - How to optimise print outs with them?
// TODO: - May have to redesign the Debug function?

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
