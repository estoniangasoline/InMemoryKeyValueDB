package replication

import (
	"errors"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage/filesystem"
	"inmemorykvdb/internal/database/storage/replication/protocol"
	"os"

	"go.uber.org/zap"
)

const (
	extraBuffer      = 10
	defaultDirectory = "C:/go/InMemoryKeyValueDB/test/"
)

type server interface {
	HandleConnections(func([]byte) []byte)
}

type Master struct {
	directory    string
	masterServer server
	logger       *zap.Logger
}

func (m *Master) DataChan() chan *request.Batch {
	return nil
}

func NewMaster(listener server, logger *zap.Logger, options ...MasterOption) (*Master, error) {
	if listener == nil {
		return nil, errors.New("server could not be nil")
	}

	if logger == nil {
		return nil, errors.New("logger could not be nil")
	}

	master := &Master{masterServer: listener, logger: logger}

	for _, option := range options {
		err := option(master)

		if err != nil {
			return nil, err
		}
	}

	if master.directory == "" {
		master.directory = defaultDirectory
	}

	master.start()

	return master, nil
}

func (m *Master) IsMaster() bool {
	return true
}

func (m *Master) start() {
	go m.masterServer.HandleConnections(func(data []byte) []byte {
		req := &protocol.Request{}
		err := protocol.Unmarshal(req, data)

		if err != nil {
			m.logger.Error(err.Error())
			return m.errorResp(err)
		}

		resp, err := m.createResponse(req)

		if err != nil {
			m.logger.Error(err.Error())
			return m.errorResp(err)
		}

		marshaledResp, err := protocol.Marshal(resp)

		if err != nil {
			m.logger.Error(err.Error())
			return m.errorResp(err)
		}

		return marshaledResp
	})
}

func (m *Master) errorResp(err error) []byte {
	r, _ := protocol.Marshal(protocol.ErrorResponse(err))
	return r
}

func (m *Master) createResponse(req *protocol.Request) (*protocol.Response, error) {
	fileNames, err := filesystem.MakeFileNames(m.directory)

	if err != nil {
		m.logger.Error("could not read file names")
		return nil, errors.New("could not read file names")
	}

	switch req.Type {
	case protocol.ReadLast:
		return m.readLast(fileNames, req.LastFileName)
	case protocol.ReadAll:
		return m.readAll(fileNames)
	}

	return nil, errors.New("unexpected request type")
}

func (m *Master) readLast(fileNames []string, lastFileName string) (*protocol.Response, error) {
	index := filesystem.FindFile(fileNames, lastFileName)

	if index == filesystem.NotFound || index == len(fileNames)-1 {
		return protocol.UnfoundResponse(), nil
	}

	targetFileName := fileNames[index+1]

	stats, err := os.Stat(m.directory + targetFileName)

	if err != nil {
		m.logger.Error("could not get stats of file")
		return nil, errors.New("could not read last file")
	}

	data := make([]byte, stats.Size()+extraBuffer)
	n, err := filesystem.ReadFile(data, m.directory+targetFileName)

	if err != nil {
		m.logger.Error(err.Error())
		return nil, errors.New("could not read last file")
	}

	data = data[:n]

	return protocol.OkResponseOneFile(targetFileName, data), nil
}

func (m *Master) readAll(fileNames []string) (*protocol.Response, error) {
	files := make([][]byte, 0, len(fileNames))
	files, err := filesystem.ReadAll(m.directory, fileNames, files)

	if err != nil {
		m.logger.Error("read all ended with problems files")
	}

	return protocol.OkResponseAllFiles(fileNames, files), err
}
