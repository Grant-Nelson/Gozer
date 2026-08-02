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
	validate(t, s)
}

func TestStack_NewWithCap(t *testing.T) {
	s := NewWithCap[int](6)
	equal(t, 0, s.Count(), `count`)
	equal(t, 6, s.Capacity(), `capacity`)
	equal(t, true, s.Empty(), `empty`)
	equal(t, 0, s.Peek(), `peek`)
	equal(t, 0, s.Pop(), `pop`)
	equal(t, ``, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)
}

func TestStack_PushOne(t *testing.T) {
	s := New[int]()
	s.PushOne(19)
	equal(t, 1, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 19, s.Peek(), `peek`)
	equal(t, `19`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	s.PushOne(28)
	equal(t, 2, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 28, s.Peek(), `peek`)
	equal(t, `28, 19`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

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
	validate(t, s)

	s.PushOne(91)
	equal(t, 9, s.Count(), `count`)
	equal(t, 16, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 91, s.Peek(), `peek`)
	equal(t, `91, 82, 73, 64, 55, 46, 37, 28, 19`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)
}

func TestStack_Push(t *testing.T) {
	s := New[int]()
	s.Push(11, 22, 33)
	equal(t, 3, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 11, s.Peek(), `peek`)
	equal(t, `11, 22, 33`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	equal(t, 11, s.Pop(), `pop`)
	equal(t, 2, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 22, s.Peek(), `peek`)
	equal(t, `22, 33`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	s.Trim()
	equal(t, 2, s.Count(), `count`)
	equal(t, 2, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 22, s.Peek(), `peek`)
	equal(t, `22, 33`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	equal(t, 22, s.Pop(), `pop`)
	equal(t, 1, s.Count(), `count`)
	equal(t, 2, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 33, s.Peek(), `peek`)
	equal(t, `33`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	equal(t, 33, s.Pop(), `pop`)
	equal(t, 0, s.Count(), `count`)
	equal(t, 2, s.Capacity(), `capacity`)
	equal(t, true, s.Empty(), `empty`)
	equal(t, 0, s.Peek(), `peek`)
	equal(t, ``, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	s.Push(44, 55, 66, 77, 88)
	equal(t, 5, s.Count(), `count`)
	equal(t, 10, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 44, s.Peek(), `peek`)
	equal(t, `44, 55, 66, 77, 88`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)
}

func TestStack_PushSeq(t *testing.T) {
	s := New[int]()
	s.PushSeq(iterator.Iterate(11, 22, 33), 5)
	equal(t, 3, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 11, s.Peek(), `peek`)
	equal(t, `11, 22, 33`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	s.PushSeq(iterator.Iterate(44, 55, 66, 77), 4)
	equal(t, 7, s.Count(), `count`)
	equal(t, 8, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 44, s.Peek(), `peek`)
	equal(t, `44, 55, 66, 77, 11, 22, 33`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	s.PushSeq(iterator.Iterate(88, 99), 1)
	equal(t, 9, s.Count(), `count`)
	equal(t, 16, s.Capacity(), `capacity`)
	equal(t, false, s.Empty(), `empty`)
	equal(t, 88, s.Peek(), `peek`)
	equal(t, `88, 99, 44, 55, 66, 77, 11, 22, 33`, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	s.Clear()
	equal(t, 0, s.Count(), `count`)
	equal(t, 16, s.Capacity(), `capacity`)
	equal(t, true, s.Empty(), `empty`)
	equal(t, 0, s.Peek(), `peek`)
	equal(t, ``, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)

	s.Grow(19)
	equal(t, 0, s.Count(), `count`)
	equal(t, 24, s.Capacity(), `capacity`)
	equal(t, true, s.Empty(), `empty`)
	equal(t, 0, s.Peek(), `peek`)
	equal(t, ``, s.Iterate().Join(`, `), `iterate`)
	validate(t, s)
}

func equal[T comparable](t *testing.T, want, got T, format string, args ...any) {
	t.Helper()
	if want != got {
		t.Errorf("Expected values to be equal:\n\twant: %v\n\tgot:  %v\n\tmsg:  %s\n",
			want, got, fmt.Sprintf(format, args...))
	}
}

func validate[T comparable](t *testing.T, sv Stack[T]) {
	t.Helper()
	s := sv.(*stackImp[T])

	cur := s.top
	count := 0
	for cur != nil {
		count++
		cur = cur.prev
	}
	equal(t, count, s.count, `count must match top stack`)

	cur = s.tombs
	count = 0
	for cur != nil {
		equal(t, s.zero, cur.value, `tomb node #%d was not zero`, count)
		count++
		cur = cur.prev
	}
	equal(t, count, s.tombCount, `count must match tomb stack`)
}
