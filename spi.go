package logs

import (
	"log"
)

const (
	ConsoleFormatter = Formatter("console")
	JsonFormatter    = Formatter("json")

	DebugLevel = Level(0)
	InfoLevel  = Level(1)
	WarnLevel  = Level(2)
	ErrorLevel = Level(3)

	defaultLogName = "AACFACTORY"

	defaultTimeFormatter = "[3:04:05 PM]"
)

type Level int

type Formatter string

type Logger interface {
	With(key string, value interface{}) Logger
	DebugEnabled() (ok bool)
	Debug() (event Event)
	InfoEnabled() (ok bool)
	Info() (event Event)
	WarnEnabled() (ok bool)
	Warn() (event Event)
	ErrorEnabled() (ok bool)
	Error() (event Event)
}

type Event interface {
	Message(message string)
	Cause(err error) Event
	Caller() Event
	CallerWithSkip(skip int) Event
	With(key string, value interface{}) Event
}

type loggerWriter struct {
	level      Level
	logger     Logger
	withCaller bool
}

func (w *loggerWriter) Write(p []byte) (n int, err error) {
	if p == nil || len(p) == 0 {
		return
	}
	if p[len(p)-1] == '\n' {
		p = p[0 : len(p)-1]
	}
	switch w.level {
	case DebugLevel:
		if w.logger.DebugEnabled() {
			e := w.logger.Debug()
			if w.withCaller {
				e = e.CallerWithSkip(4)
			}
			e.Message(string(p))
			n = len(p)
		}
		break
	case InfoLevel:
		if w.logger.InfoEnabled() {
			e := w.logger.Info()
			if w.withCaller {
				e = e.CallerWithSkip(4)
			}
			e.Message(string(p))
			n = len(p)
		}
		break
	case WarnLevel:
		if w.logger.WarnEnabled() {
			e := w.logger.Warn()
			if w.withCaller {
				e = e.CallerWithSkip(4)
			}
			e.Message(string(p))
			n = len(p)
		}
		break
	case ErrorLevel:
		if w.logger.ErrorEnabled() {
			e := w.logger.Error()
			if w.withCaller {
				e = e.CallerWithSkip(4)
			}
			e.Message(string(p))
			n = len(p)
		}
		break
	default:
		break
	}
	return
}

func MapToLogger(logger Logger, level Level, withCaller bool) (v *log.Logger) {
	v = log.New(&loggerWriter{
		level:      level,
		logger:     logger,
		withCaller: withCaller,
	}, "", 0)
	return
}
