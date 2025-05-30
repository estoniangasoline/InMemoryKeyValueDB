package readlevel

import (
	"errors"
	"os"
	"strconv"
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

func Test_findFiles(t *testing.T) {
	type testCase struct {
		name string

		pattern string

		expectedFileNames []string
		expectedErr       error
	}

	dir := "C:\\go\\InMemoryKeyValueDB\\test\\readlevel\\findfiles\\"

	fileNames := []string{"wal1.log", "wal2.log", "wal3.log", "wal4.log"}

	for _, name := range fileNames {
		os.Create(dir + name)
	}

	testCases := []testCase{
		{
			name: "correct pattern",

			pattern: "wal",

			expectedFileNames: fileNames,
			expectedErr:       nil,
		},

		{
			name: "correct pattern but empty filenames",

			pattern: "write",

			expectedFileNames: []string(nil),
			expectedErr:       nil,
		},

		{
			name: "bad pattern",

			pattern: "[anc",

			expectedFileNames: []string(nil),
			expectedErr:       errors.New("incorrect pattern to find the files"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			rl := &readLevel{pattern: test.pattern, directory: dir}

			names, err := rl.findFiles()

			for i, name := range names {
				assert.Equal(t, dir+fileNames[i], name)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}

	for _, name := range fileNames {
		os.Remove(dir + name)
	}
}

func Test_Read(t *testing.T) {
	dir := "C:\\go\\InMemoryKeyValueDB\\test\\readlevel\\read\\"

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
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			names := make([]string, len(test.dataToWrite))

			for i, data := range test.dataToWrite {
				strDig := strconv.Itoa(i)

				name := dir + test.pattern + strDig + ".log"
				names[i] = name

				fl, _ := os.Create(name)

				fl.Write(data)
				fl.Close()
			}

			rl, _ := NewReadLevel(zap.NewNop(), test.pattern, WithDirectory(dir))

			actualData, err := rl.Read()

			assert.Equal(t, test.expectedData, actualData)
			assert.Equal(t, test.expectedErr, err)

			for _, name := range names {
				os.Remove(name)
			}
		})
	}
}
