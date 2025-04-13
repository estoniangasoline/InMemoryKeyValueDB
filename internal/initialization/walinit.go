package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/wal"

	"go.uber.org/zap"
)

func createWal(cnfg *config.WalConfig, logger *zap.Logger, writeLevel writingLayer, readLevel readingLayer) (WAL, error) {
	if logger == nil {
		return nil, errors.New("nil logger")
	}

	if writeLevel == nil {
		return nil, errors.New("nil write level")
	}

	if readLevel == nil {
		return nil, errors.New("nil read level")
	}

	if cnfg == nil {
		return wal.NewWal(writeLevel, readLevel, logger)
	}

	writeAheadLog, err := wal.NewWal(
		writeLevel,
		readLevel,
		logger,
		wal.WithBatchSize(cnfg.BatchSize),
		wal.WithBatchTimeout(cnfg.BatchTimeout))

	if err != nil {
		return nil, err
	}

	return writeAheadLog, nil
}
