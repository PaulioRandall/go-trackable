package trackerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TrackedError_1(t *testing.T) {
	a := &TrackedError{
		id:             1,
		UntrackedError: *Wrap(errors.New("Root cause"), "abc"),
	}

	b := &TrackedError{
		id:             1,
		UntrackedError: *Untracked("efg"),
	}

	require.True(t, a.Is(b))
}

func Test_TrackedError_2(t *testing.T) {
	a := &TrackedError{
		id:             1,
		UntrackedError: *Wrap(errors.New("Root cause"), "abc"),
	}

	b := &TrackedError{
		id:             2,
		UntrackedError: *Wrap(errors.New("Root cause"), "abc"),
	}

	require.False(t, a.Is(b))
}
