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

func NewContext() *Context {
	c := &Context{
		cacheStore: *corestore.NewCacheStore(types.GlobalConfig.MaxSessionCacheEntries),
		nodeStore: nodestore.NewMemoryStore(),
	}
	return c
}

func InitContext() (*Context, error) {
	c := NewContext()
	return c, nil
}