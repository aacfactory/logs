package logs

import (
	"github.com/aacfactory/errors"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Field struct {
	Key   string
	Value interface{}
}

func F(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func newCodeErrorMarshalLogObject(err errors.CodeError) codeErrorMarshalLogObject {
	return codeErrorMarshalLogObject{
		err: err,
	}
}

type codeErrorMarshalLogObject struct {
	err errors.CodeError
}

func (o codeErrorMarshalLogObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	err := o.err
	fn, file, line := err.Stacktrace()

	enc.AddString("id", err.Id())
	enc.AddString("code", err.Code())
	enc.AddInt("failureCode", err.FailureCode())
	if err.Meta() != nil && len(err.Meta()) > 0 {
		meta := codeErrorMetaMarshalLogObject{
			meta: err.Meta(),
		}
		_ = enc.AddObject("meta", meta)
	}
	enc.AddString("message", err.Message())
	_ = enc.AddObject("stacktrace", codeErrorStacktraceMarshalLogObject{
		fn:   fn,
		file: file,
		line: line,
	})
	return nil
}

type codeErrorStacktraceMarshalLogObject struct {
	fn   string
	file string
	line int
}

func (o codeErrorStacktraceMarshalLogObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("fn", o.fn)
	enc.AddString("file", o.file)
	enc.AddInt("line", o.line)
	return nil
}

type codeErrorMetaMarshalLogObject struct {
	meta map[string][]string
}

func (o codeErrorMetaMarshalLogObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for key, values := range o.meta {
		enc.AddString(key, strings.Join(values, ","))
	}
	return nil
}
