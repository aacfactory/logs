package logs

import (
	"fmt"
	"go.uber.org/zap"
)

var zapLog *zap.SugaredLogger

func Log() Logs {
	return zapLog
}

func With(log Logs, fields ...Field) (_log Logs, err error) {
	if fields == nil || len(fields) == 0 {
		err = fmt.Errorf("logs with fields failed, fields are empty")
		return
	}

	sLog, ok := log.(*zap.SugaredLogger)
	if !ok {
		err = fmt.Errorf("logs with fields failed, it is not *zap.SugaredLogger")
		return
	}

	kvs := make([]zap.Field, 0, 1)
	for _, field := range fields {
		kvs = append(kvs, zap.Any(field.Key, field.Value))
	}

	_log = sLog.Desugar().With(kvs...).Sugar()
	return
}

func Sync() {
	_ = zapLog.Sync()
}
