package wad

import (
	"encoding/binary"
	"math/bits"
)

type Waddable interface {
	Wad(w *Wadder)
	UnWad(w *UnWadder)
}

const (
	version  = uint16(1)
	isUint64 = bits.UintSize > 32
)

var isBigEndian = binary.NativeEndian.Uint16([]byte{0x12, 0x48}) == 0x1248
