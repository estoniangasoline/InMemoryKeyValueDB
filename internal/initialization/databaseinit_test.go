package initialization

import (
	"errors"
	"inmemorykvdb/internal/database/compute"
	"inmemorykvdb/internal/database/storage"
	"inmemorykvdb/internal/database/storage/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_createDatabase(t *testing.T) {

	t.Parallel()

	type testCase struct {
		name string

		nilStorage bool
		nilCompute bool
		logger     *zap.Logger

		expectedNilObj bool
		expectedErr    error
	}

	testCases := []testCase{
		{
			name: "correct database",

			nilStorage: false,
			nilCompute: false,
			logger:     zap.NewNop(),

			expectedNilObj: false,
			expectedErr:    nil,
		},

		{
			name: "database without storage",

			nilStorage: true,
			nilCompute: false,
			logger:     zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("storage is nil"),
		},

		{
			name: "database without compute",

			nilStorage: false,
			nilCompute: true,
			logger:     zap.NewNop(),

			expectedNilObj: true,
			expectedErr:    errors.New("compute is nil"),
		},

		{
			name: "database without logger",

			nilStorage: false,
			nilCompute: false,
			logger:     nil,

			expectedNilObj: true,
			expectedErr:    errors.New("logger is nil"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			t.Parallel()

			var stor *storage.Storage
			var comp *compute.Compute

			if !test.nilStorage {
				eng, _ := engine.NewInMemoryEngine(zap.NewNop())
				stor, _ = storage.NewStorage(zap.NewNop(), eng)
			}

			if !test.nilCompute {
				comp, _ = compute.NewCompute(zap.NewNop())
			}

			db, err := createDatabase(stor, comp, test.logger)

			if test.expectedNilObj {
				assert.Nil(t, db)
			} else {
				assert.NotNil(t, db)
			}

			assert.Equal(t, test.expectedErr, err)
		})
	}
}
