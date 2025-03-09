package compute

import (
	"errors"
	"inmemorykvdb/internal/commands"
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

			compute, err := NewSimpleCompute(test.logger)

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

		expectedCommand   int
		expectedArguments []string
		expectedErr       error
	}

	testCases := []testCase{

		{
			name: "correct get request",

			data: "get qwerty",

			expectedCommand:   commands.GetCommand,
			expectedArguments: []string{"qwerty"},
			expectedErr:       nil,
		},

		{
			name: "correct set request",

			data: "set qwerty asdfgh",

			expectedCommand:   commands.SetCommand,
			expectedArguments: []string{"qwerty", "asdfgh"},
			expectedErr:       nil,
		},

		{
			name: "correct del request",

			data: "del qwerty",

			expectedCommand:   commands.DelCommand,
			expectedArguments: []string{"qwerty"},
			expectedErr:       nil,
		},

		{
			name: "get request in uppercase",

			data: "GET qwerty",

			expectedCommand:   commands.GetCommand,
			expectedArguments: []string{"qwerty"},
			expectedErr:       nil,
		},

		{
			name: "set request in uppercase",

			data: "SET qwerty asdfgh",

			expectedCommand:   commands.SetCommand,
			expectedArguments: []string{"qwerty", "asdfgh"},
			expectedErr:       nil,
		},

		{
			name: "del request in uppercase",

			data: "DEL qwerty",

			expectedCommand:   commands.DelCommand,
			expectedArguments: []string{"qwerty"},
			expectedErr:       nil,
		},

		{
			name: "correct del request",

			data: "del qwerty",

			expectedCommand:   commands.DelCommand,
			expectedArguments: []string{"qwerty"},
			expectedErr:       nil,
		},

		{
			name: "get request with args more than needed",

			data: "GET qwerty poiuy",

			expectedCommand:   commands.GetCommand,
			expectedArguments: []string{"qwerty"},
			expectedErr:       nil,
		},

		{
			name: "set request with args more than needed",

			data: "SET qwerty asdfgh [poiuyt]",

			expectedCommand:   commands.SetCommand,
			expectedArguments: []string{"qwerty", "asdfgh"},
			expectedErr:       nil,
		},

		{
			name: "del request with args more than needed",

			data: "DEL qwerty poiuy oiuy",

			expectedCommand:   commands.DelCommand,
			expectedArguments: []string{"qwerty"},
			expectedErr:       nil,
		},

		{
			name: "request with one arg",

			data: "DEL",

			expectedCommand:   commands.UncorrectCommand,
			expectedArguments: []string{},
			expectedErr:       errors.New("could not to parse less than two arguments"),
		},

		{
			name: "unknown commands",

			data: "LOL boba",

			expectedCommand:   commands.UncorrectCommand,
			expectedArguments: []string{},
			expectedErr:       errors.New("unknown command"),
		},

		{
			name: "set command with less than two args",

			data: "SET POP",

			expectedCommand:   commands.UncorrectCommand,
			expectedArguments: []string{},
			expectedErr:       errors.New("set command could not be realized with one argument"),
		},

		{
			name: "request with end trash",

			data: "set biba boba\r\n",

			expectedCommand:   commands.SetCommand,
			expectedArguments: []string{"biba", "boba"},
			expectedErr:       nil,
		},
	}

	compute, _ := NewSimpleCompute(zap.NewNop())

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			command, args, err := compute.Parse(test.data)

			assert.Equal(t, test.expectedCommand, command)
			assert.Equal(t, test.expectedArguments, args)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
