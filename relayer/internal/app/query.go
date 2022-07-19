package app

import (
	"github.com/vitelabs/vite-portal/internal/generics"
	"github.com/vitelabs/vite-portal/internal/node/types"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
)

func (app *RelayerCoreApp) QueryChains() []string {
	res := app.nodeService.GetChains()
	return res
}

func (app *RelayerCoreApp) QueryNodes(p nodetypes.QueryNodesParams) (res generics.GenericPage[types.Node], err error) {
	p.Offset, p.Limit = checkPagination(p.Offset, p.Limit)
	return app.nodeService.GetNodes(p.Chain, p.Offset, p.Limit)
}

func checkPagination(offset, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 1000 {
		limit = 1000
	}
	return offset, limit
}