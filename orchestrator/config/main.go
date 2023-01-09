package main

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
)

// called by Makefile during build process
func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to get the current filename")
	}
	dirname := filepath.Dir(filename)
	path := path.Join(dirname, types.DefaultConfigFilename)
	c := types.NewDefaultConfig()
	configutil.WriteConfigFile(path, &c)
	fmt.Printf("config written to: %s\n", path)
}
