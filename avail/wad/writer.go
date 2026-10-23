package wad

import (
	"go/token"
	"io"
	"reflect"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Writer struct {
	w          io.Writer
	bc         bitCaster
	regLookup  map[string]uint16
	regWriters []WriteFn
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
	return newWriterNoHeader(w, isUint64, isBigEndian, regSingleton)
}

func newWriterNoHeader(w io.Writer, useUint64, useBigEndian bool, reg *registry) *Writer {
	wad := &Writer{
		w:  w,
		bc: newBitCaster(useUint64, useBigEndian),
	}
	wad.regLookup, wad.regWriters = reg.cloneWriters()
	return wad
}

func (wad *Writer) Bool(v bool) {
	if v {
		wad.bc.writeUint8(wad.w, 0xFF)
	} else {
		wad.bc.writeUint8(wad.w, 0x00)
	}
}

func (wad *Writer) Int(v int)         { wad.bc.writeUint(wad.w, uint(v)) }
func (wad *Writer) Int8(v int8)       { wad.bc.writeUint8(wad.w, uint8(v)) }
func (wad *Writer) Int16(v int16)     { wad.bc.writeUint16(wad.w, uint16(v)) }
func (wad *Writer) Int32(v int32)     { wad.bc.writeUint32(wad.w, uint32(v)) }
func (wad *Writer) Int64(v int64)     { wad.bc.writeUint64(wad.w, uint64(v)) }
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
	wad.bc.writeFloat32(wad.w, real(v))
	wad.bc.writeFloat32(wad.w, imag(v))
}

func (wad *Writer) Complex128(v complex128) {
	wad.bc.writeFloat64(wad.w, real(v))
	wad.bc.writeFloat64(wad.w, imag(v))
}

func (wad *Writer) Bytes(v []byte) {
	wad.bc.writeUint32(wad.w, uint32(len(v)))
	write(wad.w, v)
}

func (wad *Writer) Wad[T any](v T) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer && rv.IsNil() {
		wad.bc.writeUint16(wad.w, 0)
		return
	}

	name := typeName(rv.Type())
	if index, ok := wad.regLookup[name]; ok {
		wad.bc.writeUint16(wad.w, index)
		wad.regWriters[index](wad, v)
		return
	}

	panic(faults.New(`Unable to write an unregistered Waddable`).
		With(`type`, name))
}

func (wad *Writer) slice[T any](elem func(T), values []T) {
	wad.bc.writeUint32(wad.w, uint32(len(values)))
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

func (wad *Writer) WadSlice[T any, S ~[]T](v S) { wad.slice(wad.Wad, v) }
