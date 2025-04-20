package wal

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage/wal/readlevel"
	"inmemorykvdb/internal/database/storage/wal/writelevel"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewWal(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		writer  writingLayer
		reader  readingLayer
		logger  *zap.Logger
		options []WalOptions

		expectedNilObj bool
		expectedErr    error
	}

	wl, _ := writelevel.NewWriteLevel(zap.NewNop())
	rl, _ := readlevel.NewReadLevel(zap.NewNop(), "wal")

	testCases := []testCase{
		{
			name: "correct wal",

			writer: wl,
			reader: rl,
			logger: zap.NewNop(),
			options: []WalOptions{
				WithBatchSize(100),
				WithBatchTimeout(30),
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "default wal",

			writer:  wl,
			reader:  rl,
			logger:  zap.NewNop(),
			options: []WalOptions{},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "wal without writer",

			writer:  nil,
			reader:  rl,
			logger:  zap.NewNop(),
			options: []WalOptions{},

			expectedNilObj: true,
			expectedErr:    errors.New("wal writer could not be a nil"),
		},
		{
			name: "wal without reader",

			writer:  wl,
			reader:  nil,
			logger:  zap.NewNop(),
			options: []WalOptions{},

			expectedNilObj: true,
			expectedErr:    errors.New("wal reader could not be a nil"),
		},
		{
			name: "wal without logger",

			writer: wl,
			reader: rl,

			logger:  nil,
			options: []WalOptions{},

			expectedNilObj: true,
			expectedErr:    errors.New("logger could not be nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			wal, actualErr := NewWal(test.writer, test.reader, test.logger, test.options...)

			if test.expectedNilObj {
				assert.Nil(t, wal)
			} else {
				assert.NotNil(t, wal)
			}

			assert.Equal(t, test.expectedErr, actualErr)
		})
	}
}

func Test_WriteToWal(t *testing.T) {
	wl, _ := writelevel.NewWriteLevel(zap.NewNop())
	rl, _ := readlevel.NewReadLevel(zap.NewNop(), "wal")
	logger := zap.NewNop()

	wal, _ := NewWal(wl, rl, logger)

	go func() {
		<-wal.requestChannel
		wal.blockChannel <- struct{}{}
	}()

	wal.Write(request.Request{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}})
}

func Test_writeOnDisK(t *testing.T) {
	wl, _ := writelevel.NewWriteLevel(zap.NewNop())
	rl, _ := readlevel.NewReadLevel(zap.NewNop(), "wal")
	logger := zap.NewNop()

	wal, _ := NewWal(wl, rl, logger, WithBatchSize(100))

	wal.batch.Data = append(wal.batch.Data, &request.Request{RequestType: commands.DelCommand, Args: []string{"biba"}})
	expectedData := "DEL biba\n"

	wal.writeOnDisk()

	data, _ := os.ReadFile(wl.CurrentFile.Name())
	assert.Equal(t, expectedData, string(data))
}

func Test_handleEvents(t *testing.T) {
	wl, _ := writelevel.NewWriteLevel(zap.NewNop())
	rl, _ := readlevel.NewReadLevel(zap.NewNop(), "wal")
	logger := zap.NewNop()

	wal, _ := NewWal(wl, rl, logger, WithBatchSize(100))

	wal.StartWAL()

	wal.Write(request.Request{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}})

	time.Sleep(wal.Timeout * 2)
	assert.Equal(t, 0, wal.batch.ByteSize)

	wal.Write(request.Request{RequestType: commands.DelCommand, Args: []string{string(strings.Repeat("a", wal.batch.MaxSize))}})
	assert.Equal(t, 0, wal.batch.ByteSize)
}

func Test_Read(t *testing.T) {

	type testCase struct {
		name string

		data []byte

		expectedRequests []*request.Request
	}

	testCases := []testCase{
		{
			name: "correct data",

			data: []byte("SET" + request.DelimElement + "BIBA" + request.DelimElement +
				"BOBA" + request.EndElement + "DEL" + request.DelimElement + "BIBA" + request.EndElement),

			expectedRequests: []*request.Request{
				{
					RequestType: commands.SetCommand,
					Args:        []string{"BIBA", "BOBA"},
				},

				{
					RequestType: commands.DelCommand,
					Args:        []string{"BIBA"},
				},
			},
		},

		{
			name: "reading with invalid data",

			data: []byte("SET" + request.DelimElement + "BIBA" + request.DelimElement +
				"BOBA" + request.EndElement + "LOL" + request.DelimElement + "KEK" + request.EndElement),

			expectedRequests: []*request.Request{
				{
					RequestType: commands.SetCommand,
					Args:        []string{"BIBA", "BOBA"},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			wl, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFileName("law"))
			rl, _ := readlevel.NewReadLevel(zap.NewNop(), "law")
			logger := zap.NewNop()

			os.WriteFile(wl.CurrentFile.Name(), test.data, 0644)

			defer func() {
				wl.CurrentFile.Close()
				os.Remove(wl.CurrentFile.Name())
			}()

			wal, _ := NewWal(wl, rl, logger)

			batch := wal.Read()

			assert.Equal(t, test.expectedRequests, batch.Data)
		})
	}
}
