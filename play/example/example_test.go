package example

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/go-trackerr"

	"github.com/PaulioRandall/go-trackerr/play/example/read"
)

func TestExexcuteWorkflow(t *testing.T) {
	badFilename := "play/example/NOT-DATA/acid-rain.csv"
	e := executeWorkflow(badFilename)

	trackerr.Debug(e)

	require.True(t, trackerr.AllOrdered(e,
		ErrExeWorkflow,
		read.ErrReadPkg,
		read.ErrReadingCSV,
	))
}
