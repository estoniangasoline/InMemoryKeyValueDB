package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewHashTable(t *testing.T) {
	cap := 10
	ht := NewHashTable(cap)

	require.NotNil(t, ht)

	assert.NotNil(t, ht.pairs)
	assert.NotNil(t, ht.mutex)
}

func Test_set(t *testing.T) {
	cap := 10
	ht := NewHashTable(cap)

	tests := map[string]struct {
		key   string
		value string
	}{
		"new key": {
			key:   "key1",
			value: "value1",
		},

		"existing key": {
			key:   "key1",
			value: "value2",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ht.set(test.key, test.value)

			val := ht.pairs[test.key]

			assert.Equal(t, test.value, val)
		})
	}
}

func Test_del(t *testing.T) {
	cap := 10
	ht := NewHashTable(cap)

	tests := map[string]struct {
		keyToDel string

		keyToSet string
		valToSet string
	}{
		"del existing key": {
			keyToDel: "key1",
			keyToSet: "key1",

			valToSet: "val1",
		},

		"del not existing key": {
			keyToDel: "key2",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ht.set(test.keyToSet, test.valToSet)
			ht.del(test.keyToDel)

			assert.Equal(t, "", ht.pairs[test.keyToDel])
		})
	}
}

func Test_get(t *testing.T) {
	cap := 10
	ht := NewHashTable(cap)

	tests := map[string]struct {
		keyToGet    string
		expectedVal string
		isFound     bool

		keyToSet string
		valToSet string
	}{
		"get existing key": {
			keyToGet:    "key1",
			expectedVal: "val1",
			isFound:     true,

			keyToSet: "key1",
			valToSet: "val1",
		},

		"gel not existing key": {
			keyToGet:    "key2",
			expectedVal: "",
			isFound:     false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ht.set(test.keyToSet, test.valToSet)

			val, found := ht.get(test.keyToGet)

			assert.Equal(t, test.expectedVal, val)
			assert.Equal(t, test.isFound, found)
		})
	}
}
