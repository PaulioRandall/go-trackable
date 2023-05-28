package trackerr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsTracked_1(t *testing.T) {
	a := New("a")
	require.True(t, IsTracked(a))

	b := Untracked("b")
	require.False(t, IsTracked(b))
}

func Test_HasTracked_1(t *testing.T) {
	a := Untracked("a")
	b := Untracked("b")
	c := New("c")
	d := Untracked("d")

	e := a.CausedBy(b.CausedBy(c))
	require.True(t, HasTracked(e))

	e = a.CausedBy(b.CausedBy(d))
	require.False(t, HasTracked(e))
}

func Test_HasTracked_2(t *testing.T) {
	a := Untracked("a")
	b := Untracked("b")
	c := Untracked("c")

	e := a.CausedBy(b.CausedBy(c))

	require.False(t, HasTracked(e))
}

func Test_All_1(t *testing.T) {
	a := New("a")
	b := New("b")
	c := New("c")
	d := New("d")

	e := a.CausedBy(b.CausedBy(c))

	require.True(t, All(e))
	require.True(t, All(e, a, b, c))

	require.False(t, All(e, d))
	require.False(t, All(e, a, b, d))
}

func Test_Allordered_1(t *testing.T) {
	a := New("a")
	b := New("b")
	c := New("c")
	d := New("d")

	e := a.CausedBy(b.CausedBy(c))

	require.True(t, AllOrdered(e))

	require.True(t, AllOrdered(e, a, b))
	require.False(t, AllOrdered(e, b, a))

	require.True(t, AllOrdered(e, a, b, c))
	require.False(t, AllOrdered(e, a, c, b))
	require.False(t, AllOrdered(e, a, b, c, d))
}

func Test_Any_1(t *testing.T) {
	a := New("a")
	b := New("b")
	c := New("c")

	x := New("x")
	y := New("y")
	z := New("z")

	e := a.CausedBy(b.CausedBy(c))

	require.True(t, Any(e, a))
	require.True(t, Any(e, a, b))
	require.True(t, Any(e, a, b, c))
	require.True(t, Any(e, b, c))
	require.True(t, Any(e, c))

	require.False(t, Any(e))
	require.False(t, Any(e, x, y, z))
}
