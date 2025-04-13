package parsing

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseSize(t *testing.T) {
	type testCase struct {
		name string

		unparsedSize string

		expectedSize int
		expectedErr  error
	}

	testCases := []testCase{
		{
			name: "correct message size",

			unparsedSize: "10MB",

			expectedSize: megabyteMultiply * 10,
			expectedErr:  nil,
		},
		{
			name: "correct message size with bytes",

			unparsedSize: "10B",

			expectedSize: byteMultiply * 10,
			expectedErr:  nil,
		},
		{
			name: "incorrect message size with bytes",

			unparsedSize: "bMB",

			expectedSize: -1,
			expectedErr:  errors.New("incorrect integer in string"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()

			size, err := ParseSize(test.unparsedSize)

			assert.Equal(t, test.expectedSize, size)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_parseMultiply(t *testing.T) {
	type testCase struct {
		name string

		unparsedMultiply string

		expectedMultiply int
		expectedErr      error
	}

	testCases := []testCase{
		{
			name: "megabyte",

			unparsedMultiply: "MB",

			expectedMultiply: megabyteMultiply,
			expectedErr:      nil,
		},
		{
			name: "kilobyte",

			unparsedMultiply: "KB",

			expectedMultiply: kilobyteMultiply,
			expectedErr:      nil,
		},
		{
			name: "byte",

			unparsedMultiply: "B",

			expectedMultiply: byteMultiply,
			expectedErr:      nil,
		},

		{
			name: "incorrect multiply",

			unparsedMultiply: "Q",

			expectedMultiply: -1,
			expectedErr:      errors.New("unknown or forbidden buffer size"),
		},

		{
			name: "empty string",

			unparsedMultiply: "",

			expectedMultiply: -1,
			expectedErr:      errors.New("empty string"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mult, err := ParseMultiply(test.unparsedMultiply)

			assert.Equal(t, test.expectedMultiply, mult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
