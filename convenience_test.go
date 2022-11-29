package trackable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_convenience_1(t *testing.T) {
	require.True(t, IsTracked(errReading))
}

func Test_convenience_2(t *testing.T) {
	e := Untracked("Failed to read file")
	require.False(t, IsTracked(e))
}

func Test_convenience_3(t *testing.T) {
	e := mockWorkflow()
	require.True(t, All(e, errWorkflow, errLoading, errReading))
}

func Test_convenience_4(t *testing.T) {
	e := mockDataLoad()
	require.True(t, All(e, errLoading, errReading))
	require.False(t, All(e, errWorkflow, errLoading, errReading))
}

func Test_convenience_5(t *testing.T) {
	e := mockWorkflow()
	require.True(t, Any(e, errWorkflow, errLoading, errReading))
}

func Test_convenience_6(t *testing.T) {
	e := mockFileRead()
	require.False(t, Any(e, errWorkflow, errLoading))
	require.True(t, Any(e, errWorkflow, errLoading, errReading))
}
