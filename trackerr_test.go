package trackerr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	untrackedAlpha   = &UntrackedError{msg: "untracked alpha"}
	untrackedBeta    = &UntrackedError{msg: "untracked beta"}
	untrackedCharlie = &UntrackedError{msg: "untracked charlie"}

	trackedAlpha   = &trackedError{id: 1, msg: "tracked alpha"}
	trackedBeta    = &trackedError{id: 2, msg: "tracked beta"}
	trackedCharlie = &trackedError{id: 3, msg: "tracked charlie"}
)

func Test_IsTracked_1(t *testing.T) {
	e := &trackedError{
		id:  1,
		msg: "abc",
	}

	require.True(t, IsTracked(e))
}

func Test_IsTracked_2(t *testing.T) {
	e := &UntrackedError{
		msg: "abc",
	}

	require.False(t, IsTracked(e))
}

func Test_HasTracked_1(t *testing.T) {
	c := trackedCharlie
	b := untrackedBeta.Wrap(c)
	a := untrackedAlpha.Wrap(b)

	require.True(t, HasTracked(a))
}

func Test_HasTracked_2(t *testing.T) {
	c := untrackedCharlie
	b := untrackedBeta.Wrap(c)
	a := untrackedAlpha.Wrap(b)

	require.False(t, HasTracked(a))
}

func Test_All_1(t *testing.T) {
	c := untrackedCharlie
	b := untrackedBeta.Wrap(c)
	a := untrackedAlpha.Wrap(b)

	e := a

	require.True(t, All(e))
	require.True(t, All(e, a, b, c))

	require.False(t, All(e, a, b, trackedCharlie))
}

func Test_Any_1(t *testing.T) {
	c := untrackedCharlie
	b := untrackedBeta.Wrap(c)
	a := untrackedAlpha.Wrap(b)

	e := a

	require.True(t, Any(e, a))
	require.True(t, Any(e, a, b))
	require.True(t, Any(e, a, b, c))
	require.True(t, Any(e, b, c))
	require.True(t, Any(e, c))

	require.False(t, Any(e, trackedAlpha, trackedBeta, trackedCharlie))
}
