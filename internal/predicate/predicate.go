package predicate

import "cmp"

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

func IsType[TOut, TIn any]() Predicate[TIn] {
	return func(t TIn) bool {
		_, ok := any(t).(TOut)
		return ok
	}
}

func IsZero[T comparable]() Predicate[T] {
	var zero T
	return func(t T) bool { return t == zero }
}

func IsNotZero[T comparable]() Predicate[T] {
	var zero T
	return func(t T) bool { return t != zero }
}

func GreaterThan[T cmp.Ordered](v T) Predicate[T] {
	return func(t T) bool { return t > v }
}

func GreaterOrEqual[T cmp.Ordered](v T) Predicate[T] {
	return func(t T) bool { return t >= v }
}

func LessThan[T cmp.Ordered](v T) Predicate[T] {
	return func(t T) bool { return t < v }
}

func LessOrEqual[T cmp.Ordered](v T) Predicate[T] {
	return func(t T) bool { return t <= v }
}

func InRange[T cmp.Ordered](min, max T) Predicate[T] {
	return func(t T) bool { return t >= min && t <= max }
}
