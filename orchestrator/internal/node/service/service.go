package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/interfaces"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/client"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	config  types.Config
	context interfaces.ContextI
	clients map[string]*client.ViteClient
}

// NewService creates new instances of the relayers module service
func NewService(cfg types.Config, c interfaces.ContextI) *Service {
	s := &Service{
		config:  cfg,
		context: c,
		clients: map[string]*client.ViteClient{},
	}
	for _, v := range cfg.GetChains().GetAll() {
		url := v.OfficialNodeUrl
		if url == "" {
			logger.Logger().Warn().Str("chain", v.Name).Msg("OfficialNodeUrl is empty")
		}
		s.clients[v.Name] = client.NewViteClient(url)
	}
	return s
}

func (s *Service) GetViteClient(chain string) *client.ViteClient {
	return s.clients[chain]
}
