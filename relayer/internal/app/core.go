package app

import (
	"github.com/vitelabs/vite-portal/internal/core/keeper"
	"github.com/vitelabs/vite-portal/internal/core/types"

	nodeskeeper "github.com/vitelabs/vite-portal/internal/nodes/keeper"
)

type RelayerCoreApp struct {
	coreKeeper keeper.Keeper
	nodesKeeper nodeskeeper.Keeper
}

func NewRelayerCoreApp() *RelayerCoreApp {
	app := &RelayerCoreApp{}
	nodesKeeper := nodeskeeper.NewKeeper()
	app.coreKeeper = keeper.NewKeeper(nodesKeeper)
	return app
}

func (app *RelayerCoreApp) HandleRelay(r types.Relay) (string, error) {
	res, err := app.coreKeeper.HandleRelay(r)
	if err != nil {
		return "", err
	}
	return res.Response, nil
}
