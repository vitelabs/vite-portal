package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/util/idutil"
)

func TestHandleSession_Error(t *testing.T) {
	tests := []struct {
		name          string
		header        coretypes.SessionHeader
		expectedError error
	}{
		{
			name:          "Test emtpy session",
			header:        coretypes.SessionHeader{},
			expectedError: errors.New("invalid chain: empty"),
		},
		{
			name:          "Test no nodes",
			header:        newSessionHeader("chain1"),
			expectedError: errors.New("invalid chain: no nodes"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := newDefaultTestContext()
			r, err := ctx.service.HandleSession(tc.header)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.InnerError())
				require.Empty(t, r)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandleSession_SingleNode(t *testing.T) {
	chain := "chain1"
	h := newSessionHeader(chain)
	ctx := newDefaultTestContext()
	ctx.nodeService.PutNode(newNode(chain))
	r, err := ctx.service.HandleSession(h)
	require.NoError(t, err)
	require.NotEmpty(t, r)
	require.Equal(t, 1, len(r.Nodes))
}

func TestHandleSession(t *testing.T) {
	chain := "chain1"
	h := newSessionHeader(chain)
	ctx := newDefaultTestContext()
	for i := 0; i < 2*ctx.config.SessionNodeCount; i++ {
		ctx.nodeService.PutNode(newNode(chain))
		ctx.nodeService.PutNode(newNode("chain2"))
	}
	r, err := ctx.service.HandleSession(h)
	require.NoError(t, err)
	require.NotEmpty(t, r)
	require.Equal(t, ctx.config.SessionNodeCount, len(r.Nodes))
}

func TestGetActualNodes_Empty(t *testing.T) {
	tests := []struct {
		name    string
		session coretypes.Session
	}{
		{
			name:    "Test emtpy session",
			session: coretypes.Session{},
		},
		{
			name: "Test emtpy session nodes",
			session: coretypes.Session{
				Nodes: []nodetypes.Node{},
			},
		},
		{
			name: "Test invalid session nodes",
			session: coretypes.Session{
				Nodes: []nodetypes.Node{
					newTestNode("1", "chain1"),
					newTestNode("2", "chain2"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := newDefaultTestContext()
			r := ctx.service.getActualNodes(tc.session)
			require.Equal(t, 0, len(r))
		})
	}
}

func TestGetActualNodes(t *testing.T) {
	tests := []struct {
		name     string
		session  coretypes.Session
		nodes    []nodetypes.Node
		expected []nodetypes.Node
	}{
		{
			name:     "Test emtpy",
			session:  coretypes.Session{},
			nodes:    []nodetypes.Node{},
			expected: []nodetypes.Node{},
		},
		{
			name: "Test 0 node",
			session: coretypes.Session{
				Header: newSessionHeader("chain1"),
				Nodes: []nodetypes.Node{
					newTestNode("1", "chain1"),
					newTestNode("2", "chain1"),
				},
			},
			nodes:    []nodetypes.Node{},
			expected: []nodetypes.Node{},
		},
		{
			name: "Test 1 node",
			session: coretypes.Session{
				Header: newSessionHeader("chain1"),
				Nodes: []nodetypes.Node{
					newTestNode("1", "chain1"),
					newTestNode("2", "chain1"),
					newTestNode("3", "chain2"),
				},
			},
			nodes: []nodetypes.Node{
				newTestNode("1", "chain1"),
				newTestNode("3", "chain2"),
			},
			expected: []nodetypes.Node{
				newTestNode("1", "chain1"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := newDefaultTestContext()
			for _, v := range tc.nodes {
				err := ctx.nodeService.PutNode(v)
				require.NoError(t, err)
			}
			r := ctx.service.getActualNodes(tc.session)
			require.Equal(t, len(tc.expected), len(r))
			require.Equal(t, tc.expected, r)
		})
	}
}

func newSessionHeader(chain string) coretypes.SessionHeader {
	return coretypes.SessionHeader{
		Chain:     chain,
		IpAddress: idutil.NewGuid(),
	}
}

func newNode(chain string) nodetypes.Node {
	return nodetypes.Node{
		Id: idutil.NewGuid(),
		Chain: chain,
		IpAddress: "0.0.0.0",
		RewardAddress: "vite_",
	}
}