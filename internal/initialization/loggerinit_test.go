package initialization

import (
	"inmemorykvdb/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createLogger(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		cnfg *config.LoggingConfig

		expectedNilObject bool
		expectedErr       error
	}

	testCases := []testCase{
		{
			name: "debug level logger",

			cnfg: &config.LoggingConfig{
				Level:  "debug",
				Output: defaultOutput,
			},

			expectedNilObject: false,
			expectedErr:       nil,
		},

		{
			name: "info level logger",

			cnfg: &config.LoggingConfig{
				Level:  "info",
				Output: defaultOutput,
			},

			expectedNilObject: false,
			expectedErr:       nil,
		},

		{
			name: "warn level logger",

			cnfg: &config.LoggingConfig{
				Level:  "warn",
				Output: defaultOutput,
			},

			expectedNilObject: false,
			expectedErr:       nil,
		},

		{
			name: "error level logger",

			cnfg: &config.LoggingConfig{
				Level:  "error",
				Output: defaultOutput,
			},

			expectedNilObject: false,
			expectedErr:       nil,
		},

		{
			name: "logger without level",

			cnfg: &config.LoggingConfig{
				Level:  "",
				Output: defaultOutput,
			},

			expectedNilObject: false,
			expectedErr:       nil,
		},

		{
			name: "logger without output",

			cnfg: &config.LoggingConfig{
				Level:  "debug",
				Output: "",
			},

			expectedNilObject: false,
			expectedErr:       nil,
		},

		{
			name: "nil config",

			cnfg: nil,

			expectedNilObject: false,
			expectedErr:       nil,
		},

		{
			name: "incorrect logging level config",

			cnfg: &config.LoggingConfig{
				Level:  "vasily",
				Output: defaultOutput,
			},

			expectedNilObject: false,
			expectedErr:       nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			t.Parallel()

			logger, err := createLogger(test.cnfg)

			if test.expectedNilObject {
				assert.Nil(t, logger)
			} else {
				assert.NotNil(t, logger)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
