package service

import "github.com/vitelabs/vite-portal/orchestrator/internal/node/interfaces"

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	store interfaces.StoreI
}

// NewService creates new instances of the relayers module service
func NewService(store interfaces.StoreI) *Service {
	return &Service{
		store: store,
	}
}