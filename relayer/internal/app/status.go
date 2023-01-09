package app

import (
	"fmt"
	"time"

	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/commonutil"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

func (a *RelayerApp) initOrchestrator() {
	a.handleOrchestratorStatusChange(a.orchestrator.GetStatus())
	c := a.orchestrator.SubscribeStatusChange()
	go func() {
		for {
			select {
			case status := <-c:
				_ = status
				a.handleOrchestratorStatusChange(a.orchestrator.GetStatus())
			}
		}
	}()
}

func (a *RelayerApp) handleOrchestratorStatusChange(status ws.ConnectionStatus) {
	logger.Logger().Info().Int64("status", int64(status)).Msg("orchestrator status changed")
	start := time.Now()
	if status != ws.Connected {
		return
	}
	// get all nodes for all supported chains
	chains, err := a.orchestrator.GetChains()
	if err != nil {
		return
	}
	for _, chain := range chains {
		a.getNodesRecursive(chain.Name, 0, 0)
	}
	elapsed := time.Since(start)
	logger.Logger().Info().Int64("elapsed", elapsed.Milliseconds()).Int64("status", int64(status)).Msg("orchestrator status change handled")
}

func (a *RelayerApp) getNodesRecursive(chain string, offset, limit int) {
	o, l := commonutil.CheckPagination(offset, limit)
	start := time.Now()
	res, err := a.orchestrator.GetNodes(chain, o, l)
	if err != nil {
		logger.Logger().Error().Err(err).Str("chain", chain).Int("offset", o).Int("limit", l).Msg("get nodes failed")
		return
	}

	for _, e := range res.Entries {
		n := nodetypes.Node{
			Id:         e.Id,
			Chain:      e.Chain,
			RpcHttpUrl: fmt.Sprintf("http://%s:%d", e.ClientIp, e.HTTPort),
			RpcWsUrl:   fmt.Sprintf("ws://%s:%d", e.ClientIp, e.WSPort),
		}
		a.nodeService.PutNode(n)
	}

	elapsed := time.Since(start)
	logger.Logger().Info().Int64("elapsed", elapsed.Milliseconds()).Str("chain", chain).
		Int64("offset", int64(res.Offset)).Int64("limit", int64(res.Limit)).Int64("total", int64(res.Total)).Msg("get nodes")

	newOffset := res.Offset + len(res.Entries)
	if res.Total > newOffset {
		a.getNodesRecursive(chain, newOffset, res.Limit)
	}
}
