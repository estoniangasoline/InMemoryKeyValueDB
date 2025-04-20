package network

import (
	"time"
)

type ServerOption func(*Server)

func WithServerTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.IdleTimeout = timeout
	}
}

func WithServerMaxBufferSize(bufferSize int) ServerOption {
	return func(s *Server) {
		s.MaxBufferSize = bufferSize
	}
}

func WithServerMaxConnections(maxConnections int) ServerOption {
	return func(s *Server) {
		s.MaxConnections = maxConnections
	}
}

type ClientOption func(*Client)

func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.IdleTimeout = timeout
	}
}

func WithClientMaxBufferSize(maxBufferSize int) ClientOption {
	return func(c *Client) {
		c.BufferSize = maxBufferSize
	}
}
