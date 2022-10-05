package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/util/testutil"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func newTestService(t *testing.T, nodeCount int) (*Service, []nodetypes.Node, sharedtypes.ChainConfig) {
	cfg := types.NewDefaultConfig()
	require.NoError(t, cfg.Validate())
	c := types.NewContext(cfg)
	svc := NewService(cfg, c)
	chain, found := cfg.GetChains().GetById("1")
	require.True(t, found)
	nodes := make([]nodetypes.Node, 0, nodeCount)
	for i := 0; i < nodeCount; i++ {
		node := testutil.NewNode(chain.Name)
		nodes = append(nodes, node)
		c.GetNodeStore().Add(node)
	}
	require.Equal(t, nodeCount, len(nodes))
	require.Equal(t, nodeCount, svc.context.GetNodeStore().Count(chain.Name))
	return svc, nodes, chain
}
