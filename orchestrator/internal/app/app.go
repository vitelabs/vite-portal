package app

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	nodeservice "github.com/vitelabs/vite-portal/orchestrator/internal/node/service"
	relayerservice "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/service"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/crypto"
	sharedkafka "github.com/vitelabs/vite-portal/shared/pkg/kafka"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

const (
	initializingState = iota
	runningState
	closedState
)

type OrchestratorApp struct {
	id            string
	config        types.Config
	startStopLock sync.Mutex // Start/Stop are protected by an additional lock
	state         int        // Tracks state of node lifecycle

	lock           sync.Mutex
	rpcAPIs        []rpc.API // List of APIs currently provided by the app
	rpc            *rpc.HTTPServer
	rpcAuth        *rpc.HTTPServer
	inprocHandler  *rpc.Server // In-process RPC request handler to process the API requests
	context        *types.Context
	jwtHandler     *crypto.JWTHandler
	kafka          *sharedkafka.Handler
	nodeService    *nodeservice.Service
	relayerService *relayerservice.Service
	scheduler      *gocron.Scheduler
}

func NewOrchestratorApp(cfg types.Config) *OrchestratorApp {
	defaultTimeout := time.Duration(cfg.RpcTimeout) * time.Millisecond
	c := types.NewContext(cfg)
	a := &OrchestratorApp{
		id:            idutil.NewGuid(),
		config:        cfg,
		inprocHandler: rpc.NewServer(),
		context:       c,
	}
	a.jwtHandler = crypto.NewJWTHandler([]byte(cfg.JwtSecret), time.Duration(cfg.JwtExpiryTimeout)*time.Millisecond)
	a.kafka = sharedkafka.NewHandler(defaultTimeout, cfg.Kafka)
	a.nodeService = nodeservice.NewService(cfg, a.kafka, c)
	a.relayerService = relayerservice.NewService(cfg, c.GetRelayerStore())
	a.scheduler = gocron.NewScheduler(time.UTC)
	a.InitScheduler()

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

func (a *OrchestratorApp) Start() error {
	logger.Logger().Info().Msg("Start called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()

	a.scheduler.StartAsync()

	// start RPC endpoints
	err := a.startRPC()
	if err != nil {
		a.stopRPC()
	}
	return err
}

func (a *OrchestratorApp) Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()

	a.scheduler.Stop()
	a.stopRPC()
	a.kafka.Close()
}
