package app

import (
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

func (a *OrchestratorApp) InitScheduler() {
	job1, err := a.scheduler.Every(1).Minute().Do(func() {
		a.HandleNodeStatusUpdate()
	})
	if err != nil {
		panic(err)
	}
	job1.SingletonMode()

	job2, err := a.scheduler.Every(5).Seconds().Do(func() {
		a.HandleNodeOnlineStatusUpdate()
	})
	if err != nil {
		panic(err)
	}
	job2.SingletonMode()
}

// HandleNodeStatusUpdate tries to update the local status of 1/10 of the nodes every minute.
// This means it takes max. 10 minutes or 600 seconds to update all nodes.
// Once all nodes have been updated, it starts from the beginning.
func (a *OrchestratorApp) HandleNodeStatusUpdate() {
	for _, c := range a.config.SupportedChains {
		start := time.Now()
		store, err := a.context.GetNodeStore(c.Name)
		if err != nil {
			logger.Logger().Error().Msg(err.Error())
			continue
		}
		n := store.Count() / 10
		if n < 50 {
			n = 50
		}
		a.nodeService.UpdateStatus(c.Name, n, 20)
		elapsed := time.Since(start)
		logger.Logger().Info().
			Str("chain", c.Name).
			Int("n", n).
			Int("count", store.Count()).
			Int64("elapsed", elapsed.Milliseconds()).
			Msg("node status updated")
	}
}

func (a *OrchestratorApp) HandleNodeOnlineStatusUpdate() {
	for _, c := range a.config.SupportedChains {
		a.nodeService.UpdateOnlineStatus(c.Name)
	}
}
