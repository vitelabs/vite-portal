package app

import (
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
)

var (
	CoreApp *RelayerCoreApp
)

func InitApp(debug bool, configPath string) error {
	logger.Init(debug)
	p := configPath
	if p == "" {
		p = types.DefaultConfigFilename
	}
	cfg := types.NewDefaultConfig()
	err := configutil.InitConfig(&cfg, debug, p, types.DefaultConfigVersion)
	if err != nil {
		return err
	}
	o := InitOrchestrator(cfg.OrchestratorWsUrl, cfg.RpcTimeout)
	c, err := InitContext(cfg)
	if err != nil {
		return err
	}
	CoreApp = NewRelayerCoreApp(cfg, o, c)
	return nil
}

func Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	CoreApp.context.nodeStore.Close()
}
