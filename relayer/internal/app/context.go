package app

import (
	corestore "github.com/vitelabs/vite-portal/internal/core/store"
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
	"github.com/vitelabs/vite-portal/internal/types"
)

type Context struct {
	cacheStore corestore.CacheStore
	nodeStore nodeinterfaces.StoreI
}

func NewContext(config types.Config) *Context {
	c := &Context{
		cacheStore: *corestore.NewCacheStore(config.MaxSessionCacheEntries),
		nodeStore: nodestore.NewMemoryStore(),
	}
	return c
}

func InitContext(config types.Config) (*Context, error) {
	c := NewContext(config)
	return c, nil
}