package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/runtimeutil"
)

var logger zerolog.Logger

func init() {
	Init(false)
}

func Init(debug bool) {
	// Default level is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// By default write to console at startup
	logger = zerolog.New(newConsoleWriter()).With().Timestamp().Logger()
}

func Configure(cfg *types.Config) {
	Init(cfg.Debug)

	var writers []io.Writer

	if cfg.Logging.ConsoleOutputEnabled {
		writers = append(writers, newConsoleWriter())
	}
	if cfg.Logging.FileOutputEnabled {
		writers = append(writers, newRollingFile(cfg))
	}
	mw := zerolog.MultiLevelWriter(writers...)

	logger = zerolog.New(mw).With().Timestamp().Logger()
}

func Logger() *zerolog.Logger {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return &logger
	}

	fmt.Printf("Level: %v\n", logger.GetLevel()) 

	// Append filename + line number and function name
	infos := runtimeutil.ToFuncInfos(pc)
	sublogger := logger.With().Str("file", infos.Filename).Int("line", infos.Line).Str("function", infos.Funcname).Logger()
	return &sublogger
}

func newConsoleWriter() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
}

func newRollingFile(config *types.Config) io.Writer {
	if err := os.MkdirAll(config.Logging.Directory, 0744); err != nil {
		logger.Error().Err(err).Str("path", config.Logging.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Logging.Directory, config.Logging.Filename),
		MaxBackups: config.Logging.MaxBackups,
		MaxSize:    config.Logging.MaxSize,
		MaxAge:     config.Logging.MaxAge,
	}
}
