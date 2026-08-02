package stack

import (
	"fmt"
	"testing"
)

func TestStack_Basics(t *testing.T) {
	s := New[int]()
	equal(t, 0, s.Count(), `initial count`)
	equal(t, 0, s.Capacity(), `initial capacity`)
	equal(t, true, s.Empty(), `initial empty`)
	equal(t, 0, s.Peek(), `initial peek`)
	equal(t, 0, s.Pop(), `initial pop`)

	s.Push(1, 2, 3)
	equal(t, 3, s.Count(), `initial count`)
	equal(t, 8, s.Capacity(), `initial capacity`)
	equal(t, false, s.Empty(), `empty first`)
	equal(t, 1, s.Peek(), `peek first`)

}

func equal[T comparable](t *testing.T, want, got T, format string, args ...any) {
	if want != got {
		t.Errorf("Expected values to be equal:\n\twant: %v\n\tgot:  %v\n\tmsg:  %s\n",
			want, got, fmt.Sprintf(format, args...))
	}
}
