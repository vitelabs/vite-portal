package app

import (
	"errors"

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
	app *RelayerApp
}

func (a *coreAPI) GetAppInfo() sharedtypes.RpcAppInfoResponse {
	return sharedtypes.RpcAppInfoResponse{
		Id:      a.app.id,
		Version: version.PROJECT_BUILD_VERSION,
		Name:    types.AppName,
	}
}

// adminAPI exposes API methods related to administrative tasks
type adminAPI struct {
	app *RelayerApp
}

func (a *adminAPI) GetChains() []string {
	return a.app.nodeService.GetChains()
}

func (a *adminAPI) GetNodes(chain string, offset, limit int) (generics.GenericPage[nodetypes.Node], error) {
	o, l := commonutil.CheckPagination(offset, limit)
	return a.app.nodeService.GetNodes(chain, o, l)
}

func (a *adminAPI) GetNode(id string) (nodetypes.Node, error) {
	res, found := a.app.nodeService.GetNode(id)
	if !found {
		return nodetypes.Node{}, errors.New("node does not exist")
	}
	return res, nil
}

// PutNode enables the orchestrator to add or update a node
func (a *adminAPI) PutNode(node nodetypes.Node) error {
	return a.app.nodeService.PutNode(node)
}

// DeleteNode enables the orchestrator to delete a node
func (a *adminAPI) DeleteNode(id string) error {
	return a.app.nodeService.DeleteNode(id)
}