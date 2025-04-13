package readlevel

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithFileMaxSizeRL(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		fileSize int

		expectedErr error
	}

	testCases := []testCase{
		{
			name: "correct file size",

			fileSize: 4096,

			expectedErr: nil,
		},
		{
			name: "zero file size",

			fileSize: 0,

			expectedErr: errors.New("max file size could not be a zero"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var rl readLevel

			option := WithFileMaxSize(test.fileSize)

			err := option(&rl)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
