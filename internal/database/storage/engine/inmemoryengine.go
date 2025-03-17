package engine

import (
	"errors"

	"go.uber.org/zap"
)

type InMemoryEngine struct {
	logger    zap.Logger
	hashTable map[string]string
}

func NewInMemoryEngine(logger *zap.Logger) (*InMemoryEngine, error) {
	if logger == nil {
		return nil, errors.New("engine without logger")
	}

	engine := &InMemoryEngine{logger: *logger, hashTable: make(map[string]string)}

	return engine, nil
}

func (e *InMemoryEngine) GET(key string) (string, error) {

	if len(e.hashTable) < 1 {
		e.logger.Error("could not to get value from empty engine")
		return "", errors.New("engine is empty")
	}

	value, ok := e.hashTable[key]

	if ok {
		e.logger.Debug("succesufull got value")
		return value, nil
	}

	e.logger.Debug("value not found")
	return "", errors.New("value not found")
}

func (e *InMemoryEngine) SET(key string, value string) error {

	e.hashTable[key] = value

	e.logger.Debug("succesufull set value")
	return nil
}

func (e *InMemoryEngine) DEL(key string) error {
	if len(e.hashTable) < 1 {
		e.logger.Debug("could not delete key in empty engine")
		return errors.New("engine is empty")
	}

	delete(e.hashTable, key)

	e.logger.Debug("succesufull deleted value")

	return nil
}
