package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	config types.Config
	store  *store.MemoryStore
	status *store.StatusStore
}

// NewService creates new instances of the relayers module service
func NewService(cfg types.Config, s *store.MemoryStore) *Service {
	return &Service{
		config: cfg,
		store:  s,
		status: store.NewStatusStore(),
	}
}
