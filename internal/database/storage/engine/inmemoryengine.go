package engine

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

type InMemoryEngine struct {
	logger    *zap.Logger
	mutex     sync.Mutex
	hashTable map[string]string
}

func NewInMemoryEngine(logger *zap.Logger) (*InMemoryEngine, error) {
	if logger == nil {
		return nil, errors.New("engine without logger")
	}

	engine := &InMemoryEngine{logger: logger, hashTable: make(map[string]string), mutex: sync.Mutex{}}

	return engine, nil
}

func (e *InMemoryEngine) GET(key string) (string, error) {

	value, ok := e.hashTable[key]

	if ok {
		e.logger.Debug("succesufull got value")
		return value, nil
	}

	e.logger.Debug("value not found")
	return "", errors.New("value not found")
}

func (e *InMemoryEngine) SET(key string, value string) error {

	e.mutex.Lock()
	e.hashTable[key] = value
	e.mutex.Unlock()

	e.logger.Debug("succesufull set value")
	return nil
}

func (e *InMemoryEngine) DEL(key string) error {

	delete(e.hashTable, key)

	e.logger.Debug("succesufull deleted value")

	return nil
}
