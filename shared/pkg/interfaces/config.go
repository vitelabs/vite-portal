package interfaces

import "github.com/vitelabs/vite-portal/shared/pkg/types"

type ConfigI interface {
	GetVersion() string
	GetDebug() bool
	SetDebug(debug bool)
	GetLoggingConfig() types.LoggingConfig
	Validate() error
}