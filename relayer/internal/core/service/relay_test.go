package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	coretypes "github.com/vitelabs/vite-portal/relayer/internal/core/types"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/relayer/internal/util/testutil"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
)

func TestGetConsensusNodes_Error(t *testing.T) {
	t.Parallel()
	chain := "chain1"
	ctx := newDefaultTestContext()
	h := coretypes.SessionHeader{
		IpAddress: idutil.NewGuid(),
		Chain:     chain,
	}
	nodes, err := ctx.service.getConsensusNodes(h)
	require.Error(t, err)
	require.Empty(t, nodes)
}

func TestGetConsensusNodes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		nodes    int
		expected int
	}{
		{
			name:     "Test 1",
			nodes:    1,
			expected: 1,
		},
		{
			name:     "Test default-1",
			nodes:    roottypes.DefaultConsensusNodeCount - 1,
			expected: roottypes.DefaultConsensusNodeCount - 1,
		},
		{
			name:     "Test default",
			nodes:    2 * roottypes.DefaultConsensusNodeCount,
			expected: roottypes.DefaultConsensusNodeCount,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			chain := "chain1"
			ctx := newDefaultTestContext()
			for i := 0; i < tc.nodes; i++ {
				ctx.nodeService.PutNode(testutil.NewNode(chain))
			}
			h := coretypes.SessionHeader{
				IpAddress: idutil.NewGuid(),
				Chain:     chain,
			}
			nodes, err := ctx.service.getConsensusNodes(h)
			require.NoError(t, err)
			require.Equal(t, tc.expected, len(nodes))
			// assert nodes are unique
			ids := []string{}
			for i := 0; i < len(nodes); i++ {
				ids = append(ids, nodes[i].Id)
				ids = append(ids, nodes[i].Id)
			}
			require.Equal(t, 2*tc.expected, len(ids))
			s := generics.FilterDuplicates(ids...)
			require.Equal(t, tc.expected, len(s))
		})
	}
}
