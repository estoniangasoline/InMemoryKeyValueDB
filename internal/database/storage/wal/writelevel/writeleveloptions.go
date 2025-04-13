package writelevel

import "errors"

type writeLevelOptions func(*writeLevel) error

func WithFileName(fileName string) writeLevelOptions {
	return func(wl *writeLevel) error {
		if fileName == "" {
			return errors.New("file name could not be a empty string")
		}

		wl.fileName = fileName
		return nil
	}
}

func WithFilePath(path string) writeLevelOptions {
	return func(wl *writeLevel) error {
		wl.filePath = path
		return nil
	}
}

func WithFileMaxSize(maxSize int) writeLevelOptions {
	return func(wl *writeLevel) error {
		if maxSize == 0 {
			return errors.New("max file size could not be a zero")
		}

		wl.fileMaxSize = maxSize
		return nil
	}
}
