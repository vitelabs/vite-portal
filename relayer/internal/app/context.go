package app

import (
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	nodestore "github.com/vitelabs/vite-portal/relayer/internal/node/store"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Context struct {
	sessionCacheStore *sharedtypes.TransientCache[coretypes.Session]
	nodeStore nodeinterfaces.StoreI
}

func NewContext(config types.Config) *Context {
	c := &Context{
		sessionCacheStore: sharedtypes.NewTransientCache[coretypes.Session](config.MaxSessionCacheEntries),
		nodeStore: nodestore.NewMemoryStore(),
	}
	return c
}