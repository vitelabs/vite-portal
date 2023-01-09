package handler

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	nodetypes "github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

// UpdateStatus tries to update the local status of a subset of all nodes specified by the limit parameter.
// Once all nodes have been updated, it starts from the beginning.
func (h *Handler) UpdateStatus(limit, batchSize int) {
	if limit <= 0 || batchSize <= 0 {
		return
	}
	e := h.nodeStore.GetEnumerator()
	batch := make([]nodetypes.Node, 0, batchSize)
	processed := *h.statusStore.ProcessedSet
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
			h.updateStatus(batch)
			batch = batch[:0]
		}
		count++
		if count >= limit {
			break
		}
	}
	if len(batch) > 0 {
		h.updateStatus(batch)
	}
	if count == 0 {
		processed.Clear()
	}
}

func (h *Handler) updateStatus(batch []nodetypes.Node) {
	h.updateGlobalHeight()
	var wg = sync.WaitGroup{}
	maxGoroutines := len(batch) // could be smaller than batch size if needed
	guard := make(chan struct{}, maxGoroutines)
	for _, v := range batch {
		guard <- struct{}{}
		wg.Add(1)
		go func(n nodetypes.Node) {
			start := time.Now()
			logEvent := logger.Logger().Info().Str("id", n.Id).Str("name", n.Name).Str("ip", n.ClientIp).Str("chain", n.Chain)
			runtimeInfo, err := h.getRuntimeInfo(n)
			if err != nil {
				elapsed := time.Since(start)
				logEvent.Err(err).Int64("elapsed", elapsed.Milliseconds()).Msg("update status failed")
				return
			}
			h.updateNodeStatus(n, runtimeInfo, start)
			elapsed := time.Since(start)
			logEvent.Str("height", "0").Int64("elapsed", elapsed.Milliseconds()).Msg("status updated")
			<-guard
			wg.Done()
		}(v)
	}

	wg.Wait()
}

func (h *Handler) getRuntimeInfo(node nodetypes.Node) (sharedtypes.RpcViteRuntimeInfoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()
	var runtimeInfo sharedtypes.RpcViteRuntimeInfoResponse
	err := node.RpcClient.CallContext(ctx, &runtimeInfo, "dashboard_runtimeInfo", "param1")
	return runtimeInfo, err
}

func (h *Handler) updateNodeStatus(node nodetypes.Node, runtimeInfo sharedtypes.RpcViteRuntimeInfoResponse, start time.Time) error {
	lastUpdate := node.LastUpdate
	block := runtimeInfo.LatestSnapshot
	node.LastUpdate = sharedtypes.Int64(start.UnixMilli())
	node.DelayTime = sharedtypes.Int64(time.Since(start).Milliseconds())
	node.LastBlock.Hash = block.Hash
	node.LastBlock.Height = block.Height
	node.LastBlock.Time = block.Time
	node.Status = h.getOnlineStatus(block.Height)
	err := h.nodeStore.Update(int64(lastUpdate), node)
	if err != nil {
		logger.Logger().Info().Err(err).Str("id", node.Id).Msg("update node status failed")
		return err
	}
	return nil
}

func (h *Handler) updateGlobalHeight() int {
	h.heightLock.Lock()
	defer h.heightLock.Unlock()

	current := h.statusStore.GetGlobalHeight()
	lastUpdate := h.statusStore.GetLastUpdate()
	if current != 0 && lastUpdate != 0 {
		if time.Now().UnixMilli()-lastUpdate < 500 {
			return current
		}
	}
	height, err := h.client.GetSnapshotChainHeight()
	if err != nil {
		return current
	}
	h.statusStore.SetGlobalHeight(current, height)
	return h.statusStore.GetGlobalHeight()
}

func (h *Handler) UpdateOnlineStatus() {
	e := h.nodeStore.GetEnumerator()
	for e.MoveNext() {
		n, found := e.Current()
		if !found {
			continue
		}
		status := h.getOnlineStatus(n.LastBlock.Height)
		h.nodeStore.SetStatus(n.Id, int64(n.LastUpdate), status)
	}
}

func (h *Handler) getOnlineStatus(height int) int {
	globalHeight := h.statusStore.GetGlobalHeight()
	// if the height difference is smaller than 3600 (~60 minutes) -> node is online (1)
	if mathutil.Abs(globalHeight-height) < 3600 {
		return 1
	}
	return 0
}

// SendOnlineStatus sends the local status information about every node to Apache Kafka
func (h *Handler) SendOnlineStatus() {
	round := time.Now().UnixMilli() / 1000 / 60
	logger.Logger().Info().Int64("round", round).Msg("send status started")
	e := h.nodeStore.GetEnumerator()
	totalCount, sendCount, successCount := 0, 0, 0
	// limit concurrency to 10
	semaphore := make(chan struct{}, 10)
	for e.MoveNext() {
		totalCount++
		n, found := e.Current()
		if !found {
			continue
		}
		if n.LastBlock.Height <= 0 {
			logger.Logger().Debug().Str("name", n.Name).Str("id", n.Id).Msg("node skipped")
			continue
		}
		status := sharedtypes.KafkaNodeOnlineStatus{
			EventId:     fmt.Sprintf("%d_%s", round, n.Id),
			Timestamp:   strconv.FormatInt(time.Now().UnixMilli(), 10),
			Round:       round,
			ViteAddress: n.RewardAddress,
			NodeName:    n.Name,
			Ip:          n.ClientIp,
			Chain:       n.Chain,
		}
		if n.Status == 1 {
			status.SuccessTime = 1
			successCount++
		}
		semaphore <- struct{}{}
		go func() {
			writeStart := time.Now()
			h.kafka.WriteDefault(status)
			writeElapsed := time.Since(writeStart).Milliseconds()
			logger.Logger().Debug().Str("name", n.Name).Str("id", n.Id).Int("status", n.Status).Str("chain", n.Chain).Int64("elapsed", writeElapsed).Msg("node status sent")
			<-semaphore
		}()
		sendCount++
	}
	logger.Logger().Info().Int64("round", round).
		Int("totalCount", totalCount).Int("sendCount", sendCount).Int("successCount", successCount).
		Msg("send status finished")
}
