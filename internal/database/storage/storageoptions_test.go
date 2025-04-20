package storage

import (
	"inmemorykvdb/internal/database/storage/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_WithEngine(t *testing.T) {
	testEngine, _ := engine.NewInMemoryEngine(zap.NewNop())
	var actualStorage Storage

	expectedStorage := Storage{engine: testEngine}

	option := WithEngine(testEngine)

	option(&actualStorage)

	assert.Equal(t, expectedStorage, actualStorage)
}
