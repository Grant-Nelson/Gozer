package stack

import "github.com/Grant-Nelson/Gozer/avail/iterator"

// Stack using a singular lined list with tombs for dead nodes
// that can be reused. This stack is designed to work very well
// with depth first traversal of a tree.
type Stack[T any] interface {

	// Count is the number of values in the stack.
	Count() int

	// Empty indicates there are no values in the stack.
	Empty() bool

	// Capacity it the number of values this stack can hold before
	// it needs to allocate more nodes.
	Capacity() int

	// Grow will increase the capacity of the stack until
	// greater than or equal to the give capacity.
	Grow(capacity int) Stack[T]

	// Trim removes any excess capacity.
	Trim() Stack[T]

	// Iterate will walk through the stack's values from the top
	// to the bottom. The iteration may not return all the values
	// correctly if the stack is mutated during iteration.
	Iterate() iterator.Iterator[T]

	// Pop will remove one value from the top of te stack and return it.
	// If the stack is empty, this will return the zero T value.
	Pop() T

	// Peek will return the top value on the stack without removing it.
	// If the stack is empty, this will return the zero T value.
	Peek() T

	// PushOne will push a single value onto the top of the stack.
	PushOne(value T) Stack[T]

	// Push will put these values into the stack in reverse order so that
	// the next pop will get the top value first.
	Push(values ...T) Stack[T]

	// PushSeq will put the first count values into the stack in reverse
	// order so that the next pop will get the top value first.
	// If the sequence ends before the count, the remaining will be zero values.
	PushSeq(s iterator.Iterator[T], count int) Stack[T]
}

const allocateSize = 8

type node[T any] struct {
	value T
	prev  *node[T]
}

type stackImp[T any] struct {
	zero      T
	count     int
	top       *node[T]
	tombCount int
	tombs     *node[T]
}

func New[T any]() Stack[T] {
	return &stackImp[T]{}
}

func NewWithCap[T any](count int) Stack[T] {
	return &stackImp[T]{
		tombs: allocateNodes[T](count, nil),
	}
}

func (s *stackImp[T]) Count() int    { return s.count }
func (s *stackImp[T]) Empty() bool   { return s.top == nil }
func (s *stackImp[T]) Capacity() int { return s.count + s.tombCount }

func (s *stackImp[T]) Grow(capacity int) Stack[T] {
	s.grow(capacity - s.count)
	return s
}

func (s *stackImp[T]) grow(totalTombCount int) {
	if grow := totalTombCount - s.tombCount; grow > 0 {
		if mod := grow % allocateSize; mod > 0 {
			grow += allocateSize - mod
		}
		s.tombs = allocateNodes(grow, s.tombs)
		s.tombCount += grow
	}
}

func (s *stackImp[T]) Trim() Stack[T] {
	s.tombs = nil
	s.tombCount = 0
	return s
}

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

func allocateNodes[T any](count int, prev *node[T]) *node[T] {
	if count <= 0 {
		return prev
	}
	nodes := make([]node[T], count)
	nodes[0].prev = prev
	for i := 1; i < count; i++ {
		nodes[i].prev = &nodes[i-1]
	}
	return &nodes[count-1]
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
	s.tombCount++
	return value
}

func (s *stackImp[T]) Peek() T {
	if s.top == nil {
		return s.zero
	}
	return s.top.value
}

func (s *stackImp[T]) PushOne(value T) Stack[T] {
	if s.tombs == nil {
		s.tombs = allocateNodes[T](allocateSize, nil)
		s.tombCount = allocateSize
	}
	n := s.tombs
	s.tombs = n.prev
	s.tombCount--
	s.top = set(n, value, s.top)
	s.count++
	return s
}

func (s *stackImp[T]) Push(values ...T) Stack[T] {
	count := len(values)
	if count <= 0 {
		return s
	}
	s.grow(count)
	// Asserts: count > 0, tomb count >= count
	cur := s.tombs
	var last *node[T]
	for _, v := range values {
		last = cur
		cur.value = v
		cur = cur.prev
	}
	last.prev = s.top
	s.top = s.tombs
	s.tombs = cur
	s.tombCount -= count
	s.count += count
	return s
}

func (s *stackImp[T]) PushSeq(it iterator.Iterator[T], count int) Stack[T] {
	if count <= 0 {
		count = allocateSize
	}
	s.grow(count)
	// Asserts: count > 0, tomb count >= count
	cur := s.tombs
	var last *node[T]
	actual := 0
	for v := range it {
		last = cur
		cur.value = v
		actual++
		cur = cur.prev
		if cur == nil {
			cur = allocateNodes[T](allocateSize, nil)
			s.tombCount += allocateSize
			last.prev = cur
		}
	}
	last.prev = s.top
	s.top = s.tombs
	s.tombs = cur
	s.tombCount -= actual
	s.count += actual
	return s
}
