package service

import (
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

func (s *Service) HandleOrchestratorStatusChange(status ws.ConnectionStatus) {
	logger.Logger().Info().Int64("status", int64(status)).Msg("orchestrator status changed")
	start := time.Now()
	if status != ws.Connected {
		return
	}
	// get all nodes for all supported chains
	chains := s.nodeService.GetChains()
	for _, chain := range chains {
		s.getNodesRecursive(chain, 0, 0)
	}
	elapsed := time.Since(start)
	logger.Logger().Info().Int64("elapsed", elapsed.Milliseconds()).Int64("status", int64(status)).Msg("orchestrator status change handled")
}

func (s *Service) getNodesRecursive(chain string, offset, limit int) {
	start := time.Now()
	res, err := s.nodeService.GetNodes(chain, offset, limit)
	if err != nil {
		logger.Logger().Error().Err(err).Str("chain", chain).Int("offset", offset).Int("limit", limit).Msg("get nodes failed")
		return
	}

	for _, n := range res.Entries {
		s.nodeService.PutNode(n)
	}

	elapsed := time.Since(start)
	logger.Logger().Info().Int64("elapsed", elapsed.Milliseconds()).Str("chain", chain).
		Int64("offset", int64(res.Offset)).Int64("limit", int64(res.Limit)).Int64("total", int64(res.Total)).Msg("get nodes")

	newOffset := res.Offset + len(res.Entries)
	if res.Total > newOffset {
		s.getNodesRecursive(chain, newOffset, res.Limit)
	}
}
