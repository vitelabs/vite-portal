package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/util/testutil"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	sharedtestutil "github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
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
	client, err := rpc.Dial(sharedtestutil.DefaultViteBuidlNodeUrl)
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
	require.Equal(t, 0, h.statusStore.GetGlobalHeight())
	require.Equal(t, int64(0), h.statusStore.GetLastUpdate())

	height := h.updateGlobalHeight()
	lastHeight := h.statusStore.GetGlobalHeight()
	lastUpdate := h.statusStore.GetLastUpdate()
	require.Greater(t, height, 0)
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

func TestUpdateNodeStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		node types.Node
		info sharedtypes.RpcViteRuntimeInfoResponse
		status int
	}{
		{
			name: "Test emtpy runtime info",
			node: testutil.NewNode("chain1"),
			info: sharedtypes.RpcViteRuntimeInfoResponse{},
			status: -1,
		},
		{
			name: "Test runtime info",
			node: testutil.NewNode("chain1"),
			info: sharedtypes.RpcViteRuntimeInfoResponse{
				LatestSnapshot: sharedtypes.RpcViteLatestSnapshotResponse{
					Hash:   "1234",
					Height: 1234,
					Time:   1234,
				},
			},
			status: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, _ := newTestHandler(t, 0)
			h.statusStore.SetGlobalHeight(0, 3600)
			n := tc.node
			r := tc.info
			delay := 2 * time.Second
			start := time.Now().Add(-delay)
			require.NotEqual(t, start.UnixMilli(), int64(n.LastUpdate))

			err := h.updateNodeStatus(n, r, start)
			require.Error(t, err)
			require.Equal(t, "node does not exist", err.Error())

			h.nodeStore.Add(n)
			err = h.updateNodeStatus(n, r, start)
			require.NoError(t, err)
			require.NotEqual(t, start.UnixMilli(), int64(n.LastUpdate))
			require.Equal(t, sharedtypes.Int64(0), n.DelayTime)

			n, _ = h.nodeStore.GetById(n.Id)
			require.Equal(t, start.UnixMilli(), int64(n.LastUpdate))
			require.Equal(t, delay.Milliseconds(), int64(n.DelayTime))
			require.Equal(t, r.LatestSnapshot.Hash, n.LastBlock.Hash)
			require.Equal(t, r.LatestSnapshot.Height, n.LastBlock.Height)
			require.Equal(t, r.LatestSnapshot.Time, n.LastBlock.Time)
			require.Equal(t, tc.status, n.Status)
		})
	}
}
