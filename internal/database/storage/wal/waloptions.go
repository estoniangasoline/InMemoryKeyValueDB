package wal

import (
	"time"
)

type WalOptions func(*WAL)

func WithBatchSize(batchSize int) WalOptions {
	return func(w *WAL) {
		w.BatchSize = batchSize
	}
}

func WithBatchTimeout(batchTimeout time.Duration) WalOptions {
	return func(w *WAL) {
		w.Timeout = batchTimeout
	}
}
