package store

import (
	"sync"

	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
)

type MemoryStore struct {
	db    collections.NameObjectCollectionI
	lock sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		lock: sync.RWMutex{},
	}
	s.Clear()
	return s
}

// ---
// Implement "StoreI" interface

func (s *MemoryStore) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.db = collections.NewNameObjectCollection()
}

func (s *MemoryStore) Close() {

}

func (s *MemoryStore) Count() int {
	return s.db.Count()
}

func (s *MemoryStore) GetByIndex(index int) (n types.Node, found bool) {
	// Assign default return values
	n = *new(types.Node)
	found = false

	e := s.db.GetByIndex(index)
	if e == nil {
		return
	}

	return e.(types.Node), true
}

func (s *MemoryStore) GetById(id string) (n types.Node, found bool) {
	e := s.db.Get(id)
	if e == nil {
		return
	}

	return e.(types.Node), true
}

func (s *MemoryStore) Upsert(n types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := n.Validate()
	if err != nil {
		return err
	}

	s.db.Set(n.Id, n)

	return nil
}

func (s *MemoryStore) Remove(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if id == "" {
		return nil
	}

	s.db.Remove(id)

	return nil
}