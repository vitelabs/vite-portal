package interfaces

import "github.com/vitelabs/vite-portal/relayer/internal/core/types"

type CollectorI interface {
	DispatchRelayResult(r types.RelayResult) error
}