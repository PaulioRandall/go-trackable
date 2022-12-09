package trackerr

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	untrackedAlpha   = &untrackedError{msg: "untracked alpha"}
	untrackedBeta    = &untrackedError{msg: "untracked beta"}
	untrackedCharlie = &untrackedError{msg: "untracked charlie"}

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
	e := &untrackedError{
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

func Test_ErrorStack_1(t *testing.T) {
	r := IntRealm{}

	te := r.Track("abc")
	cp := r.Checkpoint("hij")

	given := cp.Because("klm")
	given = te.CausedBy(given, "efg")

	act := ErrorStack(given)

	expLines := []string{
		"  abc",
		"⤷ efg",
		"——hij——",
		"⤷ klm",
		"",
	}
	exp := strings.Join(expLines, "\n")

	require.Equal(t, exp, act)
}

func Test_AsStack_1(t *testing.T) {
	klm := &untrackedError{msg: "abc", cause: nil}
	hij := &trackedError{id: 2, msg: "hij", cause: klm}
	efg := &untrackedError{msg: "efg", cause: hij}
	abc := &trackedError{id: 1, msg: "abc", cause: efg}

	act := AsStack(abc)
	exp := []error{abc, efg, hij, klm}

	require.Equal(t, exp, act)
}

func Test_DebugPanic_1(t *testing.T) {
	given := func() (e error) {
		defer DebugPanic(&e)

		if true {
			panic(trackedAlpha)
		}

		return nil
	}

	e := given()

	require.Equal(t, e, trackedAlpha)
}
