package app

import (
	"fmt"

	coreservice "github.com/vitelabs/vite-portal/internal/core/service"
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
	coreService *coreservice.Service
	nodeService nodeinterfaces.ServiceI
}

func NewRelayerCoreApp(cfg types.Config, o orchestrator.ClientI, c *Context) *RelayerCoreApp {
	app := &RelayerCoreApp{
		Config:  cfg,
		context: c,
	}
	app.nodeService = nodeservice.NewService(c.nodeStore)
	app.coreService = coreservice.NewService(cfg, &c.cacheStore, app.nodeService)
	return app
}

func (app *RelayerCoreApp) HandleRelay(r coretypes.Relay) (string, error) {
	if logger.DebugEnabled() {
		logger.Logger().Debug().Str("relay", fmt.Sprintf("%#v", r)).Msg("relay data")
	}
	app.setClientIp(&r)
	// TODO: extract chain (if empty -> set first from GetChains)
	res, err := app.coreService.HandleRelay(r)
	if err != nil {
		return "", err
	}
	return res.Response, nil
}

func (app *RelayerCoreApp) setClientIp(r *coretypes.Relay) {
	// Check if already set
	if r.ClientIp != "" {
		return
	}
	v := r.Payload.Headers[app.Config.HeaderTrueClientIp]
	if len(v) == 0 || v[0] == "" {
		r.ClientIp = types.DefaultIpAddress
	} else {
		r.ClientIp = v[0]
	}
}
