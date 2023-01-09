package store

import (
	"errors"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

type MemoryStore struct {
	db        collections.NameObjectCollectionI[types.Node]
	addresses mapset.Set[string]
	lock      sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		addresses: mapset.NewSet[string](),
		lock:      sync.Mutex{},
	}
	s.Clear()
	return s
}

// ---
// Implement "StoreI" interface

func (s *MemoryStore) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.db = collections.NewNameObjectCollection[types.Node]()
	s.addresses.Clear()
}

func (s *MemoryStore) Close() {

}

func (s *MemoryStore) Count() int {
	return s.db.Count()
}

func (s *MemoryStore) GetById(id string) (n types.Node, found bool) {
	node := s.db.Get(id)
	if commonutil.IsEmpty(node) {
		return *new(types.Node), false
	}

	return node, true
}

func (s *MemoryStore) GetByIndex(index int) (n types.Node, found bool) {
	node := s.db.GetByIndex(index)
	if commonutil.IsEmpty(node) {
		return *new(types.Node), false
	}

	return node, true
}

func (s *MemoryStore) GetPaginated(offset, limit int) (generics.GenericPage[types.Node], error) {
	total := s.Count()
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
		item, found := s.GetByIndex(i)
		if !found {
			return *generics.NewGenericPage[types.Node](), errors.New("inconsistent state")
		}
		result.Entries[current] = item
		current++
	}
	return result, nil
}

func (s *MemoryStore) GetEntries() []types.Node {
	return s.db.GetEntries()
}

func (s *MemoryStore) GetEnumerator() collections.EnumeratorI[types.Node] {
	return s.db.GetEnumerator()
}

func (s *MemoryStore) Add(n types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := n.Validate()
	if err != nil {
		return err
	}

	if _, found := s.GetById(n.Id); found {
		return errors.New("a node with the same id already exists")
	}

	addr := n.ClientIp
	if s.addresses.Contains(addr) {
		return errors.New("a node with the same ip address already exists")
	}

	s.db.Add(n.Id, n)
	s.addresses.Add(addr)

	return nil
}

func (s *MemoryStore) Update(lastUpdate int64, n types.Node) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := n.Validate()
	if err != nil {
		return err
	}
	
	existing, found := s.GetById(n.Id)
	if !found {
		return errors.New("node does not exist")
	}
	if lastUpdate != int64(existing.LastUpdate) {
		return errors.New("inconsistent state")
	}

	s.db.Set(n.Id, n)

	return nil
}

func (s *MemoryStore) SetStatus(id string, lastUpdate int64, status int) error {
	existing, found := s.GetById(id)
	if !found {
		return errors.New("node does not exist")
	}
	if lastUpdate != int64(existing.LastUpdate) {
		return errors.New("inconsistent state")
	}

	existing.Status = status
	s.db.Set(existing.Id, existing)
	
	return nil
}

func (s *MemoryStore) Remove(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	existing, found := s.GetById(id)
	if !found {
		return nil
	}

	s.db.Remove(id)
	s.addresses.Remove(existing.ClientIp)

	return nil
}
