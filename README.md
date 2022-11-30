# Go Trackable

Go Trackable is a library for creating trackable, traceable, and comparable errors.

I hope the code speaks mostly for itself so you don't have to trawl through my ramblings.

## Importing

```go
import "github.com/PaulioRandall/go-trackable"
```

## Usage

The trackable errors returned by `trackable.Track` have a unique internal ID which is used for comparison by `errors.Is` or `trackable.Is`. The other error struct fields are irrelevant for such comparisons.

It's important to define these errors as global or you won't be able to reference them. I'll talk about the `trackable.Interface` function a little later but it in the contexts of error tracking it's the same as `trackable.Track`.

```go
// Global variables
var ErrReadingCSV = trackable.Track("Failed to read CSV file")
var ErrReadingCSV = trackable.Interface("Failed to read CSV file")
```

When we want to track an error we have several options. Here is the full interface for errors returned by `trackable.Track` (and `trackable.Interface` for that matter). I don't expect this interface to be used much, if at all. But as an interface first software engineer I find them a great reference.

```go
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
  
  // Interface does the same as Because except the trackable error is marked
  // as being at the boundary of a key interface.
  //
  // This allows stack traces to be partitioned so they are more meaningful,
  // readable, and navigable.
  Interface(msg string, args ...any) error
  
  // InterfaceOf does the same as BecauseOf except the trackable error is marked
  // as being at the boundary of a key interface.
  //
  // This allows stack traces to be partitioned so they are more meaningful,
  // readable, and navigable.
  InterfaceOf(cause error, msg string, args ...any) error
  
  // IsInterface returns true if the trackable error was created at the site
  // of a key interface.
  IsInterface() bool
}
```

`Unwrap` and `Is` are receiving functions that work with Go's standard `errors` package. `Interface` and `IsInterface` are described later and are geared towards helping to create meaningful and navigable stack traces.

`Wrap`, `Because`, and `BecauseOf` are the ones we are interested in first.

### .Wrap

Wrapping is straight forward. `e` will be wrapped by a **COPY** of `ErrReadingCSV`. All these functions return copies of themselves so pointer comparisons will not work. Use the `Is` receiving function if a comparison is needed.

```go
func ReadCSV(filename string) error {
  _, e := os.Open(filename)
  if e != nil {
    return ErrReadingCSV.Wrap(e)
  }
  // ...
}

// Resultant stack trace:
//   Failed to read CSV file
// ⤷ <the wrapped error's message>
```

### .Because

We can create our own root cause. The `fmt.Sprintf` interface is used.

```go
func ReadCSV(filename string) error {
  if !isValidCSVFile(filename) {
    return ErrReadingCSV.Because("%q is not a valid CSV file", filename)
  }
  // ...
}

// Resultant stack trace:
//   Failed to read CSV file
// ⤷ '<filename>' is not a valid CSV file
```

### .BecauseOf

We also have a convenience function which wraps the cause `e` in a new error which is then wrapped by `ErrReadingCSV`. This is useful when the underlying cause does not or cannot provide enough relevant details for debugging. The `fmt.Sprintf` interface is used again.

```go
func ReadCSV(filename string) error {
  _, e := os.Open(filename)
  if e != nil {
    return ErrReadingCSV.BecauseOf(e, "Could not open %q", filename)
  }
  // ...
}

// Resultant stack trace:
//   Failed to read CSV file
// ⤷ Could not open '<filename>'
// ⤷ <the wrapped error's message>
```

### Testing

One place tracking becomes useful is when asserting errors in tests. Many programmers compare error messages but those messages are for humans and checking them in tests makes changing them more difficult and one wrong character can screw you over. I really don't like that. Using trackable errors means the messages can be freely changed and updated for human readers without screwing up tests.

```go
import (
  "errors"
  "testing"
)

func TestReadingCSV(t *testing.T) {
  e := ReadCSV("/bad/file/path")
  
  if !errors.Is(e, ErrReadingCSV) {
  t.Log("Expected reading CSV error but either got no error or a different error")
    t.Fail()
  }
}
```

### .Interface

The `Interface` and `InterfaceOf` receiving functions have the same signatures as `Because` and `BecauseOf` but flags the error as being at a key interface boundary or checkpoint. It may be used to indicate when an error has been returned from a call to another package.

It's a little nuanced but when printing the stack trace we highlight these interface error messages to indicate where the key checkpoints or interface boundaries are.

```go
var (
  ErrDoingThing = trackable.Track("Failed to do the thing")
  ErrDelegating = trackable.Track("Delegation returned an error")
)

func doThing() error {
  if e := delegateDoingTheThing(); e != nil {
    return ErrDoingThing.Wrap(e)
  }
  return nil
}

func delegateDoingTheThing() error {
  if e := UnhappyAPI(); e != nil {
    return ErrDelegating.InterfaceOf(e, "The Unhappy API returned an error")
  }
  return nil
}

func UnhappyAPI() error {
  e := errors.New("UnhappyAPI error root cause")
  e = fmt.Errorf("UnhappyAPI error that wraps the cause: %w", e)
  return fmt.Errorf("UnhappyAPI error wrapping at the package boundary: %w", e)
}

// Resultant stack trace:
//   Failed to do the thing
// ⤷ Delegation returned an error
// ⊖ The Unhappy API returned an error
// ⤷ UnhappyAPI error wrapping at the package boundary
// ⤷ UnhappyAPI error that wraps the cause
// ⤷ UnhappyAPI error root cause
```

## Checking out (in both senses)

1. Clone repo

```bash
git clone https://github.com/PaulioRandall/go-trackable.git
```

2. Enter repo

```bash
cd go-trackable
```

3. Go commands can be used from here but my ./godo script makes things easier. To see usage:

```bash
./godo
```
