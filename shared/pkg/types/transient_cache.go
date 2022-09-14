package types

import (
	"sync"
	"time"
)

type TransientContainer[T any] struct {
	Key       string
	Timestamp int64
	Item      T
}

type TransientCache[T any] struct {
	cache Cache
	lock  sync.RWMutex
}

func NewTransientCache[T any](capacity int) *TransientCache[T] {
	s := &TransientCache[T]{
		cache: *NewCache(capacity),
		lock:  sync.RWMutex{},
	}
	return s
}

func (c *TransientCache[T]) Get(key string, maxDuration int64) (item T, found bool) {
	if key == "" {
		return *new(T), false
	}
	val, found := c.cache.Get(key)
	if !found {
		return *new(T), false
	}
	container, ok := val.(TransientContainer[T])
	if !ok {
		return *new(T), false
	}
	// check if expired
	if time.Now().UnixMilli()-maxDuration > container.Timestamp {
		return *new(T), false
	}
	return container.Item, true
}

func (c *TransientCache[T]) Set(key string, timestamp int64, item T) {
	if key == "" {
		return
	}
	container := TransientContainer[T]{
		Key:       key,
		Timestamp: timestamp,
		Item:      item,
	}
	c.cache.Add(key, container)
}

func (c *TransientCache[T]) Delete(key string) {
	c.cache.Remove(key)
}

func (c *TransientCache[T]) Clear() {
	c.cache.Purge()
}
