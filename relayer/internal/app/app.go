package app

import (
	"sync"
	"time"

	coreservice "github.com/vitelabs/vite-portal/relayer/internal/core/service"
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"

	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	nodeservice "github.com/vitelabs/vite-portal/relayer/internal/node/service"
)

type RelayerAppStatus int64

const (
	Unknown RelayerAppStatus = iota
	Starting
	Started
	Stopping
	Stopped
)

type RelayerApp struct {
	id            string
	config        types.Config
	startStopLock sync.Mutex // Start/Stop are protected by an additional lock
	status        RelayerAppStatus
	lock          sync.Mutex
	rpcAPIs       []rpc.API // List of APIs currently provided by the app
	rpc           *rpc.HTTPServer
	rpcAuth       *rpc.HTTPServer
	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests
	context       *Context
	coreService   *coreservice.Service
	orchestrator  *orchestrator.Orchestrator
	nodeService   nodeinterfaces.ServiceI
}

func NewRelayerApp(cfg types.Config) *RelayerApp {
	defaultTimeout := time.Duration(cfg.RpcTimeout) * time.Millisecond
	jwtExpiryTimeout := time.Duration(cfg.JwtExpiryTimeout) * time.Millisecond
	c := NewContext(cfg)
	a := &RelayerApp{
		id:            idutil.NewGuid(),
		config:        cfg,
		inprocHandler: rpc.NewServer(),
		context:       c,
	}
	a.orchestrator = orchestrator.NewOrchestrator(a.id, cfg.OrchestratorWsUrl, cfg.JwtSecret, defaultTimeout, jwtExpiryTimeout)
	a.nodeService = nodeservice.NewService(c.nodeStore)
	a.coreService = coreservice.NewService(cfg, c.sessionCacheStore, a.nodeService)

	// Register built-in APIs.
	a.rpcAPIs = append(a.rpcAPIs, a.apis()...)

	timeouts := rpc.HTTPTimeouts{
		ReadTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultTimeout * 2,
	}

	// Configure RPC servers.
	a.rpc = rpc.NewHTTPServer(timeouts)
	a.rpcAuth = rpc.NewHTTPServer(timeouts)

	return a
}

func (a *RelayerApp) Start(profile bool) error {
	logger.Logger().Info().Msg("Start called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()

	a.status = Starting

	// start RPC endpoints
	err := a.startRPC(profile)
	if err != nil {
		a.stopRPC()
		return err
	}

	a.orchestrator.Start(a.rpcAPIs)
	a.initOrchestrator()

	a.status = Started

	return nil
}

func (a *RelayerApp) Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()

	a.status = Stopping

	a.stopRPC()
	a.context.nodeStore.Close()

	a.status = Stopped
}
