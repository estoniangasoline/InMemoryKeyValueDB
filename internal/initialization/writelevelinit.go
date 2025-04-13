package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/wal/writelevel"
	"inmemorykvdb/pkg/parsing"

	"go.uber.org/zap"
)

const (
	defaultMaxSegSize = 1000
)

func createWriteLevel(logger *zap.Logger, cnfg *config.WalConfig) (writingLayer, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if cnfg == nil {
		return writelevel.NewWriteLevel(logger)
	}

	maxSegSize := defaultMaxSegSize * defaultMultiply

	if len(cnfg.MaxSegmentSize) > 2 {

		probablyMaxSegSize, err := parsing.ParseSize(cnfg.MaxSegmentSize)

		if err != nil {
			maxSegSize = probablyMaxSegSize
		}
	}

	wl, err := writelevel.NewWriteLevel(
		logger,
		writelevel.WithFileMaxSize(maxSegSize),
		writelevel.WithFileName(cnfg.FileName),
		writelevel.WithFilePath(cnfg.DataDirectory),
	)

	if err != nil {
		return nil, err
	}

	return wl, nil
}
