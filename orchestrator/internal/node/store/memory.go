package store

import (
	"errors"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
)

type MemoryStore struct {
	db        collections.NameObjectCollectionI
	addresses mapset.Set[string]
	lock      sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		addresses: mapset.NewSet[string](),
		lock:      sync.RWMutex{},
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
	s.addresses.Clear()
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
	if _, found := s.GetById(n.Id); found {
		return errors.New("a node with the same id already exists")
	}
	// TODO: replace with "True-Client-Ip"
	addr := n.PeerInfo.RemoteAddr
	if s.addresses.Contains(addr) {
		return errors.New("a node with the same ip address already exists")
	}

	s.db.Set(n.Id, n)
	s.addresses.Add(addr)

	return nil
}

func (s *MemoryStore) Remove(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if id == "" {
		return nil
	}

	s.db.Remove(id)
	s.addresses.Remove(id)

	return nil
}
