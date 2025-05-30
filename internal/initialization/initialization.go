package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database"
	"inmemorykvdb/internal/database/compute"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage"
	"inmemorykvdb/internal/network"

	"go.uber.org/zap"
)

type engineLayer interface {
	SET(key string, value string)
	GET(key string) (string, bool)
	DEL(key string)
}

type WAL interface {
	Write(req request.Request)
	Read() *request.Batch
}

type readingLayer interface {
	Read() ([][]byte, error)
}

type writingLayer interface {
	Write([]byte) (int, error)
}

type replica interface {
	IsMaster() bool
	DataChan() chan *request.Batch
}

type Initializer struct {
	engine   engineLayer
	logger   *zap.Logger
	storage  *storage.Storage
	database *database.InMemoryKeyValueDatabase
	server   *network.Server
}

func NewInitializer(cnfg *config.Config) (*Initializer, error) {

	if cnfg == nil {
		return nil, errors.New("config is nil")
	}

	logger, err := createLogger(cnfg.Logging)

	if err != nil {
		return nil, err
	}

	engine, err := createEngine(cnfg.Engine, logger)

	if err != nil {
		return nil, err
	}

	wl, err := createWriteLevel(logger, cnfg.WalConfig, cnfg.Replication)

	if err != nil {
		return nil, err
	}

	rl, err := createReadLevel(logger, cnfg.WalConfig)

	if err != nil {
		return nil, err
	}

	writeAheadLog, err := createWal(cnfg.WalConfig, logger, wl, rl)

	if err != nil {
		return nil, err
	}

	repl, err := createReplica(logger, cnfg.Replication, cnfg.WalConfig)

	if err != nil {
		return nil, err
	}

	storage, err := createStorage(engine, writeAheadLog, logger, repl)

	if err != nil {
		return nil, err
	}

	compute, err := compute.NewCompute(logger)

	if err != nil {
		return nil, err
	}

	database, err := createDatabase(storage, compute, logger)

	if err != nil {
		return nil, err
	}

	server, err := createServer(cnfg.Network, logger)

	if err != nil {
		return nil, err
	}

	return &Initializer{engine: engine, logger: logger, storage: storage, database: database, server: server}, nil
}

func (i *Initializer) StartDatabase() {
	i.server.HandleConnections(func(request []byte) []byte {
		response, err := i.database.HandleRequest(string(request))
		if err != nil {
			return []byte(err.Error())
		}
		return []byte(response)
	})
}
