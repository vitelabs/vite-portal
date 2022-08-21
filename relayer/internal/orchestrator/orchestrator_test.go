package orchestrator

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

var timeout = 1000 * time.Millisecond

func TestInit(t *testing.T) {
	mock := ws.NewMockWsRpc(0)
	require.NotNil(t, mock)
	go mock.Serve(timeout)
	o, err := InitOrchestrator(mock.Url, timeout)
	require.Nil(t, err)
	require.NotNil(t, o)
	require.Equal(t, ws.Unknown, o.GetStatus())
	commonutil.WaitFor(timeout, o.StatusChanged, func(status ws.ConnectionStatus) bool {
		return status == ws.Connected
	})
	require.Equal(t, ws.Connected, o.GetStatus())
	require.True(t, ws.CanConnect(mock.Url, timeout))
	mock.Close()
	require.False(t, ws.CanConnect(mock.Url, timeout))
	commonutil.WaitFor(timeout, o.StatusChanged, func(status ws.ConnectionStatus) bool {
		return status == ws.Disconnected
	})
	require.Equal(t, ws.Disconnected, o.GetStatus())
}

func TestInitError(t *testing.T) {
	o, err := InitOrchestrator("http://localhost:1234", timeout)
	require.Nil(t, o)
	require.NotNil(t, err)
	require.Equal(t, "URL need to match WebSocket Protocol.", err.Error())
}

func TestKill(t *testing.T) {
	p := testutil.BuildFullPath("shared", "pkg", "ws", "main", "main.go")
	cmd := exec.Command("go", "run", p)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
	}

	out := sharedtypes.BufferChannel{}
	cmd.Stdout = &out

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	commonutil.WaitFor(timeout, out.Changed, func(p []byte) bool {
		return true
	})
	url := strings.TrimSpace(out.String())
	logger.Logger().Info().Msg(fmt.Sprintf("url: %s", url))
	require.True(t, ws.CanConnect(url, timeout))
	err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	require.NoError(t, err)
	require.False(t, ws.CanConnect(url, timeout))
}
