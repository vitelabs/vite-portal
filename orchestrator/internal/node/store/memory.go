package store

import (
	"errors"
	"sort"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
)

type MemoryStore struct {
	idMap     map[string]string
	db        map[string]collections.NameObjectCollectionI[types.Node]
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

	s.idMap = map[string]string{}
	s.db = map[string]collections.NameObjectCollectionI[types.Node]{}
	s.addresses.Clear()
}

func (s *MemoryStore) Close() {

}

func (s *MemoryStore) Count(chain string) int {
	if s.db[chain] == nil {
		return 0
	}

	return s.db[chain].Count()
}

func (s *MemoryStore) GetChains() []string {
	chains := make([]string, len(s.db))

	i := 0
	for k := range s.db {
		chains[i] = k
		i++
	}

	sort.Strings(chains)

	return chains
}

func (s *MemoryStore) Get(chain string, id string) (n types.Node, found bool) {
	if chain == "" || id == "" || s.db[chain] == nil {
		return *new(types.Node), false
	}

	node := s.db[chain].Get(id)
	if commonutil.IsEmpty(node) {
		return *new(types.Node), false
	}

	return node, true
}

func (s *MemoryStore) GetByIndex(chain string, index int) (n types.Node, found bool) {
	node := s.db[chain].GetByIndex(index)
	if commonutil.IsEmpty(node) {
		return *new(types.Node), false
	}

	return node, true
}

func (s *MemoryStore) GetById(id string) (n types.Node, found bool) {
	return s.Get(s.idMap[id], id)
}

func (s *MemoryStore) Add(n types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := n.Validate()
	if err != nil {
		return err
	}

	c := s.initChain(n.Chain)

	if _, found := s.GetById(n.Id); found {
		return errors.New("a node with the same id already exists")
	}

	addr := n.ClientIp
	if s.addresses.Contains(addr) {
		return errors.New("a node with the same ip address already exists")
	}

	c.Add(n.Id, n)
	s.idMap[n.Id] = n.Chain
	s.addresses.Add(addr)

	return nil
}

func (s *MemoryStore) Remove(chain string, id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if chain == "" || id == "" || s.db[chain] == nil {
		return nil
	}

	existing, found := s.GetById(id)
	if !found {
		return nil
	}

	s.db[chain].Remove(id)
	delete(s.idMap, id)
	s.addresses.Remove(existing.ClientIp)

	if s.Count(chain) == 0 {
		delete(s.db, chain)
	}

	return nil
}

func (s *MemoryStore) initChain(chain string) (c collections.NameObjectCollectionI[types.Node]) {
	if s.db[chain] == nil {
		s.db[chain] = collections.NewNameObjectCollection[types.Node]()
	}

	return s.db[chain]
}
