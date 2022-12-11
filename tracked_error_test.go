package trackerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_trackedError_0(t *testing.T) {
	var _ TrackedError = TrackedError{}
}

func Test_trackedError_1(t *testing.T) {
	a := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: errors.New("Root cause"),
	}

	b := &TrackedError{
		id:    1,
		msg:   "efg",
		cause: nil,
	}

	require.True(t, a.Is(b))
}

func Test_trackedError_2(t *testing.T) {
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

func Test_trackedError_3(t *testing.T) {
	rootCause := errors.New("Root cause")

	given := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.Wrap(rootCause)

	exp := &TrackedError{
		id:    given.id,
		msg:   given.msg,
		cause: rootCause,
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_4(t *testing.T) {
	given := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.Because("%d%d%d", 1, 2, 3)

	exp := &TrackedError{
		id:  given.id,
		msg: given.msg,
		cause: &UntrackedError{
			msg:   "123",
			cause: nil,
		},
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_5(t *testing.T) {
	rootCause := errors.New("Root cause")

	given := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.CausedBy(rootCause, "%d%d%d", 1, 2, 3)

	exp := &TrackedError{
		id:  given.id,
		msg: given.msg,
		cause: &UntrackedError{
			msg:   "123",
			cause: rootCause,
		},
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_6(t *testing.T) {
	rootCause := errors.New("Root cause")

	given := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.Checkpoint(rootCause, "%d%d%d", 1, 2, 3)

	exp := &TrackedError{
		id:  given.id,
		msg: given.msg,
		cause: &UntrackedError{
			isCheckpoint: true,
			msg:          "123",
			cause:        rootCause,
		},
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_7(t *testing.T) {
	given := &TrackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	cause := &TrackedError{
		id:    2,
		msg:   "efg",
		cause: nil,
	}

	act := given.BecauseOf(cause, "xyz")

	exp := &TrackedError{
		id:  given.id,
		msg: given.msg,
		cause: &TrackedError{
			id:  cause.id,
			msg: cause.msg,
			cause: &UntrackedError{
				msg: "xyz",
			},
		},
	}

	require.Equal(t, exp, act)
}
