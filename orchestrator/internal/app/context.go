package app

import (
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Context struct {
	nodeStore    *nodestore.MemoryStore
	relayerStore *relayerstore.MemoryStore
	ipBlacklist  *sharedtypes.TransientCache[types.IpBlacklistItem]
}

func NewContext(config types.Config) *Context {
	c := &Context{
		nodeStore:    nodestore.NewMemoryStore(),
		relayerStore: relayerstore.NewMemoryStore(),
		ipBlacklist:  sharedtypes.NewTransientCache[types.IpBlacklistItem](config.MaxIpBlacklistEntries),
	}
	return c
}
