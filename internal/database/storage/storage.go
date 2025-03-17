package storage

import (
	"errors"
	"inmemorykvdb/internal/database/commands"

	"go.uber.org/zap"
)

type engineLayer interface {
	SET(key string, value string) error
	GET(key string) (string, error)
	DEL(key string) error
}

type Storage struct {
	Engine engineLayer
	logger *zap.Logger
}

func (s *Storage) HandleRequest(requestType int, arg ...string) (string, error) {

	switch requestType {

	case commands.GetCommand:
		s.logger.Debug("started get command")
		return s.Get(arg[0])

	case commands.SetCommand:
		s.logger.Debug("started set command")
		return "", s.Set(arg[0], arg[1])

	case commands.DelCommand:
		s.logger.Debug("started del command")
		return "", s.Del(arg[0])

	default:
		s.logger.Error("uncorrect request type")
		return "", errors.New("uncorrect request type")
	}
}

func (s *Storage) Set(key string, value string) error {
	return s.Engine.SET(key, value)
}

func (s *Storage) Get(key string) (string, error) {
	return s.Engine.GET(key)
}

func (s *Storage) Del(key string) error {
	return s.Engine.DEL(key)
}

func NewStorage(engine engineLayer, logger *zap.Logger) (*Storage, error) {

	if engine == nil && logger == nil {
		return nil, errors.New("could not create storage without engine and logger")
	}

	if engine == nil {
		return nil, errors.New("could not create storage without engine")
	}

	if logger == nil {
		return nil, errors.New("could not create storage without logger")
	}

	return &Storage{engine, logger}, nil
}
