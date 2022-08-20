package app

import (
	"errors"
	"fmt"

	coreservice "github.com/vitelabs/vite-portal/relayer/internal/core/service"
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/sliceutil"

	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	nodeservice "github.com/vitelabs/vite-portal/relayer/internal/node/service"
)

type RelayerCoreApp struct {
	Config      types.Config
	context     *Context
	coreService *coreservice.Service
	nodeService nodeinterfaces.ServiceI
}

func NewRelayerCoreApp(cfg types.Config, c *Context) *RelayerCoreApp {
	app := &RelayerCoreApp{
		Config:  cfg,
		context: c,
	}
	app.nodeService = nodeservice.NewService(c.nodeStore)
	app.coreService = coreservice.NewService(cfg, &c.cacheStore, app.nodeService)
	return app
}

func (app *RelayerCoreApp) HandleRelay(r coretypes.Relay) (string, error) {
	app.setClientIp(&r)
	err := app.setChain(&r)
	if err != nil {
		return "", err
	}
	if logger.DebugEnabled() {
		logger.Logger().Debug().Str("relay", fmt.Sprintf("%#v", r)).Msg("relay data")
	}
	res, err1 := app.coreService.HandleRelay(r)
	if err1 != nil {
		return "", errors.New(err1.InnerError())
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

func (app *RelayerCoreApp) setChain(r *coretypes.Relay) error {
	chains := app.nodeService.GetChains()
	if len(chains) == 0 {
		return errors.New("chains are empty")
	}
	defaultError := func(chain string) error {
		return errors.New(fmt.Sprintf("the chain '%s' is not supported", chain))
	}
	if r.Chain == "" {
		r.Chain = app.Config.HostToChainMap[r.Host]
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
