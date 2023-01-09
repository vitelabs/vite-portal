package store

import (
	"errors"
	"sync"

	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/shared/pkg/collections"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

type MemoryStore struct {
	db   collections.NameObjectCollectionI[types.Relayer]
	lock sync.Mutex
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

	s.db = collections.NewNameObjectCollection[types.Relayer]()
}

func (s *MemoryStore) Close() {

}

func (s *MemoryStore) Count() int {
	return s.db.Count()
}

func (s *MemoryStore) GetByIndex(index int) (r types.Relayer, found bool) {
	e := s.db.GetByIndex(index)
	if commonutil.IsEmpty(e) {
		return *new(types.Relayer), false
	}

	return e, true
}

func (s *MemoryStore) GetById(id string) (r types.Relayer, found bool) {
	e := s.db.Get(id)
	if commonutil.IsEmpty(e) {
		return
	}

	return e, true
}

func (s *MemoryStore) GetPaginated(offset, limit int) (generics.GenericPage[types.Relayer], error) {
	total := s.Count()
	result := *generics.NewGenericPage[types.Relayer]()
	result.Offset = offset
	result.Limit = limit
	result.Total = total
	if offset >= total {
		return result, nil
	}
	result.Entries = make([]types.Relayer, mathutil.Min(total-result.Offset, limit))
	count := mathutil.Min(result.Offset+result.Limit, total)
	current := 0
	for i := result.Offset; i < count; i++ {
		item, found := s.GetByIndex(i)
		if !found {
			return *generics.NewGenericPage[types.Relayer](), errors.New("inconsistent state")
		}
		result.Entries[current] = item
		current++
	}
	return result, nil
}

func (s *MemoryStore) Upsert(r types.Relayer) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := r.Validate()
	if err != nil {
		return err
	}

	s.db.Set(r.Id, r)

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
