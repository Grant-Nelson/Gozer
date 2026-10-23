package wad

import (
	"encoding/binary"
	"io"
	"math"
)

type bitCaster interface {
	readUint(r io.Reader) uint
	readUint8(r io.Reader) uint8
	readUint16(r io.Reader) uint16
	readUint32(r io.Reader) uint32
	readUint64(r io.Reader) uint64
	readFloat32(r io.Reader) float32
	readFloat64(r io.Reader) float64

	writeUint(w io.Writer, v uint)
	writeUint8(w io.Writer, v uint8)
	writeUint16(w io.Writer, v uint16)
	writeUint32(w io.Writer, v uint32)
	writeUint64(w io.Writer, v uint64)
	writeFloat32(w io.Writer, v float32)
	writeFloat64(w io.Writer, v float64)
}

var newNativeCaster func() bitCaster

func newBitCaster(needUint64, needBigEndian bool) bitCaster {
	if newNativeCaster != nil && needUint64 == isUint64 && needBigEndian == isBigEndian {
		return newNativeCaster()
	}
	return newDefaultCaster(needUint64, needBigEndian)
}

type defaultCaster struct {
	bo            binary.ByteOrder
	buf1          []byte
	buf2          []byte
	buf4          []byte
	buf8          []byte
	bufU          []byte
	needUint64    bool
	needBigEndian bool
}

func newDefaultCaster(needUint64, needBigEndian bool) bitCaster {
	var bo binary.ByteOrder = binary.LittleEndian
	if needBigEndian {
		bo = binary.BigEndian
	}
	buf := [8]byte{}
	return &defaultCaster{
		bo:            bo,
		buf1:          buf[:1],
		buf2:          buf[:2],
		buf4:          buf[:4],
		buf8:          buf[:8],
		needUint64:    needUint64,
		needBigEndian: needBigEndian,
	}
}

func (bc *defaultCaster) readUint(r io.Reader) uint {
	if bc.needUint64 {
		return uint(bc.readUint64(r))
	}
	return uint(bc.readUint32(r))
}

func (bc *defaultCaster) readUint8(r io.Reader) uint8 {
	read(r, bc.buf1)
	return bc.buf1[0]
}

func (bc *defaultCaster) readUint16(r io.Reader) uint16 {
	read(r, bc.buf2)
	return bc.bo.Uint16(bc.buf2)
}

func (bc *defaultCaster) readUint32(r io.Reader) uint32 {
	read(r, bc.buf4)
	return bc.bo.Uint32(bc.buf4)
}

func (bc *defaultCaster) readUint64(r io.Reader) uint64 {
	read(r, bc.buf8)
	return bc.bo.Uint64(bc.buf8)
}

func (bc *defaultCaster) readFloat32(r io.Reader) float32 {
	return math.Float32frombits(bc.readUint32(r))
}

func (bc *defaultCaster) readFloat64(r io.Reader) float64 {
	return math.Float64frombits(bc.readUint64(r))
}

func (bc *defaultCaster) writeUint(w io.Writer, v uint) {
	if bc.needUint64 {
		bc.writeUint64(w, uint64(v))
		return
	}
	bc.writeUint32(w, uint32(v))
}

func (bc *defaultCaster) writeUint8(w io.Writer, v uint8) {
	bc.buf1[0] = v
	write(w, bc.buf1)
}

func (bc *defaultCaster) writeUint16(w io.Writer, v uint16) {
	bc.bo.PutUint16(bc.buf2, v)
	write(w, bc.buf2)
}

func (bc *defaultCaster) writeUint32(w io.Writer, v uint32) {
	bc.bo.PutUint32(bc.buf4, v)
	write(w, bc.buf4)
}

func (bc *defaultCaster) writeUint64(w io.Writer, v uint64) {
	bc.bo.PutUint64(bc.buf8, v)
	write(w, bc.buf8)
}

func (bc *defaultCaster) writeFloat32(w io.Writer, v float32) {
	bc.writeUint32(w, math.Float32bits(v))
}

func (bc *defaultCaster) writeFloat64(w io.Writer, v float64) {
	bc.writeUint64(w, math.Float64bits(v))
}
