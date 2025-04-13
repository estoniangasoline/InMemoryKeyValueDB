package database

import (
	"errors"
	"inmemorykvdb/internal/database/request"

	"go.uber.org/zap"
)

type computeLayer interface {
	Parse(data string) (request.Request, error)
}

type storageLayer interface {
	HandleRequest(request.Request) (string, error)
}

type InMemoryKeyValueDatabase struct {
	compute computeLayer
	storage storageLayer
	logger  *zap.Logger
}

func NewInMemoryKvDb(compute computeLayer, storage storageLayer, logger *zap.Logger) (*InMemoryKeyValueDatabase, error) {

	if compute == nil || storage == nil || logger == nil {
		return nil, errors.New("could not to create db without any of arguments")
	}

	return &InMemoryKeyValueDatabase{compute: compute, storage: storage, logger: logger}, nil
}

func (db *InMemoryKeyValueDatabase) HandleRequest(data string) (string, error) {

	db.logger.Debug("request started, send data to compute")

	req, err := db.compute.Parse(data)
	db.logger.Debug("data parsed")

	if err != nil {
		db.logger.Error("data parsed with error")
		return "", err
	}

	db.logger.Debug("send request to storage")
	resp, err := db.storage.HandleRequest(req)
	db.logger.Debug("storage returned a response")

	if err != nil {
		db.logger.Error("storage responsed with error")
		return "", err
	}

	db.logger.Debug("request to db is done")

	return resp, nil
}
