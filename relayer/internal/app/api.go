package app

import (
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	"github.com/vitelabs/vite-portal/shared/pkg/version"
)

// GetAPIs return two sets of APIs, both the ones that do not require
// authentication, and the complete set
func (a *RelayerApp) GetAPIs() (unauthenticated, all []rpc.API) {
	for _, api := range a.rpcAPIs {
		if !api.Authenticated {
			unauthenticated = append(unauthenticated, api)
		}
	}
	return unauthenticated, a.rpcAPIs
}

// apis returns the collection of built-in RPC APIs.
func (a *RelayerApp) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: RpcCoreModule,
			Service:   &coreAPI{a},
		},
		{
			Namespace:     RpcNodesModule,
			Authenticated: true,
			Service:       &nodesAPI{a},
		},
		// {
		// 	Namespace: "debug",
		// 	Service:   debug.Handler,
		// },
	}
}

// coreAPI exposes API methods related to core
type coreAPI struct {
	app *RelayerApp
}

func (a *coreAPI) GetAppInfo() types.AppInfo {
	return types.AppInfo{
		Id:      a.app.id,
		Version: version.PROJECT_BUILD_VERSION,
		Name:    types.AppName,
	}
}

// nodesAPI exposes API methods related to nodes
type nodesAPI struct {
	app *RelayerApp
}

func (a *nodesAPI) GetSecret() string {
	return "secret1234"
}
