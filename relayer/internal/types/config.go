package types

import (
	"errors"
	"fmt"

	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
)

const (
	DefaultConfigVersion               = "v0.1"
	DefaultDebug                       = false
	DefaultRpcPort                     = 56331
	DefaultRpcAuthPort                 = 56332
	DefaultRpcRelayHttpPort            = 56333
	DefaultRpcRelayWsPort              = 56334
	DefaultRpcTimeout                  = 10000
	DefaultRpcNodeTimeout              = 5000
	DefaultJwtSecret                   = "secret1234"
	DefaultJwtExpiryTimeout            = 0
	DefaultUserAgent                   = ""
	DefaultSortJsonResponse            = false
	DefaultConsensusNodeCount          = 5
	DefaultSessionNodeCount            = 24
	DefaultMaxSessionCacheEntries      = 10000
	DefaultMaxSessionDuration          = 3600000
	DefaultHeaderTrueClientIp          = "CF-Connecting-IP"
	DefaultOrchestratorWsUrl           = ""
	DefaultLoggingConsoleOutputEnabled = true
	DefaultLoggingFileOutputEnabled    = true
	DefaultLoggingDirectory            = "logs"
	DefaultLoggingFilename             = "relayer.log"
	DefaultLoggingMaxSize              = 100
	DefaultLoggingMaxBackups           = 10
	DefaultLoggingMaxAge               = 28
)

type Config struct {
	// Version of the configuration schema
	Version string `json:"version"`
	// Enable debug mode
	Debug bool `json:"debug"`
	// Port number for the unauthenticated RPC
	RpcPort int32 `json:"rpcPort"`
	// Port number for the authenticated RPC
	RpcAuthPort int32 `json:"rpcAuthPort"`
	// Port number for the relay HTTP RPC
	RpcRelayHttpPort int32 `json:"rpcRelayHttpPort"`
	// Port number for the relay WebSocket RPC
	RpcRelayWsPort int32 `json:"rpcRelayWsPort"`
	// The time in milliseconds before a RPC request times out
	RpcTimeout int64 `json:"rpcTimeout"`
	// The time in milliseconds before a RPC request to a node times out
	RpcNodeTimeout int64 `json:"rpcNodeTimeout"`
	// The secret used for JSON Web Tokens
	JwtSecret string `json:"jwtSecret"`
	// The expiry timeout in milliseconds of JSON Web Tokens
	JwtExpiryTimeout int64 `json:"jwtExpiryTimeout"`
	// The user agent used when sending RPC requests to nodes
	UserAgent string `json:"userAgent"`
	// Whether the JSON response from nodes should be sorted
	SortJsonResponse bool `json:"sortJsonResponse"`
	// The number of nodes a request will be forwarded to
	ConsensusNodeCount int `json:"consensusNodeCount"`
	// The number of nodes a relay request will be matched within a session
	SessionNodeCount int `json:"sessionNodeCount"`
	// The maximum session entries in the cache
	MaxSessionCacheEntries int `json:"maxSessionCacheEntries"`
	// The maximum session duration in milliseconds
	MaxSessionDuration int64 `json:"maxSessionDuration"`
	// The true IP address of the client
	HeaderTrueClientIp string `json:"headerTrueClientIp"`
	// The WebSocket URL of the orchestrator
	OrchestratorWsUrl string `json:"orchestratorWsUrl"`
	// The optional HttpCollector URL to which all relay results will be sent to
	HttpCollectorUrl string `json:"httpCollectorUrl,omitempty"`
	// The entries to map a host to the respective chain
	HostToChainMap map[string]string `json:"hostToChainMap"`
	// Logging related configurtion
	Logging sharedtypes.LoggingConfig `json:"logging"`
}

func NewDefaultConfig() Config {
	c := Config{
		Version:                DefaultConfigVersion,
		Debug:                  DefaultDebug,
		RpcPort:                DefaultRpcPort,
		RpcAuthPort:            DefaultRpcAuthPort,
		RpcRelayHttpPort:       DefaultRpcRelayHttpPort,
		RpcRelayWsPort:         DefaultRpcRelayWsPort,
		RpcTimeout:             DefaultRpcTimeout,
		RpcNodeTimeout:         DefaultRpcNodeTimeout,
		JwtSecret:              DefaultJwtSecret,
		JwtExpiryTimeout:       DefaultJwtExpiryTimeout,
		UserAgent:              DefaultUserAgent,
		SortJsonResponse:       DefaultSortJsonResponse,
		ConsensusNodeCount:     DefaultConsensusNodeCount,
		SessionNodeCount:       DefaultSessionNodeCount,
		MaxSessionCacheEntries: DefaultMaxSessionCacheEntries,
		MaxSessionDuration:     DefaultMaxSessionDuration,
		HeaderTrueClientIp:     DefaultHeaderTrueClientIp,
		OrchestratorWsUrl:      DefaultOrchestratorWsUrl,
		HostToChainMap:         map[string]string{},
		Logging: sharedtypes.LoggingConfig{
			ConsoleOutputEnabled: DefaultLoggingConsoleOutputEnabled,
			FileOutputEnabled:    DefaultLoggingFileOutputEnabled,
			Directory:            DefaultLoggingDirectory,
			Filename:             DefaultLoggingFilename,
			MaxSize:              DefaultLoggingMaxSize,
			MaxBackups:           DefaultLoggingMaxBackups,
			MaxAge:               DefaultLoggingMaxAge,
		},
	}
	return c
}

func (c *Config) GetVersion() string {
	return c.Version
}

func (c *Config) GetDebug() bool {
	return c.Debug
}

func (c *Config) SetDebug(debug bool) {
	c.Debug = debug
}

func (c *Config) GetLoggingConfig() sharedtypes.LoggingConfig {
	return c.Logging
}

func (c *Config) Validate() error {
	prefix := "Config error: "
	if c.SessionNodeCount <= 0 {
		return errors.New(fmt.Sprintf("%s SessionNodeCount must be greater than 0", prefix))
	}
	return configutil	.ValidateJwt(c.JwtSecret, c.JwtExpiryTimeout)
}
