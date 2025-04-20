package readlevel

import "errors"

type readLevelOptions func(rl *readLevel) error

func WithFileMaxSize(maxSize int) readLevelOptions {
	return func(rl *readLevel) error {
		if maxSize == 0 {
			return errors.New("max file size could not be a zero")
		}

		rl.fileMaxSize = maxSize
		return nil
	}
}
