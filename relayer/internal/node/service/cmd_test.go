package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
	"github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/util/testutil"
)

func TestPutNodeInvalid(t *testing.T) {
	tests := []struct {
		name string
		node *types.Node
	}{
		{
			name: "Test insert emtpy node",
			node: &types.Node{
			},
		},
		{
			name: "Test insert node with id only",
			node: &types.Node{
				Id: "1234",
			},
		},
		{
			name: "Test insert node with id and chain",
			node: &types.Node{
				Id: "1234",
				Chain: "chain1234",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := nodestore.NewMemoryStore()
			service := NewService(store)
			err := service.PutNode(*tc.node)
			require.Error(t, err)
			require.Equal(t, "node is invalid", err.Error())
			require.Equal(t, int64(0), service.LastActivityTimestamp(tc.node.Chain, types.Put))
			require.Equal(t, int64(0), service.LastActivityTimestamp(tc.node.Chain, types.Delete))
		})
	}
}

func TestLastActivityTimestamp(t *testing.T) {
	store := nodestore.NewMemoryStore()
	service := NewService(store)
	chain1 := "chain1"
	require.Equal(t, int64(0), service.LastActivityTimestamp(chain1, types.Put))
	require.Equal(t, int64(0), service.LastActivityTimestamp(chain1, types.Delete))

	node := testutil.NewNode(chain1)

	err := service.PutNode(node)
	require.NoError(t, err)
	require.Greater(t, service.LastActivityTimestamp(chain1, types.Put), int64(0))
	require.Equal(t, int64(0), service.LastActivityTimestamp(chain1, types.Delete))

	err = service.DeleteNode(node.Id)
	require.NoError(t, err)
	require.Greater(t, service.LastActivityTimestamp(chain1, types.Put), int64(0))
	require.Greater(t, service.LastActivityTimestamp(chain1, types.Delete), int64(0))
}