package keeper

import (
	"github.com/vitelabs/vite-portal/internal/core/generics"
	"github.com/vitelabs/vite-portal/internal/nodes/types"
)

func paginate(page, limit int, nodes []types.Node, MaxNodes int) generics.GenericPage[types.Node] {
	return generics.GenericPage[types.Node]{}
}