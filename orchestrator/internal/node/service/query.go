package service

import (
	"errors"

	"github.com/vitelabs/vite-portal/orchestrator/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

func (s *Service) GetNodes(chain string, offset, limit int) (generics.GenericPage[types.Node], error) {
	total := s.store.Count(chain)
	result := *generics.NewGenericPage[types.Node]()
	result.Offset = offset
	result.Limit = limit
	result.Total = total
	if offset >= total {
		return result, nil
	}
	result.Entries = make([]types.Node, mathutil.Min(total-result.Offset, limit))
	count := mathutil.Min(result.Offset+result.Limit, total)
	current := 0
	for i := result.Offset; i < count; i++ {
		item, found := s.store.GetByIndex(chain, i)
		if !found {
			return *generics.NewGenericPage[types.Node](), errors.New("inconsistent state")
		}
		result.Entries[current] = item
		current++
	}
	return result, nil
}