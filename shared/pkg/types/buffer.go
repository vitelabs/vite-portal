package types

import "bytes"

type BufferChannel struct {
	bytes.Buffer
	Changed chan []byte
}

func (b *BufferChannel) Write(p []byte) (n int, err error) {
	b.Changed<-p
	return b.Buffer.Write(p)
}