package network

import (
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testAddress    = "localhost:7777"
	serverResponse = "response"
	testBufferSize = 4096
	testTimeout    = 30
)

func Test_NewClient(t *testing.T) {

	listener, _ := net.Listen(tcp, testAddress)

	defer listener.Close()

	go func() {
		for {
			listener.Accept()
		}
	}()

	type testCase struct {
		name string

		address string
		options []ClientOption

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct client",

			address: "localhost:7777",
			options: []ClientOption{WithClientMaxBufferSize(testBufferSize), WithClientTimeout(testTimeout)},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		{
			name: "uncorrect client",

			address: "localhost:2222",
			options: []ClientOption{WithClientMaxBufferSize(testBufferSize), WithClientTimeout(testTimeout)},

			expectedNilObj: true,
			expectedErr:    errors.New("failed to connect"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			client, err := NewClient(test.address, test.options...)

			if test.expectedNilObj {
				assert.Nil(t, client)
			} else {
				assert.NotNil(t, client)
			}

			assert.Equal(t, test.expectedErr, err)

		})
	}
}

func Test_Send(t *testing.T) {

	startServer := func(response string) {

		listener, _ := net.Listen(tcp, testAddress)
		defer listener.Close()

		conn, _ := listener.Accept()

		conn.Read(make([]byte, testBufferSize))

		conn.Write([]byte(response))
	}

	type testCase struct {
		name string

		request  string
		response string

		expectedResponse string
		expectedErr      error
	}

	testCases := []testCase{
		{
			name: "correct request",

			request:  "client request",
			response: serverResponse,

			expectedResponse: serverResponse,
			expectedErr:      nil,
		},

		{
			name: "buffer overflow",

			request:  "client request",
			response: strings.Repeat("A", testBufferSize),

			expectedResponse: "",
			expectedErr:      errors.New("response is bigger than buffer size"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			go startServer(test.response)

			client, _ := NewClient(testAddress, WithClientMaxBufferSize(testBufferSize), WithClientTimeout(testBufferSize))

			resp, err := client.Send([]byte(test.request))

			assert.Equal(t, test.expectedResponse, string(resp))
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
