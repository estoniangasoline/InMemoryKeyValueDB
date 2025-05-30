package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewInitalizer(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		cnfg *config.Config

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct initializer",

			cnfg: &config.Config{
				Engine: &config.EngineConfig{
					EngineType: "in_memory",
				},

				Network: &config.NetworkConfig{
					Address:        "127.0.0.1:3223",
					MaxConnections: 100,
					MaxMessageSize: "4KB",
					IdleTimeout:    time.Minute * 5,
				},

				Logging: &config.LoggingConfig{
					Level:  "info",
					Output: "C:/go/InMemoryKeyValueDB/test/log/logging.txt",
				},

				WalConfig: &config.WalConfig{
					BatchSize:      4096,
					BatchTimeout:   time.Second,
					MaxSegmentSize: "2MB",
					DataDirectory:  "C:/go/InMemoryKeyValueDB/test/init/",
					FileName:       "wal",
				},

				Replication: &config.ReplicaConfig{
					ReplicaType:   master,
					MasterAddress: ":8080",
				},
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "nil config",

			cnfg: nil,

			expectedNilObj: true,
			expectedErr:    errors.New("config is nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			initializer, err := NewInitializer(test.cnfg)

			if test.expectedNilObj {
				assert.Nil(t, initializer)
			} else {
				assert.NotNil(t, initializer)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
