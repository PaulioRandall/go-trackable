package example

import (
	"github.com/PaulioRandall/trackable/rewrite"

	"github.com/PaulioRandall/trackable/rewrite/play/example/clean"
	"github.com/PaulioRandall/trackable/rewrite/play/example/read"
)

var ErrExeWorkflow = track.Error("Executing workflow failed")

// Run provides example usage of go-track. Try breaking some of the parameters
// and logic to see the error stack traces that get produced.
//
// It's verbose in terms of errors on purpose to show off the various features.
// In reality you only want to put in the error handling you need.
func Run() {
	e := executeWorkflow("rewrite/play/example/data/acid-rain.csv")
	if e != nil {
		track.Debug(e)
	}
}

func executeWorkflow(filename string) error {
	data, e := read.CSV(filename)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	data, e = clean.CSV(data)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	if e = printData(data); e != nil {
		return ErrExeWorkflow.BecauseOf(e, "Failed to print out data")
	}

	return nil
}
