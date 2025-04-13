package storage

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewStorage(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		logger  *zap.Logger
		options []StorageOption

		expectedErr error
	}

	testEngine, _ := engine.NewInMemoryEngine(zap.NewNop())

	testCases := []testCase{

		{
			name: "valid storage",

			options: []StorageOption{WithEngine(testEngine)},
			logger:  zap.NewNop(),

			expectedErr: nil,
		},

		{
			name: "storage without logger",

			options: []StorageOption{WithEngine(testEngine)},
			logger:  nil,

			expectedErr: errors.New("could not create storage without logger"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			_, err := NewStorage(test.logger, test.options...)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_HandleRequest(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		request request.Request

		expectStr   string
		expectedErr error
	}

	testCases := []testCase{

		{
			name: "set request",

			request: request.Request{RequestType: commands.SetCommand, Args: []string{"asdfg", "qwert"}},

			expectStr:   okAnswer,
			expectedErr: nil,
		},

		{
			name: "get request",

			request: request.Request{RequestType: commands.GetCommand, Args: []string{"asdfg"}},

			expectStr:   "qwert",
			expectedErr: nil,
		},

		{
			name: "del request",

			request: request.Request{RequestType: commands.DelCommand, Args: []string{"asdfg"}},

			expectStr:   okAnswer,
			expectedErr: nil,
		},

		{
			name: "not a correct request",

			request: request.Request{RequestType: 10, Args: []string{"asdfg"}},

			expectStr:   "",
			expectedErr: errors.New("uncorrect request type"),
		},
	}

	testEngine, _ := engine.NewInMemoryEngine(zap.NewNop())
	logger := zap.NewNop()

	storage, _ := NewStorage(logger, WithEngine(testEngine))

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			actualValue, actualErr := storage.HandleRequest(test.request)

			assert.Equal(t, test.expectStr, actualValue)
			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

func Test_recoverData(t *testing.T) {
	eng, _ := engine.NewInMemoryEngine(zap.NewNop())

	stor, _ := NewStorage(zap.NewNop(), WithEngine(eng))

	batch := request.Batch{Data: []*request.Request{
		{
			RequestType: commands.SetCommand,
			Args:        []string{"biba", "boba"},
		},

		{
			RequestType: commands.SetCommand,
			Args:        []string{"boba", "biba"},
		},
	}}

	stor.recoverData(&batch)

	answer, _ := stor.get("biba")
	assert.Equal(t, "boba", answer)

	answer, _ = stor.get("boba")
	assert.Equal(t, "biba", answer)

	batch = request.Batch{Data: []*request.Request{
		{
			RequestType: commands.DelCommand,
			Args:        []string{"biba"},
		},
	}}

	stor.recoverData(&batch)

	answer, _ = stor.get("biba")
	assert.Equal(t, "", answer)

	answer, _ = stor.get("boba")
	assert.Equal(t, "biba", answer)
}
