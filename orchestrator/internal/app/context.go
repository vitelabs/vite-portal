package app

import (
	relayerinterfaces "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/interfaces"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
)

type Context struct {
	relayerStore relayerinterfaces.StoreI
}

func NewContext(config types.Config) *Context {
	c := &Context{
		relayerStore: relayerstore.NewMemoryStore(),
	}
	return c
}