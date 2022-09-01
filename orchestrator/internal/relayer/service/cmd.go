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

func (s *Service) HandleConnect(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var resp sharedtypes.RpcAppInfoResponse
	err := c.CallContext(ctx, &resp, "core_getAppInfo")
	if err != nil {
		logger.Logger().Error().Err(err).Msg("calling context failed")
		return err
	}
	if err := s.validateRelayerResponse(resp); err != nil {
		logger.Logger().Warn().Err(err).Msg("invalid relayer response")
		return err
	}

	logger.Logger().Debug().Str("resp", fmt.Sprintf("%#v", resp)).Msg("handle connect result")
	return nil
}

func (s *Service) HandleDisconnect(peerInfo rpc.PeerInfo) {
	logger.Logger().Debug().Str("info", fmt.Sprintf("%#v", peerInfo)).Msg("handle disconnect called")
}

func (s *Service) validateRelayerResponse(r sharedtypes.RpcAppInfoResponse) error {
	// TODO: verify if relayer id matches with JWT subject
	if r.Id == "" {
		return errors.New("invalid relayer id")
	}
	if r.Name != "vite-portal-relayer" {
		return errors.New("invalid relayer name")
	}
	return nil
}