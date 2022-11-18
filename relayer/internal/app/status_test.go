package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	nodestore "github.com/vitelabs/vite-portal/relayer/internal/node/store"
	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/relayer/internal/util/testutil"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func TestGetNodesRecursive(t *testing.T) {
	t.Parallel()
	config := roottypes.NewDefaultConfig()
	app := NewRelayerApp(config)
	store := nodestore.NewMemoryStore()
	app.orchestrator = orchestrator.NewOrchestratorMock(store)
	chain1 := sharedtypes.DefaultSupportedChains[0].Name
	chain2 := sharedtypes.DefaultSupportedChains[1].Name
	limit := 2
	nodes, err := app.nodeService.GetNodes(chain1, 0, limit)
	require.NoError(t, err)
	require.Equal(t, 0, nodes.Total)
	app.getNodesRecursive(chain1, 0, 0)
	nodes, err = app.nodeService.GetNodes(chain1, 0, limit)
	require.NoError(t, err)
	require.Equal(t, 0, nodes.Total)
	err = store.UpsertMany([]nodetypes.Node{
		testutil.NewNode(chain1),
		testutil.NewNode(chain1),
		testutil.NewNode(chain1),
		testutil.NewNode(chain2),
	})
	require.NoError(t, err)
	app.getNodesRecursive(chain1, 0, 0)
	nodes, err = app.nodeService.GetNodes(chain1, 0, limit)
	require.NoError(t, err)
	require.Equal(t, 3, nodes.Total)
	require.Equal(t, 2, len(nodes.Entries))
	nodes, err = app.nodeService.GetNodes(chain1, 2, limit)
	require.NoError(t, err)
	require.Equal(t, 3, nodes.Total)
	require.Equal(t, 1, len(nodes.Entries))
}
