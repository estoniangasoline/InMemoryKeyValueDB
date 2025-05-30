package replication

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WithDirectoryMaster(t *testing.T) {
	dir := "/biba"
	expectedMaster := &Master{directory: dir}

	actualMaster := &Master{}

	option := WithDirectoryMaster(dir)

	option(actualMaster)

	assert.Equal(t, expectedMaster, actualMaster)
}

func Test_WithDirectorySlave(t *testing.T) {
	dir := "/biba"
	expectedSlave := &Slave{directory: dir}

	actualSlave := &Slave{}

	option := WithDirectorySlave(dir)

	option(actualSlave)

	assert.Equal(t, expectedSlave, actualSlave)
}

func Test_WithInterval(t *testing.T) {
	type testCase struct {
		name string

		option SlaveOption

		expectedSlave *Slave
		expectedErr   error
	}

	testCases := []testCase{
		{
			name: "correct option",

			option: WithInterval(time.Second * 10),

			expectedSlave: &Slave{requestInterval: time.Second * 10},
			expectedErr:   nil,
		},
		{
			name: "zero interval",

			option: WithInterval(0),

			expectedSlave: &Slave{},
			expectedErr:   errors.New("time interval could not be a zero"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			slave := &Slave{}

			option := test.option

			err := option(slave)

			assert.Equal(t, test.expectedSlave, slave)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
