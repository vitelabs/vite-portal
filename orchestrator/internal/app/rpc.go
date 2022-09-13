package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
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
			Prefix:    "/", // accept all URLs (e.g. /ws/gvite/...)
			JwtSecret: secret,
		}, a.BeforeConnect, a.OnConnect, a.OnDisconnect); err != nil {
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

func (a *OrchestratorApp) BeforeConnect(w http.ResponseWriter, r *http.Request) error {
	// TODO: check temporary blacklist
	if false {
		msg := "too many requests"
		http.Error(w, msg, http.StatusTooManyRequests)
		return errors.New(msg)
	}
	return nil
}

func (a *OrchestratorApp) OnConnect(c *rpc.Client, peerInfo rpc.PeerInfo) (sharedtypes.Connection, error) {
	defaultConn := sharedtypes.Connection{}
	timeout := time.Duration(a.config.RpcTimeout) * time.Millisecond
	if a.relayerService.IsRelayerConnection(peerInfo) {
		id, err := a.relayerService.HandleConnect(timeout, c, peerInfo)
		if err != nil {
			a.HandleOnConnectError(timeout, c.WriteConn, err)
			return defaultConn, err
		}
		return *sharedtypes.NewConnection(types.ConnectionTypes.Relayer, id), nil
	}
	// By default it is assumed the connection has been initiated by a node
	id, err := a.nodeService.HandleConnect(timeout, c, peerInfo)
	if err != nil {
		// TODO: add to temporary blacklist
		a.HandleOnConnectError(timeout, c.WriteConn, err)
		return defaultConn, err
	}
	return *sharedtypes.NewConnection(types.ConnectionTypes.Node, id), nil
}

func (a *OrchestratorApp) OnDisconnect(c sharedtypes.Connection) {
	switch c.Type {
	case types.ConnectionTypes.Relayer:
		a.relayerService.HandleDisconnect(c.Id)
	}
}

func (a *OrchestratorApp) HandleOnConnectError(timeout time.Duration, w rpc.JSONWriter, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	w.WriteJSON(ctx, rpc.NewJSONRPCErrorMessage(err))
}
