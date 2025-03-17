package engine

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewInMemoryEngine(t *testing.T) {

	type testCase struct {
		testName string

		logger *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			testName: "Correct engine",

			logger: &zap.Logger{},

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			testName: "Engine without logger",
			logger:   nil,

			expectedNilObj: true,
			expectedErr:    errors.New("engine without logger"),
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			engine, err := NewInMemoryEngine(test.logger)

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

		expectedErr error
	}

	testCases := []testCase{
		{
			testName: "Set value in correct engine",

			key:   "Qwerty",
			value: "Asdfgh",

			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop())

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

			expectedErr: nil,
		},

		{
			name: "Get value from not equal setkey and getkey",

			setKey: "Asdfgh",
			getKey: "Zxcvbn",

			setValue:      "Qwerty",
			expectedValue: "",

			expectedErr: errors.New("value not found"),
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop())

			engine.SET(test.setKey, test.setValue)

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

		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Delete key that is in engine",

			setKey:   "Qwerty",
			setValue: "ASDFG",

			delKey: "Qwerty",

			expectedErr: nil,
		},

		{
			name: "Delete key that is not in engine",

			setKey:   "Qwerty",
			setValue: "ASDFG",

			delKey: "Poiu",

			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop())

			engine.SET(test.setKey, test.setValue)

			actualErr := engine.DEL(test.delKey)

			assert.Equal(t, actualErr, test.expectedErr)
		})
	}
}
