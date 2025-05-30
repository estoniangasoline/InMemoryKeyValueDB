package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createReplication(t *testing.T) {
	listener, _ := net.Listen("tcp", ":8080")

	defer listener.Close()

	go func() {
		for {
			listener.Accept()
		}
	}()

	type testCase struct {
		name string

		logger *zap.Logger
		cnfg   *config.ReplicaConfig

		expectedNilObj bool
		expectedErr    error
	}

	walCnfg := &config.WalConfig{
		DataDirectory: "C:/go/InMemoryKeyValueDB/test/init/",
	}

	testCases := []testCase{
		{
			name: "correct slave replica",

			logger: zap.NewNop(),
			cnfg: &config.ReplicaConfig{
				ReplicaType:   slave,
				MasterAddress: "localhost:8080",
				SyncInterval:  time.Second,
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "replica without logger",

			logger: nil,
			cnfg: &config.ReplicaConfig{
				ReplicaType:   slave,
				MasterAddress: "localhost:8080",
				SyncInterval:  time.Second,
			},

			expectedNilObj: true,
			expectedErr:    errors.New("could not create replica without logger"),
		},

		{
			name: "replica without config",

			logger: zap.NewNop(),
			cnfg:   nil,

			expectedNilObj: true,
			expectedErr:    nil,
		},

		{
			name: "unknown replica type in config",

			logger: zap.NewNop(),
			cnfg: &config.ReplicaConfig{
				ReplicaType:   "tsar",
				MasterAddress: "localhost:8080",
				SyncInterval:  time.Second,
			},

			expectedNilObj: true,
			expectedErr:    errors.New("unknown replica type"),
		},

		{
			name: "correct master replica",

			logger: zap.NewNop(),
			cnfg: &config.ReplicaConfig{
				ReplicaType:   master,
				MasterAddress: ":2222",
				SyncInterval:  time.Second,
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			repl, err := createReplica(test.logger, test.cnfg, walCnfg)

			if test.expectedNilObj {
				assert.Nil(t, repl)
			} else {
				assert.NotNil(t, repl)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
