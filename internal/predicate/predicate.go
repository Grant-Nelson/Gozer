package predicate

type Predicate[T any] func(T) bool

func Not[T any](p Predicate[T]) Predicate[T] {
	return func(t T) bool { return !p(t) }
}

func And[T any](ps ...Predicate[T]) Predicate[T] {
	return func(t T) bool {
		for _, p := range ps {
			if !p(t) {
				return false
			}
		}
		return true
	}
}

func Or[T any](ps ...Predicate[T]) Predicate[T] {
	return func(t T) bool {
		for _, p := range ps {
			if p(t) {
				return true
			}
		}
		return false
	}
}

func Is[TOut, TIn any]() Predicate[TIn] {
	return func(t TIn) bool {
		_, ok := any(t).(TOut)
		return ok
	}
}
