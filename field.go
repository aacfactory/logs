package logs

import (
	"fmt"
	"github.com/valyala/bytebufferpool"
	"time"
)

type Field struct {
	Key   string
	Value any
}

func (field Field) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lb)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(keyIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(field.Key)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(comma)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(valueIdent)
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(colon)
	_, _ = buf.Write(dqm)
	_, _ = buf.WriteString(string(field.ValueBytes()))
	_, _ = buf.Write(dqm)
	_, _ = buf.Write(rb)
	p = buf.Bytes()
	bytebufferpool.Put(buf)
	return
}

func (field Field) ValueBytes() (p []byte) {
	vp := ""
	switch value := field.Value.(type) {
	case nil:
		vp = "<nil>"
		break
	case string:
		vp = value
		break
	case []byte:
		vp = string(value)
		break
	case time.Time:
		vp = value.Format(time.RFC3339)
		break
	default:
		vp = fmt.Sprintf("%v", value)
		break
	}
	p = []byte(vp)
	return
}

type Fields []Field

func (fields Fields) Add(key string, value any) Fields {
	for i, field := range fields {
		if field.Key == key {
			field.Value = value
			fields[i] = field
			return fields
		}
	}
	return append(fields, Field{
		Key:   key,
		Value: value,
	})
}

func (fields Fields) MarshalJSON() (p []byte, err error) {
	buf := bytebufferpool.Get()
	_, _ = buf.Write(lqb)
	if len(fields) == 0 {
		_, _ = buf.Write(rqb)
		p = buf.Bytes()
		bytebufferpool.Put(buf)
		return
	}
	for i, field := range fields {
		if i > 0 {
			_, _ = buf.Write(comma)
		}
		b, _ := field.MarshalJSON()
		_, _ = buf.Write(b)
	}
	_, _ = buf.Write(rqb)
	p = buf.Bytes()
	bytebufferpool.Put(buf)
	return
}
