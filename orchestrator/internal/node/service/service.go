package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/orchestrator/internal/interfaces"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/handler"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	sharedclients "github.com/vitelabs/vite-portal/shared/pkg/client"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

// Service maintains the link to storage and exposes getter/setter methods for handling relayers
type Service struct {
	config   types.Config
	context  interfaces.ContextI
	handlers map[string]*handler.Handler
}

// NewService creates new instances of the relayers module service
func NewService(cfg types.Config, c interfaces.ContextI) *Service {
	s := &Service{
		config:   cfg,
		context:  c,
		handlers: map[string]*handler.Handler{},
	}
	timeout := time.Duration(cfg.RpcTimeout) * time.Millisecond
	for _, v := range cfg.GetChains().GetAll() {
		nodeStore, err := c.GetNodeStore(v.Name)
		if err != nil {
			panic(err)
		}
		statusStore, err := c.GetStatusStore(v.Name)
		if err != nil {
			panic(err)
		}
		url := v.OfficialNodeUrl
		if url == "" {
			logger.Logger().Warn().Str("chain", v.Name).Msg("OfficialNodeUrl is empty")
		}
		client := sharedclients.NewViteClient(url, timeout)
		s.handlers[v.Name] = handler.NewHandler(cfg, client, nodeStore, statusStore)
	}
	return s
}

func (s *Service) GetHandler(chain string) (*handler.Handler, error) {
	h := s.handlers[chain]
	if h == nil {
		return nil, errors.New(fmt.Sprintf("handler not found for chain '%s'", chain))
	}
	return h, nil
}
