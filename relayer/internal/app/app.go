package app

import (
	"errors"
	"fmt"
	"sync"
	"time"

	coreservice "github.com/vitelabs/vite-portal/relayer/internal/core/service"
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/sliceutil"

	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	nodeservice "github.com/vitelabs/vite-portal/relayer/internal/node/service"
)

type RelayerApp struct {
	id            string
	config        types.Config
	startStopLock sync.Mutex // Start/Stop are protected by an additional lock
	state         int        // Tracks state of node lifecycle
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
	c := NewContext(cfg)
	a := &RelayerApp{
		id:            idutil.NewGuid(),
		config:        cfg,
		inprocHandler: rpc.NewServer(),
		context:       c,
	}
	a.orchestrator = orchestrator.NewOrchestrator(cfg.OrchestratorWsUrl, defaultTimeout)
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

	// start RPC endpoints
	err := a.startRPC(profile)
	if err != nil {
		a.stopRPC()
		return err
	}

	a.orchestrator.Start(a.inprocHandler)

	return nil
}

func (a *RelayerApp) Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()

	a.stopRPC()
	a.context.nodeStore.Close()
}

func (a *RelayerApp) HandleRelay(r coretypes.Relay) (string, error) {
	a.setClientIp(&r)
	err := a.setChain(&r)
	if err != nil {
		return "", err
	}
	if logger.DebugEnabled() {
		logger.Logger().Debug().Str("relay", fmt.Sprintf("%#v", r)).Msg("relay data")
	}
	res, err1 := a.coreService.HandleRelay(r)
	if err1 != nil {
		return "", errors.New(err1.InnerError())
	}
	return res.Response, nil
}

func (a *RelayerApp) setClientIp(r *coretypes.Relay) {
	// Check if already set
	if r.ClientIp != "" {
		return
	}
	v := r.Payload.Headers[a.config.HeaderTrueClientIp]
	if len(v) == 0 || v[0] == "" {
		r.ClientIp = types.DefaultIpAddress
	} else {
		r.ClientIp = v[0]
	}
}

func (a *RelayerApp) setChain(r *coretypes.Relay) error {
	chains := a.nodeService.GetChains()
	if len(chains) == 0 {
		return errors.New("chains are empty")
	}
	defaultError := func(chain string) error {
		return errors.New(fmt.Sprintf("the chain '%s' is not supported", chain))
	}
	if r.Chain == "" {
		r.Chain = a.config.HostToChainMap[r.Host]
	}
	// Check if chain exists
	if r.Chain != "" {
		if !sliceutil.Contains(chains, r.Chain) {
			return defaultError(r.Chain)
		}
		return nil
	}
	// Set default chain
	r.Chain = chains[0]
	return nil
}
