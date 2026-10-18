package wad

import (
	"encoding/binary"
	"io"
	"math/bits"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Waddable interface {
	Wad(*Writer)
	UnWad(*Reader)
}

const (
	version  = uint8(1)
	isUint64 = bits.UintSize > 32
)

var isBigEndian = binary.NativeEndian.Uint16([]byte{0x12, 0x48}) == 0x1248

func write(w io.Writer, data []byte) {
	if _, err := w.Write(data); err != nil {
		panic(faults.New(`wad failed to write`, err))
	}
}

func read(r io.Reader, data []byte) {
	if _, err := r.Read(data); err != nil {
		panic(faults.New(`wad failed to read`, err))
	}
}
