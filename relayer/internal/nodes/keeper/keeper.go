package keeper

import (
	"github.com/vitelabs/vite-portal/internal/core/generics"
	"github.com/vitelabs/vite-portal/internal/nodes/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for handling nodes
type Keeper struct {
}

// NewKeeper creates new instances of the nodes module keeper
func NewKeeper() Keeper {
	return Keeper{}
}

// ---
// Implement "KeeperI" interface

func (k Keeper) GetNodes() generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}
