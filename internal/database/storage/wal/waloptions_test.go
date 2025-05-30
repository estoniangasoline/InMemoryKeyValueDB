package wal

import (
	"inmemorykvdb/internal/database/storage/wal/readlevel"
	"inmemorykvdb/internal/database/storage/wal/writelevel"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func Test_WithWriter(t *testing.T) {
	t.Parallel()

	dir := "C:/go/InMemoryKeyValueDB/test/wal/options/"

	writer, _ := writelevel.NewWriteLevel(zap.NewNop(), writelevel.WithFilePath(dir))

	expectedWal := &WAL{writer: writer}

	wal := &WAL{}

	option := WithWriter(writer)
	option(wal)

	assert.Equal(t, expectedWal, wal)
}

func Test_WithReader(t *testing.T) {
	t.Parallel()

	dir := "C:/go/InMemoryKeyValueDB/test/wal/options/"

	reader, _ := readlevel.NewReadLevel(zap.NewNop(), "wal", readlevel.WithDirectory(dir))

	expectedWal := &WAL{reader: reader}

	wal := &WAL{}

	option := WithReader(reader)
	option(wal)

	assert.Equal(t, expectedWal, wal)
}
