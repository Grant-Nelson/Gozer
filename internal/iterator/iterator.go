package iterator

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/Grant-Nelson/Gozer/internal/predicate"
)

// Iterator is an iterator over a sequence of values
// with additional methods added to it for convenance.
type Iterator[T any] iter.Seq[T]

// Iterate will create an iterator for all the given values.
func Iterate[T any](values ...T) Iterator[T] {
	return Iterator[T](slices.Values(values))
}

// Index returns the sequence values paired with a zero based index.
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

// Count iterates over the sequence of values to determine
// how many values there are.
func (it Iterator[T]) Count() int {
	count := 0
	for range it {
		count++
	}
	return count
}

// Empty will determine if there are no values in the sequence.
//
// This will try to consume one value from the sequence.
func (it Iterator[T]) Empty() bool {
	for range it {
		return false
	}
	return true
}

// LessThan will determine if there are less than the given
// count of values in the sequence.
//
// This will try to remove up to the given count from the sequence.
func (it Iterator[T]) LessThan(count int) bool {
	for range it {
		if count <= 0 {
			return false
		}
		count--
	}
	return true
}

// GreaterThan will determine if there are greater than the given
// count of values in the sequence.
//
// This will try to remove one more than the given count from the sequence.
func (it Iterator[T]) GreaterThan(count int) bool {
	for range it {
		if count <= 0 {
			return true
		}
		count--
	}
	return false
}

// This will return the first value in the sequence with a true
// or zero with false if the sequence is empty.
//
// This will try to consume one value from the sequence.
func (it Iterator[T]) First() (T, bool) {
	for v := range it {
		return v, true
	}
	var zero T
	return zero, false
}

// This will return the last value in the sequence with true
// or zero with false if the sequence is empty.
//
// This will consume all the values in the sequence.
func (it Iterator[T]) Last() (T, bool) {
	var t T
	found := false
	for v := range it {
		t = v
		found = true
	}
	return t, found
}

// Skip will skip over the given count of values in the sequence
// then will return all remaining values.
//
// This will consume all the values in the sequence.
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

// Limit will return up to the given count of values in the sequence
// before stopping and not returning any following values.
//
// This will try to consume up to the limit.
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

// SkipWhile will skip all values until the predicate returns true,
// then it will return the value that the predicate returned true for
// and all following values.
//
// This will consume all the values in the sequence but will only
// call the predicate until the first time the predicate returns true.
func (it Iterator[T]) SkipWhile(p predicate.Predicate[T]) Iterator[T] {
	return func(yield func(T) bool) {
		skipping := true
		for v := range it {
			if skipping {
				if p(v) {
					continue
				}
				skipping = false
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Until will returns values from the sequence until the predicate returns
// true then it will stop without returning anymore values.
// The value that caused the predicate to return true is not returned.
//
// This will consume all the values until the predicate returns true.
func (it Iterator[T]) Until(p predicate.Predicate[T]) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if p(v) || !yield(v) {
				return
			}
		}
	}
}

// While will returns values from the sequence while the predicate returns
// true. As soon at predicate returns false it will stop without returning
// anymore values.
// The value that caused the predicate to return false is not returned.
//
// This will consume all the values until the predicate returns false.
func (it Iterator[T]) While(p predicate.Predicate[T]) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !p(v) || !yield(v) {
				return
			}
		}
	}
}

// OnlyOne returns the only value in the sequence with true,
// otherwise if there are no values or more than one value
// then zero with false is returned.
//
// This will try to consume up to two values from the sequence.
func (it Iterator[T]) OnlyOne() (T, bool) {
	var t T
	found := false
	for v := range it {
		if found {
			var zero T
			return zero, false
		}
		found = true
		t = v
	}
	return t, found
}

// Any will return true if there are any values that cause the predicate
// to return true, otherwise false if the predicate never returns true.
//
// This will consume up to the first value the predicate returns true for.
func (it Iterator[T]) Any(p predicate.Predicate[T]) bool {
	for v := range it {
		if p(v) {
			return true
		}
	}
	return false
}

// All will return true if the predicate returns true for all values,
// otherwise false if the predicate returns false for any value.
//
// This will consume up to the first value the predicate returns false for.
func (it Iterator[T]) All(p predicate.Predicate[T]) bool {
	for v := range it {
		if !p(v) {
			return false
		}
	}
	return true
}

// Where will return any value where the predicate returns true.
//
// This will consume all the values in the sequence.
func (it Iterator[T]) Where(p predicate.Predicate[T]) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if p(v) && !yield(v) {
				return
			}
		}
	}
}

// WhereNot will return any value where the predicate returns false.
//
// This will consume all the values in the sequence.
func (it Iterator[T]) WhereNot(p predicate.Predicate[T]) Iterator[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !p(v) && !yield(v) {
				return
			}
		}
	}
}

// Append will concatenate the given tails onto the end of the
// given iterator while iterating. They will be concatenated
// in the order they are given.
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

// ToSlice will create a slice containing all the values in the sequence.
func (it Iterator[T]) ToSlice() []T {
	return slices.Collect(iter.Seq[T](it))
}

// ToStrings will return all the values stringified.
//
// The values are stringified via `fmt.Sprint(value)`.
func (it Iterator[T]) ToStrings() Iterator[string] {
	return func(yield func(string) bool) {
		for v := range it {
			if !yield(fmt.Sprint(v)) {
				return
			}
		}
	}
}

// ToStringF will return all the values stringified with the given format.
//
// The values are stringified via `fmt.Sprintf(format, value)`.
func (it Iterator[T]) ToStringsF(format string) Iterator[string] {
	return func(yield func(string) bool) {
		for v := range it {
			if !yield(fmt.Sprintf(format, v)) {
				return
			}
		}
	}
}

// Join will create a string for this sequence with the given
// separator between all the stringified values.
//
// The values are stringified via `fmt.Sprint(value)`.
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

// UntilError will return the given function for each value
// until the end of the iteration or the function returns an error.
// Returns nil if no error occurred, or the error that was returned
// from the given function.
func (it Iterator[T]) UntilError(f func(v T) error) error {
	for v := range it {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

// Empty will create an empty iterator.
func Empty[T any]() Iterator[T] {
	return func(yield func(T) bool) {}
}

// NotZero will return any value that is not zero.
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

// Expand will return all the values inside the iterators that are
// returned from this iterator.
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

// Appends will concatenate the given iterators in the order they are given.
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

// Select will change all the values in the given iterator into the
// values of the returned iterator with the given selector function.
func Select[TIn, TOut any](it Iterator[TIn], s func(TIn) TOut) Iterator[TOut] {
	return func(yield func(TOut) bool) {
		for v := range it {
			if !yield(s(v)) {
				return
			}
		}
	}
}

// Cast will cast all the values in the input iterator and return
// only the values that could be cast into the output type.
func Cast[TIn, TOut any](it Iterator[TIn]) Iterator[TOut] {
	return func(yield func(TOut) bool) {
		for v := range it {
			if t, ok := any(v).(TOut); ok && !yield(t) {
				return
			}
		}
	}
}

// Dedup will only return unique values and skip over any values that were
// already seen. This uses the values of a key in a map to determine if the
// value has been seen before.
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

// Zip will return the two iterators merged into one sequence.
// The returned sequence will be the same length as the shortest iterator.
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

// Aggregate will run all the values through the given aggregator function
// to reduce the value. The given initial value is used with the first value
// from the iterator and the result is used with the next value, and so on.
// The result value will be returned.
func Aggregate[T1, T2 any](it Iterator[T1], init T2, ag func(T1, T2) T2) T2 {
	cur := init
	for v := range it {
		cur = ag(v, cur)
	}
	return cur
}

// Reduce will run all the values through the given reduce function
// to reduce the value. The first value and the second value is passed into
// the reduction function, the result and third value is passed into the
// reduction function, and son on. The last result will be returned.
//
// For example the reduction function could be `max` or `min` to get the
// maximum or minimum value of all the values in the iterator.
func Reduce[T1 any](it Iterator[T1], r func(T1, T1) T1) T1 {
	var cur T1
	first := true
	for v := range it {
		if first {
			cur = v
			first = false
			continue
		}
		cur = r(v, cur)
	}
	return cur
}
