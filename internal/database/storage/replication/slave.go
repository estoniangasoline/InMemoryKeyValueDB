package replication

import (
	"errors"
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage/filesystem"
	"inmemorykvdb/internal/database/storage/replication/protocol"
	"time"

	"go.uber.org/zap"
)

const (
	defaultInterval = 1 * time.Second
	maxBatchSize    = 4096
	delimElement    = ' '
)

type client interface {
	Send([]byte) ([]byte, error)
}

type Slave struct {
	directory       string
	slaveClient     client
	logger          *zap.Logger
	requestInterval time.Duration

	lastFileName string

	diskChannel chan *protocol.Response

	storageChannel chan *request.Batch
	ticker         *time.Ticker
}

func NewSlave(slaveClient client, logger *zap.Logger, options ...SlaveOption) (*Slave, error) {
	if slaveClient == nil {
		return nil, errors.New("can't create slave without client")
	}

	if logger == nil {
		return nil, errors.New("can't create slave without logger")
	}

	slave := &Slave{slaveClient: slaveClient, logger: logger}

	for _, option := range options {
		err := option(slave)

		if err != nil {
			return nil, err
		}
	}

	if slave.directory == "" {
		slave.directory = defaultDirectory
	}

	lastFileName, err := filesystem.FindLastFile(slave.directory)

	if err != nil {
		return nil, errors.New("incorrect directory")
	}

	slave.lastFileName = lastFileName

	if slave.requestInterval == 0 {
		slave.requestInterval = defaultInterval
	}

	slave.diskChannel = make(chan *protocol.Response)
	slave.ticker = time.NewTicker(slave.requestInterval)
	slave.storageChannel = make(chan *request.Batch)

	slave.start()

	return slave, nil
}

func (s *Slave) IsMaster() bool {
	return false
}

func (s *Slave) start() {
	go s.work()
	go s.writeToDisk()
}

func (s *Slave) work() {
	for range s.ticker.C {
		func() {
			defer s.ticker.Reset(s.requestInterval)

			resp, err := s.pull()

			if err != nil {
				s.logger.Error(err.Error())
				return
			}

			if s.hasNewFiles(resp) {
				s.diskChannel <- resp
				s.sendToStorage(resp)
			}
		}()
	}
}

func (s *Slave) pull() (*protocol.Response, error) {

	s.logger.Debug("started pulling a request")

	req := s.createRequest()

	s.logger.Debug("started marshalling request")

	marshaled, err := protocol.Marshal(req)

	if err != nil {
		s.logger.Error("could not pull a request because marshalling is failed")
		return nil, err
	}

	s.logger.Debug("request is complete, sending data to master")

	marshaledResp, err := s.send(marshaled)

	if err != nil {
		s.logger.Error("could not send a request")
		return nil, err
	}

	resp := &protocol.Response{}

	err = protocol.Unmarshal(resp, marshaledResp)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return resp, nil
}

func (s *Slave) send(data []byte) ([]byte, error) {

	resp, err := s.slaveClient.Send(data)

	if err != nil {
		s.logger.Error("sending data done with error")
		return nil, err
	}

	return resp, nil
}

func (s *Slave) createRequest() *protocol.Request {
	req := &protocol.Request{}

	if s.lastFileName == "" {
		s.logger.Debug("creating read all request")
		req = protocol.ReadAllRequest()
	} else {
		s.logger.Debug("creating read last request")
		req = protocol.ReadLastRequest(s.lastFileName)
	}

	return req
}

func (s *Slave) hasNewFiles(resp *protocol.Response) bool {
	if resp.Status == protocol.UnfoundStatus || len(resp.FileNames) == 0 {
		return false
	}

	s.lastFileName = resp.FileNames[len(resp.FileNames)-1]

	return true
}

/*
	start gourutine for send to server
	start gourutine for write on disk
	for timeout send req to master
	if response is nil do nothing
	if response has data send it to write module
*/

func (s *Slave) writeToDisk() {
	for resp := range s.diskChannel {
		s.logger.Debug("started write files to disk")

		err := filesystem.WriteFiles(s.directory, resp.FileNames, resp.Data)

		if err != nil {
			s.logger.Error(err.Error())
		}

		s.logger.Debug("writing is done")
	}
}

func (s *Slave) sendToStorage(resp *protocol.Response) {
	batch := request.NewBatch(maxBatchSize)

	for _, data := range resp.Data {
		batch.LoadData(data)
	}

	s.storageChannel <- batch
}

func (s *Slave) DataChan() chan *request.Batch {
	return s.storageChannel
}

/*
	if data in channel write it disk
	else block
*/

/*
	make response
	send it to server
	wait for answer
	send it back
*/
