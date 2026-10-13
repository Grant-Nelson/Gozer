package wad

import (
	"go/token"
	"io"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Reader struct {
	r   io.Reader
	reg *Registry
	bc  *bitCaster
}

// NewReader creates a new Wad reader.
//
// The Reader currently requires the uint size and endianness match.
// This is because this is created for building a cache that is
// expected to be read on the same machine as it was written.
// If the read data doesn't match, then this will panic.
func NewReader(r io.Reader) *Reader {
	wad := &Reader{
		r:   r,
		reg: DefaultRegistry,
		bc:  newBitCaster(),
	}
	// Read wad header information and check for match.
	otherVersion := wad.Uint8() // wad version
	wasUint64 := wad.Bool()     // is uint 64-bits?
	wasBigEndian := wad.Bool()  // is big endian?
	if version != otherVersion || wasUint64 != isUint64 || wasBigEndian != isBigEndian {
		panic(faults.New(`failed to create wad.Reader with different header`).
			With(`read version`, version).
			With(`read uint is 64-bits`, wasUint64).
			With(`read big endian`, wasBigEndian).
			With(`current version`, version).
			With(`current uint is 64-bits`, isUint64).
			With(`current big endian`, isBigEndian))
	}
	return wad
}

func (wad *Reader) SetRegistry(reg *Registry) { wad.reg = reg }
func (wad *Reader) Registry() *Registry       { return wad.reg }

func (wad *Reader) Int() int {
	if isUint64 {
		return int(wad.Uint64())
	}
	return int(wad.Uint32())
}

func (wad *Reader) Uint() uint {
	if isUint64 {
		return uint(wad.Uint64())
	}
	return uint(wad.Uint32())
}

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
	length := wad.Int()
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
	t, ok := w.(T)
	if !ok {
		panic(faults.New(`failed to convert Waddable into expected type`).
			With(`waddable`, w).
			WithF(`type`, `%T`, t))
	}
	t.UnWad(wad)
	return t
}
