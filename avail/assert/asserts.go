package assert

import (
	"cmp"
	"errors"
	"math"
	"sync/atomic"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

// disableAsserts is used to enable (false, default) / disable (true) the [Assert] methods.
var disableAsserts = &atomic.Bool{}

// Enabled indicates if asserts are enabled or not. Defaults to enabled.
//
// This is useful to ensure that asserts do not panic while in production
// or debugging. This will be checked first to remove as much computation
// as possible when disabled so that the asserts are similar to no-ops.
func Enabled() bool { return !disableAsserts.Load() }

// Enable will enable / disable the [Assert] methods.
//
// This will return a function that when called will reset
// to the enabled state prior to this call, meaning the following
// can be used to temporarily enable or disable with the following:
//
//	defer assert.Enable(true)()
func Enable(enable bool) func() {
	previous := disableAsserts.Swap(!enable)
	return func() { disableAsserts.Store(previous) }
}

var ErrAssert = errors.New(`assert failed`)

// Eq asserts that `a` and `b` are equal, otherwise panics.
func Eq[T comparable](a, b T) {
	if Enabled() && !(a == b) {
		panic(faults.New(`%w equality check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`first`, a).
			With(`second`, b))
	}
}

// Neq asserts that `a` and `b` are not equal, otherwise panics.
func Neq[T comparable](a, b T) {
	// Note: The `!(a op b)` for the comparable and orderable is
	// to handle NaN correctly when T is a float32 or float64.
	if Enabled() && !(a != b) {
		panic(faults.New(`%w non-equality check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`first`, a).
			With(`second`, b))
	}
}

// Lt asserts that `a` is less than `b`, otherwise panics.
func Lt[T cmp.Ordered](a, b T) {
	if Enabled() && !(a < b) {
		panic(faults.New(`%w less-than check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`first`, a).
			With(`second`, b))
	}
}

// Gt asserts that `a` is greater than `b`, otherwise panics.
func Gt[T cmp.Ordered](a, b T) {
	if Enabled() && !(a > b) {
		panic(faults.New(`%w greater-than check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`first`, a).
			With(`second`, b))
	}
}

// Leq asserts that `a` is less than or equal to `b`, otherwise panics.
func Leq[T cmp.Ordered](a, b T) {
	if Enabled() && !(a <= b) {
		panic(faults.New(`%w less-than-or-equal check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`first`, a).
			With(`second`, b))
	}
}

// Geq asserts that `a` is greater than or equal to `b`, otherwise panics.
func Geq[T cmp.Ordered](a, b T) {
	if Enabled() && !(a >= b) {
		panic(faults.New(`%w greater-than-or-equal check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`first`, a).
			With(`second`, b))
	}
}

// Range asserts that `a` is between `low` and `high` inclusively, otherwise panics.
func Range[T cmp.Ordered](a, low, high T) {
	if Enabled() && !(a >= low && a <= high) {
		panic(faults.New(`%w range check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`value`, a).
			With(`low`, low).
			With(`high`, high))
	}
}

// Epsilon asserts that `a` is less than `epsilon` away from zero, otherwise panics.
func Epsilon[T ~float32 | ~float64](a, epsilon T) {
	if Enabled() && !(math.Abs(float64(a)) < float64(epsilon)) {
		panic(faults.New(`%w epsilon check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`value`, a).
			With(`epsilon`, epsilon))
	}
}

// Inf asserts that `a` is an infinite value, otherwise panics.
func Inf[T ~float32 | ~float64](a T) {
	if Enabled() && !math.IsInf(float64(a), 0) {
		panic(faults.New(`%w infinite check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`value`, a))
	}
}

// NotInf asserts that `a` is not an infinite value, otherwise panics.
func NotInf[T ~float32 | ~float64](a T) {
	if Enabled() && math.IsInf(float64(a), 0) {
		panic(faults.New(`%w not infinite check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`value`, a))
	}
}

// NaN asserts that `a` is not a number (NaN), otherwise panics.
func NaN[T ~float32 | ~float64](a T) {
	if Enabled() && !math.IsNaN(float64(a)) {
		panic(faults.New(`%w NaN check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`value`, a))
	}
}

// NotNaN asserts that `a` is a number (not NaN), otherwise panics.
func NotNaN[T ~float32 | ~float64](a T) {
	if Enabled() && math.IsNaN(float64(a)) {
		panic(faults.New(`%w not NaN check`, ErrAssert).
			WithTypeFor[T](`type`).
			With(`value`, a))
	}
}

// Nil asserts that the pointer `a` is nil, otherwise panics.
func Nil[T any, P ~*T](p P) {
	if Enabled() && p != nil {
		panic(faults.New(`%w must be nil check`, ErrAssert).
			WithTypeFor[P](`type`).
			With(`value`, p))
	}
}

// NotNil asserts that the pointer `a` is not nil, otherwise panics.
func NotNil[T any, P ~*T](p P) {
	if Enabled() && p == nil {
		panic(faults.New(`%w not nil check`, ErrAssert).
			WithTypeFor[P](`type`))
	}
}

// Zero asserts that `a` is the zero instance of `a`'s type, otherwise panics.
func Zero[T comparable](a T) {
	if Enabled() {
		var zero T
		if a != zero {
			panic(faults.New(`%w must be zero check`, ErrAssert).
				WithTypeFor[T](`type`).
				With(`value`, a))
		}
	}
}

// Nonzero asserts that `a` is not the zero instance of `a`'s type, otherwise panics.
func Nonzero[T comparable](a T) {
	if Enabled() {
		var zero T
		if a == zero {
			panic(faults.New(`%w not zero check`, ErrAssert).
				WithTypeFor[T](`type`).
				With(`value`, a))
		}
	}
}

// Fn asserts that `a` returns true when called, otherwise panics.
// This will run `a` once so will block until `a` finishes.
// If asserts are disabled, then `a` will not be run.
func Fn(a func() bool) {
	if Enabled() && !a() {
		panic(faults.New(`%w function check`, ErrAssert))
	}
}

// True asserts that `a` is true, otherwise panics.
func True(a bool) {
	if Enabled() && !a {
		panic(faults.New(`%w true check`, ErrAssert))
	}
}

// False asserts that `a` is false, otherwise panics.
func False(a bool) {
	if Enabled() && a {
		panic(faults.New(`%w false check`, ErrAssert))
	}
}

// EmptyStr asserts that `s` is an empty string, otherwise panics.
func EmptyStr[T ~string](s T) {
	if Enabled() {
		if ln := len(s); ln != 0 {
			panic(faults.New(`%w empty string check`, ErrAssert).
				WithTypeFor[T](`type`).
				With(`length`, ln))
		}
	}
}

// NotEmptyStr asserts that `s` is not empty string, otherwise panics.
func NotEmptyStr[T ~string](s T) {
	if Enabled() {
		if ln := len(s); ln == 0 {
			panic(faults.New(`%w not empty string check`, ErrAssert).
				WithTypeFor[T](`type`).
				With(`length`, ln))
		}
	}
}

// StrRange asserts that `s` has a length between `low` and `high` inclusively, otherwise panics.
func StrRange[T ~string](s T, low, high int) {
	if Enabled() {
		if ln := len(s); ln < low || ln > high {
			panic(faults.New(`%w string length range check`, ErrAssert).
				WithTypeFor[T](`type`).
				With(`length`, ln).
				With(`low`, low).
				With(`high`, high))
		}
	}
}

// EmptySlice asserts that `s` is empty, otherwise panics.
func EmptySlice[T any, S ~[]T](s S) {
	if Enabled() {
		if ln := len(s); ln != 0 {
			panic(faults.New(`%w empty slice check`, ErrAssert).
				WithTypeFor[S](`type`).
				With(`length`, ln))
		}
	}
}

// NotEmptySlice asserts that `s` is not empty, otherwise panics.
func NotEmptySlice[T any, S ~[]T](s S) {
	if Enabled() {
		if ln := len(s); ln == 0 {
			panic(faults.New(`%w not empty slice check`, ErrAssert).
				WithTypeFor[S](`type`).
				With(`length`, ln))
		}
	}
}

// SliceRange asserts that `s` has a length between `low` and `high` inclusively, otherwise panics.
func SliceRange[T any, S ~[]T](s S, low, high int) {
	if Enabled() {
		if ln := len(s); ln < low || ln > high {
			panic(faults.New(`%w slice length range check`, ErrAssert).
				WithTypeFor[S](`type`).
				With(`length`, ln).
				With(`low`, low).
				With(`high`, high))
		}
	}
}

// EmptyMap asserts that `m` is empty, otherwise panics.
func EmptyMap[K comparable, V any, M ~map[K]V](m M) {
	if Enabled() {
		if ln := len(m); ln != 0 {
			panic(faults.New(`%w empty map check`, ErrAssert).
				WithTypeFor[M](`type`).
				With(`length`, ln))
		}
	}
}

// NotEmptyMap asserts that `m` is not empty, otherwise panics.
func NotEmptyMap[K comparable, V any, M ~map[K]V](m M) {
	if Enabled() {
		if ln := len(m); ln == 0 {
			panic(faults.New(`%w not empty map check`, ErrAssert).
				WithTypeFor[M](`type`).
				With(`length`, ln))
		}
	}
}

// MapRange asserts that `m` has a length between `low` and `high` inclusively, otherwise panics.
func MapRange[K comparable, V any, M ~map[K]V](m M, low, high int) {
	if Enabled() {
		if ln := len(m); ln < low || ln > high {
			panic(faults.New(`%w map length range check`, ErrAssert).
				WithTypeFor[M](`type`).
				With(`length`, ln).
				With(`low`, low).
				With(`high`, high))
		}
	}
}
