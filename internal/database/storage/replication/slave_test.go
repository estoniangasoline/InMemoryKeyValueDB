package replication

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage/filesystem"
	"inmemorykvdb/internal/database/storage/replication/protocol"
	"inmemorykvdb/internal/network"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewSlave(t *testing.T) {
	type testCase struct {
		name string

		nilClient bool
		logger    *zap.Logger
		options   []SlaveOption

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct slave",

			nilClient: false,
			logger:    zap.NewNop(),
			options:   []SlaveOption{},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "correct slave with options",

			nilClient: false,
			logger:    zap.NewNop(),
			options:   []SlaveOption{WithDirectorySlave("c:/go/InMemoryKeyValueDB/test/slave/"), WithInterval(time.Second * 2)},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "slave without logger",

			nilClient: false,
			logger:    nil,
			options:   []SlaveOption{},

			expectedNilObj: true,
			expectedErr:    errors.New("can't create slave without logger"),
		},
		{
			name: "slave without client",

			nilClient: true,
			logger:    zap.NewNop(),
			options:   []SlaveOption{},

			expectedNilObj: true,
			expectedErr:    errors.New("can't create slave without client"),
		},
		{
			name: "incorrect directory",

			nilClient: false,
			logger:    zap.NewNop(),
			options:   []SlaveOption{WithDirectorySlave("./biba")},

			expectedNilObj: true,
			expectedErr:    errors.New("incorrect directory"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			var slaveClient client

			if !test.nilClient {
				slaveClient, _ = network.NewClient(":8080")
			}

			slave, err := NewSlave(slaveClient, test.logger, test.options...)

			if test.expectedNilObj {
				assert.Nil(t, slave)
			} else {
				assert.NotNil(t, slave)
			}
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_writeToDisk(t *testing.T) {
	dir := "C:/go/InMemoryKeyValueDB/test/slave/writetodisk/"
	resps := []*protocol.Response{
		{
			Status:    protocol.OkStatus,
			FileNames: []string{"wal.log"},
			Data: [][]byte{
				[]byte("set biba boba"),
			},
		},
		{
			Status:    protocol.OkStatus,
			FileNames: []string{"wal1.log"},
			Data: [][]byte{
				[]byte("del biba"),
			},
		},
		{
			Status:    protocol.OkStatus,
			FileNames: []string{"wal2.log"},
			Data: [][]byte{
				[]byte("set boba biba"),
			},
		},
	}

	client, _ := network.NewClient(":8080")
	slave, _ := NewSlave(client, zap.NewNop(), WithDirectorySlave(dir))

	block := make(chan struct{})

	go func() {
		defer func() {
			block <- struct{}{}
		}()

		slave.writeToDisk()
	}()

	wg := &sync.WaitGroup{}

	wg.Add(len(resps))

	for i := range len(resps) {
		go func() {
			defer wg.Done()
			slave.diskChannel <- resps[i]
		}()
	}

	wg.Wait()
	close(slave.diskChannel)
	<-block

	fileNames := make([]string, len(resps))
	expectedData := make([][]byte, 0, len(resps))

	for i, resp := range resps {
		fileNames[i] = resp.FileNames[0]
		expectedData = append(expectedData, resp.Data...)
	}

	actualData := make([][]byte, 0, len(resps))

	actualData, _ = filesystem.ReadAll(dir, fileNames, actualData)

	assert.Equal(t, expectedData, actualData)
}

func Test_hasNewFiles(t *testing.T) {
	type testCase struct {
		name string

		resp *protocol.Response

		expectedAnswer       bool
		expectedLastFileName string
	}

	testCases := []testCase{
		{
			name: "has new files",

			resp: &protocol.Response{Status: protocol.OkStatus, FileNames: []string{"wal.log"}},

			expectedAnswer:       true,
			expectedLastFileName: "wal.log",
		},
		{
			name: "has not new files",
			resp: &protocol.Response{Status: protocol.UnfoundStatus},

			expectedAnswer:       false,
			expectedLastFileName: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			client, _ := network.NewClient(":8080")
			slave, _ := NewSlave(client, zap.NewNop())

			ans := slave.hasNewFiles(test.resp)

			assert.Equal(t, test.expectedAnswer, ans)
			assert.Equal(t, test.expectedLastFileName, slave.lastFileName)
		})
	}
}

func Test_createRequest(t *testing.T) {
	type testCase struct {
		name string

		directory string

		expectedRequest *protocol.Request
	}

	testCases := []testCase{
		{
			name: "read all req",

			directory: "C:/go/InMemoryKeyValueDB/test/slave/createrequest/empty/",

			expectedRequest: &protocol.Request{Type: protocol.ReadAll},
		},
		{
			name: "read last req",

			directory: "C:/go/InMemoryKeyValueDB/test/slave/createrequest/notempty/",

			expectedRequest: &protocol.Request{
				Type:         protocol.ReadLast,
				LastFileName: "wal.log",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			client, _ := network.NewClient(":8080")
			slave, _ := NewSlave(client, zap.NewNop(), WithDirectorySlave(test.directory))

			req := slave.createRequest()

			assert.Equal(t, test.expectedRequest, req)
		})
	}
}

func Test_send(t *testing.T) {
	testData := []byte("set biba boba")
	expectedResponse := []byte("get boba")

	server, _ := network.NewServer(":8080", zap.NewNop())

	go server.HandleConnections(func(data []byte) []byte {
		return expectedResponse
	})

	client, _ := network.NewClient(":8080")
	slave, _ := NewSlave(client, zap.NewNop())

	resp, err := slave.send(testData)

	assert.Equal(t, expectedResponse, resp)
	assert.Nil(t, err)
}

func Test_pull(t *testing.T) {
	masterDir := "C:/go/InMemoryKeyValueDB/test/slave/pull/masterdir/"
	fileNames := []string{"wal0.log", "wal1.log", "wal2.log"}
	fileData := [][]byte{[]byte("set biba boba"), []byte("del biba"), []byte("set boba biba")}

	i := 0
	filesystem.ForEach(masterDir, fileNames, func(file *os.File) error {
		file.Write(fileData[i])
		i++
		return nil
	})

	server, _ := network.NewServer(":8080", zap.NewNop())
	master, _ := NewMaster(server, zap.NewNop(), WithDirectoryMaster(masterDir))
	master.start()

	type testCase struct {
		name string

		dir          string
		lastFileName string

		expectedResponse *protocol.Response
		expectedErr      error
	}

	testCases := []testCase{
		{
			name: "read all req",

			dir:          "C:/go/InMemoryKeyValueDB/test/slave/pull/slavedirreadall/",
			lastFileName: "",

			expectedResponse: &protocol.Response{
				Status:    protocol.OkStatus,
				FileNames: fileNames,
				Data:      fileData},
			expectedErr: nil,
		},

		{
			name: "read last req",

			dir:          "C:/go/InMemoryKeyValueDB/test/slave/pull/slavedirreadlast/",
			lastFileName: fileNames[0],
			expectedResponse: &protocol.Response{
				Status:    protocol.OkStatus,
				FileNames: []string{fileNames[1]},
				Data:      [][]byte{fileData[1]}},
			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			os.Create(test.dir + test.lastFileName)

			client, _ := network.NewClient(":8080")
			slave, _ := NewSlave(client, zap.NewNop(), WithDirectorySlave(test.dir))

			resp, err := slave.pull()

			assert.Equal(t, test.expectedResponse, resp)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_StartSlave(t *testing.T) {
	masterDir := "C:/go/InMemoryKeyValueDB/test/slave/start/masterdir/"
	fileNames := []string{"wal0.log", "wal1.log", "wal2.log"}
	fileData := [][]byte{[]byte("SET biba boba\n"), []byte("DEL biba\n"), []byte("SET boba biba\n")}

	i := 0
	filesystem.ForEach(masterDir, fileNames, func(file *os.File) error {
		file.Write(fileData[i])
		i++
		return nil
	})

	server, _ := network.NewServer(":8080", zap.NewNop())
	master, _ := NewMaster(server, zap.NewNop(), WithDirectoryMaster(masterDir))

	master.start()

	slaveDir := "C:/go/InMemoryKeyValueDB/test/slave/start/slavedir/"
	client, _ := network.NewClient(":8080")
	slave, _ := NewSlave(client, zap.NewNop(), WithDirectorySlave(slaveDir))

	slave.start()

	time.Sleep(2 * time.Second)
	slave.ticker.Stop()

	i = 0
	filesystem.ForEach(slaveDir, fileNames, func(file *os.File) error {
		buf := make([]byte, 100)
		n, err := file.Read(buf)

		assert.Nil(t, err)
		buf = buf[:n]

		assert.Equal(t, fileData[i], buf)
		i++

		return nil
	})

	for _, name := range fileNames {
		os.Remove(slaveDir + name)
	}
}

func Test_sendToStorage(t *testing.T) {
	expectedBatch := &request.Batch{
		Data: []*request.Request{
			{
				RequestType: commands.SetCommand,

				Args: []string{"biba", "boba"},
			},
			{
				RequestType: commands.DelCommand,

				Args: []string{"biba", "boba"},
			},
			{
				RequestType: commands.SetCommand,
				Args:        []string{"BOBA", "BABA"},
			},
		},
	}

	resp := &protocol.Response{
		Status:    protocol.OkStatus,
		FileNames: []string{"wal.log"},
		Data: [][]byte{
			[]byte("SET biba boba\n"),
			[]byte("DEL biba boba\n"),
			[]byte("SET BOBA BABA\n"),
		},
	}

	client, _ := network.NewClient(":8080")
	slave, _ := NewSlave(client, zap.NewNop())

	var batch *request.Batch
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		batch = <-slave.storageChannel
	}()

	slave.sendToStorage(resp)

	wg.Wait()

	assert.Equal(t, expectedBatch.Data, batch.Data)
}
