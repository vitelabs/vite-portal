package keeper

import (
	"github.com/vitelabs/vite-portal/internal/core/types"
	ct "github.com/vitelabs/vite-portal/internal/core/types"
	rt "github.com/vitelabs/vite-portal/internal/types"
)

// HandleDispatch handles the session information for a client request
func (k Keeper) HandleDispatch(header types.SessionHeader) (*ct.DispatchResponse, rt.Error) {
	return nil, nil
}