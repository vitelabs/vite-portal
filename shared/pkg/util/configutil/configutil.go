package configutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/vitelabs/vite-portal/shared/pkg/interfaces"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/jsonutil"
)

func LoadConfigFromFile[T interfaces.ConfigI](p, defaultVersion string, c *T) {
	var jsonFile *os.File
	defer jsonFile.Close()
	// if file does not exist -> create, otherwise open and compare version
	if _, err := os.Stat(p); os.IsNotExist(err) {
		WriteConfigFile(p, c)
		return
	}
	jsonFile, err := os.OpenFile(p, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot open config json file")
	}
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot read config file")
	}
	loaded := new(T)
	jsonutil.FromByteOrExit(b, &loaded)
	if (*loaded).GetVersion() != defaultVersion {
		dir := filepath.Dir(p)
		name := filepath.Base(p)
		// config schema versions do not match -> write backup
		WriteConfigFile(path.Join(dir, fmt.Sprintf("%s_%d", name, time.Now().UnixMilli())), loaded)
		// write new config with default values
		WriteConfigFile(p, c)
		return
	}
	// config schema versions do match -> set config
	jsonutil.FromByteOrExit(b, &c)
}

func WriteConfigFile(p string, c any) {
	dir := filepath.Dir(p)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logger.Logger().Fatal().Err(err).Msg("cannot create directory")
		}
	}
	jsonFile, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot open/create config json file")
	}
	b := jsonutil.ToByteIndentOrExit(c)
	_, err = jsonFile.Write(b)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("cannot write default config to json file")
	}
}