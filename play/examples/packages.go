package examples

import (
	"errors"
	"fmt"

	"github.com/PaulioRandall/trackable"
)

var (
	ErrExePackage = trackable.Interface("Failed to execuate packages")
	ErrDoingThing = trackable.New("Failed to do a thing")
	ErrCallingAPI = trackable.New("API returned an error")
)

func Packages() {
	e := packages()

	fmt.Printf("\nExample stack trace...\n\n")
	trackable.Debug(e)
}

func packages() error {
	if e := doThing(); e != nil {
		return ErrExePackage.BecauseOf(e, "Could not do that thing")
	}

	return nil
}

func doThing() error {
	if e := useOtherPackage(); e != nil {
		return ErrDoingThing.Wrap(e)
	}

	return nil
}

func useOtherPackage() error {
	if e := unhappyAPI(); e != nil {
		return ErrCallingAPI.Interface(e, "UnhappyAPI returned an error")
	}

	return nil
}

func unhappyAPI() error {
	e := errors.New("This is the root cause")
	e = fmt.Errorf("This is the parent cause: %w", e)
	e = fmt.Errorf("This is the grand parent cause: %w", e)
	return fmt.Errorf("This is the error wrapped at the API boundary: %w", e)
}
