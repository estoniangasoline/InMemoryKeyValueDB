package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/wal/readlevel"
	"inmemorykvdb/internal/database/storage/wal/writelevel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createWal(t *testing.T) {
	type testCase struct {
		name string

		nilReadLevel  bool
		nilWriteLevel bool
		logger        *zap.Logger
		cnfg          *config.WalConfig

		expectedNilObj bool
		expectedErr    error
	}

	dataDir := "C:/go/InMemoryKeyValueDB/test/init/wal/"

	testCases := []testCase{
		{
			name: "nil logger",

			nilReadLevel:  false,
			nilWriteLevel: false,
			logger:        nil,
			cnfg: &config.WalConfig{
				BatchSize:      100,
				BatchTimeout:   time.Millisecond,
				MaxSegmentSize: "10MB",
				DataDirectory:  dataDir,
				FileName:       "wrahlo",
			},

			expectedNilObj: true,
			expectedErr:    errors.New("nil logger"),
		},

		{
			name: "default wal",

			nilReadLevel:  true,
			nilWriteLevel: true,
			logger:        zap.NewNop(),
			cnfg:          nil,

			expectedNilObj: true,
			expectedErr:    nil,
		},

		{
			name: "wal with options",

			nilReadLevel:  false,
			nilWriteLevel: false,
			logger:        zap.NewNop(),
			cnfg: &config.WalConfig{
				BatchSize:      100,
				BatchTimeout:   time.Millisecond,
				MaxSegmentSize: "10MB",
				DataDirectory:  dataDir,
				FileName:       "wrahlo",
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			var wl writingLayer

			if !test.nilWriteLevel {
				wl, _ = writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFilePath(test.cnfg.DataDirectory))
			}

			var rl readingLayer

			if !test.nilReadLevel {
				rl, _ = readlevel.NewReadLevel(zap.NewNop(), defaultPattern, readlevel.WithDirectory(test.cnfg.DataDirectory))
			}

			testwal, err := createWal(test.cnfg, test.logger, wl, rl)

			if test.expectedNilObj {
				assert.Nil(t, testwal)
			} else {
				assert.NotNil(t, testwal)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
