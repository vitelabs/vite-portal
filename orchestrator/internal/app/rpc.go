package app

import (
	"github.com/gorilla/websocket"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

// Attach creates an RPC client attached to an in-process API handler.
func (a *OrchestratorApp) Attach() (*rpc.Client, error) {
	return rpc.DialInProc(a.inprocHandler), nil
}

// RPCHandler returns the in-process RPC request handler.
func (a *OrchestratorApp) RPCHandler() (*rpc.Server, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.state == closedState {
		return nil, ErrOrchestratorStopped
	}
	return a.inprocHandler, nil
}

// startRPC is a helper method to configure all the various RPC endpoints during app
// startup. It's not meant to be called at any time afterwards as it makes certain
// assumptions about the state of the app.
func (a *OrchestratorApp) startRPC() error {
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
		// Enable WS
		if err := server.EnableWS(apis, rpc.WSConfig{
			Modules:   DefaultModules,
			Origins:   DefaultAllowedOrigins,
			Prefix:    "",
			JwtSecret: secret,
		}, a.OnConnect); err != nil {
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

func (a *OrchestratorApp) stopRPC() {
	a.rpc.Stop()
	a.rpcAuth.Stop()
	a.stopInProc()
}

// startInProc registers all RPC APIs on the inproc server.
func (a *OrchestratorApp) startInProc() error {
	for _, api := range a.rpcAPIs {
		if err := a.inprocHandler.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
	}
	return nil
}

// stopInProc terminates the in-process RPC endpoint.
func (a *OrchestratorApp) stopInProc() {
	a.inprocHandler.Stop()
}

func (a *OrchestratorApp) OnConnect(c *websocket.Conn) error {
	logger.Logger().Debug().Interface("conn", c).Msg(a.config.Version)
	return nil
}
