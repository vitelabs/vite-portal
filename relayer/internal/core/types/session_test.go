package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nodeservice "github.com/vitelabs/vite-portal/internal/node/service"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
	"github.com/vitelabs/vite-portal/internal/util/testutil"
)

func TestNewSessionKey(t *testing.T) {
	tests := []struct {
		name      string
		chain     string
		ipAddress string
		expected  string
	}{
		{
			name:      "Test 1",
			chain:     "",
			ipAddress: "",
			expected:  "3545e7f0086f5fa7183af60c8e54778f",
		},
		{
			name:      "Test 2",
			chain:     "chain1",
			ipAddress: "0.0.0.0",
			expected:  "cb8b30cecd1857c59530f8bda15fab91",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			header := NewSessionHeader(tc.ipAddress, tc.chain)
			assert.Equal(t, tc.expected, header.HashString())
		})
	}
}

func TestNewSessionNodes(t *testing.T) {
	tests := []struct {
		name             string
		sessionNodeCount int
		nodeCount        int
		expected         int
		expectedError    error
	}{
		{
			name:             "Test sessionNodeCount=0 and nodeCount=0",
			sessionNodeCount: 0,
			nodeCount:        0,
			expected:         0,
			expectedError:    errors.New("insufficient nodes available to create a session"),
		},
		{
			name:             "Test sessionNodeCount=0 and nodeCount=1",
			sessionNodeCount: 0,
			nodeCount:        1,
			expected:         0,
			expectedError:    nil,
		},
		{
			name:             "Test sessionNodeCount=1 and nodeCount=1",
			sessionNodeCount: 1,
			nodeCount:        1,
			expected:         1,
			expectedError:    nil,
		},
		{
			name:             "Test sessionNodeCount=2 and nodeCount=1",
			sessionNodeCount: 2,
			nodeCount:        1,
			expected:         1,
			expectedError:    nil,
		},
		{
			name:             "Test sessionNodeCount=1 and nodeCount=2",
			sessionNodeCount: 1,
			nodeCount:        2,
			expected:         1,
			expectedError:    nil,
		},
		{
			name:             "Test sessionNodeCount=2 and nodeCount=2",
			sessionNodeCount: 2,
			nodeCount:        2,
			expected:         2,
			expectedError:    nil,
		},
		{
			name:             "Test sessionNodeCount=24 and nodeCount=100",
			sessionNodeCount: 24,
			nodeCount:        100,
			expected:         24,
			expectedError:    nil,
		},
		{
			name:             "Test sessionNodeCount=100 and nodeCount=24",
			sessionNodeCount: 100,
			nodeCount:        24,
			expected:         24,
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store := nodestore.NewMemoryStore()
			service := nodeservice.NewService(store)
			chain := "chain1"
			testutil.PutNodes(t, service, chain, tc.nodeCount)
			testutil.PutNodes(t, service, "chain2", tc.nodeCount)
			require.Equal(t, tc.nodeCount, service.GetNodeCount(chain))
			nodes, err := NewSessionNodes(service, chain, tc.sessionNodeCount)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.InnerError())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expected, len(nodes))
		})
	}
}
