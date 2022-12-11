package trackerr

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
