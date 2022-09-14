package service

import (
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/core/interfaces"
	"github.com/vitelabs/vite-portal/relayer/internal/core/store"
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

// Service maintains the link to storage and exposes getter/setter methods related to core functionalities
type Service struct {
	config        types.Config
	sessionCache  *sharedtypes.TransientCache[coretypes.Session]
	nodeService   nodeinterfaces.ServiceI
	httpCollector interfaces.CollectorI
}

// NewService creates new instances of the core module service
func NewService(config types.Config, sessionCache *sharedtypes.TransientCache[coretypes.Session], nodeService nodeinterfaces.ServiceI) *Service {
	svc := &Service{
		config:       config,
		sessionCache: sessionCache,
		nodeService:  nodeService,
	}
	if config.HttpCollectorUrl != "" {
		timeout := time.Duration(config.RpcNodeTimeout) * time.Millisecond
		svc.httpCollector = store.NewHttpCollector(config.HttpCollectorUrl, config.UserAgent, timeout)
	}
	return svc
}
