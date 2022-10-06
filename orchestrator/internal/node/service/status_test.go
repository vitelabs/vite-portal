package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/client"
	"github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
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
	store := svc.context.GetNodeStore(chain.Name)
	processed := *svc.context.GetStatusStore(chain.Name).ProcessedSet
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
		n, _ := store.GetByIndex(i)
		require.Equal(t, uint32(2), n.RpcClient.GetID())
		require.Equal(t, uint32(2), nodes[i].RpcClient.GetID())
	}
}

func TestGetChainHeight(t *testing.T) {
	t.Parallel()
	start := time.Now().UnixMilli()
	svc, _, chain := newTestService(t, 0)
	svc.clients[chain.Name] = client.NewViteClient(testutil.DefaultViteMainNodeUrl)
	store := svc.context.GetStatusStore(chain.Name)
	require.Equal(t, int64(0), store.GlobalHeight.Height)
	require.Equal(t, int64(0), store.GlobalHeight.LastUpdate)

	height := svc.GetChainHeight(chain.Name)
	lastHeight := store.GlobalHeight.Height
	lastUpdate := store.GlobalHeight.LastUpdate
	require.Greater(t, height, int64(0))
	require.Equal(t, height, lastHeight)
	require.GreaterOrEqual(t, lastUpdate, start)

	time.Sleep(time.Millisecond * 5)
	height = svc.GetChainHeight(chain.Name)
	require.Equal(t, lastHeight, store.GlobalHeight.Height)
	require.Equal(t, lastUpdate, store.GlobalHeight.LastUpdate)

	time.Sleep(time.Second * 2)
	height = svc.GetChainHeight(chain.Name)
	require.Greater(t, store.GlobalHeight.Height, lastHeight)
	require.Greater(t, store.GlobalHeight.LastUpdate, lastUpdate)
}
