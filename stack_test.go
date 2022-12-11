package trackerr

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrorStack_1(t *testing.T) {
	r := IntRealm{}

	te := r.Track("abc")
	cp := r.Checkpoint("hij")

	given := cp.Because("klm")
	given = te.CausedBy(given, "efg")

	act := ErrorStack(given)

	expLines := []string{
		"abc",
		"⤷ efg",
		"——hij——",
		"⤷ klm",
		"",
	}
	exp := strings.Join(expLines, "\n")

	require.Equal(t, exp, act)
}

func Test_AsStack_1(t *testing.T) {
	klm := Untracked("klm")
	hij := &TrackedError{id: 2, UntrackedError: *Wrap(klm, "hij")}
	efg := Wrap(hij, "efg")
	abc := &TrackedError{id: 1, UntrackedError: *Wrap(efg, "abc")}

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
