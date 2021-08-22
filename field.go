package logs

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

func withEventField(core *zerolog.Event, key string, value interface{}) (n *zerolog.Event) {
	key = strings.TrimSpace(key)
	if key == "" {
		n = core
		return
	}
	if value == nil {
		n = core
		return
	}
	switch value.(type) {
	case string:
		n = core.Str(key, value.(string))
	case int:
		n = core.Int(key, value.(int))
	case int8:
		n = core.Int8(key, value.(int8))
	case int16:
		n = core.Int16(key, value.(int16))
	case int32:
		n = core.Int32(key, value.(int32))
	case int64:
		n = core.Int64(key, value.(int64))
	case uint:
		n = core.Uint(key, value.(uint))
	case uint8:
		n = core.Uint8(key, value.(uint8))
	case uint16:
		n = core.Uint16(key, value.(uint16))
	case uint32:
		n = core.Uint32(key, value.(uint32))
	case uint64:
		n = core.Uint64(key, value.(uint64))
	case float32:
		n = core.Float32(key, value.(float32))
	case float64:
		n = core.Float64(key, value.(float64))
	case bool:
		n = core.Bool(key, value.(bool))
	case []byte:
		n = core.Bytes(key, value.([]byte))
	case time.Time:
		o := value.(time.Time)
		n = core.Str(key, o.Format(time.RFC3339))
	case time.Duration:
		o := value.(time.Duration)
		n = core.Str(key, o.String())
	case error:
		o := value.(error)
		n = core.Err(o)
	default:
		data, encodeErr := json.Marshal(value)
		if encodeErr == nil {
			n = core.RawJSON(key, data)
		}
	}

	return
}

func withContextField(core zerolog.Context, key string, value interface{}) (n zerolog.Context) {
	key = strings.TrimSpace(key)
	if key == "" {
		n = core
		return
	}
	if value == nil {
		n = core
		return
	}
	switch value.(type) {
	case string:
		n = core.Str(key, value.(string))
	case int:
		n = core.Int(key, value.(int))
	case int8:
		n = core.Int8(key, value.(int8))
	case int16:
		n = core.Int16(key, value.(int16))
	case int32:
		n = core.Int32(key, value.(int32))
	case int64:
		n = core.Int64(key, value.(int64))
	case uint:
		n = core.Uint(key, value.(uint))
	case uint8:
		n = core.Uint8(key, value.(uint8))
	case uint16:
		n = core.Uint16(key, value.(uint16))
	case uint32:
		n = core.Uint32(key, value.(uint32))
	case uint64:
		n = core.Uint64(key, value.(uint64))
	case float32:
		n = core.Float32(key, value.(float32))
	case float64:
		n = core.Float64(key, value.(float64))
	case bool:
		n = core.Bool(key, value.(bool))
	case []byte:
		n = core.Bytes(key, value.([]byte))
	case time.Time:
		o := value.(time.Time)
		n = core.Str(key, o.Format(time.RFC3339))
	case time.Duration:
		o := value.(time.Duration)
		n = core.Str(key, o.String())
	case error:
		o := value.(error)
		n = core.Err(o)
	default:
		data, encodeErr := json.Marshal(value)
		if encodeErr == nil {
			n = core.RawJSON(key, data)
		}
	}

	return
}
