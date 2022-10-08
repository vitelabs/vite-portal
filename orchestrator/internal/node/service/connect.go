package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/httputil"
)

func (s *Service) HandleConnect(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var nodeInfo sharedtypes.RpcViteNodeInfoResponse
	if err := c.CallContext(ctx, &nodeInfo, "net_nodeInfo"); err != nil {
		return s.returnConnectError("failed to call 'net_nodeInfo'", err)
	}
	logger.Logger().Debug().Str("nodeInfo", fmt.Sprintf("%#v", nodeInfo)).Msg("handle connect response")
	chain, found := s.config.GetChains().GetById(strconv.Itoa(nodeInfo.NetID))
	if !found {
		return s.returnConnectError(fmt.Sprintf("chain id '%d' is not supported", nodeInfo.NetID), nil)
	}
	var processInfo sharedtypes.RpcViteProcessInfoResponse
	if err := c.CallContext(ctx, &processInfo, "dashboard_processInfo", "param1"); err != nil {
		return s.returnConnectError("failed to call 'dashboard_processInfo'", err)
	}
	logger.Logger().Debug().Str("processInfo", fmt.Sprintf("%#v", processInfo)).Msg("handle connect response")
	clientIp := httputil.GetClientIp(peerInfo.HTTP.Header, s.config.HeaderTrueClientIp)
	n := types.Node{
		Id:            nodeInfo.ID,
		Name:          nodeInfo.Name,
		Chain:         chain.Name,
		Version:       processInfo.BuildVersion,
		Commit:        processInfo.CommitVersion,
		RewardAddress: processInfo.RewardAddress,
		Transport:     peerInfo.Transport,
		RemoteAddress: peerInfo.RemoteAddr,
		ClientIp:      clientIp,
		Status:        0,
		LastUpdate:    0,
		DelayTime:     0,
		RpcClient:     c,
		HTTPInfo: sharedtypes.HTTPInfo{
			Version:   peerInfo.HTTP.Version,
			UserAgent: peerInfo.HTTP.UserAgent,
			Origin:    peerInfo.HTTP.Origin,
			Host:      peerInfo.HTTP.Host,
		},
	}
	msg := "failed to add node"
	store, err := s.context.GetNodeStore(chain.Name)
	if err != nil {
		return s.returnConnectError(msg, err)
	}
	if err := store.Add(n); err != nil {
		return s.returnConnectError(msg, err)
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
