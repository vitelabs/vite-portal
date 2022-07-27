package service

import (
	"github.com/vitelabs/vite-portal/internal/core/store"
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
)

// Service maintains the link to storage and exposes getter/setter methods related to core functionalities
type Service struct {
	Cache *store.CacheStore
	NodeService nodeinterfaces.ServiceI
}

// NewService creates new instances of the core module service
func NewService(cache *store.CacheStore, nodeService nodeinterfaces.ServiceI) *Service {
	return &Service{
		Cache: cache,
		NodeService: nodeService,
	}
}
