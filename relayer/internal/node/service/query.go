package service

import (
	"errors"

	"github.com/vitelabs/vite-portal/relayer/internal/generics"
	"github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/mathutil"
)

func (s *Service) GetChains() []string {
	return s.store.GetChains()
}

func (s *Service) GetNodeCount(chain string) int {
	return s.store.Count(chain)
}

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

func (s *Service) GetNode(id string) (n types.Node, found bool) {
	return s.store.GetById(id)
}

func (s *Service) GetNodeByIndex(chain string, index int) (n types.Node, found bool) {
	return s.store.GetByIndex(chain, index)
}

func paginate(page, limit int, nodes []types.Node, MaxNodes int) generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}
