package trackerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UntrackedError_1(t *testing.T) {
	rootCause := errors.New("Root cause")

	given := &UntrackedError{
		msg: "abc",
	}

	act := given.CausedBy(rootCause)

	exp := &UntrackedError{
		msg:   "abc",
		cause: rootCause,
	}

	require.Equal(t, exp, act)
}

func Test_UntrackedError_2(t *testing.T) {
	given := &UntrackedError{
		msg: "abc",
	}

	act := given.Because("%d%d%d", 1, 2, 3)

	exp := &UntrackedError{
		msg: "abc",
		cause: &UntrackedError{
			msg: "123",
		},
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_3(t *testing.T) {
	given := &UntrackedError{
		msg: "abc",
	}

	cause := &UntrackedError{
		msg: "efg",
	}

	act := given.BecauseOf(cause, "xyz")

	exp := &UntrackedError{
		msg: "abc",
		cause: &UntrackedError{
			msg: "efg",
			cause: &UntrackedError{
				msg: "xyz",
			},
		},
	}

	require.Equal(t, exp, act)
}
