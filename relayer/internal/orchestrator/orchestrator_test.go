package orchestrator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/util/idutil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
	wstest "github.com/vitelabs/vite-portal/shared/pkg/ws/test"
)

var timeout = 1000 * time.Millisecond

func TestInit(t *testing.T) {
	r := wstest.NewTestWsRpc(timeout)
	r.Start()
	o := NewOrchestrator(idutil.NewGuid(), r.Url, types.DefaultJwtSecret, timeout)
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
	o := NewOrchestrator(idutil.NewGuid(), mock.Url, types.DefaultJwtSecret, timeout)
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
	o := NewOrchestrator(idutil.NewGuid(), "http://localhost:1234", types.DefaultJwtSecret, timeout)
	require.NotNil(t, o)
}

func TestStop(t *testing.T) {
	r := wstest.NewTestWsRpc(timeout)
	r.Start()
	o := NewOrchestrator(idutil.NewGuid(), r.Url, types.DefaultJwtSecret, timeout)
	require.NotNil(t, o)
	require.Equal(t, ws.Unknown, o.GetStatus())
	o.Start(rpc.NewServer())
	require.Equal(t, ws.Connected, o.GetStatus())
	require.True(t, ws.CanConnect(r.Url, timeout))
	o.Stop()
	require.True(t, ws.CanConnect(r.Url, timeout))
	require.Equal(t, ws.Disconnected, o.GetStatus())
}
