package example

import (
	"encoding/csv"
	"os"

	"github.com/PaulioRandall/trackable/rewrite"
)

var ErrReadingCSV = track.Error("Error handling CSV file")

func lineNumber(i int) int {
	i++ // Convert from index to count
	i++ // Skip the header
	return i
}

func readCSV(filename string) ([][]string, error) {
	f, e := os.Open(filename)
	if e != nil {
		return nil, ErrReadingCSV.BecauseOf(e, "File could not be opened %q", filename)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, e := r.ReadAll()
	records = records[1:] // Remove header

	if e != nil {
		return nil, ErrReadingCSV.BecauseOf(e, "File could not be read %q", filename)
	}

	return records, nil
}
