package wal

import (
	"errors"
	"fmt"
	"inmemorykvdb/internal/database/request"
	"time"

	"go.uber.org/zap"
)

type readingLayer interface {
	Read() ([][]byte, error)
}

type writingLayer interface {
	Write([]byte) (int, error)
}

const (
	defaultTickerTime = 10 * time.Millisecond
)

type WAL struct {
	BatchSize int
	Timeout   time.Duration

	ticker         *time.Ticker
	requestChannel chan request.Request
	blockChannel   chan struct{}

	writer writingLayer
	reader readingLayer

	batch *request.Batch

	logger *zap.Logger
}

func NewWal(logger *zap.Logger, walOptions ...WalOptions) (*WAL, error) {
	if logger == nil {
		return nil, errors.New("logger could not be nil")
	}

	wal := &WAL{}

	wal.logger = logger

	for _, option := range walOptions {
		option(wal)
	}

	wal.batch = request.NewBatch(wal.BatchSize)

	if wal.Timeout == 0 {
		wal.Timeout = defaultTickerTime
	}

	wal.ticker = time.NewTicker(wal.Timeout)

	wal.blockChannel = make(chan struct{})
	wal.requestChannel = make(chan request.Request)

	wal.startWAL()

	return wal, nil
}

func (w *WAL) startWAL() {
	if w.writer == nil {
		w.logger.Debug("could not start wal without write layer")
		return
	}

	w.logger.Debug("started handle events of wal")
	go w.handleEvents()
}

func (w *WAL) handleEvents() {
	for {
		select {
		case <-w.ticker.C:
			if w.batch.ByteSize != 0 {
				w.writeOnDisk()
			}

			w.ticker.Reset(w.Timeout)

		case request := <-w.requestChannel:
			w.batch.Add(&request)

			if w.batch.IsFilled() {
				w.writeOnDisk()
			}

			w.blockChannel <- struct{}{}
		}
	}
}

func (w *WAL) writeOnDisk() {
	w.logger.Debug("started write to disk")
	batchInBytes, err := w.batch.ParseBatch()
	w.batch.Clear()

	if err != nil {
		w.logger.Error(err.Error())
	} else {
		w.logger.Debug("parsing requests is complete")
	}

	count, err := w.writer.Write(batchInBytes)

	if err != nil {
		w.logger.Error(fmt.Sprintf("%s: written %d bytes", err.Error(), count))
	} else {
		w.logger.Debug("successful writed on disk")
	}
}

func (w *WAL) Write(req request.Request) {
	w.requestChannel <- req
	<-w.blockChannel
}

func (w *WAL) Read() *request.Batch {
	if w.reader == nil {
		w.logger.Debug("could not read without reader")
		return nil
	}

	w.logger.Debug("started read from wal")
	data, err := w.reader.Read()

	if err != nil {
		w.logger.Error(err.Error())
	}

	if len(data) == 0 {
		return nil
	}

	batch := request.NewBatch(len(data) * len((data)[0]))

	for _, arr := range data {
		err := batch.LoadData(arr)
		if err != nil {
			w.logger.Error(err.Error())
		}
	}

	return batch
}
