package wad

import (
	"go/token"
	"io"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Reader struct {
	r          io.Reader
	bc         bitCaster
	regReaders []ReadFn
}

// NewReader creates a new Wad reader.
//
// When reading a Wad file, the calls to the Wad must match
// the same order as they were written.
func NewReader(r io.Reader) *Reader {
	header := []byte{0, 0}
	read(r, header)
	if version != header[0] {
		panic(faults.New(`failed to create wad.Reader for the data's version`).
			With(`read version`, header[0]).
			With(`current version`, version))
	}
	useUint64 := header[1]&0x01 != 0x00
	useBigEndian := header[1]&0x02 != 0x00
	return newReaderNoHeader(r, useUint64, useBigEndian, regSingleton)
}

func newReaderNoHeader(r io.Reader, useUint64, useBigEndian bool, reg *registry) *Reader {
	return &Reader{
		r:          r,
		bc:         newBitCaster(useUint64, useBigEndian),
		regReaders: reg.cloneReaders(),
	}
}

func (wad *Reader) Bool() bool       { return wad.bc.readUint8(wad.r) > 0 }
func (wad *Reader) Int() int         { return int(wad.bc.readUint(wad.r)) }
func (wad *Reader) Int8() int8       { return int8(wad.bc.readUint8(wad.r)) }
func (wad *Reader) Int16() int16     { return int16(wad.bc.readUint16(wad.r)) }
func (wad *Reader) Int32() int32     { return int32(wad.bc.readUint32(wad.r)) }
func (wad *Reader) Int64() int64     { return int64(wad.bc.readUint64(wad.r)) }
func (wad *Reader) String() string   { return string(wad.Bytes()) }
func (wad *Reader) Uint() uint       { return wad.bc.readUint(wad.r) }
func (wad *Reader) Uint8() uint8     { return wad.bc.readUint8(wad.r) }
func (wad *Reader) Uint16() uint16   { return wad.bc.readUint16(wad.r) }
func (wad *Reader) Uint32() uint32   { return wad.bc.readUint32(wad.r) }
func (wad *Reader) Uint64() uint64   { return wad.bc.readUint64(wad.r) }
func (wad *Reader) Float32() float32 { return wad.bc.readFloat32(wad.r) }
func (wad *Reader) Float64() float64 { return wad.bc.readFloat64(wad.r) }
func (wad *Reader) Pos() token.Pos   { return token.Pos(wad.bc.readUint(wad.r)) }

func (wad *Reader) Complex64() complex64 {
	return complex(wad.bc.readFloat32(wad.r), wad.bc.readFloat32(wad.r))
}

func (wad *Reader) Complex128() complex128 {
	return complex(wad.bc.readFloat64(wad.r), wad.bc.readFloat64(wad.r))
}

func (wad *Reader) Bytes() []byte {
	length := wad.bc.readUint32(wad.r)
	data := make([]byte, length)
	read(wad.r, data)
	return data
}

func (wad *Reader) Wad[T any]() T {
	i := wad.bc.readUint16(wad.r)
	if i <= 0 {
		var zero T
		return zero
	}
	if int(i) > len(wad.regReaders) {
		panic(faults.New(`Unable to read Waddable from registry with given index`).
			With(`index`, i))
	}

	w := wad.regReaders[i-1](wad)
	if t, ok := w.(T); ok {
		return t
	}

	panic(faults.New(`unable to cast Waddable to target type`).
		WithTypeFor[T](`target`).
		WithTypeOf(`type`, w))
}

func (wad *Reader) slice[T any](elem func() T) []T {
	count := wad.bc.readUint32(wad.r)
	values := make([]T, count)
	for i := range count {
		values[i] = elem()
	}
	return values
}

func (wad *Reader) IntSlice() []int         { return wad.slice(wad.Int) }
func (wad *Reader) BoolSlice() []bool       { return wad.slice(wad.Bool) }
func (wad *Reader) Int8Slice() []int8       { return wad.slice(wad.Int8) }
func (wad *Reader) Int16Slice() []int16     { return wad.slice(wad.Int16) }
func (wad *Reader) Int32Slice() []int32     { return wad.slice(wad.Int32) }
func (wad *Reader) Int64Slice() []int64     { return wad.slice(wad.Int64) }
func (wad *Reader) StringSlice() []string   { return wad.slice(wad.String) }
func (wad *Reader) UintSlice() []uint       { return wad.slice(wad.Uint) }
func (wad *Reader) Uint8Slice() []uint8     { return wad.Bytes() }
func (wad *Reader) Uint16Slice() []uint16   { return wad.slice(wad.Uint16) }
func (wad *Reader) Uint32Slice() []uint32   { return wad.slice(wad.Uint32) }
func (wad *Reader) Uint64Slice() []uint64   { return wad.slice(wad.Uint64) }
func (wad *Reader) Float32Slice() []float32 { return wad.slice(wad.Float32) }
func (wad *Reader) Float64Slice() []float64 { return wad.slice(wad.Float64) }
func (wad *Reader) PosSlice() []token.Pos   { return wad.slice(wad.Pos) }

func (wad *Reader) WadSlice[T any]() []T { return wad.slice(wad.Wad[T]) }
