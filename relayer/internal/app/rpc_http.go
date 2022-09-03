package app

import (
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/rpc"
)

func (a *RelayerApp) StartHttpRpc(profile bool) {
	routes := []rpc.Route{
		{Name: "Default", Method: "GET", Path: "/", HandlerFunc: a.AppInfo},
		{Name: "Relay", Method: "POST", Path: "/api/v1/client/relay", HandlerFunc: a.Relay},
	}

	timeout := time.Duration(a.config.RpcTimeout) * time.Millisecond
	rpc.StartHttpRpc(routes, a.config.RpcRelayHttpPort, timeout, a.config.Debug, profile)
}

func (a *RelayerApp) StartWsRpc() {
	timeout := time.Duration(a.config.RpcTimeout) * time.Millisecond
	rpc.StartWsRpc(a.config.RpcRelayWsPort, timeout)
}
