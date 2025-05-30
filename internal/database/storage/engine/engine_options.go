package engine

import "errors"

type EngineOption func(engine *InMemoryEngine) error

func WithPartitions(count int, cap int) EngineOption {
	return func(engine *InMemoryEngine) error {
		if count <= 0 {
			return errors.New("count could not be equal or less than zero")
		}

		if cap <= 0 {
			return errors.New("cap could not be equal or less than zero")
		}

		engine.partitions = make([]*hashTable, count)

		for i := range count {
			engine.partitions[i] = NewHashTable(cap)
		}

		return nil
	}
}
