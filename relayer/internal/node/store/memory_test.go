package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

func newTestNode() *types.Node {
	return &types.Node{
		Id: "1",
		Chain: "chain1",
		IpAddress: "0.0.0.0",
		RewardAddress: "vite_1",
	}
}

func TestGetChains(t *testing.T) {
	node := newTestNode()
	s := NewMemoryStore()

	c := s.GetChains()
	assert.Empty(t, c)
	assert.Equal(t, 0, len(c))
	require.NoError(t, s.Upsert(*node))

	c = s.GetChains()
	assert.Equal(t, 1, len(c))
	node.Chain = "chain2"
	require.NoError(t, s.Upsert(*node))

	c = s.GetChains()
	assert.Equal(t, 2, len(c))
}

func TestGet(t *testing.T) {
	node := newTestNode()
	s := NewMemoryStore()

	n, found := s.Get(node.Chain, node.Id)
	assert.Empty(t, n)
	assert.False(t, found)
	require.NoError(t, s.Upsert(*node))

	n, found = s.Get(node.Chain, node.Id)
	assert.NotEmpty(t, n)
	assert.True(t, found)

	assert.Equal(t, node.Id, n.Id)
	assert.Equal(t, node.Chain, n.Chain)
	assert.Equal(t, node.IpAddress, n.IpAddress)
	assert.Equal(t, node.RewardAddress, n.RewardAddress)

	node.IpAddress = "1.1.1.1"
	assert.NotEqual(t, node.IpAddress, n.IpAddress)
}

func TestGetById(t *testing.T) {
	node := newTestNode()
	s := NewMemoryStore()

	n, found := s.GetById(node.Id)
	assert.Empty(t, n)
	assert.False(t, found)
	require.NoError(t, s.Upsert(*node))

	n, found = s.GetById(node.Id)
	assert.NotEmpty(t, n)
	assert.True(t, found)
}

func TestCount(t *testing.T) {
	node := newTestNode()
	s := NewMemoryStore()

	assert.Equal(t, 0, s.Count(""))
	assert.Equal(t, 0, s.Count(node.Chain))

	require.Error(t, s.Upsert(*new(types.Node)))
	assert.Equal(t, 0, s.Count(""))
	assert.Equal(t, 0, s.Count(node.Chain))

	require.NoError(t, s.Upsert(*node))
	assert.Equal(t, 0, s.Count(""))
	assert.Equal(t, 1, s.Count(node.Chain))

	require.NoError(t, s.Upsert(*node))
	assert.Equal(t, 0, s.Count(""))
	assert.Equal(t, 1, s.Count(node.Chain))
}

func TestUpsertInvalid(t *testing.T) {
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
			err := s.Upsert(*tc.node)
			require.Error(t, err)
			require.Equal(t, "node is invalid", err.Error())
		})
	}
}

func TestRemove(t *testing.T) {
	node1 := newTestNode()
	s := NewMemoryStore()

	assert.Equal(t, 0, s.Count(node1.Chain))
	assert.Equal(t, 0, len(s.GetChains()))
	require.NoError(t, s.Upsert(*node1))
	assert.Equal(t, 1, s.Count(node1.Chain))
	assert.Equal(t, 1, len(s.GetChains()))

	node2 := newTestNode()
	node2.Id = "2"

	require.NoError(t, s.Upsert(*node2))
	assert.Equal(t, 2, s.Count(node2.Chain))
	assert.Equal(t, 1, len(s.GetChains()))

	s.Remove(node1.Chain, node1.Id)
	assert.Equal(t, 1, s.Count(node1.Chain))
	assert.Equal(t, 1, len(s.GetChains()))
	s.Remove(node2.Chain, node2.Id)
	assert.Equal(t, 0, s.Count(node2.Chain))
	assert.Equal(t, 0, len(s.GetChains()))
}

func TestClear(t *testing.T) {
	node := newTestNode()
	s := NewMemoryStore()

	s.Clear()
	assert.Equal(t, 0, s.Count(node.Chain))

	require.NoError(t, s.Upsert(*node))
	assert.Equal(t, 1, s.Count(node.Chain))

	node.Id = "2"

	require.NoError(t, s.Upsert(*node))
	assert.Equal(t, 2, s.Count(node.Chain))

	s.Clear()
	assert.Equal(t, 0, s.Count(node.Chain))
}
