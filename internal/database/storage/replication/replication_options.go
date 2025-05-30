package replication

import (
	"errors"
	"time"
)

type MasterOption func(*Master) error
type SlaveOption func(*Slave) error

func WithDirectoryMaster(directory string) MasterOption {
	return func(m *Master) error {
		m.directory = directory
		return nil
	}
}

func WithDirectorySlave(directory string) SlaveOption {
	return func(s *Slave) error {
		s.directory = directory
		return nil
	}
}

func WithInterval(interval time.Duration) SlaveOption {
	return func(s *Slave) error {
		if interval == 0 {
			return errors.New("time interval could not be a zero")
		}

		s.requestInterval = interval

		return nil
	}
}
