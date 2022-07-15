package app

import (
	"io/ioutil"
	"os"

	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/jsonutil"
)

const (
	AppName               = "vite-portal-relayer"
	AppVersion            = "0.0.1"
	defaultConfigFileName = "relayer_config.json"
)

var (
	CoreApp      *RelayerCoreApp
	GlobalConfig types.Config
)

func InitApp(debug bool) error {
	logger.Init(debug)
	err := InitConfig(debug)
	if err != nil {
		return err
	}
	orchestrator, err := InitOrchestrator()
	if err != nil {
		return err
	}
	CoreApp = NewRelayerCoreApp(orchestrator)
	return nil
}

func InitConfig(debug bool) error {
	c := types.NewDefaultConfig()
	logger.Logger().Info().RawJSON("config", jsonutil.ToByteOrExit(c)).Msg("DefaultConfig")

	// 1. Load config file
	loadConfigFromFile(&c)
	logger.Logger().Info().RawJSON("config", jsonutil.ToByteOrExit(c)).Msg("After loading config file")

	// 2. Apply flags, overwrite the loaded file configuration
	c.Debug = debug

	// 3. Configure logger
	logger.Configure(&c)

	// 4. Set global configuration
	GlobalConfig = c
	logger.Logger().Info().RawJSON("config", jsonutil.ToByteOrExit(GlobalConfig)).Msg("GlobalConfig")

	return nil
}

func loadConfigFromFile(c *types.Config) {
	var jsonFile *os.File
	defer jsonFile.Close()
	// if file exists open, else create and open
	if _, err := os.Stat(defaultConfigFileName); err == nil {
		jsonFile, err = os.OpenFile(defaultConfigFileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			logger.Logger().Fatal().Err(err).Msg("cannot open config json file")
		}
		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			logger.Logger().Fatal().Err(err).Msg("cannot read config file")
		}
		jsonutil.FromByteOrExit(b, &c)
	} else if os.IsNotExist(err) {
		// if does not exist create one
		jsonFile, err = os.OpenFile(defaultConfigFileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			logger.Logger().Fatal().Err(err).Msg("canot open/create config json file")
		}
		b := jsonutil.ToByteIndentOrExit(c)
		// write to the file
		_, err = jsonFile.Write(b)
		if err != nil {
			logger.Logger().Fatal().Err(err).Msg("cannot write default config to json file")
		}
	}
}