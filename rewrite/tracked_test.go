package track

import (
	//"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_trackedError_1(t *testing.T) {
	r := NewIntRealm()

	act := r.Error("Abc%d%d%d", 1, 2, 3)
	exp := &trackedError{
		untrackedError: &untrackedError{
			msg:   "Abc123",
			cause: nil,
		},
		id: 1,
	}

	require.Equal(t, exp, act)
}
