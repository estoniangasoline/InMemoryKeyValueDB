package readlevel

import (
	"errors"
	"inmemorykvdb/internal/database/storage/filesystem"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	defaultFileMaxSize = 4096
	defaultDir         = "C:/go/InMemoryKeyValueDB/test/readlevel/"
)

type readLevel struct {
	logger      *zap.Logger
	directory   string
	pattern     string
	fileMaxSize int
}

func NewReadLevel(logger *zap.Logger, pattern string, options ...readLevelOptions) (*readLevel, error) {

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	rl := &readLevel{logger: logger, pattern: pattern}

	for _, option := range options {
		err := option(rl)

		if err != nil {
			return nil, err
		}
	}

	if rl.directory == "" {
		rl.directory = defaultDir
	}

	if rl.fileMaxSize == 0 {
		rl.fileMaxSize = defaultFileMaxSize
	}

	return rl, nil
}

func (rl *readLevel) findFiles() ([]string, error) {
	fileNames, err := filepath.Glob(rl.directory + rl.pattern + "*")

	if err != nil {
		return nil, errors.New("incorrect pattern to find the files")
	}

	return fileNames, nil
}

func (rl *readLevel) Read() ([][]byte, error) {
	rl.logger.Debug("started reading files")

	names, err := rl.findFiles()

	if err != nil {
		return nil, err
	}

	files := make([][]byte, 0, len(names))

	files, err = filesystem.ReadAll("", names, files)

	if err != nil {
		rl.logger.Error(err.Error())
	}

	return files, nil
}
