package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/wal"

	"go.uber.org/zap"
)

func createWal(cnfg *config.WalConfig, logger *zap.Logger, writeLevel writingLayer, readLevel readingLayer) (WAL, error) {
	if cnfg == nil {
		return nil, nil
	}

	if logger == nil {
		return nil, errors.New("nil logger")
	}

	return wal.NewWal(
		logger,
		wal.WithBatchSize(cnfg.BatchSize),
		wal.WithBatchTimeout(cnfg.BatchTimeout),
		wal.WithWriter(writeLevel),
		wal.WithReader(readLevel),
	)
}
