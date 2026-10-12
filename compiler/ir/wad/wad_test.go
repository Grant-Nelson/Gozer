package wad

import (
	"bytes"
	"testing"
)

func Test_Basic(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWadder(buf)
	w.Int8(12)
	w.Int16(34)
	w.Int32(56)
	w.Int64(78)
	w.String(`hello`)
	w.Float64(0.128)

	u := NewUnWadder(buf)
	check(t, u.Int8(), 12, `int8`)
	check(t, u.Int16(), 34, `int16`)
	check(t, u.Int32(), 56, `int32`)
	check(t, u.Int64(), 78, `int64`)
	check(t, u.String(), `hello`, `string`)
	check(t, u.Float64(), 0.128, `float64`)
}

func check[T comparable](t *testing.T, get, want T, name string) {
	if get != want {
		t.Errorf("unexpected %s:\n\tget:  %v\n\twant: %v", name, get, want)
	}
}
