package service

import (
	corestore "github.com/vitelabs/vite-portal/internal/core/store"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/testutil"
)

type testContext struct {
	config      roottypes.Config
	cache       *corestore.CacheStore
	nodeStore   *nodestore.MemoryStore
	nodeService *nodeservice.Service
	service     *Service
}

func newDefaultTestContext() *testContext {
	config := roottypes.NewDefaultConfig()
	return newTestContext(config)
}

func newTestContext(config roottypes.Config) *testContext {
	cache := corestore.NewCacheStore(config.MaxSessionCacheEntries)
	store := nodestore.NewMemoryStore()
	nodesvc := nodeservice.NewService(store)
	svc := NewService(config, cache, nodesvc)
	return &testContext{
		config:      config,
		cache:       cache,
		nodeStore:   store,
		nodeService: nodesvc,
		service:     svc,
	}
}

func newTestNode(id string, chain string) nodetypes.Node {
	node := testutil.NewNode(chain)
	node.Id = id
	return node
}
