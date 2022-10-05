package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/util/testutil"
)

func TestGet(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	n, found := s.GetById(node.Id)
	assert.Empty(t, n)
	assert.False(t, found)
	require.NoError(t, s.Add(node))

	n, found = s.GetById(node.Id)
	assert.NotEmpty(t, n)
	assert.True(t, found)

	assert.Equal(t, node.Id, n.Id)
	assert.Equal(t, node.Chain, n.Chain)
	assert.Equal(t, node.Commit, n.Commit)
	assert.Equal(t, node.RpcClient.GetID(), n.RpcClient.GetID())
	assert.Equal(t, uint32(0), node.RpcClient.GetID())

	n.Commit = "1234"
	n.RpcClient.Notify(nil, "")
	
	n1, found := s.GetById(node.Id)
	assert.Equal(t, n.Id, n1.Id)
	assert.Equal(t, n.Chain, n1.Chain)
	assert.NotEqual(t, n.Commit, n1.Commit)
	assert.Equal(t, node.RpcClient.GetID(), n.RpcClient.GetID())
	assert.Equal(t, uint32(1), node.RpcClient.GetID())
}

func TestGetById(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	n, found := s.GetById(node.Id)
	assert.Empty(t, n)
	assert.False(t, found)
	require.NoError(t, s.Add(node))

	n, found = s.GetById(node.Id)
	assert.NotEmpty(t, n)
	assert.True(t, found)
}

func TestCount(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	assert.Equal(t, 0, s.Count())

	require.Error(t, s.Add(*new(types.Node)))
	assert.Equal(t, 0, s.Count())

	require.NoError(t, s.Add(node))
	assert.Equal(t, 1, s.Count())

	err := s.Add(node)
	require.Error(t, err)
	assert.Equal(t, "a node with the same id already exists", err.Error())
	assert.Equal(t, 1, s.Count())
}

func TestUpsertInvalid(t *testing.T) {
	t.Parallel()
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
			s := NewMemoryStore()
			err := s.Add(*tc.node)
			require.Error(t, err)
			require.Equal(t, "node is invalid", err.Error())
		})
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()
	node1 := testutil.NewNode("chain1")
	s := NewMemoryStore()

	assert.Equal(t, 0, s.Count())
	require.NoError(t, s.Add(node1))
	assert.Equal(t, 1, s.Count())

	node2 := testutil.NewNode("chain1")

	require.NoError(t, s.Add(node2))
	assert.Equal(t, 2, s.Count())

	s.Remove(node1.Id)
	assert.Equal(t, 1, s.Count())
	s.Remove(node2.Id)
	assert.Equal(t, 0, s.Count())
}

func TestClear(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	s.Clear()
	assert.Equal(t, 0, s.Count())

	require.NoError(t, s.Add(node))
	assert.Equal(t, 1, s.Count())

	node.Id = "2"

	err := s.Add(node)
	require.Error(t, err)
	assert.Equal(t, "a node with the same ip address already exists", err.Error())
	assert.Equal(t, 1, s.Count())

	s.Clear()
	assert.Equal(t, 0, s.Count())
}
