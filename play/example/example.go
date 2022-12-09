package example

import (
	"github.com/PaulioRandall/go-trackerr"

	"github.com/PaulioRandall/go-trackerr/play/example/clean"
	"github.com/PaulioRandall/go-trackerr/play/example/format"
	"github.com/PaulioRandall/go-trackerr/play/example/read"
)

var ErrExeWorkflow = trackerr.Checkpoint("Workflow error")

// Run provides example usage of go-trackerr.
//
// Try breaking some of the parameters, logic, or data to see the error stack
// traces that get produced.
//
// The example is verbose in terms of errors on purpose to show off the various
// features. In real usage I'd recommend maximising relevant information while
// minimising tracked errors.
func Run() {
	defer trackerr.DebugPanic(nil)

	e := executeWorkflow("splay/example/data/acid-rain.csv")
	if e != nil {
		trackerr.Debug(e)
	}
}

func executeWorkflow(filename string) error {
	data, e := read.Read(filename)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	data, e = clean.Clean(data)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	readings, e := format.Format(data)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	printReadings(readings)
	return nil
}
