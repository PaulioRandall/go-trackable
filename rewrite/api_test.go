package track

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrorStack_1(t *testing.T) {
	r := IntRealm{}

	te := r.Error("abc")
	cp := r.Checkpoint("hij")

	given := cp.Because("klm")
	given = te.BecauseOf(given, "efg")

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
