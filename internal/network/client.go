package network

import (
	"errors"
	"net"
	"time"
)

const (
	defaultMessageSize = 4096
)

type Client struct {
	Connection  net.Conn
	IdleTimeout time.Duration
	BufferSize  int
}

func NewClient(address string, options ...ClientOption) (*Client, error) {
	conn, err := net.Dial(tcp, address)

	if err != nil {
		return nil, errors.New("failed to connect")
	}

	client := &Client{Connection: conn}

	for _, option := range options {
		option(client)
	}

	if client.IdleTimeout != 0 {
		conn.SetDeadline(time.Now().Add(client.IdleTimeout))
	}

	if client.BufferSize == 0 {
		client.BufferSize = defaultMessageSize
	}

	return client, nil
}

func (c *Client) Send(message []byte) ([]byte, error) {
	_, err := c.Connection.Write(message)

	if err != nil {
		return nil, errors.New("failed to write data")
	}

	response := make([]byte, c.BufferSize)

	responeSize, err := c.Connection.Read(response)

	if err != nil {
		return nil, errors.New("failed to read data")
	}

	if responeSize == c.BufferSize {
		return nil, errors.New("response is bigger than buffer size")
	}

	return response[:responeSize], nil
}

func (c *Client) Close() {
	if c.Connection != nil {
		c.Connection.Close()
	}
}
