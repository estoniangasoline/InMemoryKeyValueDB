package writelevel

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_checkFileSize(t *testing.T) {
	t.Parallel()

	wl, _ := NewWriteLevel(zap.NewNop())

	err := wl.checkFileSize()
	assert.Nil(t, err)

	wl.CurrentFile.Write(make([]byte, wl.fileMaxSize))

	err = wl.checkFileSize()
	assert.Nil(t, err)
}

func Test_createFile(t *testing.T) {
	t.Parallel()

	wl, _ := NewWriteLevel(zap.NewNop())

	err := wl.createFile()

	assert.Nil(t, err)

	var stringInd = strconv.Itoa(wl.nextFileIndex - 1)
	_, err = os.Stat(wl.fileName + stringInd + fileExtension)

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
			name: "zero data",

			data: []byte{},

			expectedErr: errors.New("data is empty"),
		},
		{
			name: "correct data",

			data: []byte("SET BIBA BOBA"),

			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			wl, _ := NewWriteLevel(zap.NewNop())

			_, err := wl.Write(&test.data)

			actualData, _ := os.ReadFile(wl.CurrentFile.Name())

			assert.Equal(t, test.data, actualData)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
