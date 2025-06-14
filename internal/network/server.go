package network

import (
	"errors"
	"fmt"
	"inmemorykvdb/pkg/concurrency/serversync"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	tcp = "tcp"

	defaultMaxBufferSize = 4096
)

type HandleRequest = func([]byte) []byte

type Server struct {
	Listener net.Listener

	IdleTimeout    time.Duration
	MaxBufferSize  int
	MaxConnections int
	Logger         *zap.Logger

	semaphore *serversync.Semaphore
}

func NewServer(address string, logger *zap.Logger, options ...ServerOption) (*Server, error) {

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	listener, err := net.Listen(tcp, address)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to address: %s", address)
	}

	server := &Server{Listener: listener, Logger: logger}

	for _, option := range options {
		option(server)
	}

	if server.MaxBufferSize == 0 {
		server.MaxBufferSize = defaultMaxBufferSize
	}

	if server.MaxConnections != 0 {
		semaphore, err := serversync.NewSemaphore(server.MaxConnections)

		if err != nil {
			return nil, fmt.Errorf("failed to create semaphore")
		}

		server.semaphore = semaphore
	}

	return server, nil
}

func (s *Server) HandleConnections(handleFunc HandleRequest) {

	defer s.Listener.Close()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {

		defer wg.Done()

		for {
			conn, err := s.Listener.Accept()

			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				continue
			}

			go s.handleConnection(conn, handleFunc)
		}
	}()

	wg.Wait()
}

func (s *Server) handleConnection(conn net.Conn, handleFunc HandleRequest) {

	defer func() {

		v := recover()

		if v != nil {
			s.Logger.Error("connection had a panic %v", zap.Any("panic", v))
		}

		err := conn.Close()

		if err != nil {
			s.Logger.Error("failed to close the connection")
		}

		if s.MaxConnections != 0 {
			s.semaphore.Release()
		}
	}()

	if s.MaxConnections != 0 {
		s.semaphore.Acquire()
	}

	request := make([]byte, s.MaxBufferSize)

	for {
		if s.IdleTimeout != 0 {
			err := conn.SetReadDeadline(time.Now().Add(s.IdleTimeout))
			if err != nil {
				s.Logger.Warn("failed to set read deadline")
			}
		}

		size, err := conn.Read(request)

		if err != nil {
			s.Logger.Error("failed to read data")
			break
		}

		if size == s.MaxBufferSize {
			s.Logger.Warn("buffer size got maximum on reading")
		}

		if s.IdleTimeout != 0 {
			err := conn.SetWriteDeadline(time.Now().Add(s.IdleTimeout))
			if err != nil {
				s.Logger.Warn("failed to set write deadline")
			}
		}

		size, err = conn.Write(handleFunc(request[:size]))

		if err != nil {
			s.Logger.Error("failed to write data")
			break
		}

		if size == s.MaxBufferSize {
			s.Logger.Warn("buffer size got maximum on writing")
		}
	}
}
