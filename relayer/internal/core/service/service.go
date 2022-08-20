package service

import (
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/core/interfaces"
	"github.com/vitelabs/vite-portal/relayer/internal/core/store"
	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
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
		timeout := time.Duration(config.RpcNodeTimeout) * time.Millisecond
		svc.httpCollector = store.NewHttpCollector(config.HttpCollectorUrl, config.UserAgent, timeout)
	}
	return svc
}
