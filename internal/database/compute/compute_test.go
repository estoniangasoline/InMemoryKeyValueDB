package compute

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"inmemorykvdb/internal/database/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewSimpleCompute(t *testing.T) {

	type testCase struct {
		name string

		logger *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{

		{
			name: "correct compute",

			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "compute without logger",

			logger: nil,

			expectedNilObj: true,
			expectedErr:    errors.New("could not make compute without logger"),
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {

			compute, err := NewCompute(test.logger)

			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, compute)
			} else {
				assert.NotNil(t, compute)
			}

		})
	}
}

func Test_Parse(t *testing.T) {

	type testCase struct {
		name string

		data string

		expectedRequest request.Request
		expectedErr     error
	}

	testCases := []testCase{

		{
			name: "correct get request",

			data: "get qwerty",

			expectedRequest: request.Request{RequestType: commands.GetCommand, Args: []string{"qwerty"}},
			expectedErr:     nil,
		},

		{
			name: "correct set request",

			data: "set qwerty asdfgh",

			expectedRequest: request.Request{RequestType: commands.SetCommand, Args: []string{"qwerty", "asdfgh"}},
			expectedErr:     nil,
		},

		{
			name: "correct del request",

			data: "del qwerty",

			expectedRequest: request.Request{RequestType: commands.DelCommand, Args: []string{"qwerty"}},
			expectedErr:     nil,
		},

		{
			name: "get request in uppercase",

			data: "GET qwerty",

			expectedRequest: request.Request{RequestType: commands.GetCommand, Args: []string{"qwerty"}},
			expectedErr:     nil,
		},

		{
			name: "set request in uppercase",

			data: "SET qwerty asdfgh",

			expectedRequest: request.Request{RequestType: commands.SetCommand, Args: []string{"qwerty", "asdfgh"}},
			expectedErr:     nil,
		},

		{
			name: "del request in uppercase",

			data: "DEL qwerty",

			expectedRequest: request.Request{RequestType: commands.DelCommand, Args: []string{"qwerty"}},
			expectedErr:     nil,
		},

		{
			name: "get request with args more than needed",

			data: "GET qwerty poiuy",

			expectedRequest: request.Request{RequestType: commands.GetCommand, Args: []string{"qwerty"}},
			expectedErr:     nil,
		},

		{
			name: "set request with args more than needed",

			data: "SET qwerty asdfgh [poiuyt]",

			expectedRequest: request.Request{RequestType: commands.SetCommand, Args: []string{"qwerty", "asdfgh"}},
			expectedErr:     nil,
		},

		{
			name: "del request with args more than needed",

			data: "DEL qwerty poiuy oiuy",

			expectedRequest: request.Request{RequestType: commands.DelCommand, Args: []string{"qwerty"}},
			expectedErr:     nil,
		},

		{
			name: "request with one arg",

			data: "DEL",

			expectedRequest: request.Request{RequestType: commands.IncorrectCommand},
			expectedErr:     errors.New("could not to parse less than two arguments"),
		},

		{
			name: "incorrect commands",

			data: "LOL boba",

			expectedRequest: request.Request{RequestType: commands.IncorrectCommand},
			expectedErr:     errors.New("incorrect command"),
		},

		{
			name: "set command with less than two args",

			data: "SET POP",

			expectedRequest: request.Request{RequestType: commands.SetCommand},
			expectedErr:     errors.New("set command has two arguments"),
		},

		{
			name: "request with end enter symbols",

			data: "set biba boba\r\n",

			expectedRequest: request.Request{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}},
			expectedErr:     nil,
		},
	}

	compute, _ := NewCompute(zap.NewNop())

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			request, err := compute.Parse(test.data)

			assert.Equal(t, test.expectedRequest, request)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_parseCommands(t *testing.T) {

	type testCase struct {
		name string

		stringCommand string

		expectedParsedCommand int
		expectedErr           error
	}

	testCases := []testCase{
		{
			name: "set command",

			stringCommand: "set",

			expectedParsedCommand: commands.SetCommand,
			expectedErr:           nil,
		},

		{
			name: "get command",

			stringCommand: "GET",

			expectedParsedCommand: commands.GetCommand,
			expectedErr:           nil,
		},

		{
			name: "del command",

			stringCommand: "deL",

			expectedParsedCommand: commands.DelCommand,
			expectedErr:           nil,
		},

		{
			name: "incorrect command",

			stringCommand: "barkhat",

			expectedParsedCommand: commands.IncorrectCommand,
			expectedErr:           errors.New("incorrect command"),
		},
	}

	compute, _ := NewCompute(zap.NewNop())

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {

			parsedCommand, err := compute.parseCommand(test.stringCommand)

			assert.Equal(t, test.expectedParsedCommand, parsedCommand)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_parseArguments(t *testing.T) {

	type testCase struct {
		name string

		command   int
		arguments []string

		expectedParsedArgs []string
		expectedErr        error
	}

	testCases := []testCase{
		{
			name: "correct set request",

			command:   commands.SetCommand,
			arguments: []string{"biba", "boba"},

			expectedParsedArgs: []string{"biba", "boba"},
			expectedErr:        nil,
		},

		{
			name: "correct get request",

			command:   commands.GetCommand,
			arguments: []string{"biba"},

			expectedParsedArgs: []string{"biba"},
			expectedErr:        nil,
		},

		{
			name: "correct del request",

			command:   commands.DelCommand,
			arguments: []string{"biba"},

			expectedParsedArgs: []string{"biba"},
			expectedErr:        nil,
		},

		{
			name: "correct del request with enter symbols",

			command:   commands.DelCommand,
			arguments: []string{"biba\r\n"},

			expectedParsedArgs: []string{"biba"},
			expectedErr:        nil,
		},

		{
			name: "uncorrect set request",

			command:   commands.SetCommand,
			arguments: []string{"biba\r\n"},

			expectedParsedArgs: nil,
			expectedErr:        errors.New("set command has two arguments"),
		},

		{
			name: "uncorrect set request with enter symbols as last args",

			command:   commands.SetCommand,
			arguments: []string{"biba", "\r\n"},

			expectedParsedArgs: nil,
			expectedErr:        errors.New("set command has two arguments"),
		},
	}

	compute, _ := NewCompute(zap.NewNop())

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {

			parsedArgs, err := compute.parseArguments(test.command, test.arguments)

			assert.Equal(t, test.expectedParsedArgs, parsedArgs)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
