package replication

import (
	"errors"
	"inmemorykvdb/internal/database/storage/replication/protocol"
	"inmemorykvdb/internal/network"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewMaster(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		nilServer bool
		logger    *zap.Logger
		options   MasterOption

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct master",

			nilServer: false,
			logger:    zap.NewNop(),
			options:   WithDirectoryMaster("/wal"),

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "master with nil server",

			nilServer: true,
			logger:    zap.NewNop(),
			options:   WithDirectoryMaster("/wal"),

			expectedNilObj: true,
			expectedErr:    errors.New("server could not be nil"),
		},
		{
			name: "master with nil logger",

			nilServer: false,
			logger:    nil,
			options:   WithDirectoryMaster("/wal"),

			expectedNilObj: true,
			expectedErr:    errors.New("logger could not be nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			var listener server

			if !test.nilServer {
				listener, _ = network.NewServer(":8080", zap.NewNop())
			}

			master, err := NewMaster(listener, test.logger, test.options)

			if test.expectedNilObj {
				assert.Nil(t, master)
			} else {
				assert.NotNil(t, master)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_readLast(t *testing.T) {
	type testCase struct {
		name string

		fileNames    []string
		lastFileName string

		targetFileName string
		targetData     []byte

		expectedResp *protocol.Response
		expectedErr  error
	}

	testCases := []testCase{
		{
			name: "unfound file",

			fileNames:    []string{"wal0.log", "wal1.log"},
			lastFileName: "wal2.log",

			expectedResp: protocol.UnfoundResponse(),
			expectedErr:  nil,
		},
		{
			name: "correct reading",

			fileNames:    []string{"wal0.log", "wal1.log", "wal2.log", "wal3.log"},
			lastFileName: "wal2.log",

			targetFileName: "wal3.log",
			targetData:     []byte("set biba boba"),

			expectedResp: protocol.OkResponseOneFile("wal3.log", []byte("set biba boba")),
			expectedErr:  nil,
		},
	}

	directory := "C:/go/InMemoryKeyValueDB/test/master/readfile/"
	serv, _ := network.NewServer(":8080", zap.NewNop())

	master, _ := NewMaster(serv, zap.NewNop(), WithDirectoryMaster(directory))

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			os.WriteFile(directory+test.targetFileName, test.targetData, 0644)

			resp, err := master.readLast(test.fileNames, test.lastFileName)

			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_readAll(t *testing.T) {
	directory := "C:/go/InMemoryKeyValueDB/test/master/readall/"
	serv, _ := network.NewServer(":8080", zap.NewNop())

	master, _ := NewMaster(serv, zap.NewNop(), WithDirectoryMaster(directory))

	data := [][]byte{[]byte("set biba boba"), []byte("get biba"), []byte("del biba"), []byte("set boba biba")}
	fileNames := []string{"wal0.log", "wal1.log", "wal2.log", "wal3.log"}

	for i, name := range fileNames {
		os.WriteFile(directory+name, data[i], 0644)
	}

	expectedResp := protocol.OkResponseAllFiles(fileNames, data)

	resp, _ := master.readAll(fileNames)

	assert.Equal(t, expectedResp, resp)
}

func Test_createResponse(t *testing.T) {
	type testCase struct {
		name string

		req *protocol.Request

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "read last",

			req: &protocol.Request{Type: protocol.ReadLast, LastFileName: "wal1.log"},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "read all",

			req: &protocol.Request{Type: protocol.ReadAll, LastFileName: "wal1.log"},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "unexpected request type",

			req: &protocol.Request{Type: -1, LastFileName: "wal1.log"},

			expectedNilObj: true,
			expectedErr:    errors.New("unexpected request type"),
		},
	}

	serv, _ := network.NewServer(":8080", zap.NewNop())

	master, _ := NewMaster(serv, zap.NewNop(), WithDirectoryMaster("C:/go/InMemoryKeyValueDB/test/master/createresponse/"))

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			resp, err := master.createResponse(test.req)

			if test.expectedNilObj {
				assert.Nil(t, resp)
			} else {
				assert.NotNil(t, resp)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_Start(t *testing.T) {
	type testCase struct {
		name string

		req *protocol.Request

		expectedResp *protocol.Response
	}

	directory := "C:/go/InMemoryKeyValueDB/test/master/start/"

	files := map[string]string{
		"wal0.log": "set biba boba",
		"wal1.log": "del bibo",
		"wal2.log": "set boba biba",
		"wal3.log": "del boba",
	}

	for name, data := range files {
		os.WriteFile(directory+name, []byte(data), 0644)
	}

	testCases := []testCase{
		{
			name: "read last req",

			req: &protocol.Request{Type: protocol.ReadLast, LastFileName: "wal0.log"},

			expectedResp: &protocol.Response{Status: protocol.OkStatus,
				FileNames: []string{"wal1.log"},
				Data:      [][]byte{[]byte(files["wal1.log"])}},
		},

		{
			name: "read all req",

			req: &protocol.Request{Type: protocol.ReadAll, LastFileName: ""},

			expectedResp: &protocol.Response{Status: protocol.OkStatus,
				FileNames: []string{"wal0.log", "wal1.log", "wal2.log", "wal3.log"},
				Data: [][]byte{[]byte(files["wal0.log"]),
					[]byte(files["wal1.log"]),
					[]byte(files["wal2.log"]),
					[]byte(files["wal3.log"]),
				},
			},
		},

		{
			name: "incorrect req",

			req: &protocol.Request{Type: -1},

			expectedResp: protocol.ErrorResponse(errors.New("unexpected request type")),
		},
	}
	serv, _ := network.NewServer(":8080", zap.NewNop())

	master, _ := NewMaster(serv, zap.NewNop(), WithDirectoryMaster(directory))
	client, _ := network.NewClient(":8080")

	master.start()

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			data, _ := protocol.Marshal(test.req)

			byteResp, _ := client.Send(data)

			resp := &protocol.Response{}

			protocol.Unmarshal(resp, byteResp)

			assert.Equal(t, test.expectedResp, resp)
		})
	}
}
