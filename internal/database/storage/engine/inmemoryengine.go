package engine

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

type InMemoryEngine struct {
	logger    *zap.Logger
	mutex     sync.RWMutex
	hashTable map[string]string
}

func NewInMemoryEngine(logger *zap.Logger) (*InMemoryEngine, error) {
	if logger == nil {
		return nil, errors.New("engine without logger")
	}

	engine := &InMemoryEngine{logger: logger, hashTable: make(map[string]string), mutex: sync.RWMutex{}}

	return engine, nil
}

func (e *InMemoryEngine) GET(key string) (string, error) {

	e.mutex.RLock()
	value, ok := e.hashTable[key]
	e.mutex.RUnlock()

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

	e.mutex.Lock()
	delete(e.hashTable, key)
	e.mutex.Unlock()

	e.logger.Debug("succesufull deleted value")

	return nil
}
