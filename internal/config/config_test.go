package config

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	correctConfig = `
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:3223"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: 5m
  is_sync: true
logging:
  level: "info"
  output: "logging.txt"
wal:
  flushing_batch_size: 100
  flushing_batch_timeout: "10ms"
  max_segment_size: "10MB"
  data_directory: "/data/spider/wal"
  file_name: "write_ahead_log"
replication:
  replica_type: "slave"
  master_address: "127.0.0.1:3232"
  sync_interval: "1s"
`
)

func Test_NewConfig(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		stringConfig string

		expectedConfig *Config
	}

	testCases := []testCase{

		{
			name: "correct config",

			stringConfig: correctConfig,

			expectedConfig: &Config{
				Engine: &EngineConfig{
					EngineType: "in_memory",
				},

				Network: &NetworkConfig{
					Address:        "127.0.0.1:3223",
					MaxConnections: 100,
					MaxMessageSize: "4KB",
					IdleTimeout:    time.Minute * 5,
					IsSync:         true,
				},

				Logging: &LoggingConfig{
					Level:  "info",
					Output: "logging.txt",
				},
				WalConfig: &WalConfig{
					BatchSize:      100,
					BatchTimeout:   10 * time.Millisecond,
					MaxSegmentSize: "10MB",
					DataDirectory:  "/data/spider/wal",
					FileName:       "write_ahead_log",
				},
				Replication: &ReplicaConfig{
					ReplicaType:   "slave",
					MasterAddress: "127.0.0.1:3232",
					SyncInterval:  time.Second,
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reader := strings.NewReader(test.stringConfig)
			cnfg, err := NewConfig(reader)

			assert.NoError(t, err)
			assert.Equal(t, test.expectedConfig, cnfg)
		})
	}
}
