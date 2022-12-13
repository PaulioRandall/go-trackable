# Trackerr

Trackerr is a library for creating trackable, traceable, and comparable errors.

I crafted this package in frustration trying to decypher Go's printed error stack traces and the challenge of reliably asserting specific error types while testing.

I hope the code speaks mostly for itself so you don't have to trawl through my ramblings.

## Usage

It's important to define errors created via `trackerr.Track` and `trackerr.Checkpoint` as global or you won't be able to reference them.

You can return a tracked error directly but it's recommended to call one of the receiving functions `Wrap`, `CausedBy`, `Because`, `BecauseOf`, or `Checkpoint` with a cause or additional information; sometimes both.

For manual debugging there's `trackerr.Debug` and deferable `trackerr.DebugPanic` which will print a readable stack trace of an error.

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

## Composition > Framing

Many programmers test assert using error messages (strings) but I've found this to be unreliable and leave me feeling less than confident. trackerr attempts to rectify this by assigning errors there own unique identifiers which can be checked using errors.Is or one of trackerr's utility functions.

The package is designed to work in compositional manner such that `trackerr.Track` and `errors.new` can be exchanged incrementally. Engineers may compose all their errors using trackerr or just a handful they want tracking support for. Many of trackerr's utility functions work irrespective of the underlying error types.

Composition is much better for keeping the power to change and adapt in the hands of the programmer. I'm trying to minimise _my way or the highway_ mentality and _vendor lock-in_ in so much as I can. If my package no longer provides adequate value or there is something better then it should be easy to replace or remove.

This approach also helps keep batch sizes small to better support the few teams that use Continuous Integration (CI) and Continuous Delivery (CD).

## Testing

One place trackerr becomes useful is when asserting errors in tests.

### Comparing error messages

Many programmers compare error messages but those messages are written for humans and by validating them in tests they become more cumbersome change. Furthermore, one wrong character can screw you over. I'd much prefer to separate the concerns of communicating with humans and asserting specific errors. Communicating correct and relevant information to human programmers is far harder so I'd like to improve error messages without having to worry about breaking tests. 

## Comparing error pointers

Another problem is pointer equality. Comparing pointers is better than comparing text but this means package scooped errors must be immutable, thus cannot have a cause attached to them or be wrapped. Wrapping trackerr errors produces copies allowing causes to be attached. However, calling `errors.Is(copy, original)` will still return true as internal IDs are used for equality checks and not string messages or pointers.

Unfortunately, this means `copy == original` will return false, but this is not much of a sacrifice since error pointer comparisons lost favour with the introduction of error wrapping ([Go 1.13](https://tip.golang.org/doc/go1.13#error_wrapping)). Use `errors.Is` or one of trackerr's utility functions instead.

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

## Common errors

This package also provides a few common errors that you may want to return or panic with.

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

## Checking out (in both senses)

1. Clone repo

```bash
git clone https://github.com/PaulioRandall/go-trackerr.git
```

2. Enter repo

```bash
cd go-trackerr
```

3. Go commands can be used from here but my `./godo` script eases things:

```bash
./godo [help]   # Print usage
./godo doc[s]   # Fire up documentation server
./godo clean    # Clean Go caches and bin folder
./godo build    # fmt -> build -> test -> vet
./godo play     # fmt -> build -> test -> vet -> play
```
