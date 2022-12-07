package track

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IntRealm_untracked_1(t *testing.T) {
	r := IntRealm{}

	act := r.Untracked("abc%d%d%d", 1, 2, 3)
	exp := &untrackedError{
		msg:   "abc123",
		cause: nil,
	}

	require.Equal(t, exp, act)
}

func Test_IntRealm_Error_1(t *testing.T) {
	r := IntRealm{}

	act := r.Error("abc%d%d%d", 1, 2, 3)
	exp := &trackedError{
		id:    1,
		msg:   "abc123",
		cause: nil,
	}

	require.Equal(t, exp, act)
}

func Test_IntRealm_Error_2(t *testing.T) {
	r := IntRealm{}

	_ = r.Error("abc%d%d%d", 1, 2, 3)
	act := r.Error("efg%d%d%d", 4, 5, 6)

	exp := &trackedError{
		id:    2,
		msg:   "efg456",
		cause: nil,
	}

	require.Equal(t, exp, act)
}

func Test_IntRealm_Checkpoint_1(t *testing.T) {
	r := IntRealm{}

	act := r.Checkpoint("abc%d%d%d", 1, 2, 3)
	exp := &checkpointError{
		id:    1,
		msg:   "abc123",
		cause: nil,
	}

	require.Equal(t, exp, act)
}
