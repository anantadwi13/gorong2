package pool

import "sync"

var (
	poolBuffer512 sync.Pool
	poolBuffer1k  sync.Pool
	poolBuffer2k  sync.Pool
	poolBuffer4k  sync.Pool
	poolBuffer8k  sync.Pool
	poolBuffer16k sync.Pool
	poolBufferAny sync.Pool
)

type Buffer struct {
	Buf []byte
}

func GetBuffer(size int) *Buffer {
	var (
		x      interface{}
		bufCap int
		buf    *Buffer
	)

	switch {
	case size <= 512:
		bufCap = 512
		x = poolBuffer512.Get()
	case size <= 1*1024:
		bufCap = 1 * 1024
		x = poolBuffer1k.Get()
	case size <= 2*1024:
		bufCap = 2 * 1024
		x = poolBuffer2k.Get()
	case size <= 4*1024:
		bufCap = 4 * 1024
		x = poolBuffer4k.Get()
	case size <= 8*1024:
		bufCap = 8 * 1024
		x = poolBuffer8k.Get()
	case size <= 16*1024:
		bufCap = 16 * 1024
		x = poolBuffer16k.Get()
	default:
		bufCap = size
		x = poolBufferAny.Get()
	}

	if x == nil {
		return &Buffer{make([]byte, size, bufCap)}
	}
	buf, ok := x.(*Buffer)
	if !ok || buf == nil {
		return &Buffer{make([]byte, size, bufCap)}
	}
	if cap(buf.Buf) < size {
		PutBuffer(buf)
		return &Buffer{make([]byte, size, bufCap)}
	}
	buf.Buf = buf.Buf[:size]
	return buf
}

func PutBuffer(buf *Buffer) {
	if buf == nil {
		return
	}

	size := cap(buf.Buf)
	switch {
	case size == 512:
		poolBuffer512.Put(buf)
	case size == 1*1024:
		poolBuffer1k.Put(buf)
	case size == 2*1024:
		poolBuffer2k.Put(buf)
	case size == 4*1024:
		poolBuffer4k.Put(buf)
	case size == 8*1024:
		poolBuffer8k.Put(buf)
	case size == 16*1024:
		poolBuffer16k.Put(buf)
	default:
		poolBufferAny.Put(buf)
	}
}
