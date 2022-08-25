package app

import (
	"sync"

	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

type OrchestratorApp struct {
	config        *types.Config
	startStopLock sync.Mutex // Start/Stop are protected by an additional lock
	rpcAPIs       []rpc.API   // List of APIs currently provided by the app
}

func NewOrchestratorApp(cfg *types.Config) *OrchestratorApp {
	app := &OrchestratorApp{
		config: cfg,
	}

	// Register built-in APIs.
	app.rpcAPIs = append(app.rpcAPIs, app.apis()...)

	return app
}

func (a *OrchestratorApp) Start() {
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()
}

func (a *OrchestratorApp) Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	a.startStopLock.Lock()
	defer a.startStopLock.Unlock()
}

// Config returns the configuration of app.
func (a *OrchestratorApp) Config() *types.Config {
	return a.config
}