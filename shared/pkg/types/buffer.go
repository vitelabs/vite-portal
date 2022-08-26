package types

import "bytes"

type BufferChannel struct {
	buffer bytes.Buffer
	Changed chan []byte
}

func NewBufferChannel() *BufferChannel {
	return &BufferChannel{
		buffer: bytes.Buffer{},
		Changed: make(chan []byte),
	}
}

func (b *BufferChannel) Write(p []byte) (n int, err error) {
	n, err = b.buffer.Write(p)
	b.Changed<-p
	return 
}

func (b *BufferChannel) String() string {
	return b.buffer.String()
}