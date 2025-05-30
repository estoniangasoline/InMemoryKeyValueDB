package storage

import (
	"inmemorykvdb/internal/database/request"
	"inmemorykvdb/internal/database/storage/replication"
	"inmemorykvdb/internal/database/storage/wal"
	"inmemorykvdb/internal/network"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_WithWal(t *testing.T) {
	testWal, _ := wal.NewWal(zap.NewNop())

	expectedStorage := &Storage{wal: testWal}

	option := WithWal(testWal)

	actualStorage := &Storage{}

	option(actualStorage)

	assert.Equal(t, expectedStorage, actualStorage)
}

func Test_WithReplica(t *testing.T) {
	logger := zap.NewNop()
	server, _ := network.NewServer(":8080", logger)
	master, _ := replication.NewMaster(server, logger)

	option := WithReplica(master)

	stor := &Storage{}

	option(stor)

	expectedStor := &Storage{replica: master}

	assert.Equal(t, expectedStor, stor)
}

func Test_WithDataChan(t *testing.T) {
	dataChan := make(chan *request.Batch)

	option := WithDataChan(dataChan)

	stor := &Storage{}

	option(stor)

	expectedStor := &Storage{dataChan: dataChan}

	assert.Equal(t, expectedStor, stor)
}
