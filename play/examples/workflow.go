package examples

import (
	"fmt"

	"github.com/PaulioRandall/trackable"
)

func Workflow() {
	e := workflow("/path/to/file")

	if trackable.Is(e, ErrCleaningData) {
		// ...cleaning issue...
	}

	fmt.Printf("\nExample stack trace...\n\n")
	trackable.Debug(e)
}

var (
	ErrExeWorkflow  = trackable.Track("Failed to execute workflow")
	ErrLoadingData  = trackable.Track("Failed to load data")
	ErrCleaningData = trackable.Track("Failed to clean data")
	ErrReadingCSV   = trackable.Track("Failed to read CSV file")
)

func workflow(filename string) error {
	rows, e := loadData(filename)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	if len(rows) == 0 {
		return ErrExeWorkflow.Because("There are no data rows")
	}

	cleaned, e := cleanRows(rows)
	if e != nil {
		return ErrExeWorkflow.Wrap(e)
	}

	_ = cleaned
	return nil
}

func loadData(filename string) ([][]string, error) {
	rows, e := readCSV(filename)
	if e != nil {
		return nil, ErrLoadingData.Wrap(e)
	}

	return rows, nil
}

func readCSV(filename string) ([][]string, error) {
	// Imagine some CSV file reading here...

	e := trackable.Untracked("Could not find file %q", filename)

	return nil, ErrReadingCSV.Wrap(e)
}

func cleanRows(rows [][]string) ([][]string, error) {
	// Imagine some data cleaning here...

	rowNum := 7
	e := trackable.Untracked(
		"Invalid field value or something..." +
			"I don't know I'm just making this up as I go along",
	)

	return nil, ErrCleaningData.BecauseOf(e, "Bad data in row %d", rowNum)
}
