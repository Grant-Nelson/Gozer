package assert

import (
	"cmp"
	"errors"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

// disableAsserts is used to disable the [Assert] methods.
// This is useful to ensure that asserts doesn't panic while in production
// or debugging. This will be checked first to remove as much computation
// as possible when disabled so that the asserts are similar to no-ops.
//
// TODO: Should make this able to be switched on and off via an env var.
const enableAsserts = true

var ErrAssert = errors.New(`assert failed`)

func Eq[T comparable](a, b T) {
	if enableAsserts && a != b {
		panic(faults.New(`%w equality check`, ErrAssert).
			WithF(`type`, `%T`, a).
			With(`first`, a).
			With(`second`, b))
	}
}

func Neq[T comparable](a, b T) {
	if enableAsserts && a == b {
		panic(faults.New(`%w non-equality check`, ErrAssert).
			WithF(`type`, `%T`, a).
			With(`first`, a).
			With(`second`, b))
	}
}

func Lt[T cmp.Ordered](a, b T) {
	if enableAsserts && a >= b {
		panic(faults.New(`%w less-than check`, ErrAssert).
			WithF(`type`, `%T`, a).
			With(`first`, a).
			With(`second`, b))
	}
}

func Gt[T cmp.Ordered](a, b T) {
	if enableAsserts && a <= b {
		panic(faults.New(`%w greater-than check`, ErrAssert).
			WithF(`type`, `%T`, a).
			With(`first`, a).
			With(`second`, b))
	}
}

func Leq[T cmp.Ordered](a, b T) {
	if enableAsserts && a > b {
		panic(faults.New(`%w less-than-or-equal check`, ErrAssert).
			WithF(`type`, `%T`, a).
			With(`first`, a).
			With(`second`, b))
	}
}

func Geq[T cmp.Ordered](a, b T) {
	if enableAsserts && a < b {
		panic(faults.New(`%w greater-than-or-equal check`, ErrAssert).
			WithF(`type`, `%T`, a).
			With(`first`, a).
			With(`second`, b))
	}
}

func Range[T cmp.Ordered](a, low, high T) {
	if enableAsserts && (a < low || a > high) {
		panic(faults.New(`%w range check`, ErrAssert).
			WithF(`type`, `%T`, a).
			With(`value`, a).
			With(`low`, low).
			With(`high`, high))
	}
}

func Nil[T any, P ~*T](p P) {
	if enableAsserts && p != nil {
		panic(faults.New(`%w must be nil check`, ErrAssert).
			WithF(`type`, `%T`, p).
			With(`value`, p))
	}
}

func NotNil[T any, P ~*T](p P) {
	if enableAsserts && p == nil {
		panic(faults.New(`%w not nil check`, ErrAssert).
			WithF(`type`, `%T`, p))
	}
}

func Zero[T comparable](a T) {
	if enableAsserts {
		var zero T
		if a != zero {
			panic(faults.New(`%w must be nil check`, ErrAssert).
				WithF(`type`, `%T`, a).
				With(`value`, a))
		}
	}
}

func Nonzero[T comparable](a T) {
	if enableAsserts {
		var zero T
		if a == zero {
			panic(faults.New(`%w not nil check`, ErrAssert).
				WithF(`type`, `%T`, a).
				With(`value`, a))
		}
	}
}

func Fn(a func() bool) {
	if enableAsserts && !a() {
		panic(faults.New(`%w function check`, ErrAssert))
	}
}

func True(a bool) {
	if enableAsserts && !a {
		panic(faults.New(`%w true check`, ErrAssert))
	}
}

func False(a bool) {
	if enableAsserts && a {
		panic(faults.New(`%w false check`, ErrAssert))
	}
}

func EmptyStr(s string) {
	if enableAsserts {
		if ln := len(s); ln != 0 {
			panic(faults.New(`%w empty check`, ErrAssert).
				WithF(`type`, `%T`, s).
				With(`length`, ln))
		}
	}
}

func NotEmptyStr(s string) {
	if enableAsserts {
		if ln := len(s); ln == 0 {
			panic(faults.New(`%w not empty check`, ErrAssert).
				WithF(`type`, `%T`, s).
				With(`length`, ln))
		}
	}
}

func StrRange(s string, low, high int) {
	if enableAsserts {
		if ln := len(s); ln < low || ln > high {
			panic(faults.New(`%w length range check`, ErrAssert).
				WithF(`type`, `%T`, s).
				With(`length`, ln).
				With(`low`, low).
				With(`high`, high))
		}
	}
}

func EmptySlice[T any, S ~[]T](s S) {
	if enableAsserts {
		if ln := len(s); ln != 0 {
			panic(faults.New(`%w empty check`, ErrAssert).
				WithF(`type`, `%T`, s).
				With(`length`, ln))
		}
	}
}

func NotEmptySlice[T any, S ~[]T](s S) {
	if enableAsserts {
		if ln := len(s); ln == 0 {
			panic(faults.New(`%w not empty check`, ErrAssert).
				WithF(`type`, `%T`, s).
				With(`length`, ln))
		}
	}
}

func SliceRange[T any, S ~[]T](s S, low, high int) {
	if enableAsserts {
		if ln := len(s); ln < low || ln > high {
			panic(faults.New(`%w length range check`, ErrAssert).
				WithF(`type`, `%T`, s).
				With(`length`, ln).
				With(`low`, low).
				With(`high`, high))
		}
	}
}

func EmptyMap[K comparable, V any, M ~map[K]V](m M) {
	if enableAsserts {
		if ln := len(m); ln != 0 {
			panic(faults.New(`%w empty check`, ErrAssert).
				WithF(`type`, `%T`, m).
				With(`length`, ln))
		}
	}
}

func NotEmptyMap[K comparable, V any, M ~map[K]V](m M) {
	if enableAsserts {
		if ln := len(m); ln == 0 {
			panic(faults.New(`%w not empty check`, ErrAssert).
				WithF(`type`, `%T`, m).
				With(`length`, ln))
		}
	}
}

func MapRange[K comparable, V any, M ~map[K]V](m M, low, high int) {
	if enableAsserts {
		if ln := len(m); ln < low || ln > high {
			panic(faults.New(`%w length range check`, ErrAssert).
				WithF(`type`, `%T`, m).
				With(`length`, ln).
				With(`low`, low).
				With(`high`, high))
		}
	}
}
