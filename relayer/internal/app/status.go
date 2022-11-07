package app

import (
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
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
	start := time.Now()
	res, err := a.orchestrator.GetNodes(chain, offset, limit)
	if err != nil {
		logger.Logger().Error().Err(err).Str("chain", chain).Int("offset", offset).Int("limit", limit).Msg("get nodes failed")
		return
	}

	// for _, n := range res.Entries {
	// 	a.nodeService.PutNode(n)
	// }

	elapsed := time.Since(start)
	logger.Logger().Info().Int64("elapsed", elapsed.Milliseconds()).Str("chain", chain).
		Int64("offset", int64(res.Offset)).Int64("limit", int64(res.Limit)).Int64("total", int64(res.Total)).Msg("get nodes")

	newOffset := res.Offset + len(res.Entries)
	if res.Total > newOffset {
		a.getNodesRecursive(chain, newOffset, res.Limit)
	}
}