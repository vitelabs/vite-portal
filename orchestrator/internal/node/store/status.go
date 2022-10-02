package store

import (
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
)

type StatusStore struct {
	processedSets map[string]*mapset.Set[string]
	lock        sync.RWMutex
}

func NewStatusStore() *StatusStore {
	s := &StatusStore{
		lock: sync.RWMutex{},
	}
	s.Clear()
	return s
}

func (s *StatusStore) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.processedSets = map[string]*mapset.Set[string]{}
}

func (s *StatusStore) GetProcessedSet(chain string) *mapset.Set[string] {
	s.lock.Lock()
	defer s.lock.Unlock()

	if existing := s.processedSets[chain]; existing == nil {
		set := mapset.NewSet[string]()
		s.processedSets[chain] = &set
	}

	return s.processedSets[chain]
}