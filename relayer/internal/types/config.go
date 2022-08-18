package types

import (
	"errors"
	"fmt"

	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

const (
	DefaultConfigVersion               = "v0.3"
	DefaultDebug                       = false
	DefaultRpcHttpPort                 = 56331
	DefaultRpcWsPort                   = 56332
	DefaultRpcTimeout                  = 10000
	DefaultRpcNodeTimeout              = 5000
	DefaultUserAgent                   = ""
	DefaultSortJsonResponse            = false
	DefaultConsensusNodeCount          = 5
	DefaultSessionNodeCount            = 24
	DefaultMaxSessionCacheEntries      = 100
	DefaultMaxSessionDuration          = 60000000
	DefaultHeaderTrueClientIp          = "True-Client-Ip"
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
	// Port number for the HTTP RPC
	RpcHttpPort int32 `json:"rpcHttpPort"`
	// Port number for the WebSocket RPC
	RpcWsPort int32 `json:"rpcWsPort"`
	// The time in milliseconds before a RPC request times out
	RpcTimeout int64 `json:"rpcTimeout"`
	// The time in milliseconds before a RPC request to a node times out
	RpcNodeTimeout int64 `json:"rpcNodeTimeout"`
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
		RpcHttpPort:            DefaultRpcHttpPort,
		RpcWsPort:              DefaultRpcWsPort,
		RpcTimeout:             DefaultRpcTimeout,
		RpcNodeTimeout:         DefaultRpcNodeTimeout,
		UserAgent:              DefaultUserAgent,
		SortJsonResponse:       DefaultSortJsonResponse,
		ConsensusNodeCount:     DefaultConsensusNodeCount,
		SessionNodeCount:       DefaultSessionNodeCount,
		MaxSessionCacheEntries: DefaultMaxSessionCacheEntries,
		MaxSessionDuration:     DefaultMaxSessionDuration,
		HeaderTrueClientIp:     DefaultHeaderTrueClientIp,
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

func (c Config) GetVersion() string {
	return c.Version
}

func (c *Config) Validate() error {
	prefix := "Config error: "
	if c.SessionNodeCount <= 0 {
		return errors.New(fmt.Sprintf("%s SessionNodeCount must be greater than 0", prefix))
	}
	return nil
}
