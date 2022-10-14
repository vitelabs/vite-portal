package types

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/orchestrator/internal/interfaces"
	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"

	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Context struct {
	relayerStore *relayerstore.MemoryStore
	ipBlacklist  *sharedtypes.TransientCache[IpBlacklistItem]
	chainContexts map[string]interfaces.ChainContextI
}

func NewContext(config Config) *Context {
	c := &Context{
		relayerStore: relayerstore.NewMemoryStore(),
		ipBlacklist:  sharedtypes.NewTransientCache[IpBlacklistItem](config.MaxIpBlacklistEntries),
		chainContexts: map[string]interfaces.ChainContextI{},
	}
	for _, v := range config.GetChains().GetAll() {
		c.chainContexts[v.Name] = NewChainContext(config)
	}
	return c
}

func (c *Context) GetRelayerStore() *relayerstore.MemoryStore {
	return c.relayerStore
}

func (c *Context) GetIpBlacklist() *sharedtypes.TransientCache[IpBlacklistItem] {
	return c.ipBlacklist
}
 
func (c *Context) GetChainContext(chain string) (interfaces.ChainContextI, error) {
	cc := c.chainContexts[chain]
	if cc == nil {
		return nil, errors.New(fmt.Sprintf("chain context not found for chain '%s'", chain))
	}
	return cc, nil
}

type ChainContext struct {
	nodeStore *nodestore.MemoryStore
	statusStore *nodestore.StatusStore
}

func NewChainContext(config Config) *ChainContext {
	return &ChainContext{
		nodeStore: nodestore.NewMemoryStore(config.AllowClientIpDuplicates),
		statusStore: nodestore.NewStatusStore(),
	}
}

func (c *ChainContext) GetNodeStore() *nodestore.MemoryStore {
	return c.nodeStore
}

func (c *ChainContext) GetStatusStore() *nodestore.StatusStore {
	return c.statusStore
}