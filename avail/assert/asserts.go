package asserts

import (
	"cmp"
	"errors"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

// disableAsserts is used to disable the [Assert] methods.
// This is useful to ensure that an assert doesn't panic
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
