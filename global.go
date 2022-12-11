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
func Initialised() {
	globalInitialised = true
}

// Untracked returns a new error without a tracking ID.
//
// It is no different than using errors.New except it has a handy fmt.Sprintf
// signiture and a few extra receiving functions for any niche use cases one
// may encounter.
func Untracked(msg string, args ...any) *untrackedError {
	return because(msg, args...)
}

// Wrap returns a new untracked error that wraps a cause.
func Wrap(cause error, msg string, args ...any) *untrackedError {
	return because(msg, args...).Wrap(cause).(*untrackedError)
}

// Track returns a new tracked error from this package's global Realm.
//
// This is recommended way to use to create all trackable errors outside of
// testing.
func Track(msg string, args ...any) *trackedError {
	checkInitState()
	return globalRealm.Track(msg, args...)
}

// Checkpoint returns a new trackable checkpoint error from this package's
// global Realm.
//
// This is recommended way to use to create all checkpoint errors outside of
// testing.
func Checkpoint(msg string, args ...any) *trackedError {
	checkInitState()
	return globalRealm.Checkpoint(msg, args...)
}

func checkInitState() {
	if globalInitialised {
		panic(Untracked("No tracked errors may be created after initialisation. Initialise a package variable using this function instead."))
	}
}
