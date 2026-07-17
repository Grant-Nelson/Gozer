package stack

import (
	"github.com/Grant-Nelson/Gozer/avail/iterator"
)

type Stack[T any] interface {
	Count() int
	Empty() bool
	Iterate() iterator.Iterator[T]
	Pop() T
	PushOne(value T) Stack[T]
	// Push will put these values into the stack in reverse order so that
	// the next pop will get the first of these values first.
	Push(values ...T) Stack[T]
	PushSeq(s iterator.Iterator[T], count int) Stack[T]
}

type node[T any] struct {
	value T
	prev  *node[T]
}

type stackImp[T any] struct {
	zero  T
	count int
	top   *node[T]
	tombs *node[T]
}

func New[T any]() Stack[T] {
	return &stackImp[T]{}
}

func NewWithCap[T any](count int) Stack[T] {
	return &stackImp[T]{
		tombs: allocateNodes[T](count, nil),
	}
}

func (s *stackImp[T]) Count() int { return s.count }

func (s *stackImp[T]) Empty() bool { return s.top == nil }

func (s *stackImp[T]) Iterate() iterator.Iterator[T] {
	return func(yield func(T) bool) {
		for n := s.top; n != nil; n = n.prev {
			if !yield(n.value) {
				return
			}
		}
	}
}

func set[T any](n *node[T], value T, prev *node[T]) *node[T] {
	n.value = value
	n.prev = prev
	return n
}

func (s *stackImp[T]) Pop() T {
	if s.top == nil {
		return s.zero
	}
	n := s.top
	value := n.value
	s.top = n.prev
	s.count--
	s.tombs = set(n, s.zero, s.tombs)
	return value
}

func (s *stackImp[T]) PushOne(value T) Stack[T] {
	if s.tombs != nil {
		n := s.tombs
		s.tombs = n.prev
		s.top = set(n, value, s.top)
	} else {
		s.top = &node[T]{value: value, prev: s.top}
	}
	s.count++
	return s
}

func (s *stackImp[T]) getTombs(count int, prev *node[T]) (*node[T], int) {
	if s.tombs == nil || count == 0 {
		return prev, 0
	}
	top := s.tombs
	bot := s.tombs
	for i := range count {
		if bot.prev == nil {
			bot.prev = prev
			s.tombs = nil
			return top, i + 1
		}
		bot = bot.prev
	}
	s.tombs = bot.prev
	bot.prev = prev
	return top, count
}

func allocateNodes[T any](count int, prev *node[T]) *node[T] {
	if count == 0 {
		return prev
	}
	nodes := make([]*node[T], count)
	nodes[0].prev = prev
	for i := 1; i < count; i++ {
		nodes[i].prev = nodes[i-1]
	}
	return nodes[count-1]
}

func (s *stackImp[T]) prepareNodes(count int, prev *node[T]) *node[T] {
	top, found := s.getTombs(count, prev)
	if found < count {
		top = allocateNodes(count-found, top)
	}
	return top
}

func (s *stackImp[T]) Push(values ...T) Stack[T] {
	count := len(values)
	if count <= 0 {
		return s
	}
	s.top = s.prepareNodes(count, s.top)
	cur := s.top
	for _, v := range values {
		cur.value = v
		cur = cur.prev
	}
	return s
}

func (s *stackImp[T]) PushSeq(it iterator.Iterator[T], count int) Stack[T] {
	if count <= 0 {
		return s
	}
	s.top = s.prepareNodes(count, s.top)
	cur := s.top
	for v := range it {
		if count <= 0 {
			return s
		}
		count--
		cur.value = v
		cur = cur.prev
		if cur == nil {
			return s
		}
	}
	return s
}
