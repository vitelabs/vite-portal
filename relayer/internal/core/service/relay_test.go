package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/generics"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/idutil"
	"github.com/vitelabs/vite-portal/internal/util/testutil"
	g "github.com/zyedidia/generic"
)

func TestGetConsensusNodes_Error(t *testing.T) {
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
			s := generics.HashsetOf(uint64(len(nodes)), g.Equals[string], g.HashString, ids...)
			require.Equal(t, tc.expected, s.Size())
		})
	}
}
