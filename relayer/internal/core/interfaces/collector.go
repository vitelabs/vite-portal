package interfaces

import "github.com/vitelabs/vite-portal/internal/core/types"

type CollectorI interface {
	DispatchRelayResult(r types.RelayResult) error
}