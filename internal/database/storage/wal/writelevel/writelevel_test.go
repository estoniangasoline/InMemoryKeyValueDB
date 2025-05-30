package writelevel

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_NewWriteLevel(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		options []writeLevelOptions
		logger  *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "default",

			options: []writeLevelOptions{},
			logger:  zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "without logger",

			options: []writeLevelOptions{},
			logger:  nil,

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},
		{
			name: "custom",

			options: []writeLevelOptions{
				WithFileMaxSize(2048),
				WithFileName("logs"),
				WithFilePath("../")},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			wl, err := NewWriteLevel(test.logger, test.options...)

			if test.expectedNilObj {
				assert.Nil(t, wl)
			} else {
				assert.NotNil(t, wl)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_createFile(t *testing.T) {
	t.Parallel()

	wl, _ := NewWriteLevel(zap.NewNop())

	file, err := wl.createFile()

	require.Nil(t, err)

	_, err = file.Stat()

	assert.NotEqual(t, os.ErrNotExist, err)
}

func Test_checkFileIsExist(t *testing.T) {

	wl, _ := NewWriteLevel(zap.NewNop())

	stringInd := strconv.Itoa(wl.nextFileIndex)
	expectedIndex := wl.nextFileIndex + 1

	os.Create(wl.filePath + wl.fileName + stringInd + fileExtension)
	wl.checkFileIsExist()

	assert.Equal(t, expectedIndex, wl.nextFileIndex)
}

func Test_Write(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		data []byte

		expectedErr error
	}

	testCases := []testCase{
		{
			name: "correct data",

			data: []byte("SET BIBA BOBA"),

			expectedErr: nil,
		},
		{
			name: "zero data",

			data: []byte{},

			expectedErr: errors.New("data is empty"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			wl, _ := NewWriteLevel(zap.NewNop(), WithFilePath("C:/go/InMemoryKeyValueDB/test/wal/wl/"))

			_, err := wl.Write(test.data)

			assert.Equal(t, test.expectedErr, err)

			if err == nil {
				actualData, _ := os.ReadFile("C:/go/InMemoryKeyValueDB/test/wal/wl/write_ahead1.log")
				assert.Equal(t, test.data, actualData)
				os.Remove("C:/go/InMemoryKeyValueDB/test/wal/wl/write_ahead1.log")
			}
		})
	}
}
