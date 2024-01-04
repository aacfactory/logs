package logs

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrCloseTimeout = errors.New("logs: close time out")
)

func newSink(minLevel Level, discardLevel Level, consumes int, buffer int, sendTimeout time.Duration, shutdownTimeout time.Duration, writers []Writer) (sink *Sink) {
	if minLevel == 0 {
		minLevel = DebugLevel
	}
	if discardLevel == 0 {
		discardLevel = DebugLevel
	}
	if consumes < 1 {
		consumes = runtime.NumCPU()
	}
	if buffer < 1 {
		buffer = 40960
	}
	if sendTimeout < 1 {
		sendTimeout = time.Duration(10) * time.Microsecond
	}
	if shutdownTimeout < 1 {
		shutdownTimeout = time.Duration(10) * time.Second
	}
	if len(writers) == 0 {
		writers = append(writers, NewConsoleWriter(TextFormatter, StdOut))
	}
	sink = &Sink{
		running:         atomic.Bool{},
		consumes:        consumes,
		minLevel:        minLevel,
		discardLevel:    discardLevel,
		entries:         make(chan Entry, buffer),
		timeouts:        make(chan Entry, buffer),
		writers:         writers,
		count:           new(sync.WaitGroup),
		timers:          NewTimers(),
		sendTimeout:     sendTimeout,
		shutdownTimeout: shutdownTimeout,
	}
	sink.Listen()
	return
}

type Sink struct {
	running         atomic.Bool
	consumes        int
	minLevel        Level
	discardLevel    Level
	entries         chan Entry
	timeouts        chan Entry
	writers         []Writer
	count           *sync.WaitGroup
	timers          *Timers
	sendTimeout     time.Duration
	shutdownTimeout time.Duration
}

func (sink *Sink) Emit(entry Entry) {
	if !sink.running.Load() {
		return
	}
	if entry.Level < sink.minLevel {
		return
	}
	timer := sink.timers.Get(sink.sendTimeout)
	sink.count.Add(1)
	ok := false
	select {
	case sink.entries <- entry:
		ok = true
		break
	case <-timer.C:
		break
	}
	sink.timers.Put(timer)
	if ok {
		return
	}
	if entry.Level <= sink.discardLevel {
		sink.count.Done()
		return
	}
	timer = sink.timers.Get(sink.sendTimeout)
	select {
	case sink.timeouts <- entry:
		ok = true
		break
	case <-timer.C:
		break
	}
	sink.timers.Put(timer)
	if ok {
		return
	}
	go func(timeout chan Entry, entry Entry) {
		timeout <- entry
	}(sink.timeouts, entry)
	return
}

func (sink *Sink) Listen() {
	for i := 0; i < sink.consumes; i++ {
		consume(sink.entries, sink.count, sink.writers...)
	}
	go func(entries chan Entry, timeouts chan Entry) {
		for {
			entry, ok := <-timeouts
			if !ok {
				break
			}
			entries <- entry
		}
	}(sink.entries, sink.timeouts)
	sink.running.Store(true)
	return
}

func (sink *Sink) Shutdown(ctx context.Context) (err error) {
	if !sink.running.Load() {
		return
	}
	tc, cancel := context.WithTimeout(ctx, sink.shutdownTimeout)
	defer cancel()
	sink.running.Store(false)
	select {
	case <-tc.Done():
		err = ErrCloseTimeout
		return
	case <-sink.wait():
		return
	}
}

func (sink *Sink) wait() (ch <-chan struct{}) {
	c := make(chan struct{}, 1)
	go func(c chan struct{}, sink *Sink) {
		sink.count.Wait()
		c <- struct{}{}
		close(c)
		close(sink.entries)
		close(sink.timeouts)
		for _, writer := range sink.writers {
			_ = writer.Close()
		}
	}(c, sink)
	ch = c
	return
}

func consume(ch <-chan Entry, count *sync.WaitGroup, writers ...Writer) {
	go func(ch <-chan Entry, count *sync.WaitGroup, writers []Writer) {
		for {
			entry, ok := <-ch
			if !ok {
				break
			}
			for _, writer := range writers {
				writer.Write(entry)
			}
			count.Done()
		}
	}(ch, count, writers)
}
