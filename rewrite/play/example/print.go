package example

import (
	"github.com/PaulioRandall/trackable/rewrite"

	"github.com/PaulioRandall/trackable/rewrite/play/example/format"
)

func printData(data []format.PhReading) error {
	panic(track.ErrTodo.Because("printData() { ??? }"))
}
