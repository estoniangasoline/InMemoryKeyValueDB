package readlevel

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	testFilesCount = 3
)

func Test_NewReadLevel(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		pattern string
		logger  *zap.Logger
		options []readLevelOptions

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct read level",

			pattern: "wal",
			logger:  zap.NewNop(),
			options: []readLevelOptions{},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "correct read level with options",

			pattern: "wal",
			logger:  zap.NewNop(),
			options: []readLevelOptions{WithFileMaxSize(1024)},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "readlevel without logger",

			pattern: "wal",
			logger:  nil,
			options: []readLevelOptions{},

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},
	}

	for i := range testFilesCount {
		dig := strconv.Itoa(i)
		os.Create("wal" + dig + ".log")
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			readLogger, err := NewReadLevel(test.logger, test.pattern, test.options...)

			if test.expectedNilObj {
				assert.Nil(t, readLogger)
			} else {
				assert.NotNil(t, readLogger)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_Read(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		pattern     string
		dataToWrite [][]byte

		expectedData [][]byte
		expectedErr  error
	}

	testCases := []testCase{
		{
			name: "correct reading",

			pattern:     "val",
			dataToWrite: [][]byte{{'1', '2', '3', '4', '5'}, {'1', '2', '3', '4', '5'}, {'1', '2', '3', '4', '5'}},

			expectedData: [][]byte{{'1', '2', '3', '4', '5'}, {'1', '2', '3', '4', '5'}, {'1', '2', '3', '4', '5'}},
			expectedErr:  nil,
		},
		{
			name: "reading out of range",

			pattern:     "law",
			dataToWrite: [][]byte{[]byte(strings.Repeat("1", 10000))},

			expectedData: [][]byte{},
			expectedErr:  errors.New("could not to read the files: law0.log "),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			for i, data := range test.dataToWrite {
				strDig := strconv.Itoa(i)
				fl, _ := os.Create(test.pattern + strDig + ".log")
				fl.Write(data)
			}

			rl, _ := NewReadLevel(zap.NewNop(), test.pattern)

			actualData, err := rl.Read()

			assert.Equal(t, test.expectedData, *actualData)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
