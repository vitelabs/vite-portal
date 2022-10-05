package types

import (
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Context struct {
	nodeStores   map[string]*nodestore.MemoryStore
	relayerStore *relayerstore.MemoryStore
	statusStores map[string]*nodestore.StatusStore
	ipBlacklist  *sharedtypes.TransientCache[IpBlacklistItem]
}

func NewContext(config Config) *Context {
	c := &Context{
		nodeStores:   map[string]*nodestore.MemoryStore{},
		relayerStore: relayerstore.NewMemoryStore(),
		statusStores: map[string]*nodestore.StatusStore{},
		ipBlacklist:  sharedtypes.NewTransientCache[IpBlacklistItem](config.MaxIpBlacklistEntries),
	}
	for _, v := range config.GetChains().GetAll() {
		c.nodeStores[v.Name] = nodestore.NewMemoryStore()
		c.statusStores[v.Name] = nodestore.NewStatusStore()
	}
	return c
}

func (c *Context) GetNodeStore(chain string) *nodestore.MemoryStore {
	return c.nodeStores[chain]
}

func (c *Context) GetRelayerStore() *relayerstore.MemoryStore {
	return c.relayerStore
}

func (c *Context) GetStatusStore(chain string) *nodestore.StatusStore {
	return c.statusStores[chain]
}

func (c *Context) GetIpBlacklist() *sharedtypes.TransientCache[IpBlacklistItem] {
	return c.ipBlacklist
}
