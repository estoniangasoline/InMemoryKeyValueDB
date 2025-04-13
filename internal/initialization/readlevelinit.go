package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/wal/readlevel"
	"inmemorykvdb/pkg/parsing"

	"go.uber.org/zap"
)

const (
	defaultPattern = "write_ahead"
)

func createReadLevel(logger *zap.Logger, cnfg *config.WalConfig) (readingLayer, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	if cnfg == nil {
		return readlevel.NewReadLevel(logger, defaultPattern)
	}

	maxSegSize := defaultMaxSegSize * defaultMultiply

	if len(cnfg.MaxSegmentSize) > 2 {

		probablyMaxSegSize, err := parsing.ParseSize(cnfg.MaxSegmentSize)

		if err != nil {
			maxSegSize = probablyMaxSegSize
		}
	}

	rl, err := readlevel.NewReadLevel(logger, cnfg.FileName, readlevel.WithFileMaxSize(maxSegSize))

	if err != nil {
		return nil, err
	}

	return rl, nil
}
