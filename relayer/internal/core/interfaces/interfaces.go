package interfaces

import (
	"github.com/vitelabs/vite-portal/relayer/internal/core/types"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
)

type CollectorI interface {
	DispatchRelayResult(r types.RelayResult) error
}

type ServiceI interface {
	HandleRelay(r types.Relay) (*types.RelayResponse, roottypes.Error)
}