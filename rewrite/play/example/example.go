package example

import (
	"github.com/PaulioRandall/trackable/rewrite"
)

var ErrExeWorkflow = track.Error("Executing workflow failed")

func Run() {
	e := executeWorkflow("/rewrite/play/example/data/acid-rain.csv")
	if e != nil {
		track.Debug(e)
	}
}

func executeWorkflow(filename string) error {
	data, e := readCSV(filename)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	data, e = cleanData(data)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	if e = printData(data); e != nil {
		return ErrExeWorkflow.BecauseOf(e, "Failed to print out data")
	}

	return nil
}
