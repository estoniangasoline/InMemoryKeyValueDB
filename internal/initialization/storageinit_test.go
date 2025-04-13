package initialization

import (
	"errors"
	"inmemorykvdb/internal/database/storage/engine"
	"inmemorykvdb/internal/database/storage/wal"
	"inmemorykvdb/internal/database/storage/wal/readlevel"
	"inmemorykvdb/internal/database/storage/wal/writelevel"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createStorage(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		nilEngine bool
		nilWal    bool
		logger    *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct storage",

			nilEngine: false,
			nilWal:    false,
			logger:    zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "storage without engine",

			nilEngine: true,
			nilWal:    false,
			logger:    zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("engine is nil"),
		},

		{
			name: "storage without wal",

			nilEngine: false,
			nilWal:    true,
			logger:    zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "storage without logger",

			nilEngine: false,
			nilWal:    false,
			logger:    nil,

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var eng engineLayer

			if !test.nilEngine {
				eng, _ = engine.NewInMemoryEngine(zap.NewNop())
			}

			var writeAheadLog WAL

			if !test.nilWal {
				wl, _ := writelevel.NewWriteLevel(zap.NewNop())
				rl, _ := readlevel.NewReadLevel(zap.NewNop(), defaultPattern)
				writeAheadLog, _ = wal.NewWal(wl, rl, zap.NewNop())
			}

			stor, err := createStorage(eng, writeAheadLog, test.logger)

			if test.expectedNilObj {
				assert.Nil(t, stor)
			} else {
				assert.NotNil(t, stor)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
