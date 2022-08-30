package app

import "github.com/vitelabs/vite-portal/shared/pkg/rpc"

func (a *RelayerApp) startRPC(profile bool) error {
	a.StartHttpRpc(profile)
	a.StartWsRpc()

	if err := a.startInProc(); err != nil {
		return err
	}

	open, all := a.GetAPIs()

	init := func(server *rpc.HTTPServer, apis []rpc.API, port int, secret []byte) error {
		if err := server.SetListenAddr("", port); err != nil {
			return err
		}

		// Enable HTTP
		if err := server.EnableRPC(apis, rpc.HTTPConfig{
			CorsAllowedOrigins: DefaultAllowedOrigins,
			Vhosts:             DefaultVhosts,
			Modules:            DefaultModules,
			Prefix:             "",
			JwtSecret:          secret,
		}); err != nil {
			return err
		}

		return nil
	}

	// Set up unauthenticated RPC.
	if err := init(a.rpc, open, int(a.config.RpcPort), nil); err != nil {
		return err
	}
	// Set up authenticated RPC.
	if err := init(a.rpcAuth, all, int(a.config.RpcAuthPort), nil); err != nil {
		return err
	}

	// Start the servers
	a.rpc.Start()
	a.rpcAuth.Start()

	return nil
}

func (a *RelayerApp) stopRPC() {
	a.rpc.Stop()
	a.rpcAuth.Stop()
	a.stopInProc()
}

// startInProc registers all RPC APIs on the inproc server.
func (a *RelayerApp) startInProc() error {
	for _, api := range a.rpcAPIs {
		if err := a.inprocHandler.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
	}
	return nil
}

// stopInProc terminates the in-process RPC endpoint.
func (a *RelayerApp) stopInProc() {
	a.inprocHandler.Stop()
}