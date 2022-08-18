package app

import (
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

var (
	CoreApp *RelayerCoreApp
)

func InitApp(debug bool, configPath string) error {
	logger.Init(debug)
	cfg, err := InitConfig(debug, configPath)
	if err != nil {
		return err
	}
	o, err := InitOrchestrator()
	if err != nil {
		return err
	}
	c, err := InitContext(cfg)
	if err != nil {
		return err
	}
	CoreApp = NewRelayerCoreApp(cfg, o, c)
	return nil
}

func InitConfig(debug bool, configPath string) (types.Config, error) {
	c := types.NewDefaultConfig()
	logger.Logger().Info().RawJSON("config", jsonutil.ToByteOrExit(c)).Msg("DefaultConfig")

	// 1. Load config file
	if configPath == "" {
		configutil.LoadConfigFromFile(types.DefaultConfigFilename, types.DefaultConfigVersion, &c)
	} else {
		configutil.LoadConfigFromFile(configPath, types.DefaultConfigVersion, &c)
	}
	logger.Logger().Info().RawJSON("config", jsonutil.ToByteOrExit(c)).Msg("After loading config file")

	// 2. Apply flags, overwrite the loaded file configuration
	if debug {
		c.Debug = debug
	}

	// 3. Configure logger
	logger.Configure(c.Debug, c.Logging)
	logger.Logger().Info().RawJSON("config", jsonutil.ToByteOrExit(c)).Msg("GlobalConfig")

	// 4. Validate
	err := c.Validate()
	if err != nil {
		return types.Config{}, err
	}

	return c, nil
}

func Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	CoreApp.context.nodeStore.Close()
}
