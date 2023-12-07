package logs

import "log"

func ConvertToStandardLogger(logger Logger, level Level, withCaller bool) (v *log.Logger) {
	v = log.New(&loggerWriter{
		level:      level,
		logger:     logger,
		withCaller: withCaller,
	}, "", 0)
	return
}

type loggerWriter struct {
	level      Level
	logger     Logger
	withCaller bool
}

func (w *loggerWriter) Write(p []byte) (n int, err error) {
	pLen := len(p)
	if pLen == 0 {
		return
	}
	if p[pLen-1] == '\n' {
		p = p[0 : pLen-1]
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
