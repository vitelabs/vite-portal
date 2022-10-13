package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/util/testutil"
	"github.com/vitelabs/vite-portal/shared/pkg/client"
	sharedhandlers "github.com/vitelabs/vite-portal/shared/pkg/handler"
	sharedtestutil "github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
)

func newTestHandler(t *testing.T, nodeCount int) (*Handler, []nodetypes.Node) {
	cfg := types.NewDefaultConfig()
	require.NoError(t, cfg.Validate())
	c := types.NewContext(cfg)
	chain, found := cfg.GetChains().GetById("1")
	require.True(t, found)
	cc, err := c.GetChainContext(chain.Name)
	require.NoError(t, err)
	nodeStore := cc.GetNodeStore()
	client := client.NewViteClient(sharedtestutil.DefaultViteMainNodeUrl, time.Duration(0) * time.Millisecond)
	kafka := sharedhandlers.NewKafkaHandler(time.Duration(0) * time.Millisecond, cfg.Kafka)
	handler := NewHandler(cfg, client, kafka, cc)
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
