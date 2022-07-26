package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

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
	context *Context
	CoreApp *RelayerCoreApp
)

func InitApp(debug bool) error {
	logger.Init(debug)
	err := InitConfig(debug)
	if err != nil {
		return err
	}
	err = InitCache(types.GlobalConfig.MaxSessionCacheEntries)
	if err != nil {
		return err
	}
	o, err := InitOrchestrator()
	if err != nil {
		return err
	}
	c, err := InitContext()
	if err != nil {
		return err
	}
	context = c
	CoreApp = NewRelayerCoreApp(o, c)
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
	types.GlobalConfig = c
	logger.Logger().Info().RawJSON("config", jsonutil.ToByteOrExit(types.GlobalConfig)).Msg("GlobalConfig")

	return nil
}

func loadConfigFromFile(c *types.Config) {
	var jsonFile *os.File
	defer jsonFile.Close()
	filename := defaultConfigFileName
	// if file does not exist -> create, otherwise open and compare version
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		writeConfigFile(filename, c)
		return
	}
	jsonFile, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot open config json file")
	}
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot read config file")
	}
	loaded := &types.Config{}
	jsonutil.FromByteOrExit(b, &loaded)
	if loaded.Version != types.DefaultConfigVersion {
		// config schema versions do not match -> write backup
		writeConfigFile(fmt.Sprintf("%s_%d", filename, time.Now().UnixMilli()), loaded)
		// write new config with default values
		writeConfigFile(filename, c)
	}
}

func writeConfigFile(name string, c *types.Config) {
	jsonFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("canot open/create config json file")
	}
	b := jsonutil.ToByteIndentOrExit(c)
	_, err = jsonFile.Write(b)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot write default config to json file")
	}
}

func Shutdown() {
	logger.Logger().Info().Msg("Shutdown called")
	context.nodeStore.Close()
}
