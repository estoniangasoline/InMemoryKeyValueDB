package storage

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage/engine"
	"inmemorykvdb/internal/database/storage/replication"
	"inmemorykvdb/internal/network"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewStorage(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		eng     engineLayer
		logger  *zap.Logger
		options []StorageOption

		expectedNilObj bool
		expectedErr    error
	}

	client, _ := network.NewClient(":8080")

	testEngine, _ := engine.NewInMemoryEngine(zap.NewNop())
	slave, _ := replication.NewSlave(client, zap.NewNop())

	testCases := []testCase{

		{
			name: "valid storage",

			eng:     testEngine,
			options: []StorageOption{},
			logger:  zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "storage without logger",

			eng:     testEngine,
			options: []StorageOption{},
			logger:  nil,

			expectedNilObj: true,
			expectedErr:    errors.New("could not create storage without logger"),
		},

		{
			name: "slave storage without datachan",

			eng:     testEngine,
			options: []StorageOption{WithReplica(slave)},
			logger:  zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("could not create slave node without data chan"),
		},

		{
			name: "storage without engine",

			eng:     nil,
			options: []StorageOption{},
			logger:  zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("could not create storage without engine"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			_, err := NewStorage(test.logger, test.eng, test.options...)

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
			expectedErr: errors.New("incorrect request type"),
		},
	}

	testEngine, _ := engine.NewInMemoryEngine(zap.NewNop())
	logger := zap.NewNop()

	storage, _ := NewStorage(logger, testEngine)

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

	stor, _ := NewStorage(zap.NewNop(), eng)

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

	answer, _ := stor.engine.GET("biba")
	assert.Equal(t, "boba", answer)

	answer, _ = stor.engine.GET("boba")
	assert.Equal(t, "biba", answer)

	batch = request.Batch{Data: []*request.Request{
		{
			RequestType: commands.DelCommand,
			Args:        []string{"biba"},
		},
	}}

	stor.recoverData(&batch)

	answer, _ = stor.engine.GET("biba")
	assert.Equal(t, "", answer)

	answer, _ = stor.engine.GET("boba")
	assert.Equal(t, "biba", answer)
}

func Test_synchronization(t *testing.T) {
	eng, _ := engine.NewInMemoryEngine(zap.NewNop())
	client, _ := network.NewClient(":8080")
	slave, _ := replication.NewSlave(client, zap.NewNop())
	dataChan := make(chan *request.Batch)

	stor, _ := NewStorage(zap.NewNop(), eng, WithReplica(slave), WithDataChan(dataChan))

	batch := &request.Batch{Data: []*request.Request{
		{
			RequestType: commands.SetCommand,

			Args: []string{"biba", "boba"},
		},

		{
			RequestType: commands.SetCommand,

			Args: []string{"boba", "biba"},
		},

		{
			RequestType: commands.SetCommand,

			Args: []string{"bib", "bob"},
		},

		{
			RequestType: commands.DelCommand,

			Args: []string{"bib", "bob"},
		},
	}}

	stor.synchronization()

	dataChan <- batch
	dataChan <- &request.Batch{}

	answ, _ := stor.engine.GET("biba")
	assert.Equal(t, "boba", answ)

	answ, _ = stor.engine.GET("boba")
	assert.Equal(t, "biba", answ)

	answ, _ = stor.engine.GET("bib")
	assert.Equal(t, "", answ)
}
