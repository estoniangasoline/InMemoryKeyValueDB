package engine

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithPartitions(t *testing.T) {
	correctEngine := func(cap, count int) *InMemoryEngine {
		eng := &InMemoryEngine{}
		eng.partitions = make([]*hashTable, count)

		for i := range count {
			eng.partitions[i] = NewHashTable(cap)
		}

		return eng
	}

	tests := map[string]struct {
		cap   int
		count int

		expectedEng *InMemoryEngine
		expectedErr error
	}{
		"correct cap and count": {
			cap:   10,
			count: 10,

			expectedEng: correctEngine(10, 10),
			expectedErr: nil,
		},

		"incorrect cap": {
			cap:   0,
			count: 10,

			expectedEng: &InMemoryEngine{},
			expectedErr: errors.New("cap could not be equal or less than zero"),
		},

		"incorrect count": {
			cap:   10,
			count: 0,

			expectedEng: &InMemoryEngine{},
			expectedErr: errors.New("count could not be equal or less than zero"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			eng := &InMemoryEngine{}

			option := WithPartitions(test.count, test.cap)

			err := option(eng)

			assert.Equal(t, test.expectedEng, eng)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
