package app

import "github.com/vitelabs/vite-portal/shared/pkg/rpc"

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
			Namespace: "admin",
			Authenticated: true,
			Service:   &adminAPI{a},
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

// adminAPI is the collection of administrative API methods exposed over
// both secure and unsecure RPC channels.
type adminAPI struct {
	app *OrchestratorApp
}

// publicAPI offers helper utils
type publicAPI struct {
	app *OrchestratorApp
}

// Version returns the app version
func (a *publicAPI) Version() string {
	return a.app.config.Version
}
