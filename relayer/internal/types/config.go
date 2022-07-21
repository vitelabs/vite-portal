package types

var GlobalConfig Config

const (
	DefaultConfigVersion               = "v0.1"
	DefaultDebug                       = false
	DefaultRpcHttpPort                 = 56331
	DefaultRpcTimeout                  = 10000
	DefaultRpcNodeTimeout              = 5000
	DefaultUserAgent                   = ""
	DefaultSortJsonResponse            = false
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
	// The time in milliseconds before a RPC request times out
	RpcTimeout int64 `json:"rpcTimeout"`
	// The time in milliseconds before a RPC request to a node times out
	RpcNodeTimeout int64 `json:"rpcNodeTimeout"`
	// The user agent used when sending RPC requests to nodes
	UserAgent string `json:"userAgent"`
	// Whether the JSON response from nodes should be sorted
	SortJsonResponse bool `json:"sortJsonResponse"`
	// Logging related configurtion
	Logging LoggingConfig `json:"logging"`
}

type LoggingConfig struct {
	// Enable console logging
	ConsoleOutputEnabled bool `json:"consoleOutputEnabled"`
	// Enable logging to a file
	FileOutputEnabled bool `json:"fileOutputEnabled"`
	// Directory to log to to when file output is enabled
	Directory string `json:"directory"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `json:"filename"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `json:"maxSize"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `json:"maxBackups"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int `json:"maxAge"`
}

func NewDefaultConfig() Config {
	c := Config{
		Version:          DefaultConfigVersion,
		Debug:            DefaultDebug,
		RpcHttpPort:      DefaultRpcHttpPort,
		RpcTimeout:       DefaultRpcTimeout,
		RpcNodeTimeout:   DefaultRpcNodeTimeout,
		UserAgent:        DefaultUserAgent,
		SortJsonResponse: DefaultSortJsonResponse,
		Logging: LoggingConfig{
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
