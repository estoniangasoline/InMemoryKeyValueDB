package filesystem

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadFile(t *testing.T) {
	type testCase struct {
		name string

		fileDoesNotExist bool
		dataLen          int
		fileName         string

		expectedData []byte
		expectedErr  error
	}

	directory := "C:/go/InMemoryKeyValueDB/test/filesystem/readfile/"

	testCases := []testCase{
		{
			name: "file does not exist",

			fileDoesNotExist: true,
			dataLen:          0,
			fileName:         "a",

			expectedData: []byte{},
			expectedErr:  errors.New("opening the file " + directory + "a has error open " + directory + "a: The system cannot find the file specified."),
		},
		{
			name: "buffer overflow",

			fileDoesNotExist: false,
			dataLen:          0,
			fileName:         "test.log",

			expectedData: []byte{},
			expectedErr:  errors.New("reading the file " + directory + "test.log in buffer is overflow"),
		},
		{
			name: "correct reading",

			fileDoesNotExist: false,
			dataLen:          6,
			fileName:         "test1.log",

			expectedData: []byte("qwert"),
			expectedErr:  nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if !test.fileDoesNotExist {
				file, _ := os.Create(directory + test.fileName)
				file.Write(test.expectedData)
			}

			actualData := make([]byte, test.dataLen)
			n, err := ReadFile(actualData, directory+test.fileName)

			assert.Equal(t, test.expectedData, actualData[:n])
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_ForEach(t *testing.T) {
	fileNames := []string{"test1.txt", "test2.txt", "test3.txt"}
	fileData := []string{"set biba boba", "get biba", "del biba"}
	directory := "C:/go/InMemoryKeyValueDB/test/filesystem/foreach/"

	for _, name := range fileNames {
		f, _ := os.Create(directory + name)
		f.Close()
	}

	i := 0

	ForEach(directory, fileNames, func(file *os.File) error {
		_, err := file.Write([]byte(fileData[i]))
		i++

		return err
	})

	actualData := make([]string, 3)

	i = 0

	ForEach(directory, fileNames, func(file *os.File) error {
		buf := make([]byte, 100)
		n, err := file.Read(buf)

		actualData[i] = string(buf[:n])
		i++

		return err
	})

	assert.Equal(t, fileData, actualData)
}

func Test_MakeFileNames(t *testing.T) {
	type testCase struct {
		name string

		directory string

		expectedNames []string
		expectedErr   error
	}

	testCases := []testCase{
		{
			name: "correct directory",

			directory: "C:/go/InMemoryKeyValueDB/test/filesystem/makefilenames/",

			expectedNames: []string{"testfile1.txt", "testfile2.txt", "testfile3.txt"},
			expectedErr:   nil,
		},
		{
			name: "incorrect directory",

			directory: "./biba",

			expectedNames: []string(nil),
			expectedErr:   errors.New("could not read directory ./biba with error"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			names, err := MakeFileNames(test.directory)

			assert.Equal(t, test.expectedNames, names)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_ReadAll(t *testing.T) {
	fileNames := []string{"test1.txt", "test2.txt", "test3.txt"}
	filesData := [][]byte{[]byte("set biba boba"), []byte("get biba"), []byte("del biba")}

	directory := "C:/go/InMemoryKeyValueDB/test/filesystem/readall/"

	for i, name := range fileNames {
		f, _ := os.Create(directory + name)
		f.Write(filesData[i])
		f.Close()
	}

	actual := make([][]byte, 0, len(fileNames))

	actual, _ = ReadAll(directory, fileNames, actual)

	assert.Equal(t, filesData, actual)
}

func Test_FindFile(t *testing.T) {
	type testCase struct {
		name string

		names  []string
		target string

		expectedIndex int
	}

	testCases := []testCase{
		{
			name: "file in list",

			names:  []string{"wal0.log", "wal1.log", "wal2.log", "wal3.log", "wal4.log", "wal5.log", "wal6.log", "wal7.log", "wal8.log", "wal9.log"},
			target: "wal5.log",

			expectedIndex: 5,
		},
		{
			name: "file is not in list",

			names:  []string{"wal0.log", "wal1.log", "wal2.log", "wal3.log", "wal4.log", "wal5.log", "wal6.log", "wal7.log", "wal8.log", "wal9.log"},
			target: "wal10.log",

			expectedIndex: NotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			index := FindFile(test.names, test.target)
			assert.Equal(t, test.expectedIndex, index)
		})
	}
}

func Test_FindLastFile(t *testing.T) {

	type testCase struct {
		name string

		names []string

		expectedName string
		expectedErr  error
	}

	testCases := []testCase{
		{
			name: "correct names",

			names: []string{"wal0.log", "wal1.log", "wal2.log"},

			expectedName: "wal2.log",
			expectedErr:  nil,
		},
		{
			name: "empty names",

			names: []string{},

			expectedName: "",
			expectedErr:  nil,
		},
	}

	dir := "C:/go/InMemoryKeyValueDB/test/filesystem/findlastfile/"

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			for _, name := range test.names {
				f, _ := os.Create(dir + name)
				f.Close()
			}

			lastFileName, _ := FindLastFile(dir)

			assert.Equal(t, test.expectedName, lastFileName)

			for _, name := range test.names {
				os.Remove(dir + name)
			}
		})
	}
}

func Test_WriteFiles(t *testing.T) {
	type testCase struct {
		name string

		fileNames []string
		data      []string

		expectedData []string
		expectedErr  error
	}

	testCases := []testCase{
		{
			name: "correct writing",

			fileNames: []string{"wal0.log", "wal1.log", "wal2.log"},
			data:      []string{"set biba boba", "del biba", "set boba"},

			expectedData: []string{"set biba boba", "del biba", "set boba"},
			expectedErr:  nil,
		},

		{
			name: "len names and len data not equal",

			fileNames: []string{"wal0.log", "wal1.log", "wal2.log"},
			data:      []string{"set biba boba", "del biba", "set boba", "del boba"},

			expectedData: []string{"", "", ""},
			expectedErr:  errors.New("could not write 3 file names with 4 file data"),
		},
	}

	dir := "C:/go/InMemoryKeyValueDB/test/filesystem/writefiles/"

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			dataInBytes := make([][]byte, len(test.data))

			for i := range dataInBytes {
				dataInBytes[i] = []byte(test.data[i])
			}

			err := WriteFiles(dir, test.fileNames, dataInBytes)

			actualData := make([][]byte, 0, len(test.expectedData))

			actualData, _ = ReadAll(dir, test.fileNames, actualData)

			expectedDataInBytes := make([][]byte, len(test.expectedData))

			for i := range test.expectedData {
				expectedDataInBytes[i] = []byte(test.expectedData[i])
			}

			assert.Equal(t, expectedDataInBytes, actualData)
			assert.Equal(t, test.expectedErr, err)

			for _, name := range test.fileNames {
				os.Remove(dir + name)
			}
		})
	}
}
