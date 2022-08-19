package types

import (
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

const (
	DefaultConfigVersion               = "v0.1"
	DefaultDebug                       = false
	DefaultRpcHttpPort                 = 57331
	DefaultRpcWsPort                   = 57332
	DefaultRpcTimeout                  = 5000
	DefaultUserAgent                   = ""
	DefaultHeaderTrueClientIp          = "True-Client-Ip"
	DefaultChain                       = "vite_mainnet"
	DefaultLoggingConsoleOutputEnabled = true
	DefaultLoggingFileOutputEnabled    = true
	DefaultLoggingDirectory            = "logs"
	DefaultLoggingFilename             = "relayer.log"
	DefaultLoggingMaxSize              = 100
	DefaultLoggingMaxBackups           = 10
	DefaultLoggingMaxAge               = 28
)

var (
	DefaultSupportedChains = []string{"vite_mainnet", "vite_testnet"}
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
	// The user agent used when sending RPC requests
	UserAgent string `json:"userAgent"`
	// The true IP address of the client
	HeaderTrueClientIp string `json:"headerTrueClientIp"`
	// The default chain name
	DefaultChain string `json:"defaultChain"`
	// A list of supported chain names
	SupportedChains []string `json:"supportedChains"`
	// Logging related configurtion
	Logging sharedtypes.LoggingConfig `json:"logging"`
}

func NewDefaultConfig() Config {
	c := Config{
		Version:            DefaultConfigVersion,
		Debug:              DefaultDebug,
		RpcHttpPort:        DefaultRpcHttpPort,
		RpcWsPort:          DefaultRpcWsPort,
		RpcTimeout:         DefaultRpcTimeout,
		UserAgent:          DefaultUserAgent,
		HeaderTrueClientIp: DefaultHeaderTrueClientIp,
		DefaultChain:       DefaultChain,
		SupportedChains:    DefaultSupportedChains,
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
	return nil
}
