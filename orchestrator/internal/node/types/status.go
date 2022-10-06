package types

import (
	"sync"
	"time"
)

type ChainHeight struct {
	Height int64
	LastUpdate int64
	lock sync.Mutex
}

func NewChainHeight() *ChainHeight {
	return &ChainHeight{
		Height: 0,
		LastUpdate: 0,
	}
}

func (c *ChainHeight) Update(height int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if height > c.Height {
		c.Height = height
		c.LastUpdate = time.Now().UnixMilli()
	}
}