# Trackerr

> TODO: This README needs updating!

Package trackerr aims to facilitate creation of referenceable errors and elegant stack traces.

It was crafted in frustration trying to navigate Go's printed error stacks and the challenge of reliably asserting specific error types while testing.

I hope the code speaks mostly for itself so you don't have to trawl through my ramblings.

## Usage

```go
package main

import (
	"os"

	"github.com/PaulioRandall/go-trackerr"
)

var (
	ErrSavingText = trackerr.Checkpoint("Could not save text")
	ErrOpeningFile = trackerr.Track("Fail to open file")
	ErrWritingToFile = trackerr.Track("Failed writing text to file")
)

func init() {
	trackerr.Initialised()
}

func main() {
	filename := "./simpsons.txt" 
	text := "Dental plan. Lisa needs braces."

	var e error
	defer trackerr.DebugPanic(&e)

	if e = saveText(filename, text); e != nil {
		trackerr.Debug(e)
	}

	if trackerr.Any(e, ErrOpeningFile, ErrWritingToFile) {
		log.Println("Did you forget file system permissions again?")
	}
}

func saveText(filename, text string) error {
	// ...

	if e := writeTextToFile(filename, text); e != nil {
		return ErrSavingText.Wrap(e)
	}

	// ...

	return nil
}

func writeTextToFile(filename, text string) error {
	f, e := os.Create(filename)
	if e != nil {
		return ErrOpeningFile.CausedBy(e, "Filename %q could not be created", filename)
	}
	defer f.Close()
	
	if _, e := f.WriteString(text); e != nil {
		return ErrWritingText.Wrap(e)
	}
	return nil
}

// Stack trace if file could not be created:
//
//		——Could not save text——
//		⤷ Failed to open file
//		⤷ Filename '/simpsons.txt' could not be created
//		⤷ <error returned by os.Create>
```

It's important to define errors created via `Track` and `Checkpoint` as package scooped (global) or you won't be able to reference them. It is not recommended to create trackable errors after initialisation but Realms do exist if such use cases appear.

It is also recommended to call `Initialised` from an init function in package main to prevent the creation of trackable errors after program initialisation.

You can return a tracked error directly but it's recommended to call one of the receiving functions `Wrap`, `CausedBy`, `Because`, `BecauseOf`, or `Checkpoint` with additional information.

For manual debugging there's `trackerr.Debug` and the deferable `trackerr.DebugPanic` which will print a readable stack trace.

You may also craft your own error types and wrap or be wrapped by trackerr errors. Any that implement the `ErrorWrapper` interface may be used as an argument to the `Tracked.BecauseOf` and `Untacked.BecauseOf`.

### Testing

One place trackerr becomes useful is when asserting errors in tests.

trackerr assigns errors there own private unique identifiers which are used for comparison by `errors.Is` and trackerr's utility functions. This separates the concerns of communicating with humans from asserting that specific errors occur when they should.

```go
import (
	"errors"
	"testing"
)

func TestReadingCSV(t *testing.T) {
	e := ReadCSV("/bad/file/path")
	
	if !errors.Is(e, ErrReadingCSV) {
		t.Log("Expected CSV read error")
		t.Fail()
	}
}
```

### Common trackable errors

This package also provides a few common trackable errors that you may want to return or panic with.

```go
var (
	// ErrTodo is a convenience tracked error for specifying a TODO.
	//
	// This can be useful if you're taking a stepwise refinement or test driven
	// approach to writing code.
	ErrTodo = Track("TODO: Implementation needed")

	// ErrBug is a convenience tracked error for use at the site of known bugs.
	ErrBug = Track("BUG: Fix needed")

	// ErrInsane is a convenience tracked error for sanity checking.
	ErrInsane = Track("Sanity check failed!!")
)
```

## Design decisions

The design is largely usage lead and thus somewhat emergent. That is, I had projects requiring trackable errors to which I crafted structures and functions based on need.

### Composition > Framing

The package is designed to work in a compositional manner such that `trackerr.Track` and `errors.new` can be exchanged incrementally. Engineers may compose all their errors using trackerr or just the few that require tracking support. Most of trackerr's utility functions work on the `error` interface so the underlying error types matter little.

Composition is favoured over framing, when feasible, so the power to change and adapt, with needs and the times, remains in the hands of the consuming engineers. Essentially minimising the _my way or the highway_ and _vendor lock-in_ mentalities in so much as possible.

If my package no longer provides value for cost or something better pops up then it should be **incrementally** removable or replacable. My preference for changability, Continuous Integration (CI), and Continuous Delivery (CD) certainly influenced this decision.

### Why not string equality?

Many programmers test assert using error messages (strings) but I've found this to be unreliable, reduces changability, and leaves me feeling less than confident. Communicating correct and relevant information to humans can be quite a fraught affair so I'd like to maximise the ease of improving and rewriting error messages without having to worry about breaking tests.

### Why not pointer equality?

Comparing pointers is better than comparing text but this means package scooped errors must be immutable, thus cannot have a cause attached to them or be wrapped. The functions attached to `TrackedError` and `UntrackedError` produce copies of themselves allowing wrapping and the attachment of causes. `errors.Is(copy, original)` still returns true as private unique identifiers are compared, not string messages or pointers.

Unfortunately, this means `copy == original` will always return false. This is not much of a sacrifice as error pointer comparisons lost favour with the introduction of error wrapping ([Go 1.13](https://tip.golang.org/doc/go1.13#error_wrapping)). Use `errors.Is` or one of trackerr's utility functions instead.

## Checking out (in both senses)

1. Clone repo

```bash
git clone https://github.com/PaulioRandall/go-trackerr.git
```

2. Enter repo

```bash
cd go-trackerr
```

3. Standard Go commands can be used from here but my `./godo` script eases things:

```bash
./godo [help]   # Print usage
./godo doc[s]   # Fire up documentation server
./godo clean    # Clean Go caches and bin folder
./godo test     # fmt -> build -> test -> vet
./godo play     # fmt -> build -> test -> vet -> play
```
