package example

import (
	"github.com/PaulioRandall/trackable/rewrite"

	"github.com/PaulioRandall/trackable/rewrite/play/example/clean"
	"github.com/PaulioRandall/trackable/rewrite/play/example/format"
	"github.com/PaulioRandall/trackable/rewrite/play/example/read"
)

var ErrExeWorkflow = track.Error("Executing workflow failed")

// Run provides example usage of go-track. Try breaking some of the parameters
// and logic to see the error stack traces that get produced.
//
// It's verbose in terms of errors on purpose to show off the various features.
// In reality you only want to put in the error handling you need.
func Run() {
	defer func() {
		if v := recover(); v != nil {
			track.Debug(v.(error))
		}
	}()

	e := executeWorkflow("rewrite/play/example/data/acid-rain.csv")
	if e != nil {
		track.Debug(e)
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

	fmtData, e := format.Format(data)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	if e = printData(fmtData); e != nil {
		return ErrExeWorkflow.CausedBy(e, "Data print failure")
	}

	return nil
}
