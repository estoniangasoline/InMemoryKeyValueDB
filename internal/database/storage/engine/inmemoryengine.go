package engine

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

type InMemoryEngine struct {
	Logger    *zap.Logger
	mutex     sync.RWMutex
	hashTable map[string]string
}

func NewInMemoryEngine(logger *zap.Logger) (*InMemoryEngine, error) {
	if logger == nil {
		return nil, errors.New("engine without logger")
	}

	engine := &InMemoryEngine{Logger: logger, hashTable: make(map[string]string), mutex: sync.RWMutex{}}

	return engine, nil
}

func (e *InMemoryEngine) GET(key string) (string, error) {

	e.mutex.RLock()
	value, ok := e.hashTable[key]
	e.mutex.RUnlock()

	if ok {
		e.Logger.Debug("succesufull got value")
		return value, nil
	}

	e.Logger.Debug("value not found")
	return "", errors.New("value not found")
}

func (e *InMemoryEngine) SET(key string, value string) error {

	e.mutex.Lock()
	e.hashTable[key] = value
	e.mutex.Unlock()

	e.Logger.Debug("succesufull set value")
	return nil
}

func (e *InMemoryEngine) DEL(key string) error {

	e.mutex.Lock()
	delete(e.hashTable, key)
	e.mutex.Unlock()

	e.Logger.Debug("succesufull deleted value")

	return nil
}
