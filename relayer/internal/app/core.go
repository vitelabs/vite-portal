package app

import (
	"github.com/vitelabs/vite-portal/internal/core/service"
	"github.com/vitelabs/vite-portal/internal/core/types"

	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	orchestrator "github.com/vitelabs/vite-portal/internal/orchestrator/interfaces"
)

type RelayerCoreApp struct {
	coreService service.Service
	nodeService nodeservice.Service
}

func NewRelayerCoreApp(o orchestrator.ClientI) *RelayerCoreApp {
	app := &RelayerCoreApp{}
	app.nodeService = nodeservice.NewService()
	app.coreService = service.NewService(app.nodeService)
	return app
}

func (app *RelayerCoreApp) HandleRelay(r types.Relay) (string, error) {
	res, err := app.coreService.HandleRelay(r)
	if err != nil {
		return "", err
	}
	return res.Response, nil
}
