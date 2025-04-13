package writelevel

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithFileName(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		fileName string

		expectedErr error
	}

	testCases := []testCase{
		{
			name: "correct file name",

			fileName: "wal",

			expectedErr: nil,
		},
		{
			name: "empty file name",

			fileName: "",

			expectedErr: errors.New("file name could not be a empty string"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var wl writeLevel

			option := WithFileName(test.fileName)

			err := option(&wl)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_WithFileMaxSize(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		fileSize int

		expectedErr error
	}

	testCases := []testCase{
		{
			name: "correct file size",

			fileSize: 4096,

			expectedErr: nil,
		},
		{
			name: "zero file size",

			fileSize: 0,

			expectedErr: errors.New("max file size could not be a zero"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var wl writeLevel

			option := WithFileMaxSize(test.fileSize)

			err := option(&wl)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_WithFilePath(t *testing.T) {
	var currentWL writeLevel
	path := "/disk/logs"
	expectedWL := writeLevel{filePath: path}

	option := WithFilePath(path)
	option(&currentWL)

	assert.Equal(t, expectedWL, currentWL)
}
