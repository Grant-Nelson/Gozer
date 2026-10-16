//go:build !js

package wad

import (
	"io"
	"unsafe"
)

// nativeCaster is the smallest part of the Wad reader/writer.
//
// This implementation of it requires pointer casting between types
// which is known not to work for languages running in an environment
// where pointer casting is not possible or is very slow.
type nativeCaster struct {
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

func init() { newNativeCaster = createNativeCaster }

func createNativeCaster() bitCaster {
	buf := [8]byte{}
	bufPtr := unsafe.Pointer(&buf[0])
	return &nativeCaster{
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

func (bc *nativeCaster) readUint8(r io.Reader) uint8     { read(r, bc.buf1); return *bc.ptrUint8 }
func (bc *nativeCaster) readUint16(r io.Reader) uint16   { read(r, bc.buf2); return *bc.ptrUint16 }
func (bc *nativeCaster) readUint32(r io.Reader) uint32   { read(r, bc.buf4); return *bc.ptrUint32 }
func (bc *nativeCaster) readUint64(r io.Reader) uint64   { read(r, bc.buf8); return *bc.ptrUint64 }
func (bc *nativeCaster) readFloat32(r io.Reader) float32 { read(r, bc.buf4); return *bc.ptrFloat32 }
func (bc *nativeCaster) readFloat64(r io.Reader) float64 { read(r, bc.buf8); return *bc.ptrFloat64 }

func (bc *nativeCaster) writeUint8(w io.Writer, v uint8)     { *bc.ptrUint8 = v; write(w, bc.buf1) }
func (bc *nativeCaster) writeUint16(w io.Writer, v uint16)   { *bc.ptrUint16 = v; write(w, bc.buf2) }
func (bc *nativeCaster) writeUint32(w io.Writer, v uint32)   { *bc.ptrUint32 = v; write(w, bc.buf4) }
func (bc *nativeCaster) writeUint64(w io.Writer, v uint64)   { *bc.ptrUint64 = v; write(w, bc.buf8) }
func (bc *nativeCaster) writeFloat32(w io.Writer, v float32) { *bc.ptrFloat32 = v; write(w, bc.buf4) }
func (bc *nativeCaster) writeFloat64(w io.Writer, v float64) { *bc.ptrFloat64 = v; write(w, bc.buf8) }
