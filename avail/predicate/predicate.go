package predicate

import (
	"cmp"
	"math"
)

type Predicate[T any] func(T) bool

func True[T any](t T) bool { return true }

func False[T any](t T) bool { return false }

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

func IsNotType[TOut, TIn any]() Predicate[TIn] {
	return func(t TIn) bool {
		_, ok := any(t).(TOut)
		return !ok
	}
}

func IsInf[T ~float32 | ~float64]() Predicate[T] {
	return func(t T) bool { return math.IsInf(float64(t), 0) }
}

func IsNotInf[T ~float32 | ~float64]() Predicate[T] {
	return func(t T) bool { return !math.IsInf(float64(t), 0) }
}

func IsNaN[T ~float32 | ~float64]() Predicate[T] {
	return func(t T) bool { return math.IsNaN(float64(t)) }
}

func IsNotNaN[T ~float32 | ~float64]() Predicate[T] {
	return func(t T) bool { return !math.IsNaN(float64(t)) }
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

// InRange checks that the value is between the given min and max values inclusively.
// If `min` is greater than `max`, then this predicate will not work correctly.
func InRange[T cmp.Ordered](min, max T) Predicate[T] {
	return func(t T) bool { return t >= min && t <= max }
}
