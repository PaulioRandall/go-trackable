package clean

import (
	"strings"

	"github.com/PaulioRandall/trackable/rewrite"
)

const (
	parseDayFmt = "13 Jan"
)

var (
	ErrCleanPkg = track.Checkpoint("play/example/clean")
)

func CSV(data [][]string) error {
	trimCellSpaces(data)
	removeNulls(data)

	if e := parseDates(data); e != nil {
		return ErrCleanPkg.Wrap(e)
	}

	if e := parsePhValues(data); e != nil {
		return ErrCleanPkg.Wrap(e)
	}

	return nil
}

func trimCellSpaces(data [][]string) {
	for i := 0; i < len(data); i++ {
		row := data[i]

		for j := 0; j < len(row); j++ {
			row[j] = strings.TrimSpace(row[j])
		}
	}
}

func removeNulls(data [][]string) {
	panic(ErrCleanPkg.BecauseOf(track.ErrTodo, "removeNulls func needs implementation"))
}

func parseDates(data [][]string) error {
	return track.ErrTodo.Because("parseDates func needs implementation")
}

func parsePhValues(data [][]string) error {
	return track.ErrTodo.Because("parsePhValues func needs implementation")
}

func lineNumber(i int) int {
	i++ // Convert from index to count
	i++ // Skip the header
	return i
}
