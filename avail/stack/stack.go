package stack

type Stack[T any] interface {
	Count() int
	Empty() bool
	Pop() T
	PushOne(value T)
	// Push will put these values into the stack in reverse order so that
	// the next pop will get the first of these values first.
	Push(values ...T)
}

type node[T any] struct {
	value T
	prev  *node[T]
}

type stackImp[T any] struct {
	count int
	top   *node[T]
}

func New[T any]() Stack[T] { return &stackImp[T]{} }

func (s *stackImp[T]) Count() int { return s.count }

func (s *stackImp[T]) Empty() bool { return s.top == nil }

func (s *stackImp[T]) Pop() T {
	if s.top == nil {
		var zero T
		return zero
	}
	value := s.top.value
	s.top = s.top.prev
	s.count--
	return value
}

func (s *stackImp[T]) PushOne(value T) {
	s.top = &node[T]{value: value, prev: s.top}
	s.count++
}

func (s *stackImp[T]) Push(values ...T) {
	for i := len(values) - 1; i >= 0; i-- {
		s.PushOne(values[i])
	}
}
