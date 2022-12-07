package read

import (
	"encoding/csv"
	"os"

	"github.com/PaulioRandall/trackable/rewrite"
)

var (
	ErrReadPkg    = track.Checkpoint("play/example/read package")
	ErrReadingCSV = track.Error("Error handling CSV file")
)

func CSV(filename string) ([][]string, error) {
	data, e := openAndReadCSV(filename)
	if e != nil {
		return nil, ErrReadPkg.Wrap(e)
	}
	return data, nil
}

func openAndReadCSV(filename string) ([][]string, error) {
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
