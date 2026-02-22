package faults

import (
	"fmt"
	"testing"
)

func TestErrGroup_AddOne(t *testing.T) {
	eg := NewErrGroup(3)
	AreEqual(t, `no errors`,
		eg.Error(), `empty error`)

	eg.Add(fmt.Errorf(`error 1`))
	AreEqual(t, `error 1`,
		eg.Error(), `one error`)

	eg.Add(fmt.Errorf(`error 2`))
	AreEqual(t,
		"multiple errors:\n"+
			"1. error 1\n"+
			"2. error 2",
		eg.Error(), `two error`)

	eg.Add(fmt.Errorf(`error 3`))
	AreEqual(t,
		"multiple errors:\n"+
			"1. error 1\n"+
			"2. error 2\n"+
			"3. error 3",
		eg.Error(), `three error`)

	eg.Add(fmt.Errorf(`error 4`))
	AreEqual(t,
		"multiple errors:\n"+
			"1. error 1\n"+
			"2. error 2\n"+
			"3. error 3\n"+
			"4. too many errors (1 discarded)",
		eg.Error(), `four error`)

	eg.Add(fmt.Errorf(`error 5`))
	AreEqual(t,
		"multiple errors:\n"+
			"1. error 1\n"+
			"2. error 2\n"+
			"3. error 3\n"+
			"4. too many errors (2 discarded)",
		eg.Error(), `five error`)
}

func TestErrGroup_AddMultiple(t *testing.T) {
	eg := NewErrGroup(3)
	eg.Add(fmt.Errorf(`error 1`),
		fmt.Errorf(`error 2`),
		fmt.Errorf(`error 3`),
		nil,
		fmt.Errorf(`error 4`),
		fmt.Errorf(`error 5`))
	AreEqual(t,
		"multiple errors:\n"+
			"1. error 1\n"+
			"2. error 2\n"+
			"3. error 3\n"+
			"4. too many errors (2 discarded)",
		eg.Error(), `five error`)
}

func AreEqual[T comparable](t testing.TB, exp, actual T, msg string) {
	t.Helper()
	if exp != actual {
		t.Errorf("Unexpected value: %s:\n"+
			"\tactual: %v\n"+
			"\texp:    %v\n", msg, actual, exp)
	}
}
