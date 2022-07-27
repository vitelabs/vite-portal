package service

import (
	"sync"
	"time"

	"github.com/vitelabs/vite-portal/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling nodes
type Service struct {
	store    interfaces.StoreI
	activity map[string]*types.NodeActivityEntry
	lock     sync.RWMutex
}

// NewService creates new instances of the nodes module service
func NewService(store interfaces.StoreI) *Service {
	return &Service{
		store:    store,
		activity: make(map[string]*types.NodeActivityEntry),
		lock:     sync.RWMutex{},
	}
}

func (s *Service) LastActivityTimestamp(chain string, a types.NodeActivity) int64 {
	if s.activity[chain] == nil {
		return 0
	}
	switch a {
	case types.Put:
		return s.activity[chain].Put
	case types.Delete:
		return s.activity[chain].Delete
	}
	return 0
}

func (s *Service) updateLastActivityTimestamp(chain string, a types.NodeActivity) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.activity[chain] == nil {
		s.activity[chain] = &types.NodeActivityEntry{}
	}
	now := time.Now().UnixMilli()
	switch a {
	case types.Put:
		s.activity[chain].Put = now
		return
	case types.Delete:
		s.activity[chain].Delete = now
		return
	}
}
