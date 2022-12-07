package track

import (
	//"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IntRealm_untracked_1(t *testing.T) {
	r := IntRealm{}

	act := r.Untracked("Abc%d%d%d", 1, 2, 3)
	exp := &untrackedError{
		msg:   "Abc123",
		cause: nil,
	}

	require.Equal(t, exp, act)
}

func Test_IntRealm_Error_1(t *testing.T) {
	r := IntRealm{}

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

func Test_IntRealm_Error_2(t *testing.T) {
	r := IntRealm{}

	_ = r.Error("Abc%d%d%d", 1, 2, 3)
	act := r.Error("Efg%d%d%d", 4, 5, 6)

	exp := &trackedError{
		untrackedError: &untrackedError{
			msg:   "Efg456",
			cause: nil,
		},
		id: 2,
	}

	require.Equal(t, exp, act)
}

func Test_IntRealm_Checkpoint_1(t *testing.T) {
	r := IntRealm{}

	act := r.Checkpoint("Abc%d%d%d", 1, 2, 3)
	exp := &checkpointError{
		untrackedError: &untrackedError{
			msg:   "Abc123",
			cause: nil,
		},
		id: 1,
	}

	require.Equal(t, exp, act)
}
