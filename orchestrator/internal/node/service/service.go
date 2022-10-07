package service

import (
	"errors"
	"fmt"

	"github.com/vitelabs/vite-portal/orchestrator/internal/interfaces"
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/handler"
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
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
	for _, v := range cfg.GetChains().GetAll() {
		nodeStore, err := c.GetNodeStore(v.Name)
		if err != nil {
			panic(err)
		}
		statusStore, err := c.GetStatusStore(v.Name)
		if err != nil {
			panic(err)
		}
		s.handlers[v.Name] = handler.NewHandler(cfg, nodeStore, statusStore)
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
