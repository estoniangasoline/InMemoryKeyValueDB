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

		logger  *zap.Logger
		options []WalOptions

		expectedNilObj bool
		expectedErr    error
	}

	dir := "C:\\go\\InMemoryKeyValueDB\\test\\wal\\newwal\\"

	wl, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFilePath(dir))
	rl, _ := readlevel.NewReadLevel(zap.NewNop(), "wal", readlevel.WithDirectory(dir))

	testCases := []testCase{
		{
			name:   "correct wal",
			logger: zap.NewNop(),
			options: []WalOptions{
				WithBatchSize(100),
				WithBatchTimeout(30),
				WithReader(rl),
				WithWriter(wl),
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name:    "default wal",
			logger:  zap.NewNop(),
			options: []WalOptions{},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "wal without logger",

			logger:  nil,
			options: []WalOptions{},

			expectedNilObj: true,
			expectedErr:    errors.New("logger could not be nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			wal, actualErr := NewWal(test.logger, test.options...)

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
	dir := "C:\\go\\InMemoryKeyValueDB\\test\\wal\\writetowal\\"
	wl, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFilePath(dir))

	wal, _ := NewWal(zap.NewNop(), WithWriter(wl))

	go func() {
		<-wal.requestChannel
		wal.blockChannel <- struct{}{}
	}()

	wal.Write(request.Request{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}})
}

func Test_writeOnDisk(t *testing.T) {
	dir := "C:\\go\\InMemoryKeyValueDB\\test\\wal\\writeondisk\\"

	wl, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFilePath(dir))

	wal, _ := NewWal(zap.NewNop(), WithBatchSize(100), WithWriter(wl))

	wal.batch.Data = append(wal.batch.Data, &request.Request{RequestType: commands.DelCommand, Args: []string{"biba"}})
	expectedData := "DEL biba\n"

	wal.writeOnDisk()

	data, _ := os.ReadFile(wl.LastFileName)
	assert.Equal(t, expectedData, string(data))
}

func Test_handleEvents(t *testing.T) {
	dir := "C:\\go\\InMemoryKeyValueDB\\test\\wal\\handleevents\\"

	wl, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFilePath(dir))

	wal, _ := NewWal(zap.NewNop(), WithBatchSize(100), WithWriter(wl))

	wal.startWAL()

	wal.Write(request.Request{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}})

	time.Sleep(wal.Timeout * 2)
	assert.Equal(t, 0, wal.batch.ByteSize)

	wal.Write(request.Request{RequestType: commands.DelCommand, Args: []string{string(strings.Repeat("a", wal.batch.MaxSize))}})
	assert.Equal(t, 0, wal.batch.ByteSize)
}

func Test_Read(t *testing.T) {

	dir := "C:\\go\\InMemoryKeyValueDB\\test\\wal\\read\\"

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

			wl, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFileName("wal"), writelevel.WithFilePath(dir))
			rl, _ := readlevel.NewReadLevel(zap.NewNop(), "wal", readlevel.WithDirectory(dir))
			logger := zap.NewNop()

			wl.Write(test.data)

			defer os.Remove(wl.LastFileName)

			wal, _ := NewWal(logger, WithReader(rl))

			batch := wal.Read()

			assert.Equal(t, test.expectedRequests, batch.Data)
		})
	}
}
