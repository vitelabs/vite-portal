package service

import (
	"github.com/vitelabs/vite-portal/internal/generics"
	"github.com/vitelabs/vite-portal/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling nodes
type Service struct {
	store interfaces.StoreI
}

// NewService creates new instances of the nodes module service
func NewService(store interfaces.StoreI) Service {
	return Service{
		store: store,
	}
}

// ---
// Implement "ServiceI" interface

func (k Service) GetNodes() generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}
