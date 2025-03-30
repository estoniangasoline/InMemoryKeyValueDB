package network

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WithServerTimeout(t *testing.T) {
	t.Parallel()

	timeOut := time.Second

	expectedServer := Server{IdleTimeout: timeOut}
	var actualServer Server

	option := WithServerTimeout(timeOut)
	option(&actualServer)

	assert.Equal(t, expectedServer, actualServer)
}

func Test_WithServerMaxBufferSize(t *testing.T) {
	t.Parallel()

	bufferSize := 1024

	expectedServer := Server{MaxBufferSize: bufferSize}
	var actualServer Server

	option := WithServerMaxBufferSize(bufferSize)
	option(&actualServer)

	assert.Equal(t, expectedServer, actualServer)
}

func Test_WithServerMaxConnections(t *testing.T) {
	t.Parallel()

	maxConnections := 100

	expectedServer := Server{MaxConnections: maxConnections}
	var actualServer Server

	option := WithServerMaxConnections(maxConnections)
	option(&actualServer)

	assert.Equal(t, expectedServer, actualServer)
}

func Test_WithSync(t *testing.T) {
	t.Parallel()

	isSync := true

	expectedServer := Server{IsSync: isSync}
	var actualServer Server

	option := WithSync(isSync)
	option(&actualServer)

	assert.Equal(t, expectedServer, actualServer)
}

func Test_WithClientTimeout(t *testing.T) {
	t.Parallel()

	timeOut := time.Second

	expectedClient := Client{IdleTimeout: timeOut}
	var actualClient Client

	option := WithClientTimeout(timeOut)
	option(&actualClient)

	assert.Equal(t, expectedClient, actualClient)
}

func Test_WithClientMaxBufferSize(t *testing.T) {
	t.Parallel()

	bufferSize := 1024

	expectedClient := Client{BufferSize: bufferSize}
	var actualClient Client

	option := WithClientMaxBufferSize(bufferSize)
	option(&actualClient)

	assert.Equal(t, expectedClient, actualClient)
}
