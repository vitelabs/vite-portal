package keeper

import (
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// HandleRelay handles a read/write request to one or multiple nodes
func (k Keeper) HandleRelay(relay coretypes.Relay) (*coretypes.RelayResponse, roottypes.Error) {
	return nil, nil
}
