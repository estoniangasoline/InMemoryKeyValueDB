package serversync

import (
	"errors"
	"sync"
)

type Semaphore struct {
	maxCount int
	curCount int
	cond     sync.Cond
}

func NewSemaphore(count int) (*Semaphore, error) {
	if count == 0 {
		return nil, errors.New("semaphore work with count more than one")
	}

	return &Semaphore{maxCount: count, curCount: 0, cond: *sync.NewCond(&sync.Mutex{})}, nil
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	for s.curCount >= s.maxCount {
		s.cond.Wait()
	}
	s.curCount++
	s.cond.L.Unlock()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	s.curCount--
	s.cond.Signal()
	s.cond.L.Unlock()
}
