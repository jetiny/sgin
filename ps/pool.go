package ps

import (
	"sync"

	"github.com/pkg/errors"
)

// Pool creates multiple stay open exiftool instances and spreads the work
// across them with a simple round robin distribution.
type Pool struct {
	sync.Mutex
	stayopens []*Stayopen
	c         int
	l         int
	stopped   bool
}

func (p *Pool) Execute(flags ...string) ([]byte, error) {
	if p.stopped {
		return nil, errors.New("Stopped")
	}
	p.Lock()
	p.c++
	key := p.c % p.l
	p.Unlock()
	return p.stayopens[key].execute(flags...)
}

func (p *Pool) Stop() {
	p.Lock()
	defer p.Unlock()
	for _, s := range p.stayopens {
		s.Stop()
	}
	p.stopped = true
}

func NewPool(h EventHandler, cmd string, num int, flags ...string) (*Pool, error) {
	p := &Pool{
		stayopens: make([]*Stayopen, num),
		l:         num,
	}
	var err error
	for i := 0; i < num; i++ {
		p.stayopens[i], err = NewStayOpen(h, cmd, flags...)
		if err != nil {
			return nil, errors.Wrap(err, "Could not create StayOpen")
		}
	}
	return p, nil
}
