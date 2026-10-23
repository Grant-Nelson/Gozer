package wad

import (
	"bytes"
	"testing"
)

func Test_Basic(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf)
	w.Int8(12)
	w.Int16(34)
	w.Int32(56)
	w.Int64(78)
	w.String(`hello`)
	w.Float64(0.128)

	u := NewReader(buf)
	check(t, u.Int8(), 12, `int8`)
	check(t, u.Int16(), 34, `int16`)
	check(t, u.Int32(), 56, `int32`)
	check(t, u.Int64(), 78, `int64`)
	check(t, u.String(), `hello`, `string`)
	check(t, u.Float64(), 0.128, `float64`)
}

type address struct {
	number int
	city   string
}

func (a address) Wad(w *Writer) {
	w.Int(a.number)
	w.String(a.city)
}

func (a *address) UnWad(r *Reader) {
	a.number = r.Int()
	a.city = r.String()
}

type person struct {
	name    string
	address address
	female  bool
}

func (p *person) Wad(w *Writer) {
	w.String(p.name)
	w.Wad(p.address)
	w.Bool(p.female)
}

func (p *person) UnWad(r *Reader) {
	p.name = r.String()
	p.address = r.Wad[address]()
	p.female = r.Bool()
}

func Test_Waddables(t *testing.T) {
	reg := newRegistry()
	reg.Register[*person]()
	reg.Register[*address]()

	p := &person{
		name: `bob`,
		address: address{
			number: 42,
			city:   `ash`,
		},
		female: false,
	}

	buf := &bytes.Buffer{}
	w := newWriterNoHeader(buf, isUint64, isBigEndian, reg)
	w.Wad(p)

	r := newReaderNoHeader(buf, isUint64, isBigEndian, reg)
	p2 := r.Wad[*person]()

	check(t, *p2, *p, `person`)
}

func check[T comparable](t *testing.T, get, want T, name string) {
	if get != want {
		t.Errorf("unexpected %s:\n\tget:  %v\n\twant: %v", name, get, want)
	}
}
