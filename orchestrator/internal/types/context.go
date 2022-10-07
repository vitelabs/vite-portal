package types

import (
	"time"

	nodestore "github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	relayerstore "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/store"
	sharedclients "github.com/vitelabs/vite-portal/shared/pkg/client"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Context struct {
	nodeStores   map[string]*nodestore.MemoryStore
	relayerStore *relayerstore.MemoryStore
	statusStores map[string]*nodestore.StatusStore
	clients      map[string]*sharedclients.ViteClient
	ipBlacklist  *sharedtypes.TransientCache[IpBlacklistItem]
}

func NewContext(config Config) *Context {
	c := &Context{
		nodeStores:   map[string]*nodestore.MemoryStore{},
		relayerStore: relayerstore.NewMemoryStore(),
		statusStores: map[string]*nodestore.StatusStore{},
		clients:      map[string]*sharedclients.ViteClient{},
		ipBlacklist:  sharedtypes.NewTransientCache[IpBlacklistItem](config.MaxIpBlacklistEntries),
	}
	timeout := time.Duration(config.RpcTimeout) * time.Millisecond
	for _, v := range config.GetChains().GetAll() {
		url := v.OfficialNodeUrl
		if url == "" {
			logger.Logger().Warn().Str("chain", v.Name).Msg("OfficialNodeUrl is empty")
		}
		client := sharedclients.NewViteClient(url, timeout)
		c.clients[v.Name] = client
		c.nodeStores[v.Name] = nodestore.NewMemoryStore()
		c.statusStores[v.Name] = nodestore.NewStatusStore(client)
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
