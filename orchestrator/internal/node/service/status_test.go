package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateStatus_Empty(t *testing.T) {
	t.Parallel()
	svc, _, chain := newTestService(t, 0)
	svc.UpdateStatus(chain.Name, 4, 2)
}

func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	nodeCount := 6
	svc, nodes, chain := newTestService(t, nodeCount)
	nodeStore, err := svc.context.GetNodeStore(chain.Name)
	require.NoError(t, err)
	statusStore, err := svc.context.GetStatusStore(chain.Name)
	require.NoError(t, err)
	processed := *statusStore.ProcessedSet
	for i := 0; i < nodeCount; i++ {
		require.Equal(t, uint32(0), nodes[i].RpcClient.GetID())
	}
	// only 4 nodes should get updated
	svc.UpdateStatus(chain.Name, 4, 2)
	require.Equal(t, 4, processed.Cardinality())
	for i := 0; i < 4; i++ {
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
	for i := 4; i < nodeCount; i++ {
		require.Equal(t, uint32(0), nodes[i].RpcClient.GetID())
	}
	// all nodes should get updated
	svc.UpdateStatus(chain.Name, 4, 2)
	require.Equal(t, nodeCount, processed.Cardinality())
	for i := 0; i < nodeCount; i++ {
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
	// no node should get updated
	svc.UpdateStatus(chain.Name, 2*nodeCount, 2*nodeCount)
	require.Equal(t, 0, processed.Cardinality())
	for i := 0; i < nodeCount; i++ {
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
	// all nodes should get updated
	svc.UpdateStatus(chain.Name, 2*nodeCount, 2*nodeCount)
	require.Equal(t, nodeCount, processed.Cardinality())
	for i := 0; i < nodeCount; i++ {
		n, _ := nodeStore.GetByIndex(i)
		require.Equal(t, uint32(2), n.RpcClient.GetID())
		require.Equal(t, uint32(2), nodes[i].RpcClient.GetID())
	}
}
