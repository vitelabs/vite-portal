package types

import (
	"errors"
	"fmt"

	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"

	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Context struct {
	relayerStore *relayerstore.MemoryStore
	nodeStores   map[string]*nodestore.MemoryStore
	statusStores map[string]*nodestore.StatusStore
	ipBlacklist  *sharedtypes.TransientCache[IpBlacklistItem]
}

func NewContext(config Config) *Context {
	c := &Context{
		relayerStore: relayerstore.NewMemoryStore(),
		nodeStores:   map[string]*nodestore.MemoryStore{},
		statusStores: map[string]*nodestore.StatusStore{},
		ipBlacklist:  sharedtypes.NewTransientCache[IpBlacklistItem](config.MaxIpBlacklistEntries),
	}
	for _, v := range config.GetChains().GetAll() {
		c.nodeStores[v.Name] = nodestore.NewMemoryStore()
		c.statusStores[v.Name] = nodestore.NewStatusStore()
	}
	return c
}

func (c *Context) GetRelayerStore() *relayerstore.MemoryStore {
	return c.relayerStore
}

func (c *Context) GetNodeStore(chain string) (*nodestore.MemoryStore, error) {
	s := c.nodeStores[chain]
	if s == nil {
		return nil, errors.New(fmt.Sprintf("node store not found for chain '%s'", chain))
	}
	return s, nil
}

func (c *Context) GetStatusStore(chain string) (*nodestore.StatusStore, error) {
	s := c.statusStores[chain]
	if s == nil {
		return nil, errors.New(fmt.Sprintf("status store not found for chain '%s'", chain))
	}
	return s, nil
}

func (c *Context) GetIpBlacklist() *sharedtypes.TransientCache[IpBlacklistItem] {
	return c.ipBlacklist
}
