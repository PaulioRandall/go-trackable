package trackerr

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrorStack_1(t *testing.T) {
	r := IntRealm{}

	abc := r.New("abc")
	efg := r.New("efg")
	hij := r.New("hij")
	klm := r.New("klm")

	hij.cause = klm
	efg.cause = hij
	abc.cause = efg

	act := ErrorStack(abc)

	expLines := []string{
		"abc",
		"⤷ efg",
		"⤷ hij",
		"⤷ klm",
		"",
	}
	exp := strings.Join(expLines, "\n")

	require.Equal(t, exp, act)
}

func Test_SliceStack_1(t *testing.T) {
	klm := &UntrackedError{msg: "klm"}
	hij := &TrackedError{id: 2, msg: "hij", cause: klm}
	efg := &UntrackedError{msg: "efg", cause: hij}
	abc := &TrackedError{id: 1, msg: "abc", cause: efg}

	act := SliceStack(abc)
	exp := []error{abc, efg, hij, klm}

	require.Equal(t, exp, act)
}

func Test_DebugPanic_1(t *testing.T) {
	a := New("a")

	given := func() (e error) {
		defer DebugPanic(&e)

		if true {
			panic(a)
		}

		return nil
	}

	e := given()

	require.Equal(t, e, a)
}

func Test_Stack_1(t *testing.T) {
	a := New("a")
	b := New("b")
	c := New("c")

	e := Stack(c, b, a)
	require.True(t, a.Is(e))

	e = Unwrap(e)
	require.True(t, b.Is(e))

	e = Unwrap(e)
	require.True(t, c.Is(e))
}
