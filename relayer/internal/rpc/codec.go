package rpc

import (
	"fmt"
	"io"
	"sync"

	"github.com/vitelabs/vite-portal/relayer/internal/util/idutil"
)

type codec struct {
	id     string                    // identifier
	closer sync.Once                 // close channel once
	closed chan interface{}          // closed on Close
	decMu  sync.Mutex                // guards the decoder
	decode func(v interface{}) error // decoder to allow multiple transports
	encMu  sync.Mutex                // guards the encoder
	encode func(v interface{}) error // encoder to allow multiple transports
	rw     io.ReadWriteCloser        // connection
}

func NewCodec(rwc io.ReadWriteCloser, encode, decode func(v interface{}) error) ServerCodec {
	return &codec{
		id:     idutil.NewGuid(),
		closed: make(chan interface{}),
		encode: encode,
		decode: decode,
		rw:     rwc,
	}
}

// Id returns the unique identifier of the codec
func (c *codec) Id() string {
	return c.id
}

// ReadRequestHeaders will read new requests without parsing the arguments. It will
// return chain collection of requests, an indication if these requests are in batch
// form or an error when the incoming message could not be read/parsed.
func (c *codec) ReadRequest() ([]byte, Error) {
	c.decMu.Lock()
	defer c.decMu.Unlock()

	var incomingMsg []byte
	if err := c.decode(&incomingMsg); err != nil {
		return nil, &invalidRequestError{err.Error()}
	}
	return incomingMsg, nil
}

// CreateErrorResponse creates an error response with the given id and error.
func (c *codec) CreateErrorResponse(err Error) interface{} {
	return fmt.Sprintf("Code: %d Message: %s", err.ErrorCode(), err.Error())
}

// Write writes the message to the client
func (c *codec) Write(res interface{}) error {
	c.encMu.Lock()
	defer c.encMu.Unlock()

	return c.encode(res)
}

// Close closes the underlying connection
func (c *codec) Close() {
	c.closer.Do(func() {
		close(c.closed)
		c.rw.Close()
	})
}

// Closed returns a channel which will be closed when Close is called
func (c *codec) Closed() <-chan interface{} {
	return c.closed
}
