package store

import (
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

type StatusStore struct {
	ProcessedSet *mapset.Set[string]
	globalHeight int64
	lastUpdate   int64
	lock         sync.Mutex
}

func NewStatusStore() *StatusStore {
	s := &StatusStore{
		globalHeight: 0,
		lastUpdate:   0,
		lock:         sync.Mutex{},
	}
	set := mapset.NewSet[string]()
	s.ProcessedSet = &set
	return s
}

func (s *StatusStore) GetGlobalHeight() int64 {
	return s.globalHeight
}

func (s *StatusStore) GetLastUpdate() int64 {
	return s.lastUpdate
}

func (s *StatusStore) SetGlobalHeight(oldValue int64, newValue int64) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.globalHeight != oldValue {
		return false
	}

	s.globalHeight = newValue
	s.lastUpdate = time.Now().UnixMilli()

	return true
}
