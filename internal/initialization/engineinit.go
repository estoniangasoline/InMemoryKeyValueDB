package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/engine"

	"go.uber.org/zap"
)

const (
	inMemoryType = "in_memory"
)

func createEngine(config *config.EngineConfig, logger *zap.Logger) (engineLayer, error) {

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if config == nil {
		return engine.NewInMemoryEngine(logger)
	}

	var initEngine engineLayer
	var err error

	switch config.EngineType {
	case inMemoryType:
		initEngine, err = engine.NewInMemoryEngine(logger) // TODO: new engine type
	default:
		initEngine, err = engine.NewInMemoryEngine(logger)
	}

	return initEngine, err
}
