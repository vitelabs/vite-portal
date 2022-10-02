package rpc

import (
	"context"
	"encoding/json"
	"strconv"
	"sync/atomic"
)

type ClientMock struct {
	idCounter uint32
}

func (c *ClientMock) GetID() uint32 {
	return c.idCounter
}

func (c *ClientMock) nextID() json.RawMessage {
	id := atomic.AddUint32(&c.idCounter, 1)
	return strconv.AppendUint(nil, uint64(id), 10)
}

func (c *ClientMock) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	c.nextID()
	return nil
}

func (c *ClientMock) Notify(ctx context.Context, method string, args ...interface{}) error {
	c.nextID()
	return nil
}