package format

import (
	"time"

	"github.com/PaulioRandall/trackable/rewrite"
)

const parseDayFmt = "13 Jan"

var ErrFormatPkg = track.Checkpoint("play/example/format")

type PhReading struct {
	Date  time.Time
	Value float32
}

func Format(data [][]string) ([]PhReading, error) {
	return nil, ErrFormatPkg.BecauseOf(track.ErrTodo, "Format() not yet implemented")
}

/*
func parseDates(data [][]string) error {
	return track.ErrTodo.Because("parseDates func needs implementation")
}

func parsePhValues(data [][]string) error {
	return track.ErrTodo.Because("parsePhValues func needs implementation")
}
*/
