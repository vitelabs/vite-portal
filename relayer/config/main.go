package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/vitelabs/vite-portal/internal/logger"
	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/jsonutil"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to get the current filename")
	}
	dirname := filepath.Dir(filename)
	path := path.Join(dirname, "relayer_config.json")
	c := types.NewDefaultConfig()
	writeConfigFile(path, &c)
	fmt.Printf("config written to: %s\n", path)
}

func writeConfigFile(name string, c *types.Config) {
	jsonFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot open/create config json file")
	}
	b := jsonutil.ToByteIndentOrExit(c)
	_, err = jsonFile.Write(b)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot write default config to json file")
	}
}
