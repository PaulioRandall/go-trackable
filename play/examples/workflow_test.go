package examples

import (
	"testing"

	"github.com/PaulioRandall/trackable"

	"github.com/stretchr/testify/require"
)

func Test_workflow_1(t *testing.T) {
	e := workflow("/test-data.csv")

	require.True(t, trackable.Is(e, ErrExeWorkflow))
	require.True(t, trackable.Is(e, ErrLoadingData))
	require.True(t, trackable.Is(e, ErrReadingCSV))

	require.True(t, trackable.Any(e, ErrLoadingData, ErrCleaningData))
	require.True(t, trackable.All(e, ErrExeWorkflow, ErrLoadingData, ErrReadingCSV))

	require.False(t, trackable.Is(e, ErrCleaningData))
	require.False(t, trackable.All(e, ErrExeWorkflow, ErrCleaningData))
}
