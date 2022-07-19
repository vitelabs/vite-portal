package service

import (
	"github.com/vitelabs/vite-portal/internal/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
)

func (k Service) GetChains() []string {
	return k.store.GetChains()
}

func (k Service) GetNodes(chain string, page, limit int) generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}

func paginate(page, limit int, nodes []types.Node, MaxNodes int) generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}