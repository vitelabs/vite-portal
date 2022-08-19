package types

import (
	"net"
	"time"
)

// ConnectionWrapper is a wrapper around net.Conn that sets read/write deadlines
// before every Read() or Write() call.
// Source: https://github.com/gobwas/ws-examples/blob/master/src/chat/main.go
type ConnectionWrapper struct {
	net.Conn
	t time.Duration
}

func (c ConnectionWrapper) Write(p []byte) (int, error) {
	if err := c.Conn.SetWriteDeadline(time.Now().Add(c.t)); err != nil {
		return 0, err
	}
	return c.Conn.Write(p)
}

func (c ConnectionWrapper) Read(p []byte) (int, error) {
	if err := c.Conn.SetReadDeadline(time.Now().Add(c.t)); err != nil {
		return 0, err
	}
	return c.Conn.Read(p)
}