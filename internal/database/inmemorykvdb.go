package database

import (
	"errors"
	"inmemorykvdb/internal/compute"
	"inmemorykvdb/internal/storage"

	"go.uber.org/zap"
)

type InMemoryKeyValueDatabase struct {
	compute compute.Compute
	storage storage.Storage
	logger  *zap.Logger
}

func NewInMemoryKvDb(compute compute.Compute, storage storage.Storage, logger *zap.Logger) (Database, error) {

	if compute == nil || storage == nil || logger == nil {
		return nil, errors.New("could not to create db without any of arguments")
	}

	return &InMemoryKeyValueDatabase{compute: compute, storage: storage, logger: logger}, nil
}

func (db *InMemoryKeyValueDatabase) Request(data string) (string, error) {

	db.logger.Debug("request started")

	db.logger.Debug("send data to compute")
	command, args, err := db.compute.Parse(data)
	db.logger.Debug("data parsed")

	if err != nil {
		db.logger.Error("data parsed with error")
		return "", err
	}

	db.logger.Debug("send request to storage")
	resp, err := db.storage.Request(command, args...)
	db.logger.Debug("storage returned a response")

	if err != nil {
		db.logger.Error("storage responsed with error")
		return "", err
	}

	db.logger.Debug("request to db is done")

	return resp, nil
}
