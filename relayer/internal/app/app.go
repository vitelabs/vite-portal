package app

import (
	"errors"
	"fmt"

	coreservice "github.com/vitelabs/vite-portal/relayer/internal/core/service"
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/sliceutil"

	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	nodeservice "github.com/vitelabs/vite-portal/relayer/internal/node/service"
)

type RelayerApp struct {
	Config      types.Config
	context     *Context
	coreService *coreservice.Service
	orchestrator *orchestrator.Orchestrator
	nodeService nodeinterfaces.ServiceI
}

func NewRelayerApp(cfg types.Config, o *orchestrator.Orchestrator, c *Context) *RelayerApp {
	a := &RelayerApp{
		Config:  cfg,
		context: c,
		orchestrator: o,
	}
	a.nodeService = nodeservice.NewService(c.nodeStore)
	a.coreService = coreservice.NewService(cfg, &c.cacheStore, a.nodeService)
	return a
}

func (a *RelayerApp) Start(profile bool) error {
	a.startRPC(profile)
	return nil
}

func (a *RelayerApp) Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
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

func (app *RelayerApp) setClientIp(r *coretypes.Relay) {
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

func (app *RelayerApp) setChain(r *coretypes.Relay) error {
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
