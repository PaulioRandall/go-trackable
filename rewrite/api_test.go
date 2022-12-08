package track

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
	c := &trackedError{
		id:  1,
		msg: "c",
	}

	b := &untrackedError{
		msg:   "b",
		cause: c,
	}

	a := &untrackedError{
		msg:   "a",
		cause: b,
	}

	require.True(t, HasTracked(a))
}

func Test_HasTracked_2(t *testing.T) {
	c := &untrackedError{
		msg: "c",
	}

	b := &untrackedError{
		msg:   "b",
		cause: c,
	}

	a := &untrackedError{
		msg:   "a",
		cause: b,
	}

	require.False(t, HasTracked(a))
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
