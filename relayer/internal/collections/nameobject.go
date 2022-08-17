package collections

import (
	"sync"

	"github.com/vitelabs/vite-portal/relayer/internal/util/sliceutil"
)

type NameObjectCollectionI interface {
	Add(name string, obj any)
	Set(name string, obj any)
	Remove(name string)
	RemoveAt(i int)
	Get(name string) any
	GetByIndex(i int) any
	GetNameByIndex(i int) string
	Count() int
}

type NameObjectCollection struct {
	NameObjectCollectionI
	entriesMap   map[string]*nameObjectEntry
	entriesSlice []*nameObjectEntry
	lock         sync.RWMutex
}

// NewNameObjectCollection creates a new collection
func NewNameObjectCollection() *NameObjectCollection {
	return &NameObjectCollection{
		entriesMap:   map[string]*nameObjectEntry{},
		entriesSlice: []*nameObjectEntry{},
		lock:         sync.RWMutex{},
	}
}

// Add adds an entry with the specified name and object.
func (c *NameObjectCollection) Add(name string, obj any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.add(name, obj)
}

// AddMany adds many entries to the collection.
func (c *NameObjectCollection) AddMany(entries map[string]any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for name, obj := range entries {
		c.add(name, obj)
	}
}

// Set sets the object of the entry with the specified name, if found; otherwise, adds an entry with the specified name and object.
func (c *NameObjectCollection) Set(name string, obj any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	existing := c.findEntry(name)

	if existing == nil {
		c.add(name, obj)
	} else {
		c.entriesMap[name].obj = obj
	}
}

// Remove removes the entries with the specified name.
func (c *NameObjectCollection) Remove(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if name == "" {
		return
	}

	existing := c.findEntry(name)

	if existing == nil {
		return
	}

	delete(c.entriesMap, name)

	for i := len(c.entriesSlice) - 1; i >= 0; i-- {
		if c.GetNameByIndex(i) == name {
			c.entriesSlice = sliceutil.RemoveAt(c.entriesSlice, i)
			break
		}
	}
}

// RemoveAt removes the entry at the specified index.
func (c *NameObjectCollection) RemoveAt(i int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	name := c.GetNameByIndex(i)

	if name == "" {
		return
	}

	delete(c.entriesMap, name)
	c.entriesSlice = sliceutil.RemoveAt(c.entriesSlice, i)
}

// Get gets the object of the entry with the specified name.
func (c *NameObjectCollection) Get(name string) any {
	entry := c.findEntry(name)
	if entry == nil {
		return nil
	}
	return entry.obj
}

// GetByIndex gets the object of the entry with the specified index.
func (c *NameObjectCollection) GetByIndex(i int) any {
	return c.Get(c.GetNameByIndex(i))
}

// GetNameByIndex gets the name of the entry at the specified index.
func (c *NameObjectCollection) GetNameByIndex(i int) string {
	if i >= len(c.entriesSlice) || i < 0 {
		return ""
	}
	return c.entriesSlice[i].name
}

// Count gets the number of entries in this collection.
func (c *NameObjectCollection) Count() int {
	return len(c.entriesSlice)
}

func (c *NameObjectCollection) add(name string, obj any) {
	if name == "" || obj == nil {
		return
	}

	if c.entriesMap[name] != nil {
		return
	}

	entry := newNameObjectEntry(name, obj)
	c.entriesMap[name] = entry
	c.entriesSlice = append(c.entriesSlice, entry)
}

func (c *NameObjectCollection) findEntry(name string) *nameObjectEntry {
	return c.entriesMap[name]
}

type nameObjectEntry struct {
	name string
	obj  any
}

func newNameObjectEntry(name string, obj any) *nameObjectEntry {
	return &nameObjectEntry{
		name: name,
		obj:  obj,
	}
}
