package service

import (
	"github.com/vitelabs/vite-portal/internal/core/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling nodes
type Service struct {
}

// NewService creates new instances of the nodes module service
func NewService() Service {
	return Service{}
}

// ---
// Implement "ServiceI" interface

func (k Service) GetNodes() generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}
