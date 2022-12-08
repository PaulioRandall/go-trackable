package clean

import (
	"strings"

	"github.com/PaulioRandall/trackable/rewrite"
)

var ErrCleanPkg = track.Checkpoint("play/example/clean")

func Clean(data [][]string) ([][]string, error) {
	trimCellSpaces(data)

	var e error
	if data, e = removeNulls(data); e != nil {
		return nil, ErrCleanPkg.Wrap(e)
	}

	return data, nil
}

func trimCellSpaces(data [][]string) {
	for i := 0; i < len(data); i++ {
		row := data[i]

		for j := 0; j < len(row); j++ {
			row[j] = strings.TrimSpace(row[j])
		}
	}
}

func removeNulls(data [][]string) ([][]string, error) {
	var results [][]string

	for i := 0; i < len(data); i++ {
		row := data[i]

		if phWithoutDate(row) { // Excuse to return an error
			return nil, track.Untracked(
				"pH reading found without a date on line %d", lineNumber(i),
			)
		}

		if isNotNullRecord(row) {
			results = append(results, row)
		}
	}

	return results, nil
}

func phWithoutDate(row []string) bool {
	return row[0] == "" && row[1] != ""
}

func isNotNullRecord(row []string) bool {
	return row[0] != "" && row[1] != ""
}

func lineNumber(i int) int {
	i++ // Convert from index to count
	i++ // Skip the header
	return i
}
