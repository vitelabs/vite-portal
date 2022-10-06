package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/interfaces"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	config  types.Config
	context interfaces.ContextI
}

// NewService creates new instances of the relayers module service
func NewService(cfg types.Config, c interfaces.ContextI) *Service {
	s := &Service{
		config:  cfg,
		context: c,
	}
	return s
}
