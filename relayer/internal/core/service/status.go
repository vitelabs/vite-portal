package service

import (
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/ws"
)

func (s *Service) HandleOrchestratorStatusChange(status ws.ConnectionStatus) {
	logger.Logger().Info().Int64("status", int64(status)).Msg("orchestrator status changed")
	now := time.Now()
	if status != ws.Connected {
		return
	}
	// time.Sleep(5 * time.Second)
	// TODO: get all nodes for all supported chains
	elapsed := time.Since(now)
	logger.Logger().Info().Int64("elapsed", elapsed.Milliseconds()).Int64("status", int64(status)).Msg("orchestrator status change handled")
}
