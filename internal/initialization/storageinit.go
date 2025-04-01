package initialization

import (
	"errors"
	"inmemorykvdb/internal/database/storage"

	"go.uber.org/zap"
)

func createStorage(engine engineLayer, logger *zap.Logger) (*storage.Storage, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if engine == nil {
		return nil, errors.New("engine is nil")
	}

	storage, err := storage.NewStorage(engine, logger)

	return storage, err
}
