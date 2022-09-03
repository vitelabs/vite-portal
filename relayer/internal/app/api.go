package app

import (
	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
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
			Namespace:     RpcDbModule,
			Authenticated: true,
			Service:       &dbAPI{a},
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

func (a *coreAPI) GetAppInfo() sharedtypes.RpcAppInfoResponse {
	return sharedtypes.RpcAppInfoResponse{
		Id:      a.app.id,
		Version: version.PROJECT_BUILD_VERSION,
		Name:    types.AppName,
	}
}

// dbAPI exposes API methods related to database
type dbAPI struct {
	app *RelayerApp
}

func (a *dbAPI) GetChains() []string {
	return a.app.nodeService.GetChains()
}

func (a *dbAPI) GetNodes(chain string, offset, limit int) (generics.GenericPage[nodetypes.Node], error) {
	o, l := commonutil.CheckPagination(offset, limit)
	return a.app.nodeService.GetNodes(chain, o, l)
}