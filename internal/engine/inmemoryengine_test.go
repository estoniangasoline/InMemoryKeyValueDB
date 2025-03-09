package engine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewInMemoryEngine(t *testing.T) {

	type testCase struct {
		testName string

		logger *zap.Logger
		size   uint

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			testName: "Correct engine",

			logger: &zap.Logger{},
			size:   1000,

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			testName: "Engine without logger",
			logger:   nil,
			size:     1000,

			expectedNilObj: true,
			expectedErr:    errors.New("engine without logger"),
		},

		{
			testName: "Engine with size greater than max",
			logger:   &zap.Logger{},
			size:     1000000,

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("engine could not be bigger than %d elements", maxEngineSize),
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			engine, err := NewInMemoryEngine(test.logger, test.size)

			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, engine)
			} else {
				assert.NotNil(t, engine)
			}
		})
	}
}

func Test_SetEngine(t *testing.T) {

	type testCase struct {
		testName string

		key   string
		value string

		engineSize uint

		expectedErr error
	}

	testCases := []testCase{
		{
			testName: "Set value in correct engine",

			key:   "Qwerty",
			value: "Asdfgh",

			engineSize: 1000,

			expectedErr: nil,
		},

		{
			testName: "Set value in uncorrect engine",

			key:   "Qwerty",
			value: "Asdfgh",

			engineSize: 0,

			expectedErr: errors.New("engine is empty"),
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop(), test.engineSize)

			err := engine.SET(test.key, test.value)

			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func Test_GetEngine(t *testing.T) {

	type testCase struct {
		name string

		setKey string
		getKey string

		setValue string

		engineSize int

		expectedValue string
		expectedErr   error
	}

	testCases := []testCase{
		{
			name: "Get value from equal setkey and getkey",

			setKey: "Asdfgh",
			getKey: "Asdfgh",

			setValue:      "Qwerty",
			expectedValue: "Qwerty",

			engineSize: 1000,

			expectedErr: nil,
		},

		{
			name: "Get value from not equal setkey and getkey",

			setKey: "Asdfgh",
			getKey: "Zxcvbn",

			setValue:      "Qwerty",
			expectedValue: "",

			engineSize: 1000,

			expectedErr: errors.New("value not found"),
		},

		{
			name: "Get value from empty engine",

			getKey:        "Zxcvbn",
			expectedValue: "",

			engineSize: 0,

			expectedErr: errors.New("engine is empty"),
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop(), uint(test.engineSize))

			if test.engineSize < maxEngineSize {
				engine.SET(test.setKey, test.setValue)
			}

			actualValue, actualErr := engine.GET(test.getKey)

			assert.Equal(t, actualValue, test.expectedValue)
			assert.Equal(t, actualErr, test.expectedErr)
		})
	}
}

func Test_DelEngine(t *testing.T) {

	type testCase struct {
		name string

		setKey   string
		setValue string

		delKey string

		engineSize int

		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Delete key that is in engine",

			setKey:   "Qwerty",
			setValue: "ASDFG",

			delKey: "Qwerty",

			engineSize: 1000,

			expectedErr: nil,
		},

		{
			name: "Delete key that is not in engine",

			setKey:   "Qwerty",
			setValue: "ASDFG",

			delKey: "Poiu",

			engineSize: 1000,

			expectedErr: nil,
		},

		{
			name: "Delete key from empty engine",

			setKey:   "Qwerty",
			setValue: "ASDFG",

			delKey: "Poiu",

			engineSize: 0,

			expectedErr: errors.New("engine is empty"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop(), uint(test.engineSize))

			if test.engineSize < maxEngineSize {
				engine.SET(test.setKey, test.setValue)
			}

			actualErr := engine.DEL(test.delKey)

			assert.Equal(t, actualErr, test.expectedErr)
		})
	}
}
