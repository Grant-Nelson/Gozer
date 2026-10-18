package wad

import (
	"go/token"
	"io"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Reader struct {
	r      io.Reader
	reg    *Registry
	bc     bitCaster
	fnInt  func() int
	fnUint func() uint
}

// NewReader creates a new Wad reader.
//
// The Reader currently requires the uint size and endianness match.
// This is because this is created for building a cache that is
// expected to be read on the same machine as it was written.
// If the read data doesn't match, then this will panic.
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
	return newReaderNoHeader(r, useUint64, useBigEndian)
}

func newReaderNoHeader(r io.Reader, useUint64, useBigEndian bool) *Reader {
	wad := &Reader{
		r:   r,
		reg: defaultRegistry.clone(),
		bc:  newBitCaster(useUint64, useBigEndian),
	}
	if useUint64 {
		wad.fnInt = wad.intVia64
		wad.fnUint = wad.uintVia64
	} else {
		wad.fnInt = wad.intVia32
		wad.fnUint = wad.uintVia32
	}
	return wad
}

func (wad *Reader) Registry() *Registry { return wad.reg }

func (wad *Reader) intVia32() int   { return int(wad.Uint32()) }
func (wad *Reader) intVia64() int   { return int(wad.Uint64()) }
func (wad *Reader) uintVia32() uint { return uint(wad.Uint32()) }
func (wad *Reader) uintVia64() uint { return uint(wad.Uint64()) }

func (wad *Reader) Int() int         { return wad.fnInt() }
func (wad *Reader) Uint() uint       { return wad.fnUint() }
func (wad *Reader) Bool() bool       { return wad.Uint8() > 0 }
func (wad *Reader) Int8() int8       { return int8(wad.Uint8()) }
func (wad *Reader) Int16() int16     { return int16(wad.Uint16()) }
func (wad *Reader) Int32() int32     { return int32(wad.Uint32()) }
func (wad *Reader) Int64() int64     { return int64(wad.Uint64()) }
func (wad *Reader) String() string   { return string(wad.Bytes()) }
func (wad *Reader) Uint8() uint8     { return wad.bc.readUint8(wad.r) }
func (wad *Reader) Uint16() uint16   { return wad.bc.readUint16(wad.r) }
func (wad *Reader) Uint32() uint32   { return wad.bc.readUint32(wad.r) }
func (wad *Reader) Uint64() uint64   { return wad.bc.readUint64(wad.r) }
func (wad *Reader) Float32() float32 { return wad.bc.readFloat32(wad.r) }
func (wad *Reader) Float64() float64 { return wad.bc.readFloat64(wad.r) }
func (wad *Reader) Pos() token.Pos   { return token.Pos(wad.Int()) }

func (wad *Reader) Complex64() complex64 {
	return complex(wad.Float32(), wad.Float32())
}

func (wad *Reader) Complex128() complex128 {
	return complex(wad.Float64(), wad.Float64())
}

func (wad *Reader) Bytes() []byte {
	length := wad.Uint32()
	data := make([]byte, length)
	read(wad.r, data)
	return data
}

func (wad *Reader) Wad[T Waddable]() T {
	i := wad.Uint16()
	if i == 0 {
		var zero T
		return zero
	}
	w := wad.reg.fromIndex(i)
	if t, ok := w.(T); ok {
		t.UnWad(wad)
		return t
	}
	var zero T
	panic(faults.New(`failed to convert Waddable into expected type`).
		With(`waddable`, w).
		WithF(`expected`, `%T`, zero).
		WithF(`type`, `%T`, w))
}

func (wad *Reader) basicSlice[T any](elem func() T) []T {
	count := wad.Uint32()
	values := make([]T, count)
	for i := range count {
		values[i] = elem()
	}
	return values
}

func (wad *Reader) IntSlice() []int         { return wad.basicSlice(wad.fnInt) }
func (wad *Reader) UintSlice() []uint       { return wad.basicSlice(wad.fnUint) }
func (wad *Reader) BoolSlice() []bool       { return wad.basicSlice(wad.Bool) }
func (wad *Reader) Int8Slice() []int8       { return wad.basicSlice(wad.Int8) }
func (wad *Reader) Int16Slice() []int16     { return wad.basicSlice(wad.Int16) }
func (wad *Reader) Int32Slice() []int32     { return wad.basicSlice(wad.Int32) }
func (wad *Reader) Int64Slice() []int64     { return wad.basicSlice(wad.Int64) }
func (wad *Reader) StringSlice() []string   { return wad.basicSlice(wad.String) }
func (wad *Reader) Uint8Slice() []uint8     { return wad.Bytes() }
func (wad *Reader) Uint16Slice() []uint16   { return wad.basicSlice(wad.Uint16) }
func (wad *Reader) Uint32Slice() []uint32   { return wad.basicSlice(wad.Uint32) }
func (wad *Reader) Uint64Slice() []uint64   { return wad.basicSlice(wad.Uint64) }
func (wad *Reader) Float32Slice() []float32 { return wad.basicSlice(wad.Float32) }
func (wad *Reader) Float64Slice() []float64 { return wad.basicSlice(wad.Float64) }
func (wad *Reader) PosSlice() []token.Pos   { return wad.basicSlice(wad.Pos) }

func (wad *Reader) WadSlice[T Waddable, S ~[]T]() S {
	count := wad.Uint32()
	values := make(S, count)
	wad.reg.Preallocate[T](int(count))
	for i := range count {
		values[i] = wad.Wad[T]()
	}
	return values
}
