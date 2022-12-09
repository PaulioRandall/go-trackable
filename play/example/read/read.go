package read

import (
	"encoding/csv"
	"os"

	"github.com/PaulioRandall/go-trackerr"
)

var (
	ErrReadPkg    = trackerr.Checkpoint("play/example/read package")
	ErrReadingCSV = trackerr.Track("Error handling CSV file")
)

func Read(filename string) ([][]string, error) {
	data, e := openAndReadCSV(filename)
	if e != nil {
		return nil, ErrReadPkg.Wrap(e)
	}
	return data, nil
}

func openAndReadCSV(filename string) ([][]string, error) {
	f, e := os.Open(filename)
	if e != nil {
		return nil, ErrReadingCSV.CausedBy(e, "File could not be opened %q", filename)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, e := r.ReadAll()
	records = records[1:] // Remove header

	if e != nil {
		return nil, ErrReadingCSV.CausedBy(e, "File could not be read %q", filename)
	}

	return records, nil
}
