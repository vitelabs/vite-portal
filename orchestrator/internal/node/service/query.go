package service

import (
	"errors"

	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

func (s *Service) Get(chain string, offset, limit int) (generics.GenericPage[types.Node], error) {
	logger.Logger().Debug().Str("chain", chain).Msg("get nodes")
	cc, err := s.context.GetChainContext(chain)
	if err != nil {
		return *generics.NewGenericPage[types.Node](), err
	}
	store := cc.GetNodeStore()
	total := store.Count()
	result := *generics.NewGenericPage[types.Node]()
	result.Offset = offset
	result.Limit = limit
	result.Total = total
	if offset >= total {
		return result, nil
	}
	count := mathutil.Min(total-result.Offset, limit)
	result.Entries = make([]types.Node, count)
	current := 0
	for i := result.Offset; i < count; i++ {
		item, found := store.GetByIndex(i)
		if !found {
			return *generics.NewGenericPage[types.Node](), errors.New("inconsistent state")
		}
		result.Entries[current] = item
		current++
	}
	return result, nil
}