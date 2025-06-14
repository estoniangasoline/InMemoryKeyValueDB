package request

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewBatch(t *testing.T) {
	t.Parallel()

	size := 1000
	batch := NewBatch(size)

	assert.Equal(t, cap(batch.Data), size)
	assert.Equal(t, batch.MaxSize, size)
}

func Test_Add(t *testing.T) {
	t.Parallel()

	firstArg := "aaaaa"
	secondArg := "zzzzzzzzzzz"
	req := &Request{RequestType: commands.SetCommand, Args: []string{firstArg, secondArg}}

	expectedData := []*Request{req}
	expectedSize := intSize + len(firstArg) + len(secondArg)

	batch := NewBatch(100)

	batch.Add(req)

	assert.Equal(t, expectedData, batch.Data)
	assert.Equal(t, expectedSize, batch.ByteSize)
}

func Test_ParseBatch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		requests []Request

		expectedData []byte
		expectedErr  error
	}

	testCases := []testCase{
		{
			name: "correct parsing",

			requests: []Request{{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}}},

			expectedData: []byte("SET biba boba\n"),
			expectedErr:  nil,
		},
		{
			name: "parsing incorrect request",

			requests: []Request{{RequestType: commands.IncorrectCommand, Args: []string{"aaaaa"}}},

			expectedData: []byte{},
			expectedErr:  errors.New("some requests has not parsed"),
		},
	}

	size := 1000

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			batch := NewBatch(size)

			for _, req := range test.requests {
				batch.Add(&req)
			}

			data, err := batch.ParseBatch()

			assert.Equal(t, test.expectedData, data)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_LoadData(t *testing.T) {
	type testCase struct {
		name string

		data []byte

		expectedRequests []*Request
		expectedErr      error
	}

	testCases := []testCase{
		{
			name: "correct data",

			data: []byte("SET BIBA BOBA\nDEL BIBA\n"),

			expectedRequests: []*Request{
				{
					RequestType: commands.SetCommand,

					Args: []string{"BIBA", "BOBA"},
				},

				{
					RequestType: commands.DelCommand,

					Args: []string{"BIBA"},
				},
			},
			expectedErr: nil,
		},
		{
			name: "uncorrect data",

			data: []byte("SET BIBA BOBA\nLOL BIBA\n"),

			expectedRequests: []*Request{
				{
					RequestType: commands.SetCommand,

					Args: []string{"BIBA", "BOBA"},
				},
			},
			expectedErr: errors.New("has unparsed requests"),
		},
	}

	testMaxSize := 1000

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			batch := NewBatch(testMaxSize)

			err := batch.LoadData(test.data)

			assert.Equal(t, test.expectedRequests, batch.Data)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_LoadFilesToBatch(t *testing.T) {
	type testCase struct {
		name string

		fileNames []string
		data      []string

		expectedBatch []*Request
		expectedErr   error
	}

	testCases := []testCase{
		{
			name: "correct loading",

			fileNames: []string{"wal0.log", "wal1.log", "wal2.log"},
			data:      []string{"SET biba boba" + EndElement, "SET biba boba" + EndElement, "DEL biba" + EndElement},

			expectedBatch: []*Request{
				{
					RequestType: commands.SetCommand,
					Args:        []string{"biba", "boba"},
				},

				{
					RequestType: commands.SetCommand,
					Args:        []string{"biba", "boba"},
				},

				{
					RequestType: commands.DelCommand,
					Args:        []string{"biba"},
				},
			},
			expectedErr: nil,
		},

		{
			name: "loading with errors",

			fileNames: []string{"wal0.log", "wal1.log", "wal2.log"},

			data: []string{"SET biba boba" + EndElement, "SET biba boba" + EndElement, "lol boba" + EndElement},

			expectedBatch: []*Request{
				{
					RequestType: commands.SetCommand,
					Args:        []string{"biba", "boba"},
				},

				{
					RequestType: commands.SetCommand,
					Args:        []string{"biba", "boba"},
				},
			},
			expectedErr: errors.New("has files with errors: wal2.log"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			dataInBytes := make([][]byte, len(test.data))

			for i := range dataInBytes {
				dataInBytes[i] = []byte(test.data[i])
			}

			batch := NewBatch(1000)

			err := batch.LoadFilesToBatch(test.fileNames, dataInBytes)

			assert.Equal(t, test.expectedBatch, batch.Data)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
