package service

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

func (s *Service) Get(chain string, offset, limit int) (generics.GenericPage[types.Node], error) {
	logger.Logger().Debug().Str("chain", chain).Msg("get nodes")
	cc, err := s.context.GetChainContext(chain)
	if err != nil {
		return *generics.NewGenericPage[types.Node](), err
	}
	store := cc.GetNodeStore()
	return store.GetPaginated(offset, limit)
}