package initialization

import (
	"errors"
	"inmemorykvdb/internal/database/storage/engine"
	"inmemorykvdb/internal/database/storage/replication"
	"inmemorykvdb/internal/database/storage/wal"
	"inmemorykvdb/internal/database/storage/wal/readlevel"
	"inmemorykvdb/internal/database/storage/wal/writelevel"
	"inmemorykvdb/internal/network"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createStorage(t *testing.T) {

	t.Parallel()

	dir := "C:/go/InMemoryKeyValueDB/tests/init/storage/"

	type testCase struct {
		name string

		nilEngine  bool
		nilWal     bool
		nilReplica bool
		logger     *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct storage",

			nilEngine:  false,
			nilWal:     false,
			nilReplica: false,
			logger:     zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "storage without engine",

			nilEngine:  true,
			nilWal:     false,
			nilReplica: false,
			logger:     zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("engine is nil"),
		},

		{
			name: "storage without wal",

			nilEngine:  false,
			nilWal:     true,
			nilReplica: false,
			logger:     zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "storage without logger",

			nilEngine:  false,
			nilWal:     false,
			nilReplica: false,
			logger:     nil,

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},

		{
			name: "storage without replica",

			nilEngine:  false,
			nilWal:     false,
			nilReplica: true,
			logger:     zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			var eng engineLayer

			if !test.nilEngine {
				eng, _ = engine.NewInMemoryEngine(zap.NewNop())
			}

			var writeAheadLog WAL

			if !test.nilWal {
				wl, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFilePath(dir))
				rl, _ := readlevel.NewReadLevel(zap.NewNop(), defaultPattern, readlevel.WithDirectory(dir))
				writeAheadLog, _ = wal.NewWal(zap.NewNop(), wal.WithReader(rl), wal.WithWriter(wl))
			}

			var slave replica

			if !test.nilReplica {
				client, _ := network.NewClient(":8080")
				slave, _ = replication.NewSlave(client, zap.NewNop())
			}

			stor, err := createStorage(eng, writeAheadLog, test.logger, slave)

			if test.expectedNilObj {
				assert.Nil(t, stor)
			} else {
				assert.NotNil(t, stor)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
