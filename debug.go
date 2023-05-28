package trackerr

import (
	"fmt"
)

// Debug pretty prints the error stack trace to terminal for debugging
// purposes.
//
// If e is nil then a message will be printed indicating so. This function is
// not designed for logging, just day to day manual debugging.
func Debug(e error) (int, error) {
	s := ErrorStack(e)

	if s == "" {
		return fmt.Print("[DEBUG ERROR] nil error")
	}

	return fmt.Print("[DEBUG ERROR]\n  ", s)
}

// DebugPanic recovers from a panic, prints out the error using the Debug
// function, and finally sets it as the catch error's pointer value.
//
// If nil is passed as the catch then the panic continues after printing.
//
// If the panic value is not an error the panic will continue!
//
// This function is not designed for logging, just day to day manual debugging.
func DebugPanic(catch *error) {
	v := recover()

	if v == nil {
		return
	}

	e, ok := v.(error)
	if !ok {
		panic(v)
	}

	Debug(e)

	if catch == nil {
		panic(e)
	}
	*catch = e
}
