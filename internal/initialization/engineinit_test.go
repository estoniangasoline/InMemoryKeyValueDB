package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createEngine(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		engconf *config.EngineConfig
		logger  *zap.Logger

		expectedNilObject bool
		expectedErr       error
	}

	testCases := []testCase{
		{
			name: "correct engine",

			engconf: &config.EngineConfig{
				EngineType: "in_memory",
			},
			logger: zap.NewNop(),

			expectedNilObject: false,
			expectedErr:       nil,
		},
		{
			name: "nil config",

			engconf: nil,
			logger:  zap.NewNop(),

			expectedNilObject: false,
			expectedErr:       nil,
		},
		{
			name: "nil logger",

			engconf: &config.EngineConfig{
				EngineType: "in_memory",
			},
			logger: nil,

			expectedNilObject: true,
			expectedErr:       errors.New("logger is nil"),
		},

		{
			name: "unknown engine type",

			engconf: &config.EngineConfig{
				EngineType: "vasily",
			},
			logger: zap.NewNop(),

			expectedNilObject: false,
			expectedErr:       nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var testEngine engineLayer

			if !test.expectedNilObject {
				if test.engconf != nil {
					switch test.engconf.EngineType {
					case inMemoryType:
						testEngine, _ = engine.NewInMemoryEngine(test.logger)
					default:
						testEngine, _ = engine.NewInMemoryEngine(test.logger)
					}
				} else {
					testEngine, _ = engine.NewInMemoryEngine(test.logger)
				}
			}

			actualEngine, actualErr := createEngine(test.engconf, test.logger)

			if test.expectedNilObject {
				assert.Nil(t, actualEngine)
			} else {
				assert.Equal(t, testEngine, actualEngine)
			}

			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}
