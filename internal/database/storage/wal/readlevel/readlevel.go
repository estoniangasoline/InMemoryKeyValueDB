package readlevel

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	defaultFileMaxSize = 4096
)

type readLevel struct {
	logger      *zap.Logger
	pattern     string
	fileMaxSize int
	fileNames   []string
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

	err := rl.findFiles()

	if err != nil {
		return nil, err
	}

	if rl.fileMaxSize == 0 {
		rl.fileMaxSize = defaultFileMaxSize
	}

	return rl, nil
}

func (rl *readLevel) findFiles() error {
	fileNames, err := filepath.Glob(rl.pattern + "*")

	if err != nil {
		return errors.New("incorrect pattern to find the files")
	}

	rl.fileNames = fileNames

	return nil
}

func (rl *readLevel) Read() (*[][]byte, error) {
	rl.logger.Debug("started reading files")
	readData := make([][]byte, 0, len(rl.fileNames))

	errorStr := "could not to read the files: "
	var hasUnreadFiles bool

	for _, name := range rl.fileNames {
		rl.logger.Debug(fmt.Sprintf("reading file %s", name))
		fl, err := os.Open(name)

		if err != nil {
			rl.logger.Error(fmt.Sprintf("file %s could not be open with error: %s", name, err.Error()))

			hasUnreadFiles = true
			errorStr += name + " "

			fl.Close()
			continue
		}

		buf := make([]byte, rl.fileMaxSize)

		n, err := fl.Read(buf)
		fl.Close()

		if !errors.Is(err, io.EOF) && (err != nil || n == rl.fileMaxSize) {

			rl.logReadErr(err, n, name)

			hasUnreadFiles = true
			errorStr += name + " "
			continue
		}

		rl.logger.Debug(fmt.Sprintf("reading file %s is success", name))

		readData = append(readData, buf[:n])
	}

	rl.logger.Debug("reading is complete")

	if hasUnreadFiles {
		return &readData, errors.New(errorStr)
	}

	return &readData, nil
}

func (rl *readLevel) logReadErr(err error, readBytes int, name string) {
	var errMsg string
	if err == nil && readBytes == rl.fileMaxSize {
		errMsg = "reading out of range"
	} else {
		errMsg = err.Error()
	}

	rl.logger.Debug(fmt.Sprintf("reading file %s is done with err %s", name, errMsg))
}
