package logger

import (
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/vitelabs/vite-portal/shared/pkg/types"
	"github.com/vitelabs/vite-portal/shared/pkg/util/runtimeutil"
)

var (
	logger zerolog.Logger
	debugEnabled bool
)

func init() {
	Init(false)
}

func Init(debug bool) {
	debugEnabled = debug
	// Default level is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// By default write to console at startup
	logger = zerolog.New(newConsoleWriter()).With().Timestamp().Logger()
}

func Configure(debug bool, cfg types.LoggingConfig) {
	Init(debug)

	var writers []io.Writer

	if cfg.ConsoleOutputEnabled {
		writers = append(writers, newConsoleWriter())
	}
	if cfg.FileOutputEnabled {
		writers = append(writers, newRollingFile(cfg))
	}
	mw := zerolog.MultiLevelWriter(writers...)

	logger = zerolog.New(mw).With().Timestamp().Logger()
}

func DebugEnabled() bool {
	return debugEnabled
}

func Logger() *zerolog.Logger {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return &logger
	}

	// Append filename + line number and function name
	infos := runtimeutil.ToFuncInfos(pc)
	sublogger := logger.With().Str("file", infos.Filename).Int("line", infos.Line).Str("function", infos.Funcname).Logger()
	return &sublogger
}

func newConsoleWriter() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
}

func newRollingFile(config types.LoggingConfig) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		logger.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
	}
}
