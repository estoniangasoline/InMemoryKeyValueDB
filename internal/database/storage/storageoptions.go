package storage

type StorageOption func(*Storage)

func WithEngine(engine engineLayer) StorageOption {
	return func(s *Storage) {
		s.engine = engine
	}
}

func WithWal(wal WAL) StorageOption {
	return func(s *Storage) {
		s.wal = wal
	}
}
