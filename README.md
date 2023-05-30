# Trackerr

Package trackerr aims to facilitate creation of referenceable errors and elegant stack traces.

It was crafted in frustration trying to navigate Go's printed error stacks and the challenge of reliably asserting specific error types while testing.

I hope the code speaks mostly for itself so you don't have to trawl through my ramblings.

## API

```go
import (
	// Package imported is just called 'trackerr' 
	"github.com/PaulioRandall/go-trackerr"
)
```

**Please note:** `TrackedError` and `UntrackedError` are structs but I've specified them here as interfaces for documentation purposes.

```go
var (
    // ErrTodo for specifying a TODO.
    ErrTodo = New("TODO: Implementation needed")

    // ErrBug for the site of known bugs
    ErrBug = New("BUG: Fix needed")

    // ErrInsane for sanity checking.
    ErrInsane = New("Sanity check failed!!")
)

func New(msg string, args ...any) TrackedError {}
func Track(msg string, args ...any) TrackedError {}
func Untracked(msg string, args ...any) UntrackedError {}

func All(e error, targets ...error) bool
func AllOrdered(e error, targets ...error) bool
func Any(e error, targets ...error) bool
func HasTracked(e error) bool
func Is(e, target error) bool
func IsTracked(e error) bool
func IsTrackerr(e error) bool
func Unwrap(e error) error

func Stack(errs ...ErrorThatWraps) error
func SliceStack(e error) []error
func Squash(e error) error
func Squashf(e error, f ErrorFormatter) error
func ErrorStack(e error) string
func ErrorStackf(e error, f ErrorFormatter) string
func ErrorWithoutCause(e error) string

func Debug(e error) (int, error)
func DebugPanic(catch *error)

func Initialised()

type ErrorFormatter func(errMsg string, e error, isFirst bool) string

type ErrorThatWraps interface {
	error
	CausedBy(cause error) error
}

// Actually a struct
type TrackedError interface {
	error

	CausedBy(cause error)
	Because(msg string, args ...any) error
	BecauseOf(cause error, msg string, args ...any) error
	ContextFor(cause ErrorThatWraps, rootCause error) error

	Is(error) bool
	Unwrap() error
}

// Actually a struct
type UntrackedError interface {
	error

	CausedBy(cause error)
	Because(msg string, args ...any) error
	BecauseOf(cause error, msg string, args ...any) error
	ContextFor(cause ErrorThatWraps, rootCause error) error

	Unwrap() error
}

type Realm interface {
	New(msg string, args ...any) *TrackedError
	Track(msg string, args ...any) *TrackedError
}

type IntRealm struct {}
```

**Tracked errors should be package variables**

It's important to define errors created via `New` and `Track` as package scooped (global) or you won't be able to reference them. It is not recommended to create trackable errors after initialisation but Realms exist for such cases.

**Wrapping errors**

You can return a tracked or untracked error directly but it's recommended to call one of the receiving functions `CausedBy`, `Because`, `BecauseOf`, or `ContextFor` with additional information.

```go
var (
	ErrLoadingData = trackerr.New("Failed to load data")
	ErrOpeningDatabase = trackerr.New("Could not open database")

	dbFile = "./data/db.sqlite"
)

func Err() error {
	return ErrLoadingData
}

func CausedBy() error {
	return ErrLoadingData.CausedBy(ErrOpeningDatabase)
}

func Because() error {
	return ErrLoadingData.Because("Database file '%s' not found", dbFile)
}

func BecauseOf() error {
	e := trackerr.Untracked("Database file '%s' not found", dbFile)
	return ErrLoadingData.BecauseOf(e, "Could not open database")
}

func ContextFor() error {
	e := trackerr.Untracked("Database file '%s' not found", dbFile)
	return ErrLoadingData.ContextFor(ErrOpeningDatabase, e)
}
```

**Prevent creating tracked errors after program initialisation**

It's also recommended to call `Initialised` from an init function in package main to prevent the creation of trackable errors after program initialisation.

```go
package main

import (
	"github.com/PaulioRandall/go-trackerr"
)

var ErrForNoReason = trackerr.New("Failed for no reason")

func init() {
	trackerr.Initialised()
}

func main() {
	// Bad, will panic
	e = trackerr.New("I felt like it")

	_ = e
}
```

**Debugging**

For manual debugging there's `trackerr.Debug` which will print a readable stack trace.

```go
func Debug() {
	a := trackerr.UntrackedError("Failed to load data")
	b := trackerr.UntrackedError("Could not open database")
	c := trackerr.UntrackedError("Database file not found")

	e := Stack(a, b, c)

	trackerr.Debug(e)

	// [DEBUG ERROR]
	// Failed to load data
	// ⤷ Could not open database
	// ⤷ Database file not found
}
```

Alternatively the deferable `trackerr.DebugPanic(nil)` will recover from a panic, print the error (if it is one), then resume the panic.

```go
func DebugPanic() {
	defer trackerr.DebugPanic(nil)

	a := trackerr.UntrackedError("Failed to load data")
	b := trackerr.UntrackedError("Could not open database")
	c := trackerr.UntrackedError("Database file not found")

	e := Stack(a, b, c)
	panic(e)

	// [DEBUG ERROR]
	// Failed to load data
	// ⤷ Could not open database
	// ⤷ Database file not found
}
```

Passing a pointer to an error `trackerr.DebugPanic(&e)` will prevent the panic resuming and instead set it as the value pointed to by the pointer. 

```go
func DebugPanic() (e error) {
	defer trackerr.DebugPanic(&e)

	...
}
```

**Custom errors**

You may also craft your own error types and wrap or be wrapped by trackerr errors.

```go
type myError struct {
	msg string
	cause error
}

func (e myError) CausedBy(other error) error {
	e.cause = other
	return e
}

func (e myError) Unwrap() error {
	return e.cause
}

var (
	ErrLoadingData = trackerr.New("Failed to load data")
	ErrFileNotFound = trackerr.New("Database file not found")
)

func main() {
	e := myError{ msg: "Could not open database" }
	e = ErrLoadingData.ContextFor(ErrFileNotFound, e)
	_ = e
}
```

### Testing

One place trackerr becomes useful is when asserting errors in tests.

Trackerr assigns errors there own private unique identifiers which are used for comparison by `errors.Is` and trackerr's utility functions. This separates the concerns of communicating with humans from asserting that specific errors occur when they should.

```go
// csvreader.go

import (
	"errors"
)

var ErrParsingCSV = trackerr.New("Could not parse CSV")

func ReadCSV(file string) error {
	...

	return ErrParsingCSV
}
```

```go
// csvreader_test.go

import (
	"errors"
	"testing"
)

func TestReadCSV_InvalidFormat(t *testing.T) {
	e := ReadCSV("/path/to/csv/file")
	
	if !errors.Is(e, ErrParsingCSV) {
		t.Log("Expected ErrParsingCSV error")
		t.Fail()
	}
}
```

## Design decisions

The design is largely usage lead and thus somewhat emergent. That is, I had projects requiring trackable errors to which I crafted structures and functions based on need.

### Composition > Framing

The package is designed to work in a compositional manner such that `trackerr.New`, `trackerr.Track`, and `errors.new` can be exchanged incrementally. Engineers may compose all their errors using trackerr or just the few that require tracking. Most of trackerr's utility functions work on the `error` interface so the underlying error types matter little.

Composition is favoured over framing, when feasible, so the power to change and adapt, with needs and the times, remains in the hands of the consuming engineers. In so much as possible, minimising the _my way or the highway_ mentality which is core to commercial software but also rampant in open source tooling.

If my package no longer provides value for cost or if something better appears then it should be **incrementally** removable or replacable. I find that a good design is one that can change easily. My preference for changability, Continuous Integration (CI), and Continuous Delivery (CD) certainly influenced these decisions.

### Why not string equality?

Many programmers test assert using error messages (strings) but I've found this to be unreliable, reduces changability, and leaves me feeling less than confident in my code; and testing is all about gaining confidence.

Communicating aaccurate and relevant information to humans can be quite a fraught affair so I'd like to maximise the ease of improving and rewriting error messages without having to worry about breaking tests.

### Why not pointer equality?

Comparing pointers is better than comparing text but this means package scooped errors must be immutable, thus cannot have a cause attached to them or be wrapped. The receiving functions of `TrackedError` and `UntrackedError` produce copies of themselves (including their IDs) that allows the attachment of causes while keeping the equality checking. `errors.Is(copy, original)` still returns true as private unique identifiers are compared, not string messages or pointers.

Unfortunately, this means `copy == original` will always return false. This is not much of a sacrifice as error pointer comparisons lost favour with the introduction of error wrapping ([Go 1.13](https://tip.golang.org/doc/go1.13#error_wrapping)). Use `errors.Is`, `trackerr.Is`, or one of trackerr's other utility functions instead.

## Checking out (in both senses)

```bash
git clone https://github.com/PaulioRandall/go-trackerr.git
cd go-trackerr
```

Standard Go commands can be used from here but my `./godo` script eases things:

```bash
./godo [help]   # Print usage
./godo doc[s]   # Fire up documentation server
./godo clean    # Clean Go caches
./godo test     # fmt -> test -> vet
```
