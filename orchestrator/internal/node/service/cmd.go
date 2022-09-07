package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/sliceutil"
)

func (s *Service) HandleConnect(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var nodeInfo sharedtypes.RpcViteNodeInfoResponse
	if err := c.CallContext(ctx, &nodeInfo, "net_nodeInfo"); err != nil {
		return s.returnConnectError("failed to call 'net_nodeInfo'", err)
	}
	chain := sharedtypes.Chains.GetById(nodeInfo.NetID)
	if chain.Id == sharedtypes.Chains.Unknown.Id || !sliceutil.Contains(s.config.SupportedChains, chain.Name) {
		return s.returnConnectError(fmt.Sprintf("chain id '%d' is not supported", nodeInfo.NetID), nil)
	}
	var processInfo sharedtypes.RpcViteProcessInfoResponse
	if err := c.CallContext(ctx, &processInfo, "dashboard_processInfo"); err != nil {
		return s.returnConnectError("failed to call 'dashboard_processInfo'", err)
	}
	n := types.Node{
		Id:            nodeInfo.ID,
		Name:          nodeInfo.Name,
		Chain:         chain.Name,
		Version:       processInfo.BuildVersion,
		Commit:        processInfo.CommitVersion,
		RewardAddress: processInfo.RewardAddress,
		RpcClient:     c,
		PeerInfo:      peerInfo,
	}
	if err := s.store.Upsert(n); err != nil {
		return s.returnConnectError("failed to upsert node", err)
	}
	return n.Id, nil
}

func (s *Service) returnConnectError(msg string, err error) (string, error) {
	if err != nil {
		logger.Logger().Error().Err(err).Msg(msg)
	} else {
		logger.Logger().Info().Msg(msg)
	}
	return "", errors.New(msg)
}
