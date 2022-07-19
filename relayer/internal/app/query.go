package app

func (app *RelayerCoreApp) QueryChains() []string {
	res := app.nodeService.GetChains()
	return res
}