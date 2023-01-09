package orchestrator

import (
	"strconv"

	nodeinterfaces "github.com/vitelabs/vite-portal/relayer/internal/node/interfaces"
	"github.com/vitelabs/vite-portal/relayer/internal/orchestrator/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/handler"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

type OrchestratorMock struct {
	stopped bool
	status  ws.ConnectionStatus
	ps      *handler.Pubsub
	store   nodeinterfaces.StoreI
}

func NewOrchestratorMock(store nodeinterfaces.StoreI) *OrchestratorMock {
	return &OrchestratorMock{
		stopped: false,
		status:  ws.Unknown,
		ps:      handler.NewPubsub(),
		store:   store,
	}
}

func (o *OrchestratorMock) GetStatus() ws.ConnectionStatus {
	return o.status
}

func (o *OrchestratorMock) Start(apis []rpc.API) {
	o.stopped = false
}

func (o *OrchestratorMock) Stop() {
	o.stopped = true
}

func (o *OrchestratorMock) SubscribeStatusChange() <-chan string {
	return o.ps.Subscribe()
}

func (o *OrchestratorMock) connect() {
	if o.stopped {
		return
	}
	o.setStatus(ws.Connecting)
	o.setStatus(ws.Connected)
}

func (o *OrchestratorMock) setStatus(newStatus ws.ConnectionStatus) {
	if o.status != newStatus {
		logger.Logger().Info().
			Int64("before", int64(o.status)).
			Int64("after", int64(newStatus)).
			Msg("connection status changed")
		o.status = newStatus
		o.ps.Publish(strconv.FormatInt(int64(newStatus), 10))
	}
}

func (o *OrchestratorMock) GetChains() ([]sharedtypes.ChainConfig, error) {
	var resp []sharedtypes.ChainConfig
	return resp, nil
}

func (o *OrchestratorMock) GetNodes(chain string, offset, limit int) (generics.GenericPage[types.Node], error) {
	result := *generics.NewGenericPage[types.Node]()
	r, err := o.store.GetPaginated(chain, offset, limit)
	if err != nil {
		return result, err
	}
	result.Limit = r.Limit
	result.Offset = r.Offset
	result.Total = r.Total
	result.Entries = make([]types.Node, len(r.Entries))
	for i, n := range r.Entries {
		result.Entries[i] = types.Node{
			Id: n.Id,
			Chain: n.Chain,			
		}
	}
	return result, nil
}
