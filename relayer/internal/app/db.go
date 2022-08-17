package app

import (
	"github.com/vitelabs/vite-portal/relayer/internal/node/types"
	nodetypes "github.com/vitelabs/vite-portal/relayer/internal/node/types"
	"github.com/vitelabs/vite-portal/shared/pkg/generics"
)

func (app *RelayerCoreApp) GetChains() []string {
	res := app.nodeService.GetChains()
	return res
}

func (app *RelayerCoreApp) GetNodes(p nodetypes.GetNodesParams) (res generics.GenericPage[types.Node], err error) {
	p.Offset, p.Limit = checkPagination(p.Offset, p.Limit)
	return app.nodeService.GetNodes(p.Chain, p.Offset, p.Limit)
}

func (app *RelayerCoreApp) GetNode(id string) (res types.Node, found bool) {
	return app.nodeService.GetNode(id)
}

func (app *RelayerCoreApp) PutNode(n nodetypes.Node) error {
	return app.nodeService.PutNode(n)
}

func (app *RelayerCoreApp) DeleteNode(id string) error {
	return app.nodeService.DeleteNode(id)
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