package types

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
