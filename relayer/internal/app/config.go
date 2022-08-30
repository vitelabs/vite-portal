package app

import (
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
)

const (
	RpcCoreModule  = "core"
	RpcNodesModule = "nodes"
	RpcDebugModule = "debug"
)

var (
	DefaultAllowedOrigins = []string{"*"}
	DefaultVhosts         = []string{"localhost"}
	DefaultModules        = []string{RpcCoreModule, RpcNodesModule}
)

func InitApp(debug bool, configPath string) (*RelayerApp, error) {
	logger.Init(debug)
	p := configPath
	if p == "" {
		p = types.DefaultConfigFilename
	}
	cfg := types.NewDefaultConfig()
	err := configutil.InitConfig(&cfg, debug, p, types.DefaultConfigVersion)
	if err != nil {
		return nil, err
	}
	a := NewRelayerApp(cfg)
	return a, nil
}
