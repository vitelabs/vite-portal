package service

import (
	"errors"

	"github.com/vitelabs/vite-portal/orchestrator/internal/relayer/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

func (s *Service) Get(offset, limit int) (generics.GenericPage[types.Relayer], error) {
	total := s.store.Count()
	result := *generics.NewGenericPage[types.Relayer]()
	result.Offset = offset
	result.Limit = limit
	result.Total = total
	if offset >= total {
		return result, nil
	}
	result.Entries = make([]types.Relayer, mathutil.Min(total-result.Offset, limit))
	count := mathutil.Min(result.Offset+result.Limit, total)
	current := 0
	for i := result.Offset; i < count; i++ {
		item, found := s.store.GetByIndex(i)
		if !found {
			return *generics.NewGenericPage[types.Relayer](), errors.New("inconsistent state")
		}
		result.Entries[current] = item
		current++
	}
	return result, nil
}