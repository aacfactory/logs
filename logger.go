package logs

import (
	"context"
	"errors"
	"time"
)

type Options struct {
	level                  Level
	disableConsole         bool
	consoleWriterFormatter ConsoleWriterFormatter
	consoleWriterOutType   ConsoleWriterOutType
	buffer                 int
	discardLevel           Level
	consumes               int
	writers                []Writer
	sendTimeout            time.Duration
	shutdownTimeout        time.Duration
}

type Option = func(options *Options) (err error)

func WithWriter(writers ...Writer) Option {
	return func(options *Options) (err error) {
		for _, writer := range writers {
			if writer == nil {
				err = errors.New("one of writer is nil")
				return
			}
			options.writers = append(options.writers, writer)
		}
		return
	}
}

func WithLevel(value Level) Option {
	return func(options *Options) (err error) {
		if value < DebugLevel || value > ErrorLevel {
			err = errors.New("invalid level")
			return
		}
		options.level = value
		return
	}
}

func WithTimeoutDiscardLevel(value Level) Option {
	return func(options *Options) (err error) {
		if value < DebugLevel || value > ErrorLevel {
			err = errors.New("invalid level")
			return
		}
		options.discardLevel = value
		return
	}
}

func WithConsoleWriterFormatter(formatter ConsoleWriterFormatter) Option {
	return func(options *Options) (err error) {
		options.consoleWriterFormatter = formatter
		return
	}
}

func WithConsoleWriterOutType(typ ConsoleWriterOutType) Option {
	return func(options *Options) (err error) {
		options.consoleWriterOutType = typ
		return
	}
}

func WithBuffer(size int) Option {
	return func(options *Options) (err error) {
		if size < 1 {
			err = errors.New("invalid buffer size")
			return
		}
		options.buffer = size
		return
	}
}

func WithConsumes(consumes int) Option {
	return func(options *Options) (err error) {
		if consumes < 1 {
			err = errors.New("invalid consumes")
			return
		}
		options.consumes = consumes
		return
	}
}

func WithSendTimeout(timeout time.Duration) Option {
	return func(options *Options) (err error) {
		if timeout < 1 {
			err = errors.New("invalid timeout")
			return
		}
		options.sendTimeout = timeout
		return
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(options *Options) (err error) {
		if timeout < 1 {
			err = errors.New("invalid timeout")
			return
		}
		options.shutdownTimeout = timeout
		return
	}
}

func DisableConsoleWriter() Option {
	return func(options *Options) (err error) {
		options.disableConsole = true
		return
	}
}

func New(options ...Option) (log Logger, err error) {
	// options
	opt := &Options{
		level:                  InfoLevel,
		disableConsole:         false,
		consoleWriterFormatter: TextFormatter,
		consoleWriterOutType:   StdOut,
		buffer:                 0,
		discardLevel:           DebugLevel,
		consumes:               1,
		writers:                nil,
		sendTimeout:            0,
		shutdownTimeout:        0,
	}
	if options != nil {
		errs := 0
		optErrs := errors.New("new logger failed")
		for _, option := range options {
			optErr := option(opt)
			if optErr != nil {
				optErrs = errors.Join(optErrs, optErr)
				errs++
			}
		}
		if errs > 0 {
			err = optErrs
			return
		}
	}
	if !opt.disableConsole {
		// set console writer at end
		console := NewConsoleWriter(opt.consoleWriterFormatter, opt.consoleWriterOutType)
		opt.writers = append(opt.writers, console)
	}
	if len(opt.writers) == 0 {
		err = errors.New("new logger failed, writers is required")
		return
	}
	// sink
	sink := newSink(opt.level, opt.discardLevel, opt.consumes, opt.buffer, opt.sendTimeout, opt.shutdownTimeout, opt.writers)
	// events
	events := newEvents(sink)
	// logger
	log = &logger{
		level:  opt.level,
		events: events,
		fields: nil,
	}

	return
}

type Logger interface {
	With(key string, value any) Logger
	DebugEnabled() (ok bool)
	Debug() (event Event)
	InfoEnabled() (ok bool)
	Info() (event Event)
	WarnEnabled() (ok bool)
	Warn() (event Event)
	ErrorEnabled() (ok bool)
	Error() (event Event)
	Shutdown(ctx context.Context) (err error)
}

type logger struct {
	level  Level
	events *Events
	fields Fields
}

func (log *logger) With(key string, value any) (n Logger) {
	n = &logger{
		level:  log.level,
		events: log.events,
		fields: log.fields.Add(key, value),
	}
	return
}

func (log *logger) DebugEnabled() (ok bool) {
	ok = log.level <= DebugLevel
	return
}

func (log *logger) Debug() (e Event) {
	if log.DebugEnabled() {
		e = log.events.Get(DebugLevel, log.fields)
		return
	}
	e = empty
	return
}

func (log *logger) InfoEnabled() (ok bool) {
	ok = log.level <= InfoLevel
	return
}

func (log *logger) Info() (e Event) {
	if log.InfoEnabled() {
		e = log.events.Get(InfoLevel, log.fields)
		return
	}
	e = empty
	return
}

func (log *logger) WarnEnabled() (ok bool) {
	ok = log.level <= WarnLevel
	return
}

func (log *logger) Warn() (e Event) {
	if log.WarnEnabled() {
		e = log.events.Get(WarnLevel, log.fields)
		return
	}
	e = empty
	return
}

func (log *logger) ErrorEnabled() (ok bool) {
	ok = log.level <= ErrorLevel
	return
}

func (log *logger) Error() (e Event) {
	if log.ErrorEnabled() {
		e = log.events.Get(ErrorLevel, log.fields)
		return
	}
	e = empty
	return
}

func (log *logger) Shutdown(ctx context.Context) (err error) {
	if len(log.fields) > 0 {
		return
	}
	err = log.events.Shutdown(ctx)
	return
}
