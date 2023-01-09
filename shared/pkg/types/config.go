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

type ChainConfig struct {
	// Id is the identifier of the chain
	Id string `json:"id"`
	// Name is the name of the chain
	Name string `json:"name"`
	// OfficialNodeUrl is the HTTP URL of the official node
	OfficialNodeUrl string `json:"officialNodeUrl"`
}

type KafkaConfig struct {
	Server       KafkaServerConfig `json:"servers"`
	DefaultTopic KafkaTopicConfig  `json:"defaultTopic"`
	RpcTopic     KafkaTopicConfig  `json:"rpcTopic"`
}

type KafkaServerConfig struct {
	Servers          string `json:"servers"`
	CertLocation     string `json:"certLocation"`
	CertKeyLocation  string `json:"certKeyLocation"`
	CertPoolLocation string `josn:"certPoolLocation"`
}

type KafkaTopicConfig struct {
	GroupId string `json:"groupId"`
	Topic   string `json:"topic"`
}
