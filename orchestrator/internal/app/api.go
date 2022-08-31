package app

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/version"
)

// GetAPIs return two sets of APIs, both the ones that do not require
// authentication, and the complete set
func (a *OrchestratorApp) GetAPIs() (unauthenticated, all []rpc.API) {
	for _, api := range a.rpcAPIs {
		if !api.Authenticated {
			unauthenticated = append(unauthenticated, api)
		}
	}
	return unauthenticated, a.rpcAPIs
}

// apis returns the collection of built-in RPC APIs.
func (a *OrchestratorApp) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: RpcCoreModule,
			Service:   &coreAPI{a},
		},
		{
			Namespace:     RpcAdminModule,
			Authenticated: true,
			Service:       &adminAPI{a},
		},
		// {
		// 	Namespace: "debug",
		// 	Service:   debug.Handler,
		// },
	}
}

// coreAPI exposes API methods related to core
type coreAPI struct {
	app *OrchestratorApp
}

func (a *coreAPI) GetAppInfo() sharedtypes.RpcAppInfoResponse {
	return sharedtypes.RpcAppInfoResponse{
		Id:      a.app.id,
		Version: version.PROJECT_BUILD_VERSION,
		Name:    types.AppName,
	}
}

// adminAPI is the collection of administrative API methods exposed over
// both secure and unsecure RPC channels.
type adminAPI struct {
	app *OrchestratorApp
}

// GetSecrect returns the app secret
func (a *adminAPI) GetSecret() string {
	return "secret1234"
}
