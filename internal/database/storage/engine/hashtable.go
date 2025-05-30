package engine

import (
	"inmemorykvdb/pkg/concurrency"
	"sync"
)

type hashTable struct {
	pairs map[string]string
	mutex *sync.RWMutex
}

func NewHashTable(capacity int) *hashTable {
	return &hashTable{pairs: make(map[string]string, capacity), mutex: &sync.RWMutex{}}
}

func (h *hashTable) get(key string) (string, bool) {
	var res string
	var found bool

	concurrency.WithRLock(h.mutex, func() {
		res, found = h.pairs[key]
	})

	return res, found
}

func (h *hashTable) set(key, value string) {
	concurrency.WithLock(h.mutex, func() {
		h.pairs[key] = value
	})
}

func (h *hashTable) del(key string) {
	concurrency.WithLock(h.mutex, func() {
		delete(h.pairs, key)
	})
}
