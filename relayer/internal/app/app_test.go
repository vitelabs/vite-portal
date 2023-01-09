package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/core/types"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/relayer/internal/util/testutil"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func TestSetChain(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		relay         types.Relay
		chains        []string
		expected      string
		expectedError error
	}{
		{
			name:          "Test emtpy chains",
			relay:         types.Relay{},
			chains:        []string{},
			expected:      "",
			expectedError: errors.New("chains are empty"),
		},
		{
			name: "Test unsupported",
			relay: types.Relay{
				Chain: "chain2",
			},
			chains:        []string{"chain1"},
			expected:      "",
			expectedError: errors.New("the chain 'chain2' is not supported"),
		},
		{
			name:     "Test default",
			relay:    types.Relay{},
			chains:   []string{"chain1", "chain2"},
			expected: "chain1",
		},
		{
			name: "Test chain2",
			relay: types.Relay{
				Chain: "chain2",
			},
			chains:   []string{"chain1", "chain2"},
			expected: "chain2",
		},
		{
			name: "Test host",
			relay: types.Relay{
				Host: "test.localhost",
			},
			chains:   []string{"chain1", "chain2"},
			expected: "chain2",
		},
		{
			name: "Test invalid host",
			relay: types.Relay{
				Host: "test1234",
			},
			chains:   []string{"chain1", "chain2"},
			expected: "chain1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := newRelayerApp()
			insertChains(a, tc.chains)
			err := a.setChain(&tc.relay)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
				require.Empty(t, tc.expected)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, tc.relay.Chain)
			}
		})
	}
}

func TestSetClientIp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		relay    types.Relay
		expected string
	}{
		{
			name:     "Test emtpy relay",
			relay:    types.Relay{},
			expected: roottypes.DefaultIpAddress,
		},
		{
			name: "Test already set",
			relay: types.Relay{
				ClientIp: "1.2.3.4",
			},
			expected: "1.2.3.4",
		},
		{
			name: "Test header",
			relay: types.Relay{
				Payload: types.Payload{
					Headers: map[string][]string{"test": {"1.2.3.4"}},
				},
			},
			expected: roottypes.DefaultIpAddress,
		},
		{
			name: "Test default",
			relay: types.Relay{
				Payload: types.Payload{
					Headers: map[string][]string{sharedtypes.DefaultHeaderTrueClientIp: {"1.2.3.4"}},
				},
			},
			expected: "1.2.3.4",
		},
		{
			name: "Test default multiple",
			relay: types.Relay{
				Payload: types.Payload{
					Headers: map[string][]string{sharedtypes.DefaultHeaderTrueClientIp: {"4.3.2.1", "1.2.3.4"}},
				},
			},
			expected: "4.3.2.1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := newRelayerApp()
			a.setClientIp(&tc.relay)
			require.NotEmpty(t, tc.relay.ClientIp)
			require.Equal(t, tc.expected, tc.relay.ClientIp)
		})
	}
}

func newRelayerApp() *RelayerApp {
	config := roottypes.NewDefaultConfig()
	config.RpcTimeout = 0
	config.HostToChainMap = map[string]string{
		"test.localhost": "chain2",
	}
	return NewRelayerApp(config)
}

func insertChains(a *RelayerApp, chains []string) {
	for _, c := range chains {
		a.nodeService.PutNode(testutil.NewNode(c))
	}
}
