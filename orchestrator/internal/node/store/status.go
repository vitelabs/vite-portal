package store

import (
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
)

type StatusStore struct {
	requestedSets map[string]*mapset.Set[string]
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

	s.requestedSets = map[string]*mapset.Set[string]{}
}

func (s *StatusStore) GetRequestedSet(chain string) *mapset.Set[string] {
	s.lock.Lock()
	defer s.lock.Unlock()

	if existing := s.requestedSets[chain]; existing == nil {
		set := mapset.NewSet[string]()
		s.requestedSets[chain] = &set
	}

	return s.requestedSets[chain]
}