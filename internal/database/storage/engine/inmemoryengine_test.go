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

		logger  *zap.Logger
		options []EngineOption

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			testName: "Correct engine",

			logger: zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			testName: "Correct engine with options",

			logger:  zap.NewNop(),
			options: []EngineOption{WithPartitions(10, 10)},

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

	engine, _ := NewInMemoryEngine(zap.NewNop())

	key := "key"
	value := "value"

	engine.SET(key, value)

	actual, found := engine.GET(key)

	assert.Equal(t, value, actual)
	assert.True(t, found)
}

func Test_GetEngine(t *testing.T) {

	type testCase struct {
		name string

		setKey string
		getKey string

		setValue string

		expectedValue string
		isFound       bool
	}

	testCases := []testCase{
		{
			name: "Get value from equal setkey and getkey",

			setKey: "Asdfgh",
			getKey: "Asdfgh",

			setValue:      "Qwerty",
			expectedValue: "Qwerty",

			isFound: true,
		},

		{
			name: "Get value from not equal setkey and getkey",

			setKey: "Asdfgh",
			getKey: "Zxcvbn",

			setValue:      "Qwerty",
			expectedValue: "",

			isFound: false,
		},
	}

	for _, test := range testCases {

		t.Run(test.name, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop())

			engine.SET(test.setKey, test.setValue)

			actualValue, actualFound := engine.GET(test.getKey)

			assert.Equal(t, actualValue, test.expectedValue)
			assert.Equal(t, actualFound, test.isFound)
		})
	}
}

func Test_DelEngine(t *testing.T) {

	type testCase struct {
		name string

		setKey   string
		setValue string

		delKey string
	}

	testCases := []testCase{
		{
			name: "Delete key that is in engine",

			setKey:   "Qwerty",
			setValue: "ASDFG",

			delKey: "Qwerty",
		},

		{
			name: "Delete key that is not in engine",

			setKey:   "Qwerty",
			setValue: "ASDFG",

			delKey: "Poiu",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			engine, _ := NewInMemoryEngine(zap.NewNop())

			engine.SET(test.setKey, test.setValue)

			engine.DEL(test.delKey)

			_, found := engine.GET(test.delKey)

			assert.False(t, found)
		})
	}
}
