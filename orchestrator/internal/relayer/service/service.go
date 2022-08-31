package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/interfaces"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	store interfaces.StoreI
}

// NewService creates new instances of the relayers module service
func NewService(store interfaces.StoreI) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) IsRelayerConnection(peerInfo rpc.PeerInfo) bool {
	// TODO: verify peerInfo.Auth extracted from header
	return true
}
