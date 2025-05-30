package engine

import (
	"errors"
	"fmt"
	"hash/fnv"

	"go.uber.org/zap"
)

const (
	baseCapacity = 100
)

type InMemoryEngine struct {
	Logger     *zap.Logger
	partitions []*hashTable
}

func NewInMemoryEngine(logger *zap.Logger, options ...EngineOption) (*InMemoryEngine, error) {
	if logger == nil {
		return nil, errors.New("engine without logger")
	}

	engine := &InMemoryEngine{Logger: logger}

	for _, option := range options {
		err := option(engine)

		if err != nil {
			return nil, err
		}
	}

	if len(engine.partitions) == 0 {
		engine.partitions = make([]*hashTable, 1)
		engine.partitions[0] = NewHashTable(baseCapacity)
	}

	return engine, nil
}

func (e *InMemoryEngine) GET(key string) (string, bool) {
	e.Logger.Debug(fmt.Sprintf("started get query for key: %s", key))

	value, found := e.partitions[e.makeTxId(key)].get(key)

	e.Logger.Debug("get query is done")

	return value, found
}

func (e *InMemoryEngine) SET(key string, value string) {
	e.Logger.Debug(fmt.Sprintf("started set query for key: %s; value: %s", key, value))

	e.partitions[e.makeTxId(key)].set(key, value)

	e.Logger.Debug("set query is done")
}

func (e *InMemoryEngine) DEL(key string) {
	e.Logger.Debug(fmt.Sprintf("started del query for key: %s", key))

	e.partitions[e.makeTxId(key)].del(key)

	e.Logger.Debug("del query is done")
}

func (e *InMemoryEngine) makeTxId(key string) int {
	hash := fnv.New32a()
	hash.Write([]byte(key))

	return int(hash.Sum32()) % len(e.partitions)
}
