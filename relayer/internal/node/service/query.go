package service

import (
	"errors"

	"github.com/vitelabs/vite-portal/internal/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/util/mathutil"
)

func (k Service) GetChains() []string {
	return k.store.GetChains()
}

func (k Service) GetNodes(chain string, offset, limit int) (generics.GenericPage[types.Node], error) {
	total := k.store.Count(chain)
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
		item, found := k.store.GetByIndex(chain, i)
		if !found {
			return *generics.NewGenericPage[types.Node](), errors.New("inconsistent state")
		}
		result.Entries[current] = item
		current++
	}
	return result, nil
}

func paginate(page, limit int, nodes []types.Node, MaxNodes int) generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}
