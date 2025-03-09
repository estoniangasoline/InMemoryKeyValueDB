package engine

type HashTable struct {
	hashTable map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{hashTable: make(map[string]string)}
}

func (h *HashTable) Set(key string, value string) {
	h.hashTable[key] = value
}

func (h *HashTable) Get(key string) (string, bool) {
	value, ok := h.hashTable[key]

	return value, ok
}

func (h *HashTable) Delete(key string) {
	delete(h.hashTable, key)
}
