package wad

import (
	"go/token"
	"io"
)

type Writer struct {
	w   io.Writer
	reg *Registry
	bc  bitCaster
}

func NewWriter(w io.Writer) *Writer {
	header := []byte{version, 0}
	if isUint64 {
		header[1] |= 0x01
	}
	if isBigEndian {
		header[1] |= 0x02
	}
	write(w, header)
	return newWriterNoHeader(w, isUint64, isBigEndian)
}

func newWriterNoHeader(w io.Writer, useUint64, useBigEndian bool) *Writer {
	return &Writer{
		w:   w,
		reg: defaultRegistry.clone(),
		bc:  newBitCaster(useUint64, useBigEndian),
	}
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

func (wad *Writer) Int(v int)         { wad.Uint(uint(v)) }
func (wad *Writer) Int8(v int8)       { wad.Uint8(uint8(v)) }
func (wad *Writer) Int16(v int16)     { wad.Uint16(uint16(v)) }
func (wad *Writer) Int32(v int32)     { wad.Uint32(uint32(v)) }
func (wad *Writer) Int64(v int64)     { wad.Uint64(uint64(v)) }
func (wad *Writer) String(v string)   { wad.Bytes([]byte(v)) }
func (wad *Writer) Uint(v uint)       { wad.bc.writeUint(wad.w, v) }
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
	wad.Uint32(uint32(len(v)))
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

func (wad *Writer) IntSlice(v []int)         { wad.slice(wad.Int, v) }
func (wad *Writer) Int8Slice(v []int8)       { wad.slice(wad.Int8, v) }
func (wad *Writer) Int16Slice(v []int16)     { wad.slice(wad.Int16, v) }
func (wad *Writer) Int32Slice(v []int32)     { wad.slice(wad.Int32, v) }
func (wad *Writer) Int64Slice(v []int64)     { wad.slice(wad.Int64, v) }
func (wad *Writer) StringSlice(v []string)   { wad.slice(wad.String, v) }
func (wad *Writer) UintSlice(v []uint)       { wad.slice(wad.Uint, v) }
func (wad *Writer) Uint8Slice(v []uint8)     { wad.Bytes(v) }
func (wad *Writer) Uint16Slice(v []uint16)   { wad.slice(wad.Uint16, v) }
func (wad *Writer) Uint32Slice(v []uint32)   { wad.slice(wad.Uint32, v) }
func (wad *Writer) Uint64Slice(v []uint64)   { wad.slice(wad.Uint64, v) }
func (wad *Writer) Float32Slice(v []float32) { wad.slice(wad.Float32, v) }
func (wad *Writer) Float64Slice(v []float64) { wad.slice(wad.Float64, v) }
func (wad *Writer) PosSlice(v []token.Pos)   { wad.slice(wad.Pos, v) }

func (wad *Writer) WadSlice[T Waddable, S ~[]T](v S) { wad.slice(wad.Wad, v) }
