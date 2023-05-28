package trackerr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IntRealm_0(t *testing.T) {
	var _ Realm = &IntRealm{}
}

func Test_IntRealm_1(t *testing.T) {
	r := IntRealm{}

	act := r.Track("abc%d%d%d", 1, 2, 3)
	exp := &TrackedError{
		id:  1,
		msg: "abc123",
	}

	require.Equal(t, exp, act)
}

func Test_IntRealm_2(t *testing.T) {
	r := IntRealm{}

	_ = r.Track("abc%d%d%d", 1, 2, 3)
	act := r.Track("efg%d%d%d", 4, 5, 6)

	exp := &TrackedError{
		id:  2,
		msg: "efg456",
	}

	require.Equal(t, exp, act)
}
