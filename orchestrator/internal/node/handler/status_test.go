package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
)

func TestUpdateStatus_Empty(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t, 0)
	h.UpdateStatus(4, 2)
}

func TestUpdateStatus(t *testing.T) {
	t.Parallel()
	nodeCount := 6
	h, nodes := newTestHandler(t, nodeCount)
	processed := *h.statusStore.ProcessedSet
	for i := 0; i < nodeCount; i++ {
		require.Equal(t, uint32(0), nodes[i].RpcClient.GetID())
	}
	// only 4 nodes should get updated
	h.UpdateStatus(4, 2)
	require.Equal(t, 4, processed.Cardinality())
	for i := 0; i < 4; i++ {
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
	for i := 4; i < nodeCount; i++ {
		require.Equal(t, uint32(0), nodes[i].RpcClient.GetID())
	}
	// all nodes should get updated
	h.UpdateStatus(4, 2)
	require.Equal(t, nodeCount, processed.Cardinality())
	for i := 0; i < nodeCount; i++ {
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
	// no node should get updated
	h.UpdateStatus(2*nodeCount, 2*nodeCount)
	require.Equal(t, 0, processed.Cardinality())
	for i := 0; i < nodeCount; i++ {
		require.Equal(t, uint32(1), nodes[i].RpcClient.GetID())
	}
	// all nodes should get updated
	h.UpdateStatus(2*nodeCount, 2*nodeCount)
	require.Equal(t, nodeCount, processed.Cardinality())
	for i := 0; i < nodeCount; i++ {
		n, _ := h.nodeStore.GetByIndex(i)
		require.Equal(t, uint32(2), n.RpcClient.GetID())
		require.Equal(t, uint32(2), nodes[i].RpcClient.GetID())
	}
}

func TestGetRuntimeInfo_Empty(t *testing.T) {
	t.Parallel()
	h, nodes := newTestHandler(t, 1)
	r, err := h.getRuntimeInfo(nodes[0])
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Empty(t, r)
	require.Empty(t, r.LatestSnapshot)
}

func TestGetRuntimeInfo(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t, 0)
	client, err := rpc.Dial(testutil.DefaultViteBuidlNodeUrl)
	node := types.Node{
		RpcClient: client,
	}
	r, err := h.getRuntimeInfo(node)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.NotEmpty(t, r)
	require.NotEmpty(t, r.LatestSnapshot)
}

func TestUpdateGlobalHeight(t *testing.T) {
	t.Parallel()
	start := time.Now().UnixMilli()
	h, _ := newTestHandler(t, 0)
	require.Equal(t, int64(0), h.statusStore.GetGlobalHeight())
	require.Equal(t, int64(0), h.statusStore.GetLastUpdate())

	height := h.updateGlobalHeight()
	lastHeight := h.statusStore.GetGlobalHeight()
	lastUpdate := h.statusStore.GetLastUpdate()
	require.Greater(t, height, int64(0))
	require.Equal(t, height, lastHeight)
	require.GreaterOrEqual(t, lastUpdate, start)

	time.Sleep(time.Millisecond * 5)
	height = h.updateGlobalHeight()
	require.Equal(t, lastHeight, h.statusStore.GetGlobalHeight())
	require.Equal(t, lastUpdate, h.statusStore.GetLastUpdate())

	time.Sleep(time.Second * 2)
	height = h.updateGlobalHeight()
	require.Greater(t, h.statusStore.GetGlobalHeight(), lastHeight)
	require.Greater(t, h.statusStore.GetLastUpdate(), lastUpdate)
}