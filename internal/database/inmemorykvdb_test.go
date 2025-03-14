package database

import (
	"errors"
	"inmemorykvdb/internal/database/compute"
	"inmemorykvdb/internal/database/storage"
	"inmemorykvdb/internal/database/storage/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_NewInMemoryKvDb(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		nilCompute bool
		nilStorage bool
		logger     *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct db",

			nilCompute: false,
			nilStorage: false,
			logger:     zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "db without compute",

			nilCompute: true,
			nilStorage: false,
			logger:     zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("could not to create db without any of arguments"),
		},

		{
			name: "db without storage",

			nilCompute: false,
			nilStorage: true,
			logger:     zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("could not to create db without any of arguments"),
		},

		{
			name: "db without logger",

			nilCompute: false,
			nilStorage: false,
			logger:     nil,

			expectedNilObj: true,
			expectedErr:    errors.New("could not to create db without any of arguments"),
		},

		{
			name: "db without anything",

			nilCompute: true,
			nilStorage: true,
			logger:     nil,

			expectedNilObj: true,
			expectedErr:    errors.New("could not to create db without any of arguments"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			var stor storageLayer
			var comp computeLayer

			if !test.nilStorage {
				eng, _ := engine.NewInMemoryEngine(zap.NewNop())
				stor, _ = storage.NewStorage(eng, zap.NewNop())
			}

			if !test.nilCompute {
				comp, _ = compute.NewCompute(zap.NewNop())
			}

			db, err := NewInMemoryKvDb(comp, stor, test.logger)

			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, db)
			} else {
				assert.NotNil(t, db)
			}
		})
	}
}

func Test_Request(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		data string

		expectedResp string
		expectedErr  error
	}

	testCases := []testCase{
		{
			name: "correct set request",

			data: "set biba boba",

			expectedResp: "",
			expectedErr:  nil,
		},
		{
			name: "get request with mistake",

			data: "get bibo",

			expectedResp: "",
			expectedErr:  errors.New("value not found"),
		},
		{
			name: "correct get request",

			data: "GET biba",

			expectedResp: "boba",
			expectedErr:  nil,
		},
		{
			name: "correct del request",

			data: "deL biba",

			expectedResp: "",
			expectedErr:  nil,
		},
		{
			name: "another correct set request",

			data: "set boba biba",

			expectedResp: "",
			expectedErr:  nil,
		},
		{
			name: "get request to check is deleted",

			data: "GEt biba",

			expectedResp: "",
			expectedErr:  errors.New("value not found"),
		},
		{
			name: "del request to check deleting with no errors also deleted record",

			data: "del biba",

			expectedResp: "",
			expectedErr:  nil,
		},
		{
			name: "incorrect request",

			data: "yo yo",

			expectedResp: "",
			expectedErr:  errors.New("incorrect command"),
		},
	}

	eng, _ := engine.NewInMemoryEngine(zap.NewNop())
	stor, _ := storage.NewStorage(eng, zap.NewNop())

	comp, _ := compute.NewCompute(zap.NewNop())

	db, _ := NewInMemoryKvDb(comp, stor, zap.NewNop())

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			resp, err := db.HandleRequest(test.data)

			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
