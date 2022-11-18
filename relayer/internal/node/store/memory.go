package store

import (
	"errors"
	"sort"
	"sync"

	"github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

type MemoryStore struct {
	idMap map[string]string
	db    map[string]collections.NameObjectCollectionI[types.Node]
	lock  sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		lock: sync.Mutex{},
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

func (s *MemoryStore) GetPaginated(chain string, offset, limit int) (generics.GenericPage[types.Node], error) {
	total := s.Count(chain)
	result := *generics.NewGenericPage[types.Node]()
	result.Offset = offset
	result.Limit = limit
	result.Total = total
	if offset >= total {
		return result, nil
	}
	result.Entries = make([]types.Node, mathutil.Min(total-result.Offset, limit))
	count := mathutil.Min(result.Offset+result.Limit, total)
	current := 0
	for i := result.Offset; i < count; i++ {
		item, found := s.GetByIndex(chain, i)
		if !found {
			return *generics.NewGenericPage[types.Node](), errors.New("inconsistent state")
		}
		result.Entries[current] = item
		current++
	}
	return result, nil
}

func (s *MemoryStore) Upsert(n types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.upsert(n)
}

func (s *MemoryStore) upsert(n types.Node) error {
	err := n.Validate()
	if err != nil {
		return err
	}

	c := s.initChain(n.Chain)

	c.Add(n.Id, n)
	s.idMap[n.Id] = n.Chain

	return nil
}

func (s *MemoryStore) UpsertMany(nodes []types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, n := range nodes {
		err := s.upsert(n)
		if err != nil {
			return err
		}
	}

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