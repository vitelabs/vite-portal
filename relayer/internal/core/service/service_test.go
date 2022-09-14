package service

import (
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	nodeservice "github.com/vitelabs/vite-portal/relayer/internal/node/service"
	nodestore "github.com/vitelabs/vite-portal/relayer/internal/node/store"
	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/relayer/internal/util/testutil"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type testContext struct {
	config      roottypes.Config
	cache       *sharedtypes.TransientCache[coretypes.Session]
	nodeStore   *nodestore.MemoryStore
	nodeService *nodeservice.Service
	service     *Service
}

func newDefaultTestContext() *testContext {
	config := roottypes.NewDefaultConfig()
	return newTestContext(config)
}

func newTestContext(config roottypes.Config) *testContext {
	cache := sharedtypes.NewTransientCache[coretypes.Session](config.MaxSessionCacheEntries)
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
