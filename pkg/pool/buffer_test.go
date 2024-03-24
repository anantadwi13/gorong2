package pool

import (
	"math"
	"sync"
	"testing"
)

func BenchmarkBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := GetBuffer(1024)
		if len(buf.Buf) != 1024 {
			b.Error("invalid buf length")
		}
		fillBuffer(buf.Buf)
		PutBuffer(buf)
	}
}
func BenchmarkBuffer2(b *testing.B) {
	pool := &sync.Pool{
		New: func() any {
			return make([]byte, 1024)
		},
	}
	for i := 0; i < b.N; i++ {
		buf := pool.Get().([]byte)
		if len(buf) != 1024 {
			b.Error("invalid buf length")
		}
		fillBuffer(buf)
		pool.Put(buf)
	}
}

func fillBuffer(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % math.MaxUint8)
	}
}
