package app

import (
	"fmt"

	"github.com/vitelabs/vite-portal/internal/core/service"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/types"

	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	orchestrator "github.com/vitelabs/vite-portal/internal/orchestrator/interfaces"
)

type RelayerCoreApp struct {
	Config      types.Config
	context     *Context
	coreService *service.Service
	nodeService nodeinterfaces.ServiceI
}

func NewRelayerCoreApp(cfg types.Config, o orchestrator.ClientI, c *Context) *RelayerCoreApp {
	app := &RelayerCoreApp{
		Config:  cfg,
		context: c,
	}
	app.nodeService = nodeservice.NewService(c.nodeStore)
	app.coreService = service.NewService(cfg, &c.cacheStore, app.nodeService)
	return app
}

func (app *RelayerCoreApp) HandleRelay(r coretypes.Relay) (string, error) {
	if app.Config.Debug {
		logger.Logger().Debug().Str("relay", fmt.Sprintf("%#v", r)).Msg("relay data")
	}
	res, err := app.coreService.HandleRelay(r)
	if err != nil {
		return "", err
	}
	return res.Response, nil
}
