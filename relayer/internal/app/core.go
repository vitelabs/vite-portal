package app

import (
	"errors"
	"fmt"

	coreservice "github.com/vitelabs/vite-portal/internal/core/service"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/sliceutil"

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
	err := app.setChain(&r)
	if err != nil {
		return "", err
	}
	app.setClientIp(&r)
	res, err := app.coreService.HandleRelay(r)
	if err != nil {
		return "", err
	}
	return res.Response, nil
}

func (app *RelayerCoreApp) setChain(r *coretypes.Relay) error {
	chains := app.nodeService.GetChains()
	if len(chains) == 0 {
		return errors.New("chains are empty")
	}
	defaultError := func (chain string) error {
		return errors.New(fmt.Sprintf("the chain '%s' is not supported", chain))
	}
	// Check if already set
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
