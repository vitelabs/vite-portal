package keeper

import (
	nodesinterfaces "github.com/vitelabs/vite-portal/internal/nodes/interfaces"
)

// Keeper maintains the link to storage and exposes getter/setter methods related to core functionalities
type Keeper struct {
	NodesKeeper nodesinterfaces.KeeperI
}

// NewKeeper creates new instances of the core module keeper
func NewKeeper(nodesKeeper nodesinterfaces.KeeperI) Keeper {
	return Keeper{
		NodesKeeper: nodesKeeper,
	}
}
