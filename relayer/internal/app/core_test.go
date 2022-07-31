package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/internal/core/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

func newRelayerCoreApp() *RelayerCoreApp {
	config := roottypes.NewDefaultConfig()
	o, _ := NewOrchestrator()
	c := NewContext(config)
	return NewRelayerCoreApp(config, o, c)
}

func TestSetClientIp(t *testing.T) {
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
					Headers: map[string][]string{roottypes.DefaultHeaderTrueClientIp: {"1.2.3.4"}},
				},
			},
			expected: "1.2.3.4",
		},
		{
			name: "Test default multiple",
			relay: types.Relay{
				Payload: types.Payload{
					Headers: map[string][]string{roottypes.DefaultHeaderTrueClientIp: {"4.3.2.1", "1.2.3.4"}},
				},
			},
			expected: "4.3.2.1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			coreApp := newRelayerCoreApp()
			coreApp.setClientIp(&tc.relay)
			require.NotEmpty(t, tc.relay.ClientIp)
			require.Equal(t, tc.expected, tc.relay.ClientIp)
		})
	}
}
