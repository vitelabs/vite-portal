package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corestore "github.com/vitelabs/vite-portal/internal/core/store"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/util/idutil"
)

func TestGetActualNodes_Empty(t *testing.T) {
	tests := []struct {
		name    string
		session coretypes.Session
	}{
		{
			name:    "Test emtpy session",
			session: coretypes.Session{},
		},
		{
			name: "Test emtpy session nodes",
			session: coretypes.Session{
				Nodes: []nodetypes.Node{},
			},
		},
		{
			name: "Test invalid session nodes",
			session: coretypes.Session{
				Nodes: []nodetypes.Node{
					newNode("1", "chain1"),
					newNode("2", "chain2"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := newTestContext()
			r := ctx.service.getActualNodes(tc.session)
			require.Equal(t, 0, len(r))
		})
	}
}

func TestGetActualNodes(t *testing.T) {
	tests := []struct {
		name     string
		session  coretypes.Session
		nodes    []nodetypes.Node
		expected []nodetypes.Node
	}{
		{
			name:     "Test emtpy",
			session:  coretypes.Session{},
			nodes:    []nodetypes.Node{},
			expected: []nodetypes.Node{},
		},
		{
			name:     "Test 0 node",
			session: coretypes.Session{
				Header: newSessionHeader("chain1"),
				Nodes: []nodetypes.Node{
					newNode("1", "chain1"),
					newNode("2", "chain1"),
				},
			},
			nodes:    []nodetypes.Node{},
			expected: []nodetypes.Node{},
		},
		{
			name: "Test 1 node",
			session: coretypes.Session{
				Header: newSessionHeader("chain1"),
				Nodes: []nodetypes.Node{
					newNode("1", "chain1"),
					newNode("2", "chain1"),
					newNode("3", "chain2"),
				},
			},
			nodes: []nodetypes.Node{
				newNode("1", "chain1"),
				newNode("3", "chain2"),
			},
			expected: []nodetypes.Node{
				newNode("1", "chain1"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := newTestContext()
			for _, v := range tc.nodes {
				err := ctx.nodeService.PutNode(v)
				require.NoError(t, err)
			}
			r := ctx.service.getActualNodes(tc.session)
			require.Equal(t, len(tc.expected), len(r))
			require.Equal(t, tc.expected, r)
		})
	}
}

type testContext struct {
	cache       *corestore.CacheStore
	nodeStore   *nodestore.MemoryStore
	nodeService *nodeservice.Service
	service     *Service
}

func newTestContext() *testContext {
	cache := corestore.NewCacheStore(1000)
	store := nodestore.NewMemoryStore()
	nodesvc := nodeservice.NewService(store)
	svc := NewService(cache, nodesvc)
	return &testContext{
		cache:       cache,
		nodeStore:   store,
		nodeService: nodesvc,
		service:     svc,
	}
}

func newNode(id string, chain string) nodetypes.Node {
	return nodetypes.Node{
		Id:    id,
		Chain: chain,
		IpAddress: "0.0.0.0",
	}
}

func newSessionHeader(chain string) coretypes.SessionHeader {
	return coretypes.SessionHeader{
		Chain: chain,
		IpAddress: idutil.NewGuid(),
		Timestamp: time.Now().UnixMilli(),
	}
}