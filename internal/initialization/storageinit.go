package initialization

import (
	"errors"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage"

	"go.uber.org/zap"
)

func createStorage(engine engineLayer, wal WAL, logger *zap.Logger, replication replica) (*storage.Storage, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if engine == nil {
		return nil, errors.New("engine is nil")
	}

	var dataChan chan *request.Batch

	if replication != nil && !replication.IsMaster() {
		dataChan = replication.DataChan()
	}

	storage, err := storage.NewStorage(logger, engine,
		storage.WithWal(wal),
		storage.WithReplica(replication),
		storage.WithDataChan(dataChan))

	return storage, err
}
