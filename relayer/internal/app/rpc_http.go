package app

import (
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/rpc"
)

func (a *RelayerApp) StartHttpRpc(profile bool) {
	routes := []rpc.Route{
		{Name: "Default", Method: "GET", Path: "/", HandlerFunc: a.AppInfo},
		{Name: "Relay", Method: "POST", Path: "/api/v1/client/relay", HandlerFunc: a.Relay},
		{Name: "GetNode", Method: "GET", Path: "/api/v1/db/nodes/:id", HandlerFunc: a.GetNode},
		{Name: "PutNode", Method: "PUT", Path: "/api/v1/db/nodes", HandlerFunc: a.PutNode},
		{Name: "DeleteNode", Method: "DELETE", Path: "/api/v1/db/nodes/:id", HandlerFunc: a.DeleteNode},
	}

	timeout := time.Duration(a.config.RpcTimeout) * time.Millisecond
	rpc.StartHttpRpc(routes, a.config.RpcRelayHttpPort, timeout, a.config.Debug, profile)
}

func (a *RelayerApp) StartWsRpc() {
	timeout := time.Duration(a.config.RpcTimeout) * time.Millisecond
	rpc.StartWsRpc(a.config.RpcRelayWsPort, timeout)
}
