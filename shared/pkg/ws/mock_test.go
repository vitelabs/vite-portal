package ws

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStartMockWsRpc(t *testing.T) {
	t.Parallel()
	timeout := 100 * time.Millisecond
	mock := NewMockWsRpc(0)
	require.NotNil(t, mock)
	require.Greater(t, mock.Port, 0)
	require.Equal(t, "/ws/mock", mock.Pattern)
	require.Equal(t, fmt.Sprintf("ws://localhost:%d/ws/mock", mock.Port), mock.Url)
	require.False(t, CanConnect(mock.Url, timeout))
	go mock.Serve(timeout)
	require.True(t, CanConnect(mock.Url, timeout))
}