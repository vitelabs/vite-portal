package app

import (
	"fmt"

	"github.com/vitelabs/vite-portal/internal/core/service"
	"github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/logger"

	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	orchestrator "github.com/vitelabs/vite-portal/internal/orchestrator/interfaces"
)

type RelayerCoreApp struct {
	coreService service.Service
	nodeService nodeinterfaces.ServiceI
}

func NewRelayerCoreApp(o orchestrator.ClientI, c *Context) *RelayerCoreApp {
	app := &RelayerCoreApp{}
	app.nodeService = nodeservice.NewService(c.nodeStore)
	app.coreService = service.NewService(app.nodeService)
	return app
}

func (app *RelayerCoreApp) HandleRelay(r types.Relay) (string, error) {
	if GlobalConfig.Debug {
		logger.Logger().Debug().Str("relay", fmt.Sprintf("%#v", r)).Msg("relay data")
	}
	res, err := app.coreService.HandleRelay(r)
	if err != nil {
		return "", err
	}
	return res.Response, nil
}