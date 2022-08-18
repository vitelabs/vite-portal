package app

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

type OrchestratorCoreApp struct {
	Config types.Config
}

func NewOrchestratorCoreApp(cfg types.Config) *OrchestratorCoreApp {
	app := &OrchestratorCoreApp{
		Config:  cfg,
	}
	return app
}