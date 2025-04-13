package wal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WithBatchSize(t *testing.T) {
	t.Parallel()

	testBatchSize := 100
	var actualWal WAL

	expectedWal := WAL{BatchSize: testBatchSize}

	option := WithBatchSize(testBatchSize)
	option(&actualWal)

	assert.Equal(t, expectedWal, actualWal)
}

func Test_WithBatchTimeout(t *testing.T) {
	t.Parallel()

	var testBatchTimeout time.Duration = 30
	var actualWal WAL

	expectedWal := WAL{Timeout: testBatchTimeout}

	option := WithBatchTimeout(testBatchTimeout)
	option(&actualWal)

	assert.Equal(t, expectedWal, actualWal)
}
