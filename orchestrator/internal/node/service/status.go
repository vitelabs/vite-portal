package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

// UpdateStatus tries to update the local status of a subset of all nodes specified by the limit parameter.
// Once all nodes have been updated, it starts from the beginning.
func (s *Service) UpdateStatus(chain string, limit, batchSize int) {
	if chain == "" || limit <= 0 {
		return
	}
	store := s.context.GetNodeStore(chain)
	if store == nil {
		logger.Logger().Error().Msg(fmt.Sprintf("node store not found for chain '%s'", chain))
		return
	}
	e := store.GetEnumerator()
	batch := make([]nodetypes.Node, 0, batchSize)
	processed := *s.context.GetStatusStore(chain).ProcessedSet
	count := 0
	for e.MoveNext() {
		n, found := e.Current()
		if !found {
			continue
		}
		if processed.Contains(n.Id) {
			continue
		}
		processed.Add(n.Id)
		batch = append(batch, n)
		if len(batch) >= batchSize {
			s.updateStatus(batch)
			batch = batch[:0]
		}
		count++
		if count >= limit {
			break
		}
	}
	if len(batch) > 0 {
		s.updateStatus(batch)
	}
	if count == 0 {
		processed.Clear()
	}
}

func (s *Service) updateStatus(batch []nodetypes.Node) {
	var wg = sync.WaitGroup{}
	maxGoroutines := len(batch) // could be smaller than batch size if needed
	guard := make(chan struct{}, maxGoroutines)
	timeout := time.Duration(s.config.RpcTimeout) * time.Millisecond

	for _, v := range batch {
		guard <- struct{}{}
		wg.Add(1)
		go func(n nodetypes.Node) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			logger.Logger().Info().
				Str("id", n.Id).
				Str("name", n.Name).
				Str("ip", n.ClientIp).
				Str("chain", n.Chain).
				Str("rewardAddress", n.RewardAddress).
				Msg("calling 'dashboard_runtimeInfo'")
			var runtimeInfo sharedtypes.RpcViteRuntimeInfoResponse
			if err := n.RpcClient.CallContext(ctx, &runtimeInfo, "dashboard_runtimeInfo", "param1"); err != nil {
				// not successful
			} else {
				// successful
			}
			<-guard
			wg.Done()
		}(v)
	}

	wg.Wait()
}

// SendStatus sends the local status information about every node to Apache Kafka
func (s *Service) SendStatus(chain string) {
	// round := time.Now().UnixMilli() / 1000 / 60
}

func (s *Service) GetChainHeight(chain string) int64 {
	store := s.context.GetStatusStore(chain)
	if store.GlobalHeight.Height != 0 && store.GlobalHeight.LastUpdate != 0 {
		if time.Now().UnixMilli() - store.GlobalHeight.LastUpdate < 10 {
			return store.GlobalHeight.Height
		}
	}
	h, err := s.GetViteClient(chain).GetSnapshotChainHeight()
	if err != nil {
		return 0
	}
	store.GlobalHeight.Update(h)
	return h
}