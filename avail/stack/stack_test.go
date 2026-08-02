package stack

import (
	"fmt"
	"testing"

	"github.com/Grant-Nelson/Gozer/avail/iterator"
)

func TestStack_New(t *testing.T) {
	s := New[int]()
	equal(t, 0, s.Count(), `count`)
	equal(t, 0, s.Capacity(), `capacity`)
	equal(t, true, s.Empty(), `empty`)
	equal(t, 0, s.Peek(), `peek`)
	equal(t, 0, s.Pop(), `pop`)
	equal(t, ``, s.Iterate().Join(`, `), `iterate`)
}

func TestStack_PushOne(t *testing.T) {
	s := New[int]()
	s.PushOne(19)
	equal(t, 1, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 19, s.Peek(), `peek`)
	equal(t, `19`, s.Iterate().Join(`, `), `iterate`)

	s.PushOne(28)
	equal(t, 2, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 28, s.Peek(), `peek`)
	equal(t, `28, 19`, s.Iterate().Join(`, `), `iterate`)

	s.PushOne(37)
	s.PushOne(46)
	s.PushOne(55)
	s.PushOne(64)
	s.PushOne(73)
	s.PushOne(82)
	equal(t, 8, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 82, s.Peek(), `peek`)
	equal(t, `82, 73, 64, 55, 46, 37, 28, 19`, s.Iterate().Join(`, `), `iterate`)

	s.PushOne(91)
	equal(t, 9, s.Count(), `count`)
	equal(t, 16, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 91, s.Peek(), `peek`)
	equal(t, `91, 82, 73, 64, 55, 46, 37, 28, 19`, s.Iterate().Join(`, `), `iterate`)
}

func TestStack_Push(t *testing.T) {
	s := New[int]()
	s.Push(11, 22, 33)
	equal(t, 3, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 11, s.Peek(), `peek`)
	equal(t, `11, 22, 33`, s.Iterate().Join(`, `), `iterate`)

	equal(t, 11, s.Pop(), `pop`)
	equal(t, 2, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 22, s.Peek(), `peek`)
	equal(t, `22, 33`, s.Iterate().Join(`, `), `iterate`)

	s.Trim()
	equal(t, 2, s.Count(), `count`)
	equal(t, 2, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 22, s.Peek(), `peek`)
	equal(t, `22, 33`, s.Iterate().Join(`, `), `iterate`)

	equal(t, 22, s.Pop(), `pop`)
	equal(t, 1, s.Count(), `count`)
	equal(t, 2, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 33, s.Peek(), `peek`)
	equal(t, `33`, s.Iterate().Join(`, `), `iterate`)

	equal(t, 33, s.Pop(), `pop`)
	equal(t, 0, s.Count(), `count`)
	equal(t, 2, s.Capacity(), `capacity`)
	equal(t, true, s.Empty(), `empty`)
	equal(t, 0, s.Peek(), `peek`)
	equal(t, ``, s.Iterate().Join(`, `), `iterate`)

	s.Push(44, 55, 66, 77, 88)
	equal(t, 5, s.Count(), `count`)
	equal(t, 10, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 44, s.Peek(), `peek`)
	equal(t, `44, 55, 66, 77, 88`, s.Iterate().Join(`, `), `iterate`)
}

func TestStack_PushSeq(t *testing.T) {
	s := New[int]()
	s.PushSeq(iterator.Iterate(11, 22, 33), 6)

}

func equal[T comparable](t *testing.T, want, got T, format string, args ...any) {
	if want != got {
		t.Errorf("Expected values to be equal:\n\twant: %v\n\tgot:  %v\n\tmsg:  %s\n",
			want, got, fmt.Sprintf(format, args...))
	}
}
