package app

import (
	nodeinterfaces "github.com/vitelabs/vite-portal/orchestrator/internal/node/interfaces"
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerinterfaces "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/interfaces"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

type Context struct {
	nodeStore nodeinterfaces.StoreI
	relayerStore relayerinterfaces.StoreI
}

func NewContext(config types.Config) *Context {
	c := &Context{
		nodeStore: nodestore.NewMemoryStore(),
		relayerStore: relayerstore.NewMemoryStore(),
	}
	return c
}