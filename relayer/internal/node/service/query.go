package service

import (
	"github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
)

func (s *Service) GetChains() []string {
	return s.store.GetChains()
}

func (s *Service) GetNodeCount(chain string) int {
	return s.store.Count(chain)
}

func (s *Service) GetNodes(chain string, offset, limit int) (generics.GenericPage[types.Node], error) {
	return s.store.GetPaginated(chain, offset, limit)
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
