package wad

import (
	"go/token"
	"io"
)

type Writer struct {
	w   io.Writer
	reg *Registry
	bc  *bitCaster
}

func NewWriter(w io.Writer) *Writer {
	wad := &Writer{
		w:   w,
		reg: DefaultRegistry,
		bc:  newBitCaster(),
	}
	// Write wad header information
	wad.Uint8(version)    // wad version
	wad.Bool(isUint64)    // is uint 64-bits?
	wad.Bool(isBigEndian) // is big endian?
	return wad
}

func (wad *Writer) SetRegistry(reg *Registry) { wad.reg = reg }
func (wad *Writer) Registry() *Registry       { return wad.reg }

func (wad *Writer) Bool(v bool) {
	if v {
		wad.Uint8(0xFF)
	} else {
		wad.Uint8(0x00)
	}
}

func (wad *Writer) Int(v int) {
	if isUint64 {
		wad.Uint64(uint64(v))
	} else {
		wad.Uint32(uint32(v))
	}
}

func (wad *Writer) Uint(v uint) {
	if isUint64 {
		wad.Uint64(uint64(v))
	} else {
		wad.Uint32(uint32(v))
	}
}

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

func (wad *Writer) Wad(v Waddable) {
	i := wad.reg.getIndex(v)
	wad.Uint16(i)
	if i > 0 {
		v.Wad(wad)
	}
}
