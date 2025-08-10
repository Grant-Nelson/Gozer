package iterator

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/Grant-Nelson/Gozer/internal/predicate"
)

type Iterator[T any] iter.Seq[T]

func (it Iterator[T]) Indexed() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range it {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

func (it Iterator[T]) Count() int {
	count := 0
	for range it {
		count++
	}
	return count
}

func (it Iterator[T]) Empty() bool {
	for range it {
		return false
	}
	return true
}

func (it Iterator[T]) First() T {
	for v := range it {
		return v
	}
	var zero T
	return zero
}

func (it Iterator[T]) Last() T {
	var t T
	for v := range it {
		t = v
	}
	return t
}

func (it Iterator[T]) LessThan(count int) bool {
	for range it {
		if count <= 0 {
			return false
		}
		count--
	}
	return true
}

func (it Iterator[T]) GreaterThan(count int) bool {
	for range it {
		if count <= 0 {
			return true
		}
		count--
	}
	return false
}

func (it Iterator[T]) Skip(count int) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if count <= 0 && !yield(v) {
				return
			}
			count--
		}
	}
}

func (it Iterator[T]) Limit(count int) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !yield(v) || count <= 0 {
				return
			}
			count--
		}
	}
}

func (it Iterator[T]) OnlyOne() T {
	first := true
	var t T
	for v := range it {
		if !first {
			var zero T
			return zero
		}
		first = false
		t = v
	}
	return t
}

func (it Iterator[T]) Any(p predicate.Predicate[T]) bool {
	for v := range it {
		if p(v) {
			return true
		}
	}
	return false
}

func (it Iterator[T]) All(p predicate.Predicate[T]) bool {
	for v := range it {
		if !p(v) {
			return false
		}
	}
	return true
}

func (it Iterator[T]) Where(p predicate.Predicate[T]) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if p(v) && !yield(v) {
				return
			}
		}
	}
}

func (it Iterator[T]) WhereNot(p predicate.Predicate[T]) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if p(v) && !yield(v) {
				return
			}
		}
	}
}

func (it Iterator[T]) Append(tails ...Iterator[T]) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !yield(v) {
				return
			}
		}
		for _, tail := range tails {
			for v := range tail {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (it Iterator[T]) ToSlice() []T {
	return slices.Collect(iter.Seq[T](it))
}

func (it Iterator[T]) ToStrings() Iterator[string] {
	return func(yield func(string) bool) {
		for v := range it {
			if !yield(fmt.Sprint(v)) {
				return
			}
		}
	}
}

func (it Iterator[T]) Join(sep string) string {
	buf := strings.Builder{}
	first := true
	for v := range it {
		if first {
			first = false
		} else {
			_, _ = buf.WriteString(sep)
		}
		_, _ = buf.WriteString(fmt.Sprint(v))
	}
	return buf.String()
}

func Empty[T any]() Iterator[T] {
	return func(yield func(T) bool) {}
}

func NotZero[T comparable](it Iterator[T]) Iterator[T] {
	return func(yield func(T) bool) {
		var zero T
		for v := range it {
			if v != zero && !yield(v) {
				return
			}
		}
	}
}

func Expand[T1 any, T2 iter.Seq[T1]](it Iterator[T2]) Iterator[T1] {
	return func(yield func(T1) bool) {
		for p := range it {
			for v := range p {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Append[T any](its ...Iterator[T]) Iterator[T] {
	switch len(its) {
	case 0:
		return Empty[T]()
	case 1:
		return its[0]
	default:
		return its[0].Append(its[1:]...)
	}
}

func Select[TIn, TOut any](it Iterator[TIn], s func(TIn) TOut) Iterator[TOut] {
	return func(yield func(TOut) bool) {
		for v := range it {
			if !yield(s(v)) {
				return
			}
		}
	}
}

func Dedup[T comparable](it Iterator[T]) Iterator[T] {
	return func(yield func(T) bool) {
		seen := map[T]struct{}{}
		for v := range it {
			if _, has := seen[v]; has {
				continue
			}
			seen[v] = struct{}{}
			if !yield(v) {
				return
			}
		}
	}
}

func Zip[T1, T2 any](p1 Iterator[T1], p2 Iterator[T2]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		next1, stop1 := iter.Pull(iter.Seq[T1](p1))
		next2, stop2 := iter.Pull(iter.Seq[T2](p2))
		for {
			value1, has1 := next1()
			value2, has2 := next2()
			if has1 {
				if has2 {
					if !yield(value1, value2) {
						return
					}
				} else {
					stop1()
					return
				}
			} else if has2 {
				stop2()
			}
		}
	}
}

func Aggregate[T1, T2 any](it Iterator[T1], start T2, ag func(T1, T2) T2) T2 {
	cur := start
	for v := range it {
		cur = ag(v, cur)
	}
	return cur
}
