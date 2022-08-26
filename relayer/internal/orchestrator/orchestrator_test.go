package orchestrator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
	wstest "github.com/vitelabs/vite-portal/shared/pkg/ws/test"
)

var timeout = 1000 * time.Millisecond

func TestInit(t *testing.T) {
	r := wstest.NewTestWsRpc(timeout)
	r.Start()
	o := NewOrchestrator(r.Url, timeout)
	require.NotNil(t, o)
	require.Equal(t, ws.Unknown, o.GetStatus())
	o.Start(rpc.NewServer())
	require.Equal(t, ws.Connected, o.GetStatus())
	require.True(t, ws.CanConnect(r.Url, timeout))
	r.Stop()
	require.False(t, ws.CanConnect(r.Url, timeout))
	time.Sleep(timeout)
	require.Equal(t, ws.Disconnected, o.GetStatus())
}

func TestMockInit(t *testing.T) {
	mock := ws.NewMockWsRpc(0)
	require.NotNil(t, mock)
	go mock.Serve(timeout)
	o := NewOrchestrator(mock.Url, timeout)
	require.NotNil(t, o)
	require.Equal(t, ws.Unknown, o.GetStatus())
	o.Start(rpc.NewServer())
	require.Equal(t, ws.Connected, o.GetStatus())
	require.True(t, ws.CanConnect(mock.Url, timeout))
	mock.Close()
	require.False(t, ws.CanConnect(mock.Url, timeout))
	time.Sleep(timeout)
	// Connections are not closed -> use TestWsRpc
	require.Equal(t, ws.Connected, o.GetStatus())
}

func TestInitInvalidUrl(t *testing.T) {
	t.Parallel()
	o := NewOrchestrator("http://localhost:1234", timeout)
	require.NotNil(t, o)
}