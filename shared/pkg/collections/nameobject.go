package collections

import (
	"sync"

	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/sliceutil"
)

type NameObjectCollectionI[T any] interface {
	Add(name string, obj T)
	Set(name string, obj T)
	Remove(name string)
	RemoveAt(i int)
	Get(name string) T
	GetByIndex(i int) T
	GetNameByIndex(i int) string
	GetEntries() []T
	GetEnumerator() EnumeratorI[T]
	Count() int
}

// Implementation based on:
// https://github.com/microsoft/referencesource/blob/master/System/compmod/system/collections/specialized/nameobjectcollectionbase.cs
type NameObjectCollection[T any] struct {
	entriesMap   map[string]*nameObjectEntry[T]
	entriesSlice []*nameObjectEntry[T]
	lock         sync.Mutex
}

// NewNameObjectCollection creates a new collection
func NewNameObjectCollection[T any]() *NameObjectCollection[T] {
	return &NameObjectCollection[T]{
		entriesMap:   map[string]*nameObjectEntry[T]{},
		entriesSlice: []*nameObjectEntry[T]{},
		lock:         sync.Mutex{},
	}
}

// Add adds an entry with the specified name and object.
func (c *NameObjectCollection[T]) Add(name string, obj T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.add(name, obj)
}

// AddMany adds many entries to the collection.
func (c *NameObjectCollection[T]) AddMany(entries map[string]T) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for name, obj := range entries {
		c.add(name, obj)
	}
}

// Set sets the object of the entry with the specified name, if found; otherwise, adds an entry with the specified name and object.
func (c *NameObjectCollection[T]) Set(name string, obj T) {
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
func (c *NameObjectCollection[T]) Remove(name string) {
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
func (c *NameObjectCollection[T]) RemoveAt(i int) {
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
func (c *NameObjectCollection[T]) Get(name string) T {
	entry := c.findEntry(name)
	if entry == nil {
		return *new(T)
	}
	return entry.obj
}

// GetByIndex gets the object of the entry with the specified index.
func (c *NameObjectCollection[T]) GetByIndex(i int) T {
	return c.Get(c.GetNameByIndex(i))
}

// GetNameByIndex gets the name of the entry at the specified index.
func (c *NameObjectCollection[T]) GetNameByIndex(i int) string {
	if i >= len(c.entriesSlice) || i < 0 {
		return ""
	}
	return c.entriesSlice[i].name
}

// GetEntries gets all entries of the collection.
func (c *NameObjectCollection[T]) GetEntries() []T {
	c.lock.Lock()
	defer c.lock.Unlock()

	entries := make([]T, len(c.entriesSlice))
	for i, v := range c.entriesSlice {
		entries[i] = v.obj
	}

	return entries
}

// GetEnumerator gets an enumerator that can iterate through the collection.
func (c *NameObjectCollection[T]) GetEnumerator() EnumeratorI[T] {
	tmp := make([]*nameObjectEntry[T], len(c.entriesSlice))
	copy(tmp, c.entriesSlice)
	return NewNameObjectKeysEnumerator(tmp)
}

// Count gets the number of entries in this collection.
func (c *NameObjectCollection[T]) Count() int {
	return len(c.entriesSlice)
}

func (c *NameObjectCollection[T]) add(name string, obj T) {
	if name == "" || commonutil.IsEmpty(obj) {
		return
	}

	if c.entriesMap[name] != nil {
		return
	}

	entry := newNameObjectEntry(name, obj)
	c.entriesMap[name] = entry
	c.entriesSlice = append(c.entriesSlice, entry)
}

func (c *NameObjectCollection[T]) findEntry(name string) *nameObjectEntry[T] {
	return c.entriesMap[name]
}

type nameObjectEntry[T any] struct {
	name string
	obj  T
}

func newNameObjectEntry[T any](name string, obj T) *nameObjectEntry[T] {
	return &nameObjectEntry[T]{
		name: name,
		obj:  obj,
	}
}

type nameObjectKeysEnumerator[T any] struct {
	pos     int
	entries []*nameObjectEntry[T]
}

func NewNameObjectKeysEnumerator[T any](entries []*nameObjectEntry[T]) *nameObjectKeysEnumerator[T] {
	return &nameObjectKeysEnumerator[T]{
		pos:     -1,
		entries: entries,
	}
}

func (e *nameObjectKeysEnumerator[T]) MoveNext() bool {
	if e.pos < len(e.entries)-1 {
		e.pos++
		return true
	}
	e.pos = len(e.entries)
	return false
}

func (e *nameObjectKeysEnumerator[T]) Current() (curr T, found bool) {
	if e.pos < 0 || e.pos >= len(e.entries) {
		return *new(T), false
	}
	return e.entries[e.pos].obj, true
}

func (e *nameObjectKeysEnumerator[T]) Reset() {
	e.pos = -1
}
