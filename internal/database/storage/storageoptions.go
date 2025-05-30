package storage

import "inmemorykvdb/internal/database/request"

type StorageOption func(*Storage)

func WithWal(wal WAL) StorageOption {
	return func(s *Storage) {
		s.wal = wal
	}
}

func WithReplica(repl Replica) StorageOption {
	return func(s *Storage) {
		s.replica = repl
	}
}

func WithDataChan(dataChan chan *request.Batch) StorageOption {
	return func(s *Storage) {
		s.dataChan = dataChan
	}
}
