package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createReadLevel(t *testing.T) {
	type testCase struct {
		name string

		logger *zap.Logger
		cnfg   *config.WalConfig

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "nil config",

			logger: zap.NewNop(),
			cnfg:   nil,

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "nil logger",

			logger: nil,
			cnfg:   nil,

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
				DataDirectory:  "./",
				FileName:       "wrahlo",
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			rl, err := createReadLevel(test.logger, test.cnfg)

			if test.expectedNilObj {
				assert.Nil(t, rl)
			} else {
				assert.NotNil(t, rl)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
