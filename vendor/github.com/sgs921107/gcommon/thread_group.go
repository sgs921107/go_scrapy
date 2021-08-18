package gcommon

import (
	"sync"
)

// Group wait group
type Group interface {
	Size() int
	Len() int
	Add(num int)
	Done()
	Wait()
}

// ThreadGroup thread group
type ThreadGroup struct {
	size      int
	groupChan chan struct{}
	wg      *sync.WaitGroup
}

func (p *ThreadGroup) init() {
	if p.size <= 0 {
		p.size = 1
	}
	p.groupChan = make(chan struct{}, p.size)
	p.wg = &sync.WaitGroup{}
}

// Size size
func (p *ThreadGroup) Size() int {
	return p.size
}

// Len len
func (p *ThreadGroup) Len() int {
	return len(p.groupChan)
}

// Add add thread num
func (p *ThreadGroup) Add(num int) {
	for i := 0; i < num; i++ {
		p.groupChan <- struct{}{}
		p.wg.Add(1)
	}
}

// Done done
func (p *ThreadGroup) Done() {
	<-p.groupChan
	p.wg.Add(-1)
}

// Wait wait
func (p *ThreadGroup) Wait() {
	p.wg.Wait()
	close(p.groupChan)
}

// NewThreadGroup new a thread group
func NewThreadGroup(size int) Group {
	group := &ThreadGroup{
		size: size,
	}
	group.init()
	return group
}