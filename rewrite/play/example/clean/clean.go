package clean

import (
	"github.com/PaulioRandall/trackable/rewrite"
)

const (
	parseDayFmt = "13 Jan"
)

var (
	ErrCleanPkg = track.Checkpoint("play/example/clean")
)

func CSV(data [][]string) ([][]string, error) {
	e := track.ErrTodo.Because("CSV func not implemented")
	return nil, ErrCleanPkg.Wrap(e)
}

func lineNumber(i int) int {
	i++ // Convert from index to count
	i++ // Skip the header
	return i
}
