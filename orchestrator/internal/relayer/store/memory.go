package store

import (
	"sort"
	"sync"

	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
)

type MemoryStore struct {
	db   map[string]types.Relayer
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

	s.db = map[string]types.Relayer{}
}

func (s *MemoryStore) Close() {

}

func (s *MemoryStore) Count() int {
	return len(s.db)
}

func (s *MemoryStore) GetAll() []types.Relayer {
	relayers := make([]types.Relayer, len(s.db))

	i := 0
	for _, r := range s.db {
		relayers[i] = r
		i++
	}

	sort.Slice(relayers, func(i, j int) bool {
		return relayers[i].Id < relayers[j].Id
	})

	return relayers
}

func (s *MemoryStore) GetById(id string) (r types.Relayer, found bool) {
	// Assign default return values
	r = *new(types.Relayer)
	found = false

	if id == "" {
		return
	}

	relayer := s.db[id]
	if commonutil.IsZero(relayer) {
		return
	}

	return relayer, true
}

func (s *MemoryStore) Upsert(r types.Relayer) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	err := r.Validate()
	if err != nil {
		return err
	}

	s.db[r.Id] = r

	return nil
}

func (s *MemoryStore) Remove(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if id == "" {
		return nil
	}

	delete(s.db, id)

	return nil
}