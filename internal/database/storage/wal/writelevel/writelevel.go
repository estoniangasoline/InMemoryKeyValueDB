package writelevel

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"
)

const (
	defaultFirstFileIndex = 1
	defaultFileName       = "write_ahead"
	fileExtension         = ".log"
	defaultFileMaxSize    = 4096
)

type writeLevel struct {
	LastFileName string

	filePath    string
	fileName    string
	fileMaxSize int

	nextFileIndex int

	logger *zap.Logger
}

func NewWriteLevel(logger *zap.Logger, options ...writeLevelOptions) (*writeLevel, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	wl := &writeLevel{logger: logger}
	for _, option := range options {
		err := option(wl)

		if err != nil {
			return nil, err
		}
	}

	if wl.fileMaxSize == 0 {
		wl.fileMaxSize = defaultFileMaxSize
	}

	if wl.fileName == "" {
		wl.fileName = defaultFileName
	}

	wl.nextFileIndex = defaultFirstFileIndex

	return wl, nil
}

func (wl *writeLevel) Write(data []byte) (int, error) {
	wl.logger.Debug("started write data")
	if len(data) == 0 {
		wl.logger.Error("small data size")
		return 0, errors.New("data is empty")
	}

	file, err := wl.createFile()

	if err != nil {
		wl.logger.Error(fmt.Sprintf("unable write to file with error: %s", err.Error()))
		return 0, err
	}

	defer file.Close()

	wl.logger.Debug("started write to file")

	count, err := file.Write(data)
	file.Sync()

	if err != nil {
		wl.logger.Error(fmt.Sprintf("writing file done with error: %s", err.Error()))
		return count, err
	}

	wl.logger.Debug("writing is success")

	return count, nil
}

func (wl *writeLevel) checkFileIsExist() {
	wl.logger.Debug("check existing of file")
	var stringInd = strconv.Itoa(wl.nextFileIndex)
	_, err := os.Stat(wl.filePath + wl.fileName + stringInd + fileExtension)

	for !errors.Is(err, os.ErrNotExist) {
		wl.nextFileIndex++
		var stringInd string = strconv.Itoa(wl.nextFileIndex)
		_, err = os.Stat(wl.filePath + wl.fileName + stringInd + fileExtension)
	}

	wl.logger.Debug("found free file name")
}

func (wl *writeLevel) createFile() (*os.File, error) {

	wl.checkFileIsExist()

	wl.logger.Debug("creating new file")
	var stringInd = strconv.Itoa(wl.nextFileIndex)

	probablyLastFileName := wl.filePath + wl.fileName + stringInd + fileExtension
	file, err := os.Create(probablyLastFileName)

	if err != nil {
		wl.logger.Error(err.Error())
		return nil, errors.New("could not create the file")
	}

	wl.LastFileName = probablyLastFileName

	wl.nextFileIndex++

	wl.logger.Debug("creating of new file is success")

	return file, nil
}
