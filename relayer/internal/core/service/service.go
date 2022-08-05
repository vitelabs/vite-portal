package service

import (
	"github.com/vitelabs/vite-portal/internal/core/interfaces"
	"github.com/vitelabs/vite-portal/internal/core/store"
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/internal/types"
)

// Service maintains the link to storage and exposes getter/setter methods related to core functionalities
type Service struct {
	config        types.Config
	cache         *store.CacheStore
	nodeService   nodeinterfaces.ServiceI
	httpCollector interfaces.CollectorI
}

// NewService creates new instances of the core module service
func NewService(config types.Config, cache *store.CacheStore, nodeService nodeinterfaces.ServiceI) *Service {
	svc := &Service{
		config:      config,
		cache:       cache,
		nodeService: nodeService,
	}
	if config.HttpCollectorUrl != "" {
		svc.httpCollector = store.NewHttpCollector(config.HttpCollectorUrl, config.RpcNodeTimeout)
	}
	return svc 
}
