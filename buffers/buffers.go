package buffers

import (
	"bytes"
	"sync"
)

type xBuffer struct {
	bpool *sync.Pool
	spool *sync.Pool
}

var once sync.Once
var singleton *xBuffer

func GetInstance() *xBuffer {
	once.Do(func() {
		singleton = &xBuffer{
			bpool: &sync.Pool{
				New: func() any {
					return new(bytes.Buffer)
				},
			},
			spool: &sync.Pool{
				New: func() any {
					return make([]byte, 0, 32<<10)
				},
			},
		}
	})
	return singleton
}

func (x *xBuffer) GetBufferX() (buf *bytes.Buffer, release func()) {
	buf = x.bpool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf, func() {
		x.bpool.Put(buf)
	}
}

// GetBuffer returns a bytes.Buffer from the pool and release func to return it.
func (x *xBuffer) GetBuffer() (buf *bytes.Buffer) {
	buf = x.bpool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (x *xBuffer) PutBuffer(buf *bytes.Buffer) {
	x.bpool.Put(buf)
}

func (x *xBuffer) GetBytesX() (buf []byte, release func()) {
	buf = x.spool.Get().([]byte)
	buf = buf[0:0]
	return buf, func() {
		x.spool.Put(buf)
	}
}

func (x *xBuffer) GetBytes() (buf []byte) {
	b := x.spool.Get().([]byte)
	b = b[0:0]
	return b
}

func (x *xBuffer) PutBytes(b []byte) {
	x.spool.Put(b)
}

type poolCloser struct {
	*bytes.Buffer
	closer *xBuffer
}

func (x *poolCloser) Close() error {
	x.closer.PutBuffer(x.Buffer)
	return nil
}

func (x *xBuffer) GetReader() *poolCloser {
	r := x.GetBuffer()
	return &poolCloser{r, x}
}
