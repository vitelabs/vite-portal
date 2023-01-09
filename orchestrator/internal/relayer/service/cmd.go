package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

func (s *Service) IsRelayerConnection(claims jwt.RegisteredClaims) bool {
	return claims.Issuer == sharedtypes.JWTRelayerIssuer
}

func (s *Service) HandleConnect(timeout time.Duration, c *rpc.Client, peerInfo rpc.PeerInfo, claims jwt.RegisteredClaims) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var resp sharedtypes.RpcAppInfoResponse
	if err := c.CallContext(ctx, &resp, "core_getAppInfo"); err != nil {
		logger.Logger().Error().Err(err).Msg("calling context failed")
		return "", err
	}
	logger.Logger().Debug().Str("resp", fmt.Sprintf("%#v", resp)).Msg("handle connect response")
	if err := s.validateRelayerResponse(resp, claims); err != nil {
		logger.Logger().Warn().Err(err).Msg("invalid relayer response")
		return "", err
	}
	id, err = s.insertRelayer(c, peerInfo, resp)
	if err != nil {
		logger.Logger().Warn().Err(err).Msg("insert relayer failed")
		return "", err
	}
	return id, nil
}

func (s *Service) HandleDisconnect(id string) {
	logger.Logger().Debug().Str("id", id).Msg("handle disconnect called")
	s.store.Remove(id)
}

func (s *Service) validateRelayerResponse(r sharedtypes.RpcAppInfoResponse, claims jwt.RegisteredClaims) error {
	if r.Id == "" || r.Id != claims.Subject {
		return errors.New("invalid relayer id")
	}
	if r.Name != "vite-portal-relayer" {
		return errors.New("invalid relayer name")
	}
	return nil
}

func (s *Service) insertRelayer(c *rpc.Client, peerInfo rpc.PeerInfo, r sharedtypes.RpcAppInfoResponse) (id string, err error) {
	relayer := types.Relayer{
		Id:            r.Id,
		Version:       r.Version,
		Transport:     peerInfo.Transport,
		RemoteAddress: peerInfo.RemoteAddr,
		RpcClient:     c,
		HTTPInfo: sharedtypes.HTTPInfo{
			Version:   peerInfo.HTTP.Version,
			UserAgent: peerInfo.HTTP.UserAgent,
			Origin:    peerInfo.HTTP.Origin,
			Host:      peerInfo.HTTP.Host,
		},
	}
	if err := s.store.Upsert(relayer); err != nil {
		return "", err
	}
	return relayer.Id, nil
}
