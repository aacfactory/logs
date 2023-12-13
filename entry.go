package logs

import (
	"encoding/json"
	"github.com/valyala/bytebufferpool"
	"strconv"
	"time"
	"unsafe"
)

type Caller struct {
	Fn   string
	File string
	Line int
}

func (caller Caller) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lb)
	// fn
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(fnIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(caller.Fn)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// file
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(fileIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(caller.File)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// lind
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(lineIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.WriteString(strconv.Itoa(caller.Line))
	_, _ = buf.Write(rb)
	s := buf.String()
	p = unsafe.Slice(unsafe.StringData(s), len(s))
	bytebufferpool.Put(buf)
	return
}

type Error interface {
	json.Marshaler
	Error() string
}

type jsonError struct {
	cause error
}

func (err jsonError) MarshalJSON() ([]byte, error) {
	return append(append([]byte(`{"message": "`), err.cause.Error()...), []byte(`"}`)...), nil
}

func (err jsonError) String() string {
	return err.cause.Error()
}

func (err jsonError) Error() string {
	return err.cause.Error()
}

func (err jsonError) Unwrap() error {
	return err.cause
}

type Entry struct {
	Level   Level
	Occur   time.Time
	Message string
	Fields  Fields
	Cause   Error
	Caller  Caller
}

func (e Entry) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lb)
	// Level
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(levelIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(e.Level.String())
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// occur
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(occurIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(e.Occur.Format(time.RFC3339))
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// Message
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(messageIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(e.Message)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	// Fields
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(fieldsIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	metaBytes, _ := e.Fields.MarshalJSON()
	_, _ = buf.Write(metaBytes)
	_, _ = buf.Write(comma)
	// Caller
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(callerIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	stacktraceBytes, _ := e.Caller.MarshalJSON()
	_, _ = buf.Write(stacktraceBytes)
	// Cause
	if e.Cause != nil {
		_, _ = buf.Write(comma)
		_, _ = buf.Write(dqm)
		_, _ = buf.Write(causeIdent)
		_, _ = buf.Write(dqm)
		_, _ = buf.Write(colon)
		causeBytes, _ := e.Cause.MarshalJSON()
		_, _ = buf.Write(causeBytes)
	}
	_, _ = buf.Write(rb)
	s := buf.String()
	p = unsafe.Slice(unsafe.StringData(s), len(s))
	bytebufferpool.Put(buf)
	return
}

type Writer interface {
	Write(entry Entry)
	Close() (err error)
}
