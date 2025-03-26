package storage

import (
	"errors"
	"inmemorykvdb/internal/database/commands"

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
		err := s.Set(arg[0], arg[1])
		if err == nil {
			return okAnswer, nil
		}
		return "", err

	case commands.DelCommand:
		s.logger.Debug("started del command")
		err := s.Del(arg[0])
		if err == nil {
			return okAnswer, nil
		}
		return "", err

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
