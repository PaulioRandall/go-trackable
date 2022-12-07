package track

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_trackedError_Is_1(t *testing.T) {
	a := &trackedError{
		id:    1,
		msg:   "abc",
		cause: errors.New("Root cause"),
	}

	b := &trackedError{
		id:    1,
		msg:   "efg",
		cause: nil,
	}

	require.True(t, a.Is(b))
}

func Test_trackedError_Is_2(t *testing.T) {
	a := &trackedError{
		id:    1,
		msg:   "abc",
		cause: errors.New("Root cause"),
	}

	b := &trackedError{
		id:    2,
		msg:   "abc",
		cause: errors.New("Root cause"),
	}

	require.False(t, a.Is(b))
}

func Test_trackedError_Wrap_1(t *testing.T) {
	rootCause := errors.New("Root cause")

	given := &trackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.Wrap(rootCause)

	exp := &trackedError{
		id:    given.id,
		msg:   given.msg,
		cause: rootCause,
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_Because_1(t *testing.T) {
	given := &trackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.Because("%d%d%d", 1, 2, 3)

	exp := &trackedError{
		id:  given.id,
		msg: given.msg,
		cause: &untrackedError{
			msg:   "123",
			cause: nil,
		},
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_BecauseOf_1(t *testing.T) {
	rootCause := errors.New("Root cause")

	given := &trackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.BecauseOf(rootCause, "%d%d%d", 1, 2, 3)

	exp := &trackedError{
		id:  given.id,
		msg: given.msg,
		cause: &untrackedError{
			msg:   "123",
			cause: rootCause,
		},
	}

	require.Equal(t, exp, act)
}

func Test_trackedError_Checkpoint_1(t *testing.T) {
	rootCause := errors.New("Root cause")

	given := &trackedError{
		id:    1,
		msg:   "abc",
		cause: nil,
	}

	act := given.Checkpoint(rootCause, "%d%d%d", 1, 2, 3)

	exp := &trackedError{
		id:  given.id,
		msg: given.msg,
		cause: &trackedError{
			isCheckpoint: true,
			msg:          "123",
			cause:        rootCause,
		},
	}

	require.Equal(t, exp, act)
}
