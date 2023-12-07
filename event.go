package logs

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Event interface {
	Message(message string)
	MessageF(format string, a ...any)
	Cause(err error) Event
	Caller() Event
	CallerWithSkip(skip int) Event
	With(key string, value any) Event
}

var (
	empty Event = &emptyEvent{}
)

type emptyEvent struct {
}

func (e *emptyEvent) Message(message string) {
}

func (e *emptyEvent) MessageF(format string, a ...any) {
}

func (e *emptyEvent) Cause(err error) Event {
	return e
}

func (e *emptyEvent) Caller() Event {
	return e
}

func (e *emptyEvent) CallerWithSkip(skip int) Event {
	return e
}

func (e *emptyEvent) With(key string, value any) Event {
	return e
}

type event struct {
	sink  *Sink
	reset func(event *event)
	entry Entry
}

func (e *event) Message(message string) {
	e.entry.Message = message
	e.sink.Emit(e.entry)
	e.reset(e)
}

func (e *event) MessageF(format string, a ...any) {
	e.entry.Message = fmt.Sprintf(format, a...)
	e.sink.Emit(e.entry)
	e.reset(e)
}

func (e *event) Cause(err error) Event {
	if err == nil {
		return e
	}
	cause, ok := err.(Error)
	if ok {
		e.entry.Cause = cause
		return e
	}
	e.entry.Cause = jsonError{
		cause: err,
	}
	return e
}

func (e *event) Caller() Event {
	return e.CallerWithSkip(2)
}

func (e *event) CallerWithSkip(skip int) Event {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return e
	}
	if strings.IndexByte(file, '/') == 0 || strings.IndexByte(file, ':') == 1 {
		idx := strings.Index(file, "/src/")
		if idx > 0 {
			file = file[idx+5:]
		} else {
			idx = strings.Index(file, "/pkg/mod/")
			if idx > 0 {
				file = file[idx+9:]
			}
		}
	}
	fn := runtime.FuncForPC(pc)
	e.entry.Caller.File = file
	e.entry.Caller.Line = line
	e.entry.Caller.Fn = fn.Name()
	return e
}

func (e *event) With(key string, value any) Event {
	e.entry.Fields = e.entry.Fields.Add(key, value)
	return e
}

func newEvents(sink *Sink) *Events {
	return &Events{
		pool: sync.Pool{},
		sink: sink,
	}
}

type Events struct {
	pool sync.Pool
	sink *Sink
}

func (events *Events) Get(level Level, fields Fields) (v Event) {
	c := events.pool.Get()
	if c == nil {
		v = &event{
			sink:  events.sink,
			reset: events.release,
			entry: Entry{
				Level:  level,
				Occur:  time.Now(),
				Fields: fields,
			},
		}
		return
	}
	vv := c.(*event)
	vv.entry.Level = level
	vv.entry.Occur = time.Now()
	vv.entry.Fields = fields
	v = vv
	return
}

func (events *Events) Shutdown(ctx context.Context) (err error) {
	err = events.sink.Shutdown(ctx)
	return
}

func (events *Events) release(v *event) {
	v.entry.Message = ""
	v.entry.Fields = nil
	v.entry.Cause = nil
	v.entry.Occur = time.Time{}
	v.entry.Caller.Fn = ""
	v.entry.Caller.File = ""
	v.entry.Caller.Line = 0
	events.pool.Put(v)
}
