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
