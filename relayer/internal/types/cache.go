package types

import (
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/golang-lru/simplelru"
)

var (
	GlobalSessionCache *Cache
)

// Cache is a thread-safe LRU cache with a fixed capacity.
type Cache struct {
	lru      simplelru.LRUCache
	capacity int
	lock     sync.RWMutex
}

// New creates an LRU of the given capacity.
func NewCache(capacity int) *Cache {
	lru, err := lru.New(capacity)
	if err != nil {
		panic(err)
	}
	c := &Cache{
		lru:      lru,
		capacity: capacity,
		lock:     sync.RWMutex{},
	}
	return c
}

// Get returns key's value from the cache and updates the "recently used"-ness of the key.
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	return c.lru.Get(key)
}

// Contains checks if a key exists in cache without updating the "recently used"-ness.
func (c *Cache) Contains(key string) bool {
	return c.lru.Contains(key)
}

// Peek returns key's value without updating the "recently used"-ness of the key.
func (c *Cache) Peek(key string) (value interface{}, ok bool) {
	return c.lru.Peek(key)
}

// Add adds a value to the cache, returns true if an eviction occurred and updates the "recently used"-ness of the key.
func (c *Cache) Add(key string, value interface{}) (evicted bool) {
	evicted = c.lru.Add(key, value)
	return evicted
}

// ContainsOrAdd checks if a key exists in cache without updating the "recently used"-ness.
// If not, the item is added to the cache.
// Returns whether found and whether an eviction occurred.
func (c *Cache) ContainsOrAdd(key, value interface{}) (ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	}
	evicted = c.lru.Add(key, value)
	return false, evicted
}

// PeekOrAdd checks if a key exists in cache without updating the "recently used"-ness.
// If not, the item is added to the cache.
// Returns the existing value, whether found and whether an eviction occurred.
func (c *Cache) PeekOrAdd(key, value interface{}) (existing interface{}, ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	existing, ok = c.lru.Peek(key)
	if ok {
		return existing, true, false
	}

	evicted = c.lru.Add(key, value)
	return nil, false, evicted
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) (present bool) {
	present = c.lru.Remove(key)
	return
}

// Resize changes the capacity of the cache.
func (c *Cache) Resize(size int) (evicted int) {
	evicted = c.lru.Resize(size)
	return evicted
}

// Purge clears all cache entries.
func (c *Cache) Purge() {
	c.lru.Purge()
}

// GetOldest returns the oldest entry from the cache.
func (c *Cache) GetOldest() (key string, value interface{}, ok bool) {
	k, v, ok := c.lru.GetOldest()
	if !ok {
		return
	}
	key, ok = k.(string)
	return key, v, ok
}

// RemoveOldest removes the oldest entry from the cache.
func (c *Cache) RemoveOldest() (key string, value interface{}, ok bool) {
	k, v, ok := c.lru.RemoveOldest()
	if !ok {
		return
	}
	key, ok = k.(string)
	return key, v, ok
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *Cache) Keys() []interface{} {
	keys := c.lru.Keys()
	return keys
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	length := c.lru.Len()
	return length
}

// Capacity returns the capacity of the cache.
func (c *Cache) Capacity() int {
	return c.capacity
}
