package types

import (
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

const (
	DefaultConfigVersion               = "v0.1"
	DefaultDebug                       = false
	DefaultRpcPort                     = 57331
	DefaultRpcAuthPort                 = 57332
	DefaultRpcTimeout                  = 5000
	DefaultMaxIpBlacklistEntries       = 10000
	DefaultMaxIpBlacklistDuration      = 5000
	DefaultUserAgent                   = ""
	DefaultHeaderTrueClientIp          = "CF-Connecting-IP"
	DefaultLoggingConsoleOutputEnabled = true
	DefaultLoggingFileOutputEnabled    = true
	DefaultLoggingDirectory            = "logs"
	DefaultLoggingFilename             = "orchestrator.log"
	DefaultLoggingMaxSize              = 100
	DefaultLoggingMaxBackups           = 10
	DefaultLoggingMaxAge               = 28
)

var (
	DefaultSupportedChains = []string{
		sharedtypes.Chains.ViteMain.Name,
		sharedtypes.Chains.ViteBuidl.Name,
	}
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
	// The time in milliseconds before a RPC request times out
	RpcTimeout int64 `json:"rpcTimeout"`
	// The maximum ip entries in the blacklist
	MaxIpBlacklistEntries int `json:"maxIpBlacklistEntries"`
	// The maximum ip blacklist duration in milliseconds
	MaxIpBlacklistDuration int64 `json:"maxIpBlacklistDuration"`
	// The user agent used when sending RPC requests
	UserAgent string `json:"userAgent"`
	// The true IP address of the client
	HeaderTrueClientIp string `json:"headerTrueClientIp"`
	// A list of supported chain names
	SupportedChains []string `json:"supportedChains"`
	// Logging related configurtion
	Logging sharedtypes.LoggingConfig `json:"logging"`
}

func NewDefaultConfig() Config {
	c := Config{
		Version:                DefaultConfigVersion,
		Debug:                  DefaultDebug,
		RpcPort:                DefaultRpcPort,
		RpcAuthPort:            DefaultRpcAuthPort,
		RpcTimeout:             DefaultRpcTimeout,
		MaxIpBlacklistEntries:  DefaultMaxIpBlacklistEntries,
		MaxIpBlacklistDuration: DefaultMaxIpBlacklistDuration,
		UserAgent:              DefaultUserAgent,
		HeaderTrueClientIp:     DefaultHeaderTrueClientIp,
		SupportedChains:        DefaultSupportedChains,
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
