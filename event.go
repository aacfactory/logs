package logs

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"runtime"
	"strings"
)

type event struct {
	core *zerolog.Event
}

func (e *event) Message(message string) {
	e.core.Msg(message)
}

func (e *event) Cause(err error) Event {
	if err == nil {
		return e
	}

	if isCodeError(err) {
		data, encodeErr := json.Marshal(err)
		if encodeErr != nil {
			e.core = e.core.Err(err)
		} else {
			e.core = e.core.RawJSON(zerolog.ErrorFieldName, data)
		}
	} else {
		e.core = e.core.Err(err)
	}

	return e
}

func (e *event) Caller() Event {
	return e.CallerWithSkip(2)
}

func (e *event) CallerWithSkip(skip int) Event {
	_, file, line, ok := runtime.Caller(skip)
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
	e.core = e.core.Str(zerolog.CallerFieldName, fmt.Sprintf("%s:%d", file, line))
	return e
}

func (e *event) With(key string, value interface{}) Event {
	e.core = withEventField(e.core, key, value)
	return e
}

func isCodeError(err error) (ok bool) {
	_, ok = err.(codeError)
	return
}

type codeError interface {
	Id() string
	Code() int
	Name() string
	Message() string
	Stacktrace() (fn string, file string, line int)
}
