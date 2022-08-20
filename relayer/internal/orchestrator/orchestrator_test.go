package orchestrator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

var timeout = 1000 * time.Millisecond

func TestInit(t *testing.T) {
	mock := ws.StartMockWsRpc(timeout)
	o, err := InitOrchestrator(mock.Url, timeout)
	require.Nil(t, err)
	require.NotNil(t, o)
	require.Equal(t, ws.Unknown, o.GetStatus())
	time.Sleep(100 * time.Millisecond)
	require.Equal(t, ws.Connected, o.GetStatus())
}

func TestInitError(t *testing.T) {
	o, err := InitOrchestrator("http://localhost:1234", timeout)
	require.Nil(t, o)
	require.NotNil(t, err)
	require.Equal(t, "URL need to match WebSocket Protocol.", err.Error())
}