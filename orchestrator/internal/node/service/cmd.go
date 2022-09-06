package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func (s *Service) HandleConnect(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var resp sharedtypes.RpcViteNodeInfoResponse
	if err := c.CallContext(ctx, &resp, "net_nodeInfo"); err != nil {
		logger.Logger().Error().Err(err).Msg("calling context failed")
		return "", err
	}
	chain := sharedtypes.Chains.GetById(resp.NetID)
	if chain.Id == sharedtypes.Chains.Unknown.Id {
		return "", errors.New(fmt.Sprintf("chain id '%d' is not supported", resp.NetID))
	}
	return "", errors.New("not implemented yet")
}
