package initialization

import (
	"inmemorykvdb/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"

	defaultLevel    = zapcore.InfoLevel
	defaultEncoding = "json"
	defaultOutput   = "C:/go/InMemoryKeyValueDB/test/log/pretty.log"
)

func createLogger(config *config.LoggingConfig) (*zap.Logger, error) {

	if config == nil {
		return zap.Config{
			Level:       zap.NewAtomicLevelAt(defaultLevel),
			OutputPaths: []string{defaultOutput},
			Encoding:    defaultEncoding}.Build()
	}

	output := defaultOutput

	if config.Output != "" {
		output = config.Output
	}

	var level zapcore.Level

	switch config.Level {

	case debugLevel:
		level = zapcore.DebugLevel

	case infoLevel:
		level = zapcore.InfoLevel

	case warnLevel:
		level = zapcore.WarnLevel

	case errorLevel:
		level = zapcore.ErrorLevel

	default:
		level = defaultLevel
	}

	loggerCnfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{output},
		Encoding:    defaultEncoding,
	}

	return loggerCnfg.Build()
}
