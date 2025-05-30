package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createWriteLevel(t *testing.T) {
	type testCase struct {
		name string

		logger   *zap.Logger
		cnfg     *config.WalConfig
		replCnfg *config.ReplicaConfig

		expectedNilObj bool
		expectedErr    error
	}

	dataDir := "C:/go/InMemoryKeyValueDB/tests/init/wal/"

	testCases := []testCase{
		{
			name: "nil config",

			logger:   zap.NewNop(),
			cnfg:     nil,
			replCnfg: nil,

			expectedNilObj: true,
			expectedErr:    nil,
		},

		{
			name: "nil logger",

			logger:   nil,
			cnfg:     nil,
			replCnfg: nil,

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},

		{
			name: "not nil config",

			logger: zap.NewNop(),
			cnfg: &config.WalConfig{
				BatchSize:      100,
				BatchTimeout:   time.Millisecond,
				MaxSegmentSize: "10MB",
				DataDirectory:  dataDir,
				FileName:       "wrahlo",
			},
			replCnfg: nil,

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name:   "slave node",
			logger: zap.NewNop(),
			cnfg: &config.WalConfig{
				BatchSize:      100,
				BatchTimeout:   time.Millisecond,
				MaxSegmentSize: "10MB",
				DataDirectory:  dataDir,
				FileName:       "wrahlo",
			},
			replCnfg: &config.ReplicaConfig{
				ReplicaType: slave,
			},

			expectedNilObj: true,
			expectedErr:    nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			rl, err := createWriteLevel(test.logger, test.cnfg, test.replCnfg)

			if test.expectedNilObj {
				assert.Nil(t, rl)
			} else {
				assert.NotNil(t, rl)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
