package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	config types.Config
	store interfaces.StoreI
}

// NewService creates new instances of the relayers module service
func NewService(cfg types.Config, store interfaces.StoreI) *Service {
	return &Service{
		config: cfg,
		store: store,
	}
}