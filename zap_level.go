package logs

import (
	"os"
	"strings"

	"go.uber.org/zap/zapcore"
)

func zapLogLevel(name LogLevel) zapcore.Level {
	value := strings.ToLower(string(name))
	switch value {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

var (
	_zapLoglevelToColor = map[zapcore.Level]Color{
		zapcore.DebugLevel:  magenta,
		zapcore.InfoLevel:   blue,
		zapcore.WarnLevel:   yellow,
		zapcore.ErrorLevel:  red,
		zapcore.DPanicLevel: red,
		zapcore.PanicLevel:  red,
		zapcore.FatalLevel:  red,
	}
	gopath = os.Getenv("GOPATH")
)

func zapLogCapitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s := ""
	color, ok := _zapLoglevelToColor[l]
	if ok {
		s = color.Add(l.CapitalString())
	}
	enc.AppendString(s)
}

func zapLogFullCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if gopath == "" {
		enc.AppendString(caller.FullPath())
	} else {
		enc.AppendString(caller.FullPath()[len(gopath)+5:])
	}
}
