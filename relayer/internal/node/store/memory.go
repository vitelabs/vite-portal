package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vitelabs/vite-portal/internal/collections"
	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

type MemoryStore struct {
	idMap map[string]string
	db    map[string]collections.NameObjectCollectionI
	lock  sync.RWMutex
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

func (s *MemoryStore) GetChains() []string {
	chains := make([]string, len(s.db))

	i := 0
	for k := range s.db {
		chains[i] = k
		i++
	}

	return chains
}

func (s *MemoryStore) Get(chain string, id string) (n types.Node, found bool) {
	// Assign default return values
	n = *new(types.Node)
	found = false

	if chain == "" || id == "" || s.db[chain] == nil {
		return
	}

	node := s.db[chain].Get(id)
	if node == nil {
		return
	}

	return node.(types.Node), true
}

func (s *MemoryStore) GetByIndex(chain string, index int) (n types.Node, found bool) {
	// Assign default return values
	n = *new(types.Node)
	found = false

	node := s.db[chain].GetByIndex(index)
	if node == nil {
		return
	}

	return node.(types.Node), true
}

func (s *MemoryStore) GetById(id string) (n types.Node, found bool) {
	return s.Get(s.idMap[id], id)
}

func (s *MemoryStore) Upsert(n types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := validateNode(n)
	if err != nil {
		return err
	}

	c := s.initChain(n.Chain)

	c.Add(n.Id, n)
	s.idMap[n.Id] = n.Chain

	return nil
}

func (s *MemoryStore) UpsertMany(n []types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return nil
}

func (s *MemoryStore) Remove(chain string, id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if chain == "" || id == "" || s.db[chain] == nil {
		return nil
	}

	s.db[chain].Remove(id)
	delete(s.idMap, id)

	return nil
}

func (s *MemoryStore) Count(chain string) int {
	if s.db[chain] == nil {
		return 0
	}

	return s.db[chain].Count()
}

func (s *MemoryStore) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.idMap = map[string]string{}
	s.db = map[string]collections.NameObjectCollectionI{}
}

func (s *MemoryStore) Close() {

}

func (s *MemoryStore) initChain(chain string) (c collections.NameObjectCollectionI) {
	if s.db[chain] == nil {
		s.db[chain] = collections.NewNameObjectCollection()
	}

	return s.db[chain]
}

func validateNode(n types.Node) error {
	if !n.IsValid() {
		err := errors.New("Trying to insert invalid node")
		logger.Logger().Error().Err(err).Str("node", fmt.Sprintf("%#v", n))
		return err
	}
	return nil
}
