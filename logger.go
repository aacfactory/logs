package logs

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"os"
	"strings"
)

type Options struct {
	name          string
	writer        io.Writer
	level         Level
	formatter     Formatter
	noColor       bool
	timeFormatter string
}

type Option = func(options *Options) (err error)

func Name(name string) Option {
	name = strings.TrimSpace(name)
	return func(options *Options) (err error) {
		options.name = name
		return
	}
}

func Writer(w io.Writer) Option {
	return func(options *Options) (err error) {
		if w == nil {
			w = io.Discard
		}
		options.writer = w
		return
	}
}

func WithLevel(value Level) Option {
	return func(options *Options) (err error) {
		if value < 0 || value > 3 {
			value = DebugLevel
		}
		options.level = value
		return
	}
}

func Color(enable bool) Option {
	return func(options *Options) (err error) {
		options.noColor = !enable
		return
	}
}

func TimeFormatter(formatter string) Option {
	return func(options *Options) (err error) {
		if formatter == "" {
			formatter = defaultTimeFormatter
		}
		options.timeFormatter = formatter
		return
	}
}

func WithFormatter(value Formatter) Option {
	return func(options *Options) (err error) {
		if value != ConsoleFormatter && value != JsonFormatter {
			value = JsonFormatter
		}
		options.formatter = value
		return
	}
}

func New(options ...Option) (log Logger, err error) {

	opt := &Options{
		name:          defaultLogName,
		writer:        os.Stdout,
		level:         DebugLevel,
		formatter:     ConsoleFormatter,
		noColor:       false,
		timeFormatter: defaultTimeFormatter,
	}

	if options != nil {
		for _, option := range options {
			optErr := option(opt)
			if optErr != nil {
				err = fmt.Errorf("create log failed, %v", optErr)
				return
			}
		}
	}

	core := zerolog.New(opt.writer)

	if opt.formatter == ConsoleFormatter {
		core = core.Output(buildConsoleWriter(opt.writer, !opt.noColor, opt.timeFormatter))
	}

	core = core.Level(zerolog.Level(opt.level))

	core = core.With().Timestamp().Str("app", opt.name).Logger()

	log = &logger{
		level: opt.level,
		core:  core,
	}

	return
}

type logger struct {
	level Level
	core  zerolog.Logger
}

func (l *logger) With(key string, value interface{}) (n Logger) {
	ctx := l.core.With()
	ctx = withContextField(ctx, key, value)
	n = &logger{
		level: l.level,
		core:  ctx.Logger(),
	}
	return l
}

func (l *logger) DebugEnabled() (ok bool) {
	ok = l.level <= DebugLevel
	return
}

func (l *logger) Debug() (e Event) {
	e = &event{
		core: l.core.Debug(),
	}
	return
}

func (l *logger) InfoEnabled() (ok bool) {
	ok = l.level <= InfoLevel
	return
}

func (l *logger) Info() (e Event) {
	e = &event{
		core: l.core.Info(),
	}
	return
}

func (l *logger) WarnEnabled() (ok bool) {
	ok = l.level <= WarnLevel
	return
}

func (l *logger) Warn() (e Event) {
	e = &event{
		core: l.core.Warn(),
	}
	return
}

func (l *logger) ErrorEnabled() (ok bool) {
	ok = l.level <= ErrorLevel
	return
}

func (l *logger) Error() (e Event) {
	e = &event{
		core: l.core.Error(),
	}
	return
}
