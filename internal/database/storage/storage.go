package storage

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"

	"go.uber.org/zap"
)

const (
	okAnswer = "SUCCESS"
	notFound = "NOT FOUND"
)

type engineLayer interface {
	SET(key string, value string)
	GET(key string) (string, bool)
	DEL(key string)
}

type Replica interface {
	IsMaster() bool
}

type WAL interface {
	Write(req request.Request)
	Read() *request.Batch
}

type Storage struct {
	wal     WAL
	engine  engineLayer
	logger  *zap.Logger
	replica Replica

	dataChan <-chan *request.Batch
}

func (s *Storage) HandleRequest(req request.Request) (string, error) {
	if s.wal != nil && req.RequestType != commands.GetCommand && (s.replica == nil || s.replica.IsMaster()) {
		s.wal.Write(req)
	}

	resp, err := s.requestToEngine(req, true)
	return resp, err
}

func (s *Storage) requestToEngine(req request.Request, fromClient bool) (string, error) {
	switch req.RequestType {

	case commands.GetCommand:
		s.logger.Debug("started get command")
		val, found := s.engine.GET(req.Args[0])

		if !found {
			return notFound, nil
		}

		return val, nil

	case commands.SetCommand:
		if s.isNotMutable(fromClient) {
			return "", errors.New("slave node is read-only")
		}

		s.logger.Debug("started set command")
		s.engine.SET(req.Args[0], req.Args[1])

		return okAnswer, nil

	case commands.DelCommand:
		if s.isNotMutable(fromClient) {
			return "", errors.New("slave node is read-only")
		}

		s.logger.Debug("started del command")
		s.engine.DEL(req.Args[0])

		return okAnswer, nil

	default:
		s.logger.Error("incorrect request type")
		return "", errors.New("incorrect request type")
	}
}

func (s *Storage) isNotMutable(fromClient bool) bool {
	return s.replica != nil && !s.replica.IsMaster() && fromClient
}

func NewStorage(logger *zap.Logger, engine engineLayer, options ...StorageOption) (*Storage, error) {

	if logger == nil {
		return nil, errors.New("could not create storage without logger")
	}

	if engine == nil {
		return nil, errors.New("could not create storage without engine")
	}

	storage := &Storage{logger: logger, engine: engine}

	for _, option := range options {
		option(storage)
	}

	if storage.replica != nil && !storage.replica.IsMaster() {
		if storage.dataChan == nil {
			return nil, errors.New("could not create slave node without data chan")
		}

		storage.synchronization()
	}

	if storage.wal != nil {
		recovered := storage.wal.Read()

		if recovered != nil {
			storage.recoverData(recovered)
		}
	}

	return storage, nil
}

func (s *Storage) recoverData(batch *request.Batch) {

	for _, req := range batch.Data {
		_, err := s.requestToEngine(*req, false)

		if err != nil {
			s.logger.Error(err.Error())
		}
	}
}

func (s *Storage) synchronization() {
	go func() {
		for data := range s.dataChan {
			s.recoverData(data)
		}
	}()
}
