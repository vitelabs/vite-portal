package types

const (
	DefaultDebug                       = false
	DefaultRpcHttpPort                 = 56331
	DefaultRpcTimeout                  = 10000
	DefaultLoggingConsoleOutputEnabled = true
	DefaultLoggingFileOutputEnabled    = true
	DefaultLoggingDirectory            = "logs"
	DefaultLoggingFilename             = "relayer.log"
	DefaultLoggingMaxSize              = 100
	DefaultLoggingMaxBackups           = 10
	DefaultLoggingMaxAge               = 28
)

type Config struct {
	// Enable debug mode
	Debug bool `json:"debug"`
	// Port number for the HTTP RPC
	RpcHttpPort int32 `json:"rpcHttpPort"`
	// The time in milliseconds before a RPC request times out
	RpcTimeout int64 `json:"rpcTimeout"`
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
		Debug:       DefaultDebug,
		RpcHttpPort: DefaultRpcHttpPort,
		RpcTimeout:  DefaultRpcTimeout,
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