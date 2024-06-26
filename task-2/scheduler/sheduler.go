package scheduler

import (
	"sync"
	"time"
)

type Scheduler struct {
	m            *sync.Mutex
	currentJobId int
	jobsCancel   map[int]chan bool
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		currentJobId: 0,
		jobsCancel:   make(map[int]chan bool),
		m:            new(sync.Mutex),
	}
	return s
}

func (s *Scheduler) SetTimeout(job Job, delay time.Duration) int {
	s.m.Lock()
	id := s.currentJobId
	s.currentJobId++
	s.jobsCancel[id] = make(chan bool)
	s.m.Unlock()

	go func(job Job, delay time.Duration, cancelCh chan bool) {
		timer := time.NewTimer(delay)
		for {
			select {
			case <-cancelCh:
				timer.Stop()
				job.Cancel()
				break
			case <-timer.C:
				job.Do()
			}
		}
	}(job, delay, s.jobsCancel[id])

	return id
}

func (s *Scheduler) CancelTimeout(id int) {
	s.m.Lock()
	defer s.m.Unlock()
	if ch, ok := s.jobsCancel[id]; ok {
		ch <- true
	}
}
