package types

import (
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/configutil"
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
	// The secret used for JSON Web Tokens
	JwtSecret string `json:"jwtSecret"`
	// The expiry timeout in milliseconds of JSON Web Tokens
	JwtExpiryTimeout int64 `json:"jwtExpiryTimeout"`
	// The maximum ip entries in the blacklist
	MaxIpBlacklistEntries int `json:"maxIpBlacklistEntries"`
	// The maximum ip blacklist duration in milliseconds
	MaxIpBlacklistDuration int64 `json:"maxIpBlacklistDuration"`
	// The user agent used when sending RPC requests
	UserAgent string `json:"userAgent"`
	// The true IP address of the client
	HeaderTrueClientIp string `json:"headerTrueClientIp"`
	// Apache Kafka related configuration
	Kafka sharedtypes.KafkaConfig `json:"kafka"`
	// A list of supported chains
	SupportedChains []sharedtypes.ChainConfig `json:"supportedChains"`
	supportedChains *sharedtypes.Chains
	// Logging related configuration
	Logging sharedtypes.LoggingConfig `json:"logging"`
}

func NewDefaultConfig() Config {
	c := Config{
		Version:                DefaultConfigVersion,
		Debug:                  DefaultDebug,
		RpcPort:                DefaultRpcPort,
		RpcAuthPort:            DefaultRpcAuthPort,
		RpcTimeout:             DefaultRpcTimeout,
		JwtSecret:              sharedtypes.DefaultJwtSecret,
		JwtExpiryTimeout:       sharedtypes.DefaultJwtExpiryTimeout,
		MaxIpBlacklistEntries:  DefaultMaxIpBlacklistEntries,
		MaxIpBlacklistDuration: DefaultMaxIpBlacklistDuration,
		UserAgent:              DefaultUserAgent,
		HeaderTrueClientIp:     DefaultHeaderTrueClientIp,
		Kafka:                  sharedtypes.DefaultKafkaConfig,
		SupportedChains:        sharedtypes.DefaultSupportedChains,
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

func (c *Config) GetChains() *sharedtypes.Chains {
	return c.supportedChains
}

func (c *Config) Validate() error {
	c.supportedChains = sharedtypes.NewChains(c.SupportedChains)
	err := configutil.ValidateChains(c.supportedChains)
	if err != nil {
		return err
	}
	return configutil.ValidateJwt(c.JwtSecret, c.JwtExpiryTimeout)
}
