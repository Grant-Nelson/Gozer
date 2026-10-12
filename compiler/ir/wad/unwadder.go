package wad

import (
	"io"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type UnWadder struct {
	r   io.Reader
	reg *Registry
	bitCaster
}

// NewUnWadder creates a new Wad reader.
//
// The UnWadder currently requires the uint size and endianness match.
// This is because this is created for building a cache that is
// expected to be read on the same machine as it was written.
// If the read data doesn't match, then this will panic.
func NewUnWadder(r io.Reader) *UnWadder {
	wad := &UnWadder{
		r:         r,
		reg:       DefaultRegistry,
		bitCaster: newBitCaster(),
	}
	// Read wad header information and check for match.
	otherVersion := wad.Uint16() // wad version
	wasUint64 := wad.Bool()      // is uint 64-bits?
	wasBigEndian := wad.Bool()   // is big endian?
	if version != otherVersion || wasUint64 != isUint64 || wasBigEndian != isBigEndian {
		panic(faults.New(`failed to create UnWadder with different header`).
			With(`read version`, version).
			With(`read uint is 64-bits`, wasUint64).
			With(`read big endian`, wasBigEndian).
			With(`current version`, version).
			With(`current uint is 64-bits`, isUint64).
			With(`current big endian`, isBigEndian))
	}
	return wad
}

func (wad *UnWadder) read(data []byte) {
	_, err := wad.r.Read(data)
	if err != nil {
		panic(faults.New(`wad failed to read`, err))
	}
}

func (wad *UnWadder) SetRegistry(reg *Registry) { wad.reg = reg }

func (wad *UnWadder) Bool() bool { return wad.Uint8() > 0 }

func (wad *UnWadder) Int() int {
	if isUint64 {
		return int(wad.Uint64())
	}
	return int(wad.Uint32())
}

func (wad *UnWadder) Uint() uint {
	if isUint64 {
		return uint(wad.Uint64())
	}
	return uint(wad.Uint32())
}

func (wad *UnWadder) Int8() int8     { return int8(wad.Uint8()) }
func (wad *UnWadder) Int16() int16   { return int16(wad.Uint16()) }
func (wad *UnWadder) Int32() int32   { return int32(wad.Uint32()) }
func (wad *UnWadder) Int64() int64   { return int64(wad.Uint64()) }
func (wad *UnWadder) String() string { return string(wad.Bytes()) }

func (wad *UnWadder) Uint8() uint8 {
	wad.read(wad.buf1)
	return wad.getUint8()
}

func (wad *UnWadder) Uint16() uint16 {
	wad.read(wad.buf2)
	return wad.getUint16()
}

func (wad *UnWadder) Uint32() uint32 {
	wad.read(wad.buf4)
	return wad.getUint32()
}

func (wad *UnWadder) Uint64() uint64 {
	wad.read(wad.buf8)
	return wad.getUint64()
}

func (wad *UnWadder) Float32() float32 {
	wad.read(wad.buf4)
	return wad.getFloat32()
}

func (wad *UnWadder) Float64() float64 {
	wad.read(wad.buf8)
	return wad.getFloat64()
}

func (wad *UnWadder) Complex64() complex64 {
	return complex(wad.Float32(), wad.Float32())
}

func (wad *UnWadder) Complex128() complex128 {
	return complex(wad.Float64(), wad.Float64())
}

func (wad *UnWadder) Bytes() []byte {
	length := wad.Int()
	data := make([]byte, length)
	wad.read(data)
	return data
}

func (wad *UnWadder) Wad[T Waddable]() T {
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
