package interfaces

import (
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

type ClientI interface {
	Connect() error
	Subscribe(c chan string)
}

type OrchestratorI interface {
	GetStatus() ws.ConnectionStatus
	Start(apis []rpc.API)
	Stop()
	SubscribeStatusChange() <-chan string
	GetChains() ([]sharedtypes.ChainConfig, error)
	GetNodes(chain string, offset int, limit int) (generics.GenericPage[types.Node], error)
}