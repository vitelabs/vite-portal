package app

import (
	"errors"
	"fmt"

	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/sliceutil"
)

func (a *RelayerApp) HandleRelay(r coretypes.Relay) (string, error) {
	if a.status != Started {
		return "", errors.New(fmt.Sprintf("relay not possible: %d", a.status))
	}

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

func (a *RelayerApp) setClientIp(r *coretypes.Relay) {
	// Check if already set
	if r.ClientIp != "" {
		return
	}
	v := r.Payload.Headers[a.config.HeaderTrueClientIp]
	if len(v) == 0 || v[0] == "" {
		r.ClientIp = types.DefaultIpAddress
	} else {
		r.ClientIp = v[0]
	}
}

func (a *RelayerApp) setChain(r *coretypes.Relay) error {
	chains := a.nodeService.GetChains()
	if len(chains) == 0 {
		return errors.New("chains are empty")
	}
	defaultError := func(chain string) error {
		return errors.New(fmt.Sprintf("the chain '%s' is not supported", chain))
	}
	if r.Chain == "" {
		r.Chain = a.config.HostToChainMap[r.Host]
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
