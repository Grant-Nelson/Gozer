//go:build !js

package wad

import "unsafe"

type bitCaster struct {
	buf1       []byte
	buf2       []byte
	buf4       []byte
	buf8       []byte
	ptrUint8   *uint8
	ptrUint16  *uint16
	ptrUint32  *uint32
	ptrUint64  *uint64
	ptrFloat32 *float32
	ptrFloat64 *float64
}

func newBitCaster() bitCaster {
	buf := [8]byte{}
	bufPtr := unsafe.Pointer(&buf[0])
	return bitCaster{
		buf1:       buf[:1],
		buf2:       buf[:2],
		buf4:       buf[:4],
		buf8:       buf[:8],
		ptrUint8:   (*uint8)(bufPtr),
		ptrUint16:  (*uint16)(bufPtr),
		ptrUint32:  (*uint32)(bufPtr),
		ptrUint64:  (*uint64)(bufPtr),
		ptrFloat32: (*float32)(bufPtr),
		ptrFloat64: (*float64)(bufPtr),
	}
}

func (bc bitCaster) getUint8() uint8     { return *bc.ptrUint8 }
func (bc bitCaster) getUint16() uint16   { return *bc.ptrUint16 }
func (bc bitCaster) getUint32() uint32   { return *bc.ptrUint32 }
func (bc bitCaster) getUint64() uint64   { return *bc.ptrUint64 }
func (bc bitCaster) getFloat32() float32 { return *bc.ptrFloat32 }
func (bc bitCaster) getFloat64() float64 { return *bc.ptrFloat64 }

func (bc bitCaster) setUint8(v uint8)     { *bc.ptrUint8 = v }
func (bc bitCaster) setUint16(v uint16)   { *bc.ptrUint16 = v }
func (bc bitCaster) setUint32(v uint32)   { *bc.ptrUint32 = v }
func (bc bitCaster) setUint64(v uint64)   { *bc.ptrUint64 = v }
func (bc bitCaster) setFloat32(v float32) { *bc.ptrFloat32 = v }
func (bc bitCaster) setFloat64(v float64) { *bc.ptrFloat64 = v }
