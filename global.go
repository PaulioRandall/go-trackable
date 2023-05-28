package trackerr

var (
	globalRealm       IntRealm
	globalInitialised bool
)

// Initialised causes all future calls to New or Track to panic.
//
// When called from an init function in the main package, it prevents creation
// of trackable errors after program initialisation.
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

func checkInitState() {
	if globalInitialised {
		panic(Untracked("No tracked errors may be created after initialisation."))
	}
}
