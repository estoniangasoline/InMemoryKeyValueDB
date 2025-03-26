package initialization

import (
	"errors"
	"inmemorykvdb/internal/database/storage/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createStorage(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		nilEngine bool
		logger    *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct storage",

			nilEngine: false,
			logger:    zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "storage without engine",

			nilEngine: true,
			logger:    zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("engine is nil"),
		},

		{
			name: "storage without logger",

			nilEngine: false,
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

			stor, err := createStorage(eng, test.logger)

			if test.expectedNilObj {
				assert.Nil(t, stor)
			} else {
				assert.NotNil(t, stor)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
