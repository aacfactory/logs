package logs

import (
	"sync"
	"time"
)

func NewTimers() *Timers {
	return &Timers{
		pool: sync.Pool{},
	}
}

type Timers struct {
	pool sync.Pool
}

func (timers *Timers) Get(d time.Duration) (timer *time.Timer) {
	v := timers.pool.Get()
	if v == nil {
		timer = time.NewTimer(d)
		return
	}
	timer = v.(*time.Timer)
	timer.Reset(d)
	return
}

func (timers *Timers) Put(timer *time.Timer) {
	timers.pool.Put(timer)
}
