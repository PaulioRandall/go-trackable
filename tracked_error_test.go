package trackerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TrackedError_1(t *testing.T) {
	a := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: errors.New("Root cause"),
	}

	b := &TrackedError{
		id:    1,
		msg:   "efg",
		cause: errors.New("Root cause"),
	}

	require.True(t, a.Is(b))
}

func Test_TrackedError_2(t *testing.T) {
	a := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: errors.New("Root cause"),
	}

	b := &TrackedError{
		id:    2,
		msg:   "abc",
		cause: errors.New("Root cause"),
	}

	require.False(t, a.Is(b))
}
