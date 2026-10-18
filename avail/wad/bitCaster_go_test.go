//go:build !js

package wad

import (
	"math/rand"
	"testing"
)

const valueCount = 1000

type bitBuffer struct{ data []byte }

func newBitBuffer() *bitBuffer { return &bitBuffer{data: make([]byte, 8)} }

func (bb *bitBuffer) Write(p []byte) (n int, err error) {
	copy(bb.data, p)
	return len(p), nil
}

func (bb *bitBuffer) Read(p []byte) (n int, err error) {
	copy(p, bb.data)
	return len(p), nil
}

func Benchmark_Uint(b *testing.B) {
	values := make([]uint, valueCount)
	r := rand.New(rand.NewSource(0))
	for i := range valueCount {
		values[i] = uint(r.Uint64())
	}
	bb := newBitBuffer()

	test := func(b *testing.B, bc bitCaster) {
		for b.Loop() {
			for _, v := range values {
				bc.writeUint(bb, v)
				v2 := bc.readUint(bb)
				if v != v2 {
					b.Fatalf(`incorrect value returned in: %d, out: %d`, v, v2)
				}
			}
		}
	}

	b.Run(`Native`, func(b *testing.B) { test(b, newNativeCaster()) })
	b.Run(`Default`, func(b *testing.B) { test(b, newDefaultCaster(isUint64, isBigEndian)) })
}

func Benchmark_Uint8(b *testing.B) {
	values := make([]uint8, valueCount)
	r := rand.New(rand.NewSource(0))
	for i := range valueCount {
		values[i] = uint8(r.Int31n(1 << 8))
	}
	bb := newBitBuffer()

	test := func(b *testing.B, bc bitCaster) {
		for b.Loop() {
			for _, v := range values {
				bc.writeUint8(bb, v)
				v2 := bc.readUint8(bb)
				if v != v2 {
					b.Fatalf(`incorrect value returned in: %d, out: %d`, v, v2)
				}
			}
		}
	}

	b.Run(`Native`, func(b *testing.B) { test(b, newNativeCaster()) })
	b.Run(`Default`, func(b *testing.B) { test(b, newDefaultCaster(isUint64, isBigEndian)) })
}

func Benchmark_Uint16(b *testing.B) {
	values := make([]uint16, valueCount)
	r := rand.New(rand.NewSource(0))
	for i := range valueCount {
		values[i] = uint16(r.Int31n(1 << 16))
	}
	bb := newBitBuffer()

	test := func(b *testing.B, bc bitCaster) {
		for b.Loop() {
			for _, v := range values {
				bc.writeUint16(bb, v)
				v2 := bc.readUint16(bb)
				if v != v2 {
					b.Fatalf(`incorrect value returned in: %d, out: %d`, v, v2)
				}
			}
		}
	}

	b.Run(`Native`, func(b *testing.B) { test(b, newNativeCaster()) })
	b.Run(`Default`, func(b *testing.B) { test(b, newDefaultCaster(isUint64, isBigEndian)) })
}

func Benchmark_Uint32(b *testing.B) {
	values := make([]uint32, valueCount)
	r := rand.New(rand.NewSource(0))
	for i := range valueCount {
		values[i] = uint32(r.Int31())
	}
	bb := newBitBuffer()

	test := func(b *testing.B, bc bitCaster) {
		for b.Loop() {
			for _, v := range values {
				bc.writeUint32(bb, v)
				v2 := bc.readUint32(bb)
				if v != v2 {
					b.Fatalf(`incorrect value returned in: %d, out: %d`, v, v2)
				}
			}
		}
	}

	b.Run(`Native`, func(b *testing.B) { test(b, newNativeCaster()) })
	b.Run(`Default`, func(b *testing.B) { test(b, newDefaultCaster(isUint64, isBigEndian)) })
}

func Benchmark_Uint64(b *testing.B) {
	values := make([]uint64, valueCount)
	r := rand.New(rand.NewSource(0))
	for i := range valueCount {
		values[i] = uint64(r.Int63())
	}
	bb := newBitBuffer()

	test := func(b *testing.B, bc bitCaster) {
		for b.Loop() {
			for _, v := range values {
				bc.writeUint64(bb, v)
				v2 := bc.readUint64(bb)
				if v != v2 {
					b.Fatalf(`incorrect value returned in: %d, out: %d`, v, v2)
				}
			}
		}
	}

	b.Run(`Native`, func(b *testing.B) { test(b, newNativeCaster()) })
	b.Run(`Default`, func(b *testing.B) { test(b, newDefaultCaster(isUint64, isBigEndian)) })
}

func Benchmark_Float32(b *testing.B) {
	values := make([]float32, valueCount)
	r := rand.New(rand.NewSource(0))
	for i := range valueCount {
		values[i] = r.Float32()
	}
	bb := newBitBuffer()

	test := func(b *testing.B, bc bitCaster) {
		for b.Loop() {
			for _, v := range values {
				bc.writeFloat32(bb, v)
				v2 := bc.readFloat32(bb)
				if v != v2 {
					b.Fatalf(`incorrect value returned in: %f, out: %f`, v, v2)
				}
			}
		}
	}

	b.Run(`Native`, func(b *testing.B) { test(b, newNativeCaster()) })
	b.Run(`Default`, func(b *testing.B) { test(b, newDefaultCaster(isUint64, isBigEndian)) })
}

func Benchmark_Float64(b *testing.B) {
	values := make([]float64, valueCount)
	r := rand.New(rand.NewSource(0))
	for i := range valueCount {
		values[i] = r.Float64()
	}
	bb := newBitBuffer()

	test := func(b *testing.B, bc bitCaster) {
		for b.Loop() {
			for _, v := range values {
				bc.writeFloat64(bb, v)
				v2 := bc.readFloat64(bb)
				if v != v2 {
					b.Fatalf(`incorrect value returned in: %f, out: %f`, v, v2)
				}
			}
		}
	}

	b.Run(`Native`, func(b *testing.B) { test(b, newNativeCaster()) })
	b.Run(`Default`, func(b *testing.B) { test(b, newDefaultCaster(isUint64, isBigEndian)) })
}
