package initialization

import (
	"errors"
	"inmemorykvdb/internal/database"
	"inmemorykvdb/internal/database/compute"
	"inmemorykvdb/internal/database/storage"

	"go.uber.org/zap"
)

func createDatabase(stor *storage.Storage, comp *compute.Compute, logger *zap.Logger) (*database.InMemoryKeyValueDatabase, error) {
	if stor == nil {
		return nil, errors.New("storage is nil")
	}

	if comp == nil {
		return nil, errors.New("compute is nil")
	}

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	db, err := database.NewInMemoryKvDb(comp, stor, logger)

	return db, err
}
