package concurrency

import "sync"

type mut interface {
	*sync.Mutex | *sync.RWMutex
	Lock()
	Unlock()
}

func WithLock[m mut](mutex m, action func()) {
	defer mutex.Unlock()
	mutex.Lock()
	action()
}

func WithRLock(mutex *sync.RWMutex, action func()) {
	defer mutex.RUnlock()
	mutex.RLock()

	action()
}
