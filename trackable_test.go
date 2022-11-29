package trackable

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	errWorkflow = Track("Failed to run workflow")
	errLoading  = Track("Failed to load data")
	errReading  = Track("Failed to read file")
)

func mockFileRead() error {
	return errReading.Because("There is no file")
}

func mockDataLoad() error {
	e := mockFileRead()
	return errLoading.BecauseOf(e, "Some file reading issue")
}

func mockWorkflow() error {
	e := mockDataLoad()
	return errWorkflow.Wrap(e)
}

func Test_trackable_1(t *testing.T) {
	require.True(t, errors.Is(errReading, errReading))
}

func Test_trackable_2(t *testing.T) {
	require.False(t, errors.Is(errLoading, errReading))
}

func Test_trackable_3(t *testing.T) {
	e := Wrap(errReading, "Failed to load data")

	require.True(t, errors.Is(e, errReading))
	require.False(t, errors.Is(errReading, e))
}

func Test_trackable_4(t *testing.T) {
	e := errReading.Because("There is no file")
	require.True(t, errors.Is(e, errReading))
}

func Test_trackable_5(t *testing.T) {
	e := errLoading.BecauseOf(errReading, "Some file reading issue")

	require.True(t, errors.Is(e, errLoading))
	require.True(t, errors.Is(e, errReading))
}

func Test_trackable_6(t *testing.T) {
	e := errLoading.Wrap(errReading)

	require.True(t, errors.Is(e, errLoading))
	require.True(t, errors.Is(e, errReading))
}

func Test_trackable_7(t *testing.T) {
	e := mockWorkflow()

	require.True(t, errors.Is(e, errWorkflow))
	require.True(t, errors.Is(e, errLoading))
	require.True(t, errors.Is(e, errReading))
}
