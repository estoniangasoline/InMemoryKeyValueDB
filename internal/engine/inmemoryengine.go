package engine

import (
	"errors"
	"fmt"
	"hash/fnv"

	"go.uber.org/zap"
)

const (
	maxEngineSize = 100000
)

type InMemoryEngine struct {
	logger     zap.Logger
	hashTables []*HashTable
}

func NewInMemoryEngine(logger *zap.Logger, size uint) (Engine, error) {
	if logger == nil {
		return nil, errors.New("engine without logger")
	}

	if size > maxEngineSize {
		return nil, fmt.Errorf("engine could not be bigger than %d elements", maxEngineSize)
	}

	engine := &InMemoryEngine{logger: *logger, hashTables: make([]*HashTable, size)}

	for i := 0; i < int(size); i++ {
		engine.hashTables[i] = NewHashTable()
	}

	return engine, nil
}

func (e *InMemoryEngine) GET(key string) (string, error) {

	if len(e.hashTables) < 1 {
		e.logger.Error("could not to get value from empty engine")
		return "", errors.New("engine is empty")
	}

	id := e.getHashTableId(key)

	value, ok := e.hashTables[id].Get(key)

	if ok {
		e.logger.Debug("succesufull got value")
		return value, nil
	}

	e.logger.Debug("value not found")
	return "", errors.New("value not found")
}

func (e *InMemoryEngine) SET(key string, value string) error {
	if len(e.hashTables) < 1 {
		e.logger.Error("could not set value for empty engine")

		return errors.New("engine is empty")
	}

	id := e.getHashTableId(key)
	e.hashTables[id].Set(key, value)

	e.logger.Debug("succesufull set value")
	return nil
}

func (e *InMemoryEngine) DEL(key string) error {
	if len(e.hashTables) < 1 {
		e.logger.Debug("could not delete key in empty engine")
		return errors.New("engine is empty")
	}

	id := e.getHashTableId(key)
	e.hashTables[id].Delete(key)

	e.logger.Debug("succesufull deleted value")

	return nil
}

func (e *InMemoryEngine) getHashTableId(key string) int {
	hash := fnv.New32()
	hash.Write([]byte(key))
	return int(hash.Sum32()) % len(e.hashTables)
}
