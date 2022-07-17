package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

func TestGetById(t *testing.T) {
	tests := []struct {
		name string
		id string
		node *types.Node
		found bool
	}{
		{
			name: "Test nil node",
			found: false,
		},
		{
			name: "Test getting node by id",
			id: "1234",
			node: &types.Node{
				Id: "1234",
				Chain: "chain1",
				IpAddress: "0.0.0.0",
				RewardAddress: "vite_1234",
			},
			found: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewBuntdbStore()

			if tc.node != nil {
				require.NoError(t, s.Upsert(*tc.node))
			}

			n, found := s.GetById(tc.id)
			require.Equal(t, tc.found, found)

			if !found {
				assert.Empty(t, n)
			} else {
				assert.Equal(t, tc.node.Id, n.Id)
				assert.Equal(t, tc.node.Chain, n.Chain)
				assert.Equal(t, tc.node.IpAddress, n.IpAddress)
				assert.Equal(t, tc.node.RewardAddress, n.RewardAddress)
			}
		})
	}
}

func TestUpsertEmpty(t *testing.T) {
	s := NewBuntdbStore()
	err := s.Upsert(types.Node{})
	require.Error(t, err)
	require.Equal(t, "Empty node", err.Error())
}
