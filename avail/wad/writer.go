package wad

import (
	"go/token"
	"io"
)

type Writer struct {
	w      io.Writer
	reg    *Registry
	bc     bitCaster
	fnInt  func(int)
	fnUint func(uint)
}

func NewWriter(w io.Writer) *Writer {
	// Write wad header information
	var flags byte
	if isUint64 {
		flags |= 0x01
	}
	if isBigEndian {
		flags |= 0x02
	}
	write(w, []byte{version, flags})

	// Create wad writer
	wad := &Writer{
		w:   w,
		reg: defaultRegistry.clone(),
		bc:  newBitCaster(isUint64, isBigEndian),
	}
	if isUint64 {
		wad.fnInt = wad.intVia64
		wad.fnUint = wad.uintVia64
	} else {
		wad.fnInt = wad.intVia32
		wad.fnUint = wad.uintVia32
	}
	return wad
}

func (wad *Writer) Registry() *Registry { return wad.reg }

func (wad *Writer) Bool(v bool) {
	if v {
		wad.Uint8(0xFF)
	} else {
		wad.Uint8(0x00)
	}
}

func (wad *Writer) intVia32(v int)   { wad.Uint32(uint32(v)) }
func (wad *Writer) intVia64(v int)   { wad.Uint64(uint64(v)) }
func (wad *Writer) uintVia32(v uint) { wad.Uint32(uint32(v)) }
func (wad *Writer) uintVia64(v uint) { wad.Uint64(uint64(v)) }

func (wad *Writer) Int(v int)         { wad.fnInt(v) }
func (wad *Writer) Uint(v uint)       { wad.fnUint(v) }
func (wad *Writer) Int8(v int8)       { wad.Uint8(uint8(v)) }
func (wad *Writer) Int16(v int16)     { wad.Uint16(uint16(v)) }
func (wad *Writer) Int32(v int32)     { wad.Uint32(uint32(v)) }
func (wad *Writer) Int64(v int64)     { wad.Uint64(uint64(v)) }
func (wad *Writer) String(v string)   { wad.Bytes([]byte(v)) }
func (wad *Writer) Uint8(v uint8)     { wad.bc.writeUint8(wad.w, v) }
func (wad *Writer) Uint16(v uint16)   { wad.bc.writeUint16(wad.w, v) }
func (wad *Writer) Uint32(v uint32)   { wad.bc.writeUint32(wad.w, v) }
func (wad *Writer) Uint64(v uint64)   { wad.bc.writeUint64(wad.w, v) }
func (wad *Writer) Float32(v float32) { wad.bc.writeFloat32(wad.w, v) }
func (wad *Writer) Float64(v float64) { wad.bc.writeFloat64(wad.w, v) }
func (wad *Writer) Pos(v token.Pos)   { wad.Int(int(v)) }

func (wad *Writer) Complex64(v complex64) {
	wad.Float32(real(v))
	wad.Float32(imag(v))
}

func (wad *Writer) Complex128(v complex128) {
	wad.Float64(real(v))
	wad.Float64(imag(v))
}

func (wad *Writer) Bytes(v []byte) {
	wad.Int(len(v))
	write(wad.w, v)
}

func (wad *Writer) Wad[T Waddable](v T) {
	i := wad.reg.toIndex(v)
	wad.Uint16(i)
	if i > 0 {
		v.Wad(wad)
	}
}

func (wad *Writer) slice[T any](elem func(T), values []T) {
	wad.Uint32(uint32(len(values)))
	for _, v := range values {
		elem(v)
	}
}

func (wad *Writer) IntSlice(v []int)         { wad.slice(wad.fnInt, v) }
func (wad *Writer) UintSlice(v []uint)       { wad.slice(wad.fnUint, v) }
func (wad *Writer) Int8Slice(v []int8)       { wad.slice(wad.Int8, v) }
func (wad *Writer) Int16Slice(v []int16)     { wad.slice(wad.Int16, v) }
func (wad *Writer) Int32Slice(v []int32)     { wad.slice(wad.Int32, v) }
func (wad *Writer) Int64Slice(v []int64)     { wad.slice(wad.Int64, v) }
func (wad *Writer) StringSlice(v []string)   { wad.slice(wad.String, v) }
func (wad *Writer) Uint8Slice(v []uint8)     { wad.slice(wad.Uint8, v) }
func (wad *Writer) Uint16Slice(v []uint16)   { wad.slice(wad.Uint16, v) }
func (wad *Writer) Uint32Slice(v []uint32)   { wad.slice(wad.Uint32, v) }
func (wad *Writer) Uint64Slice(v []uint64)   { wad.slice(wad.Uint64, v) }
func (wad *Writer) Float32Slice(v []float32) { wad.slice(wad.Float32, v) }
func (wad *Writer) Float64Slice(v []float64) { wad.slice(wad.Float64, v) }
func (wad *Writer) PosSlice(v []token.Pos)   { wad.slice(wad.Pos, v) }

func (wad *Writer) WadSlice[T Waddable, S ~[]T](v S) { wad.slice(wad.Wad, v) }
