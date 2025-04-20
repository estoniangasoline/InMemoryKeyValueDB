package storage

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"

	"go.uber.org/zap"
)

const (
	okAnswer = "SUCCESS"
)

type engineLayer interface {
	SET(key string, value string) error
	GET(key string) (string, error)
	DEL(key string) error
}

type WAL interface {
	StartWAL()
	Write(req request.Request)
	Read() *request.Batch
}

type Storage struct {
	wal    WAL
	engine engineLayer
	logger *zap.Logger
}

func (s *Storage) HandleRequest(req request.Request) (string, error) {
	if s.wal != nil && req.RequestType != commands.GetCommand {
		s.wal.Write(req)
	}

	resp, err := s.requestToEngine(req)
	return resp, err
}

func (s *Storage) requestToEngine(req request.Request) (string, error) {
	switch req.RequestType {

	case commands.GetCommand:
		s.logger.Debug("started get command")
		return s.engine.GET(req.Args[0])

	case commands.SetCommand:
		s.logger.Debug("started set command")
		err := s.engine.SET(req.Args[0], req.Args[1])
		if err == nil {
			return okAnswer, nil
		}
		return "", err

	case commands.DelCommand:
		s.logger.Debug("started del command")
		err := s.engine.DEL(req.Args[0])
		if err == nil {
			return okAnswer, nil
		}
		return "", err

	default:
		s.logger.Error("uncorrect request type")
		return "", errors.New("uncorrect request type")
	}
}

func NewStorage(logger *zap.Logger, options ...StorageOption) (*Storage, error) {

	if logger == nil {
		return nil, errors.New("could not create storage without logger")
	}

	storage := &Storage{logger: logger}

	for _, option := range options {
		option(storage)
	}

	if storage.wal != nil {
		recovered := storage.wal.Read()
		storage.recoverData(recovered)

		storage.wal.StartWAL()
	}

	return storage, nil
}

func (s *Storage) recoverData(batch *request.Batch) {

	for _, req := range batch.Data {
		_, err := s.requestToEngine(*req)

		if err != nil {
			s.logger.Error(err.Error())
		}
	}
}
