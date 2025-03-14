package storage

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/storage/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewStorage(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		nilEngine bool
		logger    *zap.Logger

		expectedErr error
	}

	testCases := []testCase{

		{
			name: "valid storage",

			nilEngine: false,
			logger:    zap.NewNop(),

			expectedErr: nil,
		},

		{
			name: "storage without engine",

			nilEngine: true,
			logger:    zap.NewNop(),

			expectedErr: errors.New("could not create storage without engine"),
		},

		{
			name: "storage without logger",

			nilEngine: false,
			logger:    nil,

			expectedErr: errors.New("could not create storage without logger"),
		},

		{
			name: "storage without anything",

			nilEngine: true,
			logger:    nil,

			expectedErr: errors.New("could not create storage without engine and logger"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			var testEngine engineLayer

			if test.nilEngine {
				testEngine = nil
			} else {
				testEngine, _ = engine.NewInMemoryEngine(zap.NewNop())
			}

			_, err := NewStorage(testEngine, test.logger)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_HandleRequest(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		requestType int
		args        []string

		expectStr   string
		expectedErr error
	}

	testCases := []testCase{

		{
			name: "set request",

			requestType: commands.SetCommand,
			args:        []string{"asdfg", "qwert"},

			expectStr:   "",
			expectedErr: nil,
		},

		{
			name: "get request",

			requestType: commands.GetCommand,
			args:        []string{"asdfg"},

			expectStr:   "qwert",
			expectedErr: nil,
		},

		{
			name: "del request",

			requestType: commands.DelCommand,
			args:        []string{"asdfg"},

			expectStr:   "",
			expectedErr: nil,
		},

		{
			name: "not a correct request",

			requestType: 10,
			args:        []string{"asdfg"},

			expectStr:   "",
			expectedErr: errors.New("uncorrect request type"),
		},
	}

	testEngine, _ := engine.NewInMemoryEngine(zap.NewNop())

	storage, _ := NewStorage(testEngine, zap.NewNop())

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			actualValue, actualErr := storage.HandleRequest(test.requestType, test.args...)

			assert.Equal(t, test.expectStr, actualValue)
			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}
