package wad

import (
	"io"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Wadder struct {
	w   io.Writer
	reg *Registry
	bitCaster
}

func NewWadder(w io.Writer) *Wadder {
	wad := &Wadder{
		w:         w,
		reg:       DefaultRegistry,
		bitCaster: newBitCaster(),
	}
	// Write wad header information
	wad.Uint16(version)   // wad version
	wad.Bool(isUint64)    // is uint 64-bits?
	wad.Bool(isBigEndian) // is big endian?
	return wad
}

func (wad *Wadder) write(data []byte) {
	_, err := wad.w.Write(data)
	if err != nil {
		panic(faults.New(`wad failed to write`, err))
	}
}

func (wad *Wadder) SetRegistry(reg *Registry) { wad.reg = reg }

func (wad *Wadder) Bool(v bool) {
	if v {
		wad.Uint8(0xFF)
	} else {
		wad.Uint8(0x00)
	}
}

func (wad *Wadder) Int(v int) {
	if isUint64 {
		wad.Uint64(uint64(v))
	} else {
		wad.Uint32(uint32(v))
	}
}

func (wad *Wadder) Uint(v uint) {
	if isUint64 {
		wad.Uint64(uint64(v))
	} else {
		wad.Uint32(uint32(v))
	}
}

func (wad *Wadder) Int8(v int8)     { wad.Uint8(uint8(v)) }
func (wad *Wadder) Int16(v int16)   { wad.Uint16(uint16(v)) }
func (wad *Wadder) Int32(v int32)   { wad.Uint32(uint32(v)) }
func (wad *Wadder) Int64(v int64)   { wad.Uint64(uint64(v)) }
func (wad *Wadder) String(v string) { wad.Bytes([]byte(v)) }

func (wad *Wadder) Uint8(v uint8) {
	wad.setUint8(v)
	wad.write(wad.buf1)
}

func (wad *Wadder) Uint16(v uint16) {
	wad.setUint16(v)
	wad.write(wad.buf2)
}

func (wad *Wadder) Uint32(v uint32) {
	wad.setUint32(v)
	wad.write(wad.buf4)
}

func (wad *Wadder) Uint64(v uint64) {
	wad.setUint64(v)
	wad.write(wad.buf8)
}

func (wad *Wadder) Float32(v float32) {
	wad.setFloat32(v)
	wad.write(wad.buf4)
}

func (wad *Wadder) Float64(v float64) {
	wad.setFloat64(v)
	wad.write(wad.buf8)
}

func (wad *Wadder) Complex64(v complex64) {
	wad.Float32(real(v))
	wad.Float32(imag(v))
}

func (wad *Wadder) Complex128(v complex128) {
	wad.Float64(real(v))
	wad.Float64(imag(v))
}

func (wad *Wadder) Bytes(v []byte) {
	wad.Int(len(v))
	wad.write(v)
}

func (wad *Wadder) Wad(v Waddable) {
	i := wad.reg.getIndex(v)
	wad.Uint16(i)
	if i > 0 {
		v.Wad(wad)
	}
}
