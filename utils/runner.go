package utils

import (
	"sync"
)

type Runner struct {
	channelRun chan bool
	close      chan bool
	group      sync.WaitGroup
	listIn     Array[any]
	fn         func(any)
}

func NewRunner(fn func(a any)) *Runner {
	r := &Runner{
		channelRun: make(chan bool),
		close:      make(chan bool),
		group:      sync.WaitGroup{},
		listIn:     Array[any]{},
		fn:         fn,
	}
	return r
}

func (s *Runner) SetRunner(fn func(a any)) {
	s.fn = fn
}

func (s *Runner) Push(data any) {
	s.group.Add(1)
	s.listIn.Push(data)
	go func() {
		s.channelRun <- true
	}()
}

func (s *Runner) PushList(data ...any) {
	n := len(data)
	s.group.Add(n)
	s.listIn.Push(data...)
	go func() {
		for i := 0; i < n; i++ {
			s.channelRun <- true
		}
	}()
}

func (s *Runner) Close() {
	close(s.channelRun)
	close(s.close)
}

func (s *Runner) Wait() {
	s.group.Wait()
	s.close <- true
}

func (s *Runner) Start() {
	go func() {
		for {
			select {
			case <-s.close:
				return
			case <-s.channelRun:
				v := s.listIn.Shift()
				s.fn(*v)
				s.group.Done()
			}
		}
	}()
}
