package service

import (
	corestore "github.com/vitelabs/vite-portal/internal/core/store"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

type testContext struct {
	cache       *corestore.CacheStore
	nodeStore   *nodestore.MemoryStore
	nodeService *nodeservice.Service
	service     *Service
}

func newTestContext() *testContext {
	config := roottypes.NewDefaultConfig()
	cache := corestore.NewCacheStore(1000)
	store := nodestore.NewMemoryStore()
	nodesvc := nodeservice.NewService(store)
	svc := NewService(config, cache, nodesvc)
	return &testContext{
		cache:       cache,
		nodeStore:   store,
		nodeService: nodesvc,
		service:     svc,
	}
}

func newTestNode(id string, chain string) nodetypes.Node {
	return nodetypes.Node{
		Id:    id,
		Chain: chain,
		IpAddress: "0.0.0.0",
	}
}