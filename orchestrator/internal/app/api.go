package app

import (
	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	relayertypes "github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
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

// adminAPI exposes API methods for authenticated requests only
type adminAPI struct {
	app *OrchestratorApp
}

func (a *adminAPI) GetChains() []sharedtypes.ChainConfig {
	return a.app.config.SupportedChains
}


func (a *adminAPI) GetNodes(chain string, offset int, limit int) (generics.GenericPage[nodetypes.Node], error) {
	o, l := commonutil.CheckPagination(offset, limit)
	return a.app.nodeService.Get(chain, o, l)
}

func (a *adminAPI) GetRelayers(offset int, limit int) (generics.GenericPage[relayertypes.Relayer], error) {
	return a.app.relayerService.Get(commonutil.CheckPagination(offset, limit))
}

func (a *adminAPI) GetKafkaDefaultMessages(offset, limit, timeout int) ([]string, error) {
	o, l := commonutil.CheckPagination(offset, limit)
	return a.app.kafka.ReadDefault(o, l, timeout)
}

func (a *adminAPI) GetKafkaRpcMessages(offset, limit, timeout int) ([]string, error) {
	o, l := commonutil.CheckPagination(offset, limit)
	return a.app.kafka.ReadRpc(o, l, timeout)
}

func (a *adminAPI) UpdateNodeStatus() {
	a.app.HandleNodeStatusUpdate()
}

func (a *adminAPI) UpdateNodeOnlineStatus() {
	a.app.HandleNodeOnlineStatusUpdate()
}

func (a *adminAPI) DispatchNodeStatus() {
	a.app.HandleNodeStatusDispatch()
}
