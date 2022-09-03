package app

import (
	"github.com/vitelabs/vite-portal/orchestrator/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
)

const (
	RpcCoreModule     = "core"
	RpcRelayersModule = "relayers"
	RpcDebugModule    = "debug"
)

var (
	DefaultAllowedOrigins = []string{"*"}
	DefaultVhosts         = []string{"localhost"}
	DefaultModules        = []string{RpcCoreModule, RpcRelayersModule}
)

func InitApp(debug bool, configPath string) (*OrchestratorApp, error) {
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
	app := NewOrchestratorApp(cfg)
	return app, nil
}
