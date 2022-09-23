package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/interfaces"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	store      interfaces.StoreI
}

// NewService creates new instances of the relayers module service
func NewService(cfg types.Config, store interfaces.StoreI) *Service {
	return &Service{
		store:      store,
	}
}
