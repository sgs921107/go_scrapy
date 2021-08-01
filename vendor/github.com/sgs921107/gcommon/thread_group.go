package gcommon

import (
	"time"
)

// Group wait group
type Group interface {
	Size() int
	Len() int
	Add(num int)
	Done()
	Wait(inrterval time.Duration, timeout ...time.Duration)
}

// ThreadGroup thread group
type ThreadGroup struct {
	size      int
	groupChan chan interface{}
}

func (p *ThreadGroup) init() {
	if p.size == 0 {
		p.size = 1
	}
	p.groupChan = make(chan interface{}, p.size)
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
		p.groupChan <- 1
	}
}

// Done done
func (p *ThreadGroup) Done() {
	<-p.groupChan
}

// Wait wait
// interval 指定间隔多久检测一次是否所有任务都已完成 默认：100纳秒
// timeout为可选参数，如果大于0则等待指定的时间后退出
func (p *ThreadGroup) Wait(interval time.Duration, timeout ...time.Duration) {
	if interval == 0 {
		interval = 100
	}
	if len(timeout) == 0 || timeout[0] == 0 {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				if len(p.groupChan) == 0 {
					goto Break
				}
			}
		}
	} else {
		ticker := time.NewTicker(interval)
		timer := time.NewTimer(timeout[0])
		for {
			select {
			case <-ticker.C:
				if len(p.groupChan) == 0 {
					goto Break
				}
			case <-timer.C:
				goto Break
			}
		}
	}
Break:
	close(p.groupChan)
}

// NewThreadGroup new a thread group
func NewThreadGroup(size uint) Group {
	group := &ThreadGroup{
		size: int(size),
	}
	group.init()
	return group
}