package testutil

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

const (
	DefaultViteMainNodeUrl  = "https://node.vite.net/gvite"
	DefaultViteBuidlNodeUrl = "https://buidl.vite.net/gvite"
)

func BuildFullPath(elem ...string) string {
	_, filename, _, _ := runtime.Caller(0)
	rootName := "vite-portal"
	index := strings.Index(filepath.Dir(filename), rootName)
	if index == -1 {
		logger.Logger().Fatal().Msg(fmt.Sprintf("Root '%s' not found in %s", rootName, filename))
	}
	root := filename[0 : index+len(rootName)]
	p := path.Join(root, path.Join(elem...))
	logger.Logger().Info().Msg(fmt.Sprintf("root: %s, path: %s", root, p))
	_, err := os.Stat(p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
