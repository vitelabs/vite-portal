package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/store"
	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/util/testutil"
)

func TestUpdateStatus_Empty(t *testing.T) {
	t.Parallel()
	svc, _, chain := newTestService(t, 0)
	svc.UpdateStatus(chain, 4, 2)
}

func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	nodeCount := 6
	svc, nodes, chain := newTestService(t, nodeCount)
	for i := 0; i < nodeCount; i++ {
		n, _ := svc.store.GetByIndex(chain, i)
		require.Equal(t, uint32(0), n.RpcClient.GetID())
		require.Equal(t, uint32(0), nodes[i].RpcClient.GetID())
	}
	svc.UpdateStatus(chain, 4, 2)
	// only 4 nodes should be updated
	for i := 0; i < 4; i++ {
		n, _ := svc.store.GetByIndex(chain, i)
		require.Equal(t, uint32(1), n.RpcClient.GetID())
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
	for i := 4; i < nodeCount; i++ {
		n, _ := svc.store.GetByIndex(chain, i)
		require.Equal(t, uint32(0), n.RpcClient.GetID())
		require.Equal(t, uint32(0), nodes[i].RpcClient.GetID())
	}
	svc.UpdateStatus(chain, 4, 2)
	// all nodes should be updated
	for i := 0; i < nodeCount; i++ {
		n, _ := svc.store.GetByIndex(chain, i)
		require.Equal(t, uint32(1), n.RpcClient.GetID())
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
}


func newTestService(t *testing.T, nodeCount int) (*Service, []nodetypes.Node, string) {
	s := store.NewMemoryStore()
	cfg := types.NewDefaultConfig()
	require.NoError(t, cfg.Validate())
	svc := NewService(cfg, s)
	chain, found := cfg.GetChains().GetById("1")
	require.True(t, found)
	nodes := make([]nodetypes.Node, 0, nodeCount)
	for i := 0; i < nodeCount; i++ {
		node := testutil.NewNode(chain.Name)
		nodes = append(nodes, node)
		s.Add(node)
	}
	require.Equal(t, nodeCount, len(nodes))
	require.Equal(t, nodeCount, svc.store.Count(chain.Name))
	return svc, nodes, chain.Name
}