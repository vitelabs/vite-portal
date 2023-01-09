package store

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/util/testutil"
)

func TestGet(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	n, found := s.GetById(node.Id)
	require.Empty(t, n)
	require.False(t, found)
	require.NoError(t, s.Add(node))

	n, found = s.GetById(node.Id)
	require.NotEmpty(t, n)
	require.True(t, found)

	require.Equal(t, node.Id, n.Id)
	require.Equal(t, node.Chain, n.Chain)
	require.Equal(t, node.Commit, n.Commit)
	require.Equal(t, node.RpcClient.GetID(), n.RpcClient.GetID())
	require.Equal(t, uint32(0), node.RpcClient.GetID())

	n.Commit = "1234"
	n.RpcClient.Notify(nil, "")

	n1, found := s.GetById(node.Id)
	require.Equal(t, n.Id, n1.Id)
	require.Equal(t, n.Chain, n1.Chain)
	require.NotEqual(t, n.Commit, n1.Commit)
	require.Equal(t, node.RpcClient.GetID(), n.RpcClient.GetID())
	require.Equal(t, uint32(1), node.RpcClient.GetID())
}

func TestGetById(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	n, found := s.GetById(node.Id)
	require.Empty(t, n)
	require.False(t, found)
	require.NoError(t, s.Add(node))

	n, found = s.GetById(node.Id)
	require.NotEmpty(t, n)
	require.True(t, found)
}

func TestCount(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	require.Equal(t, 0, s.Count())

	require.Error(t, s.Add(*new(types.Node)))
	require.Equal(t, 0, s.Count())

	require.NoError(t, s.Add(node))
	require.Equal(t, 1, s.Count())

	err := s.Add(node)
	require.Error(t, err)
	require.Equal(t, "a node with the same id already exists", err.Error())
	require.Equal(t, 1, s.Count())
}

func TestAddInvalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		node *types.Node
	}{
		{
			name: "Test add emtpy node",
			node: &types.Node{},
		},
		{
			name: "Test add node with id only",
			node: &types.Node{
				Id: "1234",
			},
		},
		{
			name: "Test add node with id and chain",
			node: &types.Node{
				Id:    "1234",
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

func TestUpdate(t *testing.T) {
	t.Parallel()
	n := testutil.NewNode("chain1")
	lastUpdate := int64(n.LastUpdate)
	s := NewMemoryStore()

	require.Equal(t, 0, s.Count())
	require.NoError(t, s.Add(n))
	require.Equal(t, 1, s.Count())

	n.Name = "test1"
	n.LastUpdate = 1
	err := s.Update(lastUpdate-1, n)
	require.Error(t, err)
	require.Equal(t, "inconsistent state", err.Error())
	after, found := s.GetById(n.Id)
	require.True(t, found)
	require.Equal(t, lastUpdate, int64(after.LastUpdate))
	require.NotEqual(t, n.Name, after.Name) // name should not be updated

	err = s.Update(lastUpdate, n)
	require.NoError(t, err)

	after, found = s.GetById(n.Id)
	require.True(t, found)
	require.Equal(t, n.LastUpdate, after.LastUpdate)
	require.Equal(t, n.Name, after.Name) // name should be updated
}

func TestRemove(t *testing.T) {
	t.Parallel()
	node1 := testutil.NewNode("chain1")
	s := NewMemoryStore()

	require.Equal(t, 0, s.Count())
	require.NoError(t, s.Add(node1))
	require.Equal(t, 1, s.Count())

	node2 := testutil.NewNode("chain1")

	require.NoError(t, s.Add(node2))
	require.Equal(t, 2, s.Count())

	s.Remove(node1.Id)
	require.Equal(t, 1, s.Count())
	s.Remove(node2.Id)
	require.Equal(t, 0, s.Count())
}

func TestClear(t *testing.T) {
	t.Parallel()
	node := testutil.NewNode("chain1")
	s := NewMemoryStore()

	s.Clear()
	require.Equal(t, 0, s.Count())

	require.NoError(t, s.Add(node))
	require.Equal(t, 1, s.Count())

	node.Id = "2"

	err := s.Add(node)
	require.Error(t, err)
	require.Equal(t, "a node with the same ip address already exists", err.Error())
	require.Equal(t, 1, s.Count())

	s.Clear()
	require.Equal(t, 0, s.Count())
}
