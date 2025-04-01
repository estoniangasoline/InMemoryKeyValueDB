package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createServer(t *testing.T) {

	type testCase struct {
		name string

		cnfg   *config.NetworkConfig
		logger *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct server",

			cnfg: &config.NetworkConfig{
				Address:        ":7777",
				MaxConnections: 100,
				MaxMessageSize: "4KB",
				IdleTimeout:    30,
			},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "nil config",

			cnfg:   nil,
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "nil logger",

			cnfg: &config.NetworkConfig{
				Address:        ":7777",
				MaxConnections: 100,
				MaxMessageSize: "4KB",
				IdleTimeout:    30,
			},
			logger: nil,

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},

		{
			name: "empty address",

			cnfg: &config.NetworkConfig{
				Address:        "",
				MaxConnections: 100,
				MaxMessageSize: "4KB",
				IdleTimeout:    30,
			},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "zero max connections",

			cnfg: &config.NetworkConfig{
				Address:        ":7777",
				MaxConnections: 0,
				MaxMessageSize: "4KB",
				IdleTimeout:    30,
			},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "uncorrect max message size",

			cnfg: &config.NetworkConfig{
				Address:        ":7777",
				MaxConnections: 100,
				MaxMessageSize: "4",
				IdleTimeout:    30,
			},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "unknown max message size",

			cnfg: &config.NetworkConfig{
				Address:        ":7777",
				MaxConnections: 100,
				MaxMessageSize: "4TB",
				IdleTimeout:    30,
			},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "max message size without count",

			cnfg: &config.NetworkConfig{
				Address:        ":7777",
				MaxConnections: 100,
				MaxMessageSize: "KB",
				IdleTimeout:    30,
			},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server, err := createServer(test.cnfg, test.logger)

			if server != nil {
				defer server.Listener.Close()
			}

			if test.expectedNilObj {
				assert.Nil(t, server)
			} else {
				assert.NotNil(t, server)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
