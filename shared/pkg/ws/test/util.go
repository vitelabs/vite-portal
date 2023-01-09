package test

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/testutil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

var defaultPid = -1

type TestWsRpc struct {
	Timeout time.Duration
	Url     string
	pid     int
}

func NewTestWsRpc(timeout time.Duration) *TestWsRpc {
	return &TestWsRpc{
		Timeout: timeout,
		pid: defaultPid,
	}
}

func (r *TestWsRpc) Start() {
	if r.pid != defaultPid {
		r.Stop()
	}

	p := testutil.BuildFullPath("shared", "pkg", "ws", "main", "main.go")
	cmd := exec.Command("go", "run", p)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	out := sharedtypes.NewBufferChannel()
	cmd.Stdout = out

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	r.pid = cmd.Process.Pid

	commonutil.WaitForChan(r.Timeout, out.Changed, func(p []byte) bool {
		return true
	})
	r.Url = strings.TrimSpace(out.String())
	logger.Logger().Info().Msg(fmt.Sprintf("url: %s", r.Url))

	if canConnect := ws.CanConnect(r.Url, r.Timeout); !canConnect {
		logger.Logger().Fatal().Msg(fmt.Sprintf("Failed to start TestWsRpc: %d (%s)", r.pid, r.Url))
	}
}

func (r *TestWsRpc) Stop() {
	if r.pid == defaultPid {
		return
	}

	if err := syscall.Kill(-r.pid, syscall.SIGKILL); err != nil {
		log.Fatal(err)
	}
	if canConnect := ws.CanConnect(r.Url, r.Timeout); canConnect {
		logger.Logger().Fatal().Msg(fmt.Sprintf("Failed to stop TestWsRpc: %d (%s)", r.pid, r.Url))
	}
	
	r.pid = defaultPid
}
