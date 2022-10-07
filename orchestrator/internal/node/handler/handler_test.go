package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/util/testutil"
)

func newTestHandler(t *testing.T, nodeCount int) (*Handler, []nodetypes.Node) {
	cfg := types.NewDefaultConfig()
	require.NoError(t, cfg.Validate())
	c := types.NewContext(cfg)
	chain, found := cfg.GetChains().GetById("1")
	require.True(t, found)
	nodeStore, err := c.GetNodeStore(chain.Name)
	require.NoError(t, err)
	statusStore, err := c.GetStatusStore(chain.Name)
	require.NoError(t, err)
	handler := NewHandler(cfg, nodeStore, statusStore)
	nodes := make([]nodetypes.Node, 0, nodeCount)
	for i := 0; i < nodeCount; i++ {
		node := testutil.NewNode(chain.Name)
		nodes = append(nodes, node)
		nodeStore.Add(node)
	}
	require.Equal(t, nodeCount, len(nodes))
	require.Equal(t, nodeCount, nodeStore.Count())
	return handler, nodes
}
