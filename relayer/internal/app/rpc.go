package app

import (
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/rpc"
)

func (a *RelayerApp) startRPC(profile bool) {
	routes := []rpc.Route{
		{Name: "Relay", Method: "POST", Path: "/api/v1/client/relay", HandlerFunc: a.Relay},
		{Name: "GetChains", Method: "GET", Path: "/api/v1/db/chains", HandlerFunc: a.GetChains},
		{Name: "GetNodes", Method: "GET", Path: "/api/v1/db/nodes", HandlerFunc: a.GetNodes},
		{Name: "GetNode", Method: "GET", Path: "/api/v1/db/nodes/:id", HandlerFunc: a.GetNode},
		{Name: "PutNode", Method: "PUT", Path: "/api/v1/db/nodes", HandlerFunc: a.PutNode},
		{Name: "DeleteNode", Method: "DELETE", Path: "/api/v1/db/nodes/:id", HandlerFunc: a.DeleteNode},
	}

	timeout := time.Duration(a.Config.RpcTimeout) * time.Millisecond
	rpc.StartHttpRpc(routes, a.Config.RpcHttpPort, timeout, a.Config.Debug, profile)
	rpc.StartWsRpc(a.Config.RpcWsPort, timeout)
}
