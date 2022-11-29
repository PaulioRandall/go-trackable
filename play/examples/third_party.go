package examples

import (
	"errors"
	"fmt"

	"github.com/PaulioRandall/trackable"
)

func ThirdParty() {
	e := thirdParty()

	fmt.Printf("\nExample stack trace...\n\n")
	trackable.Debug(e)
}

func thirdParty() error {
	return otherPackageAPI()
}

func otherPackageAPI() error {
	e := errors.New("This is the root cause")
	e = fmt.Errorf("This is the parent cause: %w", e)
	e = fmt.Errorf("This is the grand parent cause: %w", e)
	return fmt.Errorf("This is the error wrapped at the API boundary: %w", e)
}
