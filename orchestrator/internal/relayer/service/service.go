package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/interfaces"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/crypto"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	jwtHandler crypto.JWTHandler
	store      interfaces.StoreI
}

// NewService creates new instances of the relayers module service
func NewService(cfg types.Config, store interfaces.StoreI) *Service {
	return &Service{
		jwtHandler: *crypto.NewDefaultJWTHandler([]byte(cfg.JwtSecret)),
		store:      store,
	}
}

func (s *Service) IsRelayerConnection(peerInfo rpc.PeerInfo) bool {
	token, err := s.jwtHandler.Extract(peerInfo.HTTP.Header)
	if err != nil {
		return false
	}
	claims, err := s.jwtHandler.Validate(token)
	if err != nil {
		return false
	}
	if claims.Subject != sharedtypes.JWTRelayerSubject {
		logger.Logger().Info().Str("subject", claims.Subject).Msg("invalid JWT subject")
		return false
	}
	return true
}
