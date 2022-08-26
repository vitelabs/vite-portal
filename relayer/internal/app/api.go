package app

import "github.com/vitelabs/vite-portal/shared/pkg/rpc"

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
			Namespace:     "nodes",
			Authenticated: true,
			Service:       &nodesAPI{a},
		},
		// {
		// 	Namespace: "debug",
		// 	Service:   debug.Handler,
		// },
		{
			Namespace: "public",
			Service:   &publicAPI{a},
		},
	}
}

// nodesAPI exposes API methods related to nodes
type nodesAPI struct {
	app *RelayerApp
}

func (a *nodesAPI) GetSecret() string {
	return "secret1234"
}

// publicAPI offers helper utils
type publicAPI struct {
	app *RelayerApp
}

func (a *publicAPI) Version() string {
	return a.app.config.Version
}
