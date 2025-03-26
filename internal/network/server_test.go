package network

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	testMaxConnections = 50
)

func Test_NewServer(t *testing.T) {

	type testCase struct {
		name string

		address string
		logger  *zap.Logger
		options []ServerOption

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct server",

			address: testAddress,
			options: []ServerOption{WithServerTimeout(testTimeout),
				WithServerMaxBufferSize(testBufferSize),
				WithServerMaxConnections(testMaxConnections)},
			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "uncorrect server address",

			address: "uncorrect address",
			options: []ServerOption{WithServerTimeout(testTimeout),
				WithServerMaxBufferSize(testBufferSize),
				WithServerMaxConnections(testMaxConnections)},
			logger: zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("failed to connect to address: uncorrect address"),
		},

		{
			name: "nil logger",

			address: testAddress,
			options: []ServerOption{WithServerTimeout(testTimeout),
				WithServerMaxBufferSize(testBufferSize),
				WithServerMaxConnections(testMaxConnections)},
			logger: nil,

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server, err := NewServer(test.address, test.logger, test.options...)

			if test.expectedNilObj {
				assert.Nil(t, server)
			} else {
				assert.NotNil(t, server)
				server.Listener.Close()
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_HandleConnections(t *testing.T) {

	server, _ := NewServer(testAddress, zap.NewNop(), []ServerOption{WithServerTimeout(testTimeout),
		WithServerMaxBufferSize(testBufferSize),
		WithServerMaxConnections(testMaxConnections)}...)

	defer server.Listener.Close()

	const defaultServerResp = "your message is "

	go func() {
		server.HandleConnections(func(data []byte) []byte {
			return fmt.Append(nil, defaultServerResp, string(data))
		})
	}()

	type testCase struct {
		name string

		connectionCount int
	}

	genMessages := func(msgCount int) []string {
		messages := make([]string, msgCount)

		for i := range msgCount {
			messages[i] = string(byte(i))
		}

		return messages
	}

	genResponses := func(respCount int) []string {
		responses := make([]string, respCount)

		for i := range respCount {
			responses[i] = fmt.Sprint(defaultServerResp, string(byte(i)))
		}

		return responses
	}

	testCases := []testCase{
		{
			name: "normal load",

			connectionCount: 98,
		},

		{
			name: "load bigger than max connections count",

			connectionCount: 122,
		},
	}

	for _, test := range testCases {

		messages := genMessages(test.connectionCount)
		resps := genResponses(test.connectionCount)

		for i := range test.connectionCount {
			go func(message []byte, expectedResp []byte) {

				client, _ := NewClient(testAddress, WithClientMaxBufferSize(testBufferSize), WithClientTimeout(testTimeout))
				resp, _ := client.Send(message)

				assert.Equal(t, string(expectedResp), string(resp))

			}([]byte(messages[i]), []byte(resps[i]))
		}
	}
}
