package request

import (
	"errors"
	"inmemorykvdb/internal/database/commands"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseToBytes(t *testing.T) {
	type testCase struct {
		name string

		request *Request

		expectedArray []byte
		expectedErr   error
	}

	testCases := []testCase{
		{
			name: "correct request",

			request: &Request{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}},

			expectedArray: []byte("SET" + DelimElement + "biba" + DelimElement + "boba" + EndElement),
			expectedErr:   nil,
		},
		{
			name: "incorrect command type",

			request: &Request{RequestType: -1, Args: []string{"biba", "boba"}},

			expectedArray: []byte(nil),
			expectedErr:   errors.New("incorrect command type"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			arr, err := test.request.ParseToBytes()

			assert.Equal(t, test.expectedArray, arr)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_NewRequest(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		data string

		expectedReq *Request
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "set data",

			data: "SET biba boba\n",

			expectedReq: &Request{RequestType: commands.SetCommand, Args: []string{"biba", "boba"}},
			expectedErr: nil,
		},
		{
			name: "del data",

			data: "DEL biba\n",

			expectedReq: &Request{RequestType: commands.DelCommand, Args: []string{"biba"}},
			expectedErr: nil,
		},
		{
			name: "incorrect data",

			data: "DEL\n",

			expectedReq: nil,
			expectedErr: errors.New("incorrect data"),
		},
		{
			name: "set command with less than two arguments",

			data: "SET BUBA\n",

			expectedReq: nil,
			expectedErr: errors.New("set command in data has less than two arguments"),
		},
		{
			name: "incorrect command",

			data: "LOL biba boba\n",

			expectedReq: nil,
			expectedErr: errors.New("incorrect command"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req, err := NewRequest(test.data)

			assert.Equal(t, test.expectedReq, req)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
