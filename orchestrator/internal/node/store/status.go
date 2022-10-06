package store

import (
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/vitelabs/vite-portal/shared/pkg/client"
)

type StatusStore struct {
	ProcessedSet *mapset.Set[string]
	globalHeight int64
	lastUpdate   int64
	lock         sync.RWMutex
	client       *client.ViteClient
}

func NewStatusStore(client *client.ViteClient) *StatusStore {
	s := &StatusStore{
		globalHeight: 0,
		lastUpdate:   0,
		lock:         sync.RWMutex{},
		client:       client,
	}
	set := mapset.NewSet[string]()
	s.ProcessedSet = &set
	return s
}

func (s *StatusStore) GetGlobalHeight() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.globalHeight != 0 && s.lastUpdate != 0 {
		if time.Now().UnixMilli()-s.lastUpdate < 500 {
			return s.globalHeight
		}
	}
	h, err := s.client.GetSnapshotChainHeight()
	if err != nil {
		return 0
	}
	s.globalHeight = h
	s.lastUpdate = time.Now().UnixMilli()
	return h
}
